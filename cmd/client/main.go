package main

import (
	config "file-transfer/configs"
	"file-transfer/utils/client"
	"file-transfer/utils/utils"
	"fmt"
)

func main() {
	config.LoadEnv()
	// logic to get ip for share
	addr, err := utils.GetLocalIP()
	if err != nil {
		fmt.Printf("Error getting local IP: %v\n", err)
		return
	}

	client.SendBroadCasts()

	client.StartClient(addr)
}
