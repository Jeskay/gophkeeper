package main

import (
	"gophkeeper/config"
	"gophkeeper/internal/client"
)

func main() {
	app := client.NewApp(&config.ClientConfig{
		GRPCAddress: config.Address{
			Host: "localhost",
			Port: "8080",
		},
	})
	app.Start()
}
