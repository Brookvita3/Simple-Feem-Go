package main

import (
	"file-transfer/utils/client"
	"file-transfer/utils/utils"
	"fmt"
)

func main() {
	// logic to get ip for share
	addr, err := utils.GetLocalIP()
	if err != nil {
		fmt.Printf("Error getting local IP: %v\n", err)
		return
	}

	client.SendBroadCasts()
	client.StartClient(addr)
}
