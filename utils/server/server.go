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

var listPeers = make(chan string)

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

	writter := bufio.NewWriterSize(file, fileChannels.BufferSize)
	fileChannels.WriteChunks(writter)

	fmt.Println("Write Successfully")
	close(fileChannels.ErrorChan)
}

// sendBroadcast sends a UDP broadcast message
func SendBroadCasts(listenPort int) error {
	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:0")
	if err != nil {
		fmt.Println("Error resolving broadcast address:", err)
		return err
	}

	conn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(fmt.Sprintf("%d", listenPort)))
	if err != nil {
		fmt.Println("Error sending broadcast message:", err)
		return err
	}

	return nil
}

func ListenForBroadCasts(ctx context.Context, listPeer chan string) error {
	addr, err := net.ResolveUDPAddr("udp", ":"+config.Config.BROADCAST_PORT)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return err
	}

	broadCast, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening for broadcast:", err)
		return err
	}
	defer broadCast.Close()

	fmt.Println("Listening for broadcasts on", config.Config.BROADCAST_PORT)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Listener stopped.")
			return nil
		default:
			buffer := make([]byte, 1024)
			broadCast.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, remoteAddr, err := broadCast.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Continue listening after timeout
					continue
				}
				fmt.Println("Error reading from UDP:", err)
				continue
			}

			peerAddr := string(buffer[:n])
			listPeer <- peerAddr

			fmt.Printf("Received message: %s from %s\n", string(buffer[:n]), remoteAddr)
		}
	}

}

func receiveFile(conn net.Conn, filePath string) {

	fileChannels := *utils.NewFileChannels(
		make(chan []byte),
		make(chan error),
		config.Config.CHUNK_SIZE,
		config.Config.BUFFER_SIZE,
	)

	go ReceiveFileChunks(conn, fileChannels)
	go WriteFileChunks(filePath, fileChannels)

	for err := range fileChannels.ErrorChan {
		fmt.Printf("Error: %v\n", err)
	}
}

func listenConnection(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			fmt.Printf("Client closed connection: %s\n", conn.LocalAddr().String())
			break
		}

		// mo phong nhan file
		conn.Write([]byte("Welcome to the server!\n"))
		path := "../files/receive/file1.pdf"
		receiveFile(conn, path)
	}

}

func StartServer(ctx context.Context) {

	go ListenForBroadCasts(ctx, listPeers)

	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server listening on port: %s\n", listener.Addr().String())

	listenerPort := listener.Addr().(*net.TCPAddr).Port
	go SendBroadCasts(listenerPort)

	go listenConnection(listener)

}

func ShutdownServer(cancel context.CancelFunc) {
	fmt.Println("Shutting down the server...")
	cancel()
	fmt.Println("Server has been shut down gracefully.")
}
