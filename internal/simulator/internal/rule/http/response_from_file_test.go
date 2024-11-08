package http

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewResponseFromFile(t *testing.T) {
	filePath := "/test/filename.yaml"
	responseTime := 100 * time.Millisecond

	r := NewResponseFromFile(filePath, responseTime)

	assert.Equal(t, filePath, r.filePath)
	assert.Equal(t, responseTime, r.responseTime)
}

func TestResponseFromFile_Response(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		responseTime     time.Duration
		expectedResponse Response
	}{
		{
			name:             "Valid response",
			content:          "status: 200\nbody: |-\n    Hello, World!\n",
			responseTime:     50 * time.Millisecond,
			expectedResponse: Response{StatusCode: 200, Body: []byte("Hello, World!")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tempPath := filepath.Join(tempDir, "test_file.yaml")
			err := os.WriteFile(tempPath, []byte("status: 200\nbody: |-\n    Hello, World!\n"), 0644)
			assert.NoError(t, err)

			r := NewResponseFromFile(tempPath, tt.responseTime)
			start := time.Now()
			response, err := r.Response(Request{})
			duration := time.Since(start)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResponse, response)
			assert.GreaterOrEqual(t, duration, tt.responseTime)
			assert.LessOrEqual(t, duration, 2*tt.responseTime)
		})
	}
}
