package clientUtils

import (
	"bufio"
	config "file-transfer/configs"
	"file-transfer/utils/utils"
	"fmt"
	"net"
	"os"
)

func ReadFileChunks(filePath string, fileChannels utils.FileChannels) {

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

	reader := bufio.NewReaderSize(file, fileChannels.BufferSize)
	fileChannels.ReadChunks(reader)

	fmt.Println("Read Successfully")
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

func sendFile(conn net.Conn, filePath string) {
	fileChannels := *utils.NewFileChannels(
		make(chan []byte),
		make(chan error),
		config.Config.CHUNK_SIZE,
		config.Config.BUFFER_SIZE,
	)

	go ReadFileChunks(filePath, fileChannels)
	go SendFileChunks(conn, fileChannels)

	for err := range fileChannels.ErrorChan {
		fmt.Printf("Error: %v\n", err)
	}
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
	path := "../files/send/file1.pdf"
	sendFile(conn, path)
}
