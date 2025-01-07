package utils

import (
	"net"
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
