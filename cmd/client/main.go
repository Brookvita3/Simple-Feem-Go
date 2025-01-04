package main

import (
	"file-transfer/utils/client"
)

func main() {
	// logic to get ip for share
	addr := "192.168.1.6"
	client.StartClient(addr)
}
