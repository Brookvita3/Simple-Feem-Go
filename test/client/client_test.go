// test/client_test.go
package client_test

// import (
// 	"file-transfer/utils/client" // Replace with actual package import
// 	"fmt"
// 	"io"
// 	"net"
// 	"testing"
// )

// func TestSendFile(t *testing.T) {
// 	// Set up the server
// 	listener, err := net.Listen("tcp", "localhost:8081")
// 	if err != nil {
// 		t.Fatalf("Error setting up server: %v", err)
// 	}
// 	defer listener.Close()

// 	errorChan := make(chan error, 1)

// 	// Set up a goroutine to handle the server connection
// 	go func() {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			errorChan <- fmt.Errorf("Error accepting connection: %v", err)
// 		}
// 		defer conn.Close()

// 		// Server simply writes back the data it receives (for testing)
// 		io.Copy(conn, conn)
// 		if err != nil {
// 			errorChan <- fmt.Errorf("Error copying data: %v", err)
// 		}
// 	}()

// 	// Set up the client and dial the server
// 	conn, err := net.Dial("tcp", "localhost:8081")
// 	if err != nil {
// 		t.Fatalf("Error dialing server: %v", err)
// 	}
// 	defer conn.Close()

// 	// Use the connection to send the file
// 	err = client.SendFile("../../files/send/file1.txt", conn)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}
// }
