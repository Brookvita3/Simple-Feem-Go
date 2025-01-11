package server

import (
	"context"
	config "file-transfer/configs"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

func ReceiveFileChunks(conn net.Conn, dataChan chan<- []byte, errorChan chan<- error) {
	defer conn.Close()

	chunkSize := config.Config.CHUNK_SIZE
	buffer := make([]byte, chunkSize)
	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			dataChan <- buffer[:n]
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			errorChan <- fmt.Errorf("error receiving chunk: %v", err)
			break
		}
	}
	close(dataChan) // Close the channel when done
}

func WriteFileChunks(filePath string, dataChan <-chan []byte, errorChan chan<- error) {
	file, err := os.Create(filePath)
	if err != nil {
		errorChan <- fmt.Errorf("error creating file: %v", err)
		return
	}
	defer file.Close()

	for chunk := range dataChan {
		_, err := file.Write(chunk)
		if err != nil {
			errorChan <- fmt.Errorf("error writing to file: %v", err)
			break
		}
	}
}

func ListenForBroadCasts(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(config.Config.BROADCAST_PORT))
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

		path := "../../files/receive/file1.pdf"
		dataChan := make(chan []byte)
		errorChan := make(chan error)

		go ReceiveFileChunks(conn, dataChan, errorChan)
		go WriteFileChunks(path, dataChan, errorChan)

		// Monitor for errors
		for err := range errorChan {
			fmt.Printf("Error: %v\n", err)
		}

	}
}
