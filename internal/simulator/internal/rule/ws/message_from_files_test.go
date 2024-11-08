package ws

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewMessageFromFiles(t *testing.T) {
	dirPath := "/test/path"

	h := NewMessageFromFiles(dirPath)

	assert.Equal(t, dirPath, h.dirPath)
}

func TestMessageFromFiles_Handle(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Create test YAML files
	testFiles := []string{
		"2000-01-23T12:34:56.000000+09:00.yaml",
		"2000-01-23T12:34:56.010000+09:00.yaml",
		"2000-01-23T12:34:56.020000+09:00.yaml",
	}
	for _, file := range testFiles {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte("type: text\ndata: 'test content'"), 0644)
		assert.NoError(t, err)
	}

	// Create a non-YAML file to test filtering
	err := os.WriteFile(filepath.Join(tempDir, "non_yaml.txt"), []byte("type: text\ndata: 'test content'"), 0644)
	assert.NoError(t, err)

	// Create a mock connection
	mockConn := NewMockConnection(t)

	// Set up expectations
	mockConn.On("Write", mock.Anything, mock.AnythingOfType("Message")).Return(nil).Times(len(testFiles))

	h := NewMessageFromFiles(tempDir)

	ctx := context.Background()
	err = h.Handle(ctx, Message{}, mockConn, nil)

	assert.NoError(t, err)
}

func TestMessageFromFiles_Handle_Error(t *testing.T) {
	mockConn := NewMockConnection(t)

	h := NewMessageFromFiles("/non/existent/path")

	ctx := context.Background()
	err := h.Handle(ctx, Message{}, mockConn, nil)

	assert.Error(t, err)
}

func TestMessageFromFiles_Handle_ContextCancellation(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Create a test YAML file
	err := os.WriteFile(filepath.Join(tempDir, "2000-01-23T12:34:56.000000+09:00.yaml"), []byte("type: text\ndata: 'test content'"), 0644)
	assert.NoError(t, err)

	// Create a mock connection
	mockConn := NewMockConnection(t)

	h := NewMessageFromFiles(tempDir)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately
	err = h.Handle(ctx, Message{}, mockConn, nil)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}
