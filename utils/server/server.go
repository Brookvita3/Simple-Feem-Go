package server

import (
	"bufio"
	"context"
	config "file-transfer/configs"
	"file-transfer/utils/utils"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func ReceiveFileChunks(conn net.Conn, fileChannels utils.FileChannels) {

	defer conn.Close()

	for {
		buffer := make([]byte, fileChannels.ChunkSize)
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fileChannels.ErrorChan <- err
			break
		}
		fileChannels.DataChan <- buffer[:n]
	}

	fmt.Println("Receive Successfully")
	close(fileChannels.DataChan)
}

func WriteFileChunks(filePath string, fileChannels utils.FileChannels) {

	file, err := os.Create(filePath)
	if err != nil {
		fileChannels.ErrorChan <- err
		return
	}
	defer file.Close()

	writter := bufio.NewWriterSize(file, fileChannels.ChunkSize)
	fileChannels.WriteChunks(writter)

	fmt.Println("Write Successfully")
	close(fileChannels.ErrorChan)
}

func ListenForBroadCasts(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", ":"+config.Config.BROADCAST_PORT)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening for broadcast:", err)
		return err
	}
	defer conn.Close()

	fmt.Println("Listening for broadcasts on", config.Config.BROADCAST_PORT)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Listener stopped.")
			return nil
		default:
			buffer := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Timeout reached, continue checking for ctx.Done()
					continue
				}
				fmt.Println("Error reading from UDP:", err)
				continue
			}

			fmt.Printf("Received message: %s from %s\n", string(buffer[:n]), remoteAddr)
		}
	}

}

func StartServer() {

	// Listen for incoming connections
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			fmt.Printf("Client closed connection: %s\n", conn.LocalAddr().String())
			break
		}

		conn.Write([]byte("Welcome to the server!\n"))

		path := "../files/receive/file1.pdf"
		fileChannels := *utils.NewFileChannels(
			make(chan []byte),
			make(chan error),
			config.Config.CHUNK_SIZE,
		)

		go ReceiveFileChunks(conn, fileChannels)
		go WriteFileChunks(path, fileChannels)

		// Monitor for errors
		for err := range fileChannels.ErrorChan {
			fmt.Printf("Error: %v\n", err)
		}

	}
}

func ShutdownServer(cancel context.CancelFunc) {
	fmt.Println("Shutting down the server...")
	cancel()
	fmt.Println("Server has been shut down gracefully.")
}
