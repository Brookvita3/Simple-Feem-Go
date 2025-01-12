package main

import (
	"context"
	"file-transfer/cmd/client"
	config "file-transfer/configs"
	"file-transfer/utils/server"
	"file-transfer/utils/utils"
	"fmt"
)

func showMenu() {
	fmt.Println("\nAvailable Commands:")
	fmt.Println("1. Start")
	fmt.Println("2. Send File")
	fmt.Println("5. Exit")
}

func main() {

	config.LoadEnv()

	// logic to start
	ctx, cancel := context.WithCancel(context.Background())

	// logic to run
	for {
		showMenu()

		command := utils.GetInput()

		switch command {
		case "1":
			go server.StartServer()
			go server.ListenForBroadCasts(ctx)
		case "2":
			client.SendFile()
			utils.ClearScreen()
		case "5":
			server.ShutdownServer(cancel)
			utils.ClearScreen()
			return
		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}
