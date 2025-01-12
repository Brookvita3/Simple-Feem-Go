package clientUtils

import (
	"bufio"
	config "file-transfer/configs"
	"file-transfer/utils/utils"
	"fmt"
	"io"
	"net"
	"os"
)

func ReadFileChunks(filePath string, fileChannels utils.FileChannels, chunkSize int) {

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fileChannels.ErrorChan <- err
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fileChannels.ErrorChan <- err
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, chunkSize)
	for {
		err := fileChannels.SendChunk(reader, chunkSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			fileChannels.ErrorChan <- err
			return
		}
	}
	close(fileChannels.DataChan)
}

func SendFileChunks(conn net.Conn, fileChannels utils.FileChannels) {
	defer conn.Close()
	for chunk := range fileChannels.DataChan {
		_, err := conn.Write(chunk)
		if err != nil {
			fileChannels.ErrorChan <- err
			fmt.Printf("error sending chunk: %v\n", err)
			break
		}
	}
	fmt.Println("Send Successfully")
	close(fileChannels.ErrorChan)
}

// sendBroadcast sends a UDP broadcast message
func SendBroadCasts() error {
	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+config.Config.BROADCAST_PORT)
	fmt.Println("255.255.255.255:" + config.Config.BROADCAST_PORT)
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

	fileChannels := *utils.NewFileChannels(
		make(chan []byte),
		make(chan error),
	)

	go ReadFileChunks(path, fileChannels, config.Config.CHUNK_SIZE)
	go SendFileChunks(conn, fileChannels)

	for err := range fileChannels.ErrorChan {
		fmt.Printf("Error: %v\n", err)
	}
}
