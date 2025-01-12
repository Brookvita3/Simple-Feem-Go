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

func ReadFileChunks(filePath string, dataChan chan<- []byte, errorChan chan<- error) {
	file, err := os.Open(filePath)
	if err != nil {
		errorChan <- fmt.Errorf("error opening file: %v", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, config.Config.CHUNK_SIZE)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			dataChan <- buffer[:n]
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			errorChan <- fmt.Errorf("error reading file: %v", err)
			break
		}
	}
	fmt.Println("Read Successfully")
	close(dataChan)
}

func SendFileChunks(conn net.Conn, dataChan <-chan []byte, errorChan chan<- error) {
	defer conn.Close()
	for chunk := range dataChan {
		_, err := conn.Write(chunk)
		if err != nil {
			errorChan <- fmt.Errorf("error sending chunk: %v", err)
			break
		}
	}
	fmt.Println("Send Successfully")
	close(errorChan)
}

// sendBroadcast sends a UDP broadcast message
func SendBroadCasts() error {
	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+strconv.Itoa(config.Config.BROADCAST_PORT))
	fmt.Println("255.255.255.255:" + strconv.Itoa(config.Config.BROADCAST_PORT))
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

	// example send file
	path := "../../files/send/file1.pdf"
	dataChan := make(chan []byte)
	errorChan := make(chan error)

	go ReadFileChunks(path, dataChan, errorChan)
	go SendFileChunks(conn, dataChan, errorChan)

	for err := range errorChan {
		fmt.Println("Error:", err)
	}

}
