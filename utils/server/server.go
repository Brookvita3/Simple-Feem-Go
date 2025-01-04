package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func ReceiveFile(conn net.Conn) error {
	// Define the path to save the received file in the receive folder
	path := "../../files/receive"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating receive directory: %v", err)
	}

	// Open a new file to save the received data
	receivedFilePath := filepath.Join(path, "file_received.pdf")
	file, err := os.Create(receivedFilePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Copy the data from the connection into the file
	_, err = io.Copy(file, conn)
	if err != nil {
		return fmt.Errorf("error saving received file: %v", err)
	}

	// File received and saved successfully
	fmt.Println("File received and saved successfully!")
	return nil
}

func GetLocalIP() (string, error) {
	// Get the local IP address of the machine
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
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

		// Handle the file reception
		conn.Write([]byte("Welcome to the server!\n"))

		err = ReceiveFile(conn)
		if err != nil {
			fmt.Println("Error receiving file:", err)
		}
		conn.Close()
	}
}
