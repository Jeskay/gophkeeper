package main

import "gophkeeper/internal/client"

func main() {
	app := client.NewApp()
	app.Start()
}
