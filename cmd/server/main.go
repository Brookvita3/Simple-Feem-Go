package main

import (
	"context"
	config "file-transfer/configs"
	"file-transfer/utils/server"
)

func main() {

	config.LoadEnv()

	// logic to start
	ctx, cancel := context.WithCancel(context.Background())
	go server.ListenForBroadCasts(ctx)

	// logic to run
	server.StartServer()

	// logic to shutdown
	cancel()
}
