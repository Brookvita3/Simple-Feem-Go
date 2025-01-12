// test/client_test.go
package client_test

import (
	"bytes"
	config "file-transfer/configs"
	clientUtils "file-transfer/utils/client"
	"file-transfer/utils/utils"

	//"net"
	"os"
	"testing"
	"time"
)

func TestReadFileChunks(t *testing.T) {

	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := []byte("This is a test file content to be read in chunks.")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}
	tempFile.Close()

	testCases := []struct {
		name       string
		filePath   string
		chunkSize  int
		wantError  bool
		expected   []byte
		errorCheck func(err error) bool
	}{
		{
			name:      "valid file read",
			filePath:  tempFile.Name(),
			chunkSize: 10,
			wantError: false,
			expected:  content,
		},
		{
			name:      "file not found",
			filePath:  "nonexistentfile.txt",
			chunkSize: 10,
			wantError: true,
			errorCheck: func(err error) bool {
				return os.IsNotExist(err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			config.Config.CHUNK_SIZE = tc.chunkSize

			fileChannels := *utils.NewFileChannels(
				make(chan []byte),
				make(chan error),
			)

			go clientUtils.ReadFileChunks(tc.filePath, fileChannels, tc.chunkSize)

			var result []byte
			var receivedError error
			done := make(chan bool)

			go func() {
				for {
					select {
					case data, ok := <-fileChannels.DataChan:
						if !ok {
							done <- true
							return
						}
						result = append(result, data...)
					case err := <-fileChannels.ErrorChan:
						receivedError = err
						done <- true
						return
					}
				}
			}()

			// Wait for results or timeout
			select {
			case <-done:
				// Validate the results
				if tc.wantError {
					if receivedError == nil {
						t.Errorf("expected an error but got none")
					} else if tc.errorCheck != nil && !tc.errorCheck(receivedError) {
						t.Errorf("error did not match condition: %v", receivedError)
					}
				} else {
					if !bytes.Equal(result, tc.expected) {
						t.Errorf("expected: %s, got: %s", tc.expected, result)
					}
					if receivedError != nil {
						t.Errorf("unexpected error: %v", receivedError)
					}
				}
			case <-time.After(2 * time.Second):
				t.Error("test timed out")
			}
		})
	}
}

// func TestSendFileChunks(t *testing.T) {
// 	conn, err := net.Dial("tcp", "localhost:8081")
// 	if err != nil {
// 		t.Fatalf("Error dialing server: %v", err)
// 	}
// 	defer conn.Close()

// 	testCases := []struct {
// 		name       string
// 		chunkSize  int
// 		wantError  bool
// 		errorCheck func(err error) bool
// 	}{
// 		{
// 			name:      "valid chunk send",
// 			chunkSize: 10,
// 			wantError: false,
// 		},
// 	}

// }
