package client

import (
	"fmt"
	"io"
	"net"
	"os"
)

func SendFile(filePath string, conn net.Conn) error {
	// Open the file to send
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Send the file to the server
	_, err = io.Copy(conn, file)
	if err != nil {
		return fmt.Errorf("error sending file: %v", err)
	}

	return nil
}

func getFreePort() (int, error) {
	// Create a temporary listener on port 0, which tells the OS to find an available port.
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("could not find a free port: %v", err)
	}
	defer listener.Close()

	// Retrieve the address of the listener and extract the port number.
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

func StartClient(addr string) {
	addr = fmt.Sprintf("%s:%d", addr, 8080)
	conn, err := net.Dial("tcp", addr) // Example address and port
	if err != nil {
		fmt.Println("Error dialing server:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024) // Create a buffer to hold the response
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	fmt.Println(string(buffer))

	// Specify the file path to send
	path := "../../files/send/file1.pdf"
	err = SendFile(path, conn)
	if err != nil {
		fmt.Println("Error sending file:", err)
	}
}
