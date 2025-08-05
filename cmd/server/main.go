package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	proto "gophkeeper/api/protos"
	"gophkeeper/config"
	"gophkeeper/internal/server/auth"
	"gophkeeper/internal/server/db"
	"gophkeeper/internal/server/file"
	"gophkeeper/internal/server/keeper"
	"gophkeeper/internal/server/transport"
	"gophkeeper/internal/server/transport/interceptors"
)

var (
	cfg config.ServerConfig
)

func main() {
	grpcAddr := net.JoinHostPort(cfg.GRPCAddress.Host, cfg.GRPCAddress.Port)
	connectionStr := cfg.GetDSN()
	zapL, err := zap.NewProduction()
	if err != nil {
		return
	}
	database, err := sql.Open("pgx", connectionStr)
	if err != nil {
		zapL.Fatal("failed to connect to database")
		return
	}
	dbService, err := db.NewService(database, zapslog.NewHandler(zapL.Core()))
	if err != nil {
		zapL.Error("failed to initialize db", zap.Error(err))
		return
	}
	fileService := file.NewService("storage")

	authService := auth.NewService(dbService, []byte(cfg.SecretKey), time.Hour*5)
	keeperService := keeper.NewKeeperService(dbService, fileService)
	grpcServer := transport.NewGRPCServer(authService, keeperService)

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		zapL.Fatal("failed to start grpc server", zap.Error(err))
	}
	baseServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.NewAuthUnaryInterceptor(authService)),
		grpc.StreamInterceptor(interceptors.NewAuthStreamInterceptor(authService)),
	)
	reflection.Register(baseServer)
	proto.RegisterGophkeeperServer(baseServer, grpcServer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-c
		zapL.Info("server shutting down...")
		baseServer.GracefulStop()
		wg.Done()
	}()
	if err := baseServer.Serve(grpcListener); err != nil {
		zapL.Fatal("internal server error", zap.Error(err))
	}
	wg.Wait()
	zapL.Info("server exited")
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no configuration file was found")
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
}
