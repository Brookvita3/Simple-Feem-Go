package utils

import (
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
