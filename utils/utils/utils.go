package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
)

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

func ClearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error clearing screen:", err)
	}
}

func GetInput() string {
	var input string
	fmt.Scanln(&input)
	return input
}

type FileChannels struct {
	DataChan  chan []byte
	ErrorChan chan error
}

func (fileChannels *FileChannels) SendChunk(reader *bufio.Reader, chunkSize int) error {
	chunk := make([]byte, chunkSize)
	n, err := reader.Read(chunk)
	if n > 0 {
		fileChannels.DataChan <- chunk[:n]
	}
	return err
}

func NewFileChannels(dataChan chan []byte, errorChan chan error) *FileChannels {
	return &FileChannels{
		DataChan:  dataChan,
		ErrorChan: errorChan,
	}
}
