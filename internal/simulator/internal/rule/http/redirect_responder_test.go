package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRedirectResponder(t *testing.T) {
	targetURL := "http://example.com"
	recordDir := "test/path"

	responder := NewRedirectResponder(targetURL, recordDir)

	assert.Equal(t, targetURL, responder.targetUrl)
	assert.Equal(t, recordDir, responder.recordDir)
}

func TestRedirectResponder_Response(t *testing.T) {
	tests := []struct {
		name           string
		targetURL      string
		request        Request
		expectedStatus int
		expectedBody   string
		expectedError  string
	}{
		{
			name:           "Successful redirect",
			targetURL:      "http://example.com",
			request:        Request{Method: "GET", Path: "/test", Body: []byte(""), Header: http.Header{}},
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, World!",
		},
		{
			name:          "Invalid target URL",
			targetURL:     "://invalid-url",
			request:       Request{Method: "GET", Path: "/test", Body: []byte(""), Header: http.Header{}},
			expectedError: "invalid target URL",
		},
		{
			name:          "Target server error",
			targetURL:     "http://non-existent-server.com",
			request:       Request{Method: "GET", Path: "/test", Body: []byte(""), Header: http.Header{}},
			expectedError: "failed to reach target server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// If we're testing a successful redirect, set up a mock server
			if tt.expectedStatus == http.StatusOK {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.expectedStatus)
					w.Write([]byte(tt.expectedBody))
				}))
				defer server.Close()
				tt.targetURL = server.URL
			}

			responder := NewRedirectResponder(tt.targetURL, "")
			response, err := responder.Response(tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, response.StatusCode)
				assert.Equal(t, []byte(tt.expectedBody), response.Body)
			}
		})
	}
}

func TestRedirectResponder_saveResponseToFile(t *testing.T) {
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

			responder := NewRedirectResponder("", tempDir)

			err := responder.saveResponseToFile(tt.response)
			assert.NoError(t, err)

			files, err := os.ReadDir(tempDir)
			assert.NoError(t, err)
			require.Len(t, files, 1)

			content, err := os.ReadFile(filepath.Join(tempDir, files[0].Name()))
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedContent, string(content))
		})
	}
}
