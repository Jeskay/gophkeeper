package main

import (
	"context"
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpckit "github.com/go-kit/kit/transport/grpc"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	proto "gophkeeper/api/protos"
	"gophkeeper/config"
	"gophkeeper/internal/server/auth"
	"gophkeeper/internal/server/db"
	"gophkeeper/internal/server/endpoints"
	"gophkeeper/internal/server/transport"
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
	var authService auth.Service
	{
		authService = auth.NewService(dbService, []byte(cfg.SecretKey), time.Second*5)
	}
	eps := endpoints.NewEndpointList(authService)
	grpcServer := transport.NewGRPCServer(eps)

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		zapL.Fatal("failed to start grpc server", zap.Error(err))
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		baseServer := grpc.NewServer(grpc.UnaryInterceptor(grpckit.Interceptor))
		reflection.Register(baseServer)
		proto.RegisterGophkeeperServer(baseServer, grpcServer)
		if err := baseServer.Serve(grpcListener); err != nil {
			grpcListener.Close()
			zapL.Fatal("internal server error", zap.Error(err))
		}
	}()
	<-c
	zapL.Info("server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := grpcListener.Close(); err != nil {
		zapL.Error("failed to close grpc connection", zap.Error(err))
	}

	<-ctx.Done()
	zapL.Info("server exited")
}
