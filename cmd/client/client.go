package client

import (
	config "file-transfer/configs"
	"file-transfer/utils/client"
	"file-transfer/utils/utils"
	"fmt"
)

func showMenuSendFile() {
	fmt.Println("\nAvailable Commands For Send File:")
	fmt.Println("1. Send File")
	fmt.Println("0. Return To Menu")
}

func SendFile() {
	config.LoadEnv()

	// logic to get ip for share
	addr, err := utils.GetLocalIP()
	if err != nil {
		fmt.Printf("Error getting local IP: %v\n", err)
		return
	}

	for {
		showMenuSendFile()

		var sendFileCommand string
		fmt.Print("Enter command number: ")
		fmt.Scanln(&sendFileCommand)

		switch sendFileCommand {
		case "1":
			client.SendBroadCasts()
			client.StartClient(addr)
		case "0":
			fmt.Println("Returning to main menu...")
			utils.ClearScreen()
			return
		default:
			fmt.Println("Invalid command in Send File menu. Please try again.")
			utils.ClearScreen()
		}
	}
}
