package main

import (
	"context"
	"file-transfer/cmd/client"
	"file-transfer/utils/server"
	utils "file-transfer/utils/utils"
	"fmt"
)

func showMenu() {
	fmt.Println("\nAvailable Commands:")
	fmt.Println("1. Start")
	fmt.Println("2. Send File")
	fmt.Println("5. Exit")
}

func main() {

	// logic to start
	ctx, cancel := context.WithCancel(context.Background())

	// logic to run
	for {
		showMenu()

		// Get user input
		var command string
		fmt.Print("Enter command number: ")
		fmt.Scan(&command)

		// Handle commands
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
