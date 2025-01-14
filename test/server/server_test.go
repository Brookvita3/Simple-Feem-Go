package server_test

import (
	"bytes"
	severUltils "file-transfer/utils/server"
	"file-transfer/utils/utils"
	"fmt"
	"net"
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

func TestReceiveFileChunks(t *testing.T) {

	content := []byte("This is a test file content")
	listTestCase := []testCase{
		{
			name:      "valid file recieve",
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
				make(chan []byte),
				make(chan error),
				tc.chunkSize,
			)

			server, client := net.Pipe()

			go func() {
				for i := 0; i < len(content); i += tc.chunkSize {
					end := i + tc.chunkSize
					if end > len(content) {
						end = len(content)
					}
					chunk := content[i:end]
					n, err := server.Write(chunk)
					if err != nil || n == 0 {
						fmt.Println("Error sending chunk:", err)
						break
					}
					time.Sleep(10 * time.Millisecond) // Give receiver time to process
				}
				server.Close()
				done <- true
			}()

			go severUltils.ReceiveFileChunks(client, fileChannels)

			// Collect all data chunks
			go func() {
				for chunk := range fileChannels.DataChan {
					result = append(result, chunk...)
				}
			}()

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
