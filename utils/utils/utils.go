package utils

import (
	"bufio"
	"fmt"
	"io"
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
	DataChan   chan []byte
	ErrorChan  chan error
	ChunkSize  int
	BufferSize int
}

func (fileChannels *FileChannels) ReadChunks(reader *bufio.Reader) error {
	for {
		chunk := make([]byte, fileChannels.ChunkSize)
		n, err := reader.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			fileChannels.ErrorChan <- err
			return err
		}
		if n > 0 {
			fileChannels.DataChan <- chunk[:n]
		}
	}
	return nil
}

func (fileChannels *FileChannels) WriteChunks(writer *bufio.Writer) error {
	for chunk := range fileChannels.DataChan {
		_, err := writer.Write(chunk)
		if err != nil {
			fileChannels.ErrorChan <- err
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		fileChannels.ErrorChan <- err
		return err
	}

	return nil
}

func NewFileChannels(dataChan chan []byte, errorChan chan error, chunkSize int, bufferSize int) *FileChannels {
	return &FileChannels{
		DataChan:   dataChan,
		ErrorChan:  errorChan,
		ChunkSize:  chunkSize,
		BufferSize: bufferSize,
	}
}
