package main

import (
	"context"
	"file-transfer/utils/server"
)

func main() {
	// logic to start
	ctx, cancel := context.WithCancel(context.Background())
	go server.ListenForBroadCasts(ctx)

	// logic to run
	server.StartServer()

	// logic to shutdown
	cancel()
}
