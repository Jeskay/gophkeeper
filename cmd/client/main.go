package main

import (
	"flag"
	"log"
	"net"
	"os"
	"path"

	"github.com/caarlos0/env/v11"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gophkeeper/api/protos"
	"gophkeeper/config"
	authControl "gophkeeper/internal/client/authentication/control"
	authPresent "gophkeeper/internal/client/authentication/presentation"
	downloadControl "gophkeeper/internal/client/download/control"
	downloadPresent "gophkeeper/internal/client/download/presentation"
	menuControl "gophkeeper/internal/client/menu/control"
	menuPresent "gophkeeper/internal/client/menu/presentation"
	uploadControl "gophkeeper/internal/client/upload/control"
	uploadPresent "gophkeeper/internal/client/upload/presentation"
)

var cfg *config.ClientConfig

func main() {
	addr := net.JoinHostPort(cfg.GRPCAddress.Host, cfg.GRPCAddress.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not establish grpc connection", err)
	}
	defer conn.Close()

	grpcClient := pb.NewGophkeeperClient(conn)

	authController := authControl.NewController(grpcClient)
	menuController := menuControl.NewController(authController)
	uploadController := uploadControl.NewController(grpcClient, menuController)
	downloadController := downloadControl.NewController(cfg.DownloadDirectory, grpcClient, menuController)

	authPresenter := authPresent.NewPresenter(authController)
	uploadPresenter := uploadPresent.NewPresenter(uploadController)
	downloadPresenter := downloadPresent.NewPresenter(downloadController)
	menuPresenter := menuPresent.NewPresenter(menuController, authPresenter, uploadPresenter, downloadPresenter)

	m := menuPresenter.CreateMenu()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal("unexpected error", err)
	}
}

func init() {
	cfg = loadParams()
	if err := godotenv.Load(); err != nil {
		log.Println("no configuration file was found")
	}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	saveDir, err := os.UserHomeDir()
	if err != nil {
		saveDir = ""
	}
	saveDir = path.Join(saveDir, "Downloads")
	cfg.DownloadDirectory = saveDir
}

func loadParams() *config.ClientConfig {
	var cfg = &config.ClientConfig{}
	flag.StringVar(&cfg.GRPCAddress.Host, "host", "localhost", "grpc server hostname")
	flag.StringVar(&cfg.GRPCAddress.Port, "port", "8080", "grpc server port")
	flag.Parse()
	return cfg
}
