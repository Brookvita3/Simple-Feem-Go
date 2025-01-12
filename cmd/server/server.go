package server

import (
	"context"
	config "file-transfer/configs"
	"file-transfer/utils/server"
	"fmt"
)

func showMenu() {
	fmt.Println("\nAvailable Commands:")
	fmt.Println("1. Start")
	fmt.Println("5. Exit")
}

func Main() {

	config.LoadEnv()

	// logic to start
	ctx, cancel := context.WithCancel(context.Background())

	// logic to run
	for {
		showMenu()

		// Get user input
		var command string
		fmt.Print("Enter command number: ")
		fmt.Scanln(&command)

		// Handle commands
		switch command {
		case "1":
			go server.StartServer()
			go server.ListenForBroadCasts(ctx)
		case "5":
			server.ShutdownServer(cancel)
			return
		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}
