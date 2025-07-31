package client

import (
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
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

type App struct {
	conf *config.ClientConfig
}

func NewApp(conf *config.ClientConfig) *App {
	return &App{conf: conf}
}

func (a *App) Start() {
	addr := net.JoinHostPort(a.conf.GRPCAddress.Host, a.conf.GRPCAddress.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not establish grpc connection", err)
	}
	grpcClient := pb.NewGophkeeperClient(conn)

	authController := authControl.NewController(grpcClient)
	menuController := menuControl.NewController(authController)
	uploadController := uploadControl.NewController(grpcClient, menuController)
	downloadController := downloadControl.NewController(grpcClient)

	authPresenter := authPresent.NewPresenter(authController)
	uploadPresenter := uploadPresent.NewPresenter(uploadController)
	downloadPresenter := downloadPresent.NewPresenter(downloadController)
	menuPresenter := menuPresent.NewPresenter(menuController, authPresenter, uploadPresenter, downloadPresenter)

	m := menuPresenter.CreateMenu()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal("unexpected error", err)
	}
}
