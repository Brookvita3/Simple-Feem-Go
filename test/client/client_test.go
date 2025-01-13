// test/client_test.go
package client_test

import (
	"bytes"
	clientUtils "file-transfer/utils/client"
	"file-transfer/utils/utils"
	"net"
	"os"
	"testing"
	"time"
)

type testCase struct {
	name       string
	filePath   string
	chunkSize  int
	wantError  bool
	expected   []byte
	errorCheck func(err error) bool
}

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

	listTestCase := []testCase{
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

	for _, tc := range listTestCase {
		t.Run(tc.name, func(t *testing.T) {

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

func TestSendFileChunks(t *testing.T) {
	content := []byte("This is a test file content")
	listTestCase := []testCase{
		{
			name:      "valid file send",
			chunkSize: 10,
			wantError: false,
			expected:  content,
		},
	}

	for _, tc := range listTestCase {
		t.Run(tc.name, func(t *testing.T) {
			var result []byte
			done := make(chan bool)

			fileChannels := *utils.NewFileChannels(
				make(chan []byte, 1),
				make(chan error),
			)

			fileChannels.DataChan <- content
			close(fileChannels.DataChan) // Important to close after sending

			server, client := net.Pipe()

			go func() {
				buffer := make([]byte, len(content))
				n, err := server.Read(buffer)
				if err != nil {
					t.Errorf("Failed to read: %v", err)
				}
				result = append(result, buffer[:n]...)
				server.Close()
				done <- true
			}()

			go clientUtils.SendFileChunks(client, fileChannels)

			// Wait for completion
			select {
			case <-done:
				if !bytes.Equal(result, tc.expected) {
					t.Errorf("expected %s, got %s", string(tc.expected), string(result))
				}
			case <-time.After(2 * time.Second):
				t.Fatal("Test timed out")
			}
		})
	}
}
