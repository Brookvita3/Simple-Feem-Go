package client

import (
	config "file-transfer/configs"
	"file-transfer/utils/utils"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func ReadFile(filePath string, dataChan chan<- []byte, errorChan chan<- error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, config.Config.CHUNK_SIZE)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			dataChan <- buffer[:n] // Send chunk to data channel
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			errorChan <- fmt.Errorf("error reading file: %v", err)
			break
		}
	}
	return nil

}

func SendChunks(conn net.Conn, dataChan <-chan []byte, errorChan chan<- error) {
	defer conn.Close()
	for chunk := range dataChan {
		_, err := conn.Write(chunk)
		if err != nil {
			errorChan <- fmt.Errorf("error sending chunk: %v", err)
			break
		}
	}
}

// TODO: break file to chunk, use goroutine to handle send and receive
func SendFile(filePath string, conn net.Conn) error {
	// Open the file to send
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	// Send metadata: filename and size
	metadata := fmt.Sprintf("%s:%d\n", fileInfo.Name(), fileInfo.Size())
	_, err = conn.Write([]byte(metadata))
	if err != nil {
		return fmt.Errorf("error sending metadata: %v", err)
	}

	// Send the file to the server
	_, err = io.Copy(conn, file)
	if err != nil {
		return fmt.Errorf("error sending file: %v", err)
	}

	return nil
}

// sendBroadcast sends a UDP broadcast message
func SendBroadCasts() error {
	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+strconv.Itoa(config.Config.BROADCAST_PORT))
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

	message, err := utils.GetLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		return err
	}

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending broadcast message:", err)
		return err
	}

	fmt.Println("Broadcast message sent:", message)
	return nil
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
