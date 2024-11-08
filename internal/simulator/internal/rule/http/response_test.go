package http

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteToFile(t *testing.T) {
	tests := []struct {
		name            string
		response        Response
		expectedContent string
	}{
		{
			name:            "Basic test",
			response:        Response{StatusCode: 200, Body: []byte("Hello, World!")},
			expectedContent: "status: 200\nbody: |-\n    Hello, World!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tempPath := filepath.Join(tempDir, "test_http_response.yaml")

			err := WriteToFile(tempPath, tt.response)
			assert.NoError(t, err)

			content, err := os.ReadFile(tempPath)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedContent, string(content))
		})
	}
}

func TestReadFromFile(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		expectedResponse Response
		wantErr          bool
	}{
		{
			name:             "Valid response",
			content:          "status: 200\nbody: |-\n    Hello, World!\n",
			expectedResponse: Response{StatusCode: 200, Body: []byte("Hello, World!")},
			wantErr:          false,
		},
		{
			name:             "Invalid YAML",
			content:          "key: value\n",
			expectedResponse: Response{},
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tempPath := filepath.Join(tempDir, "test_http_response.yaml")

			err := os.WriteFile(tempPath, []byte(tt.content), 0644)
			assert.NoError(t, err)

			response, err := ReadFromFile(tempPath)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}
