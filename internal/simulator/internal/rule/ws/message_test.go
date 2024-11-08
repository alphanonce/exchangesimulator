package ws

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteToFile(t *testing.T) {
	tests := []struct {
		name            string
		message         Message
		expectedContent string
	}{
		{
			name:            "Text message",
			message:         Message{Type: MessageText, Data: []byte("Hello, World!")},
			expectedContent: "type: text\ndata: |-\n    Hello, World!\n",
		},
		{
			name:            "Binary message",
			message:         Message{Type: MessageBinary, Data: []byte{0x01, 0x02, 0x03, 0x04}},
			expectedContent: "type: binary\ndata: |-\n    01020304\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tempPath := filepath.Join(tempDir, "test_ws_message.yaml")

			err := WriteToFile(tempPath, tt.message)
			assert.NoError(t, err)

			content, err := os.ReadFile(tempPath)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedContent, string(content))
		})
	}
}

func TestReadFromFile(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		expectedMessage Message
		wantErr         bool
	}{
		{
			name:            "Valid text message",
			content:         "type: text\ndata: |-\n    Hello, World!\n",
			expectedMessage: Message{Type: MessageText, Data: []byte("Hello, World!")},
			wantErr:         false,
		},
		{
			name:            "Valid binary message",
			content:         "type: binary\ndata: |-\n    01020304\n",
			expectedMessage: Message{Type: MessageBinary, Data: []byte{0x01, 0x02, 0x03, 0x04}},
			wantErr:         false,
		},
		{
			name:            "Invalid YAML",
			content:         "key: value\n",
			expectedMessage: Message{},
			wantErr:         true,
		},
		{
			name:            "Invalid hex in binary message",
			content:         "type: binary\ndata: |-\n    0102030G\n",
			expectedMessage: Message{},
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tempPath := filepath.Join(tempDir, "test_ws_message.yaml")

			err := os.WriteFile(tempPath, []byte(tt.content), 0644)
			assert.NoError(t, err)

			message, err := ReadFromFile(tempPath)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMessage, message)
			}
		})
	}
}

func TestHexEncoding(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	encoded := hex.EncodeToString(data)
	assert.Equal(t, "01020304", encoded)

	decoded, err := hex.DecodeString(encoded)
	assert.NoError(t, err)
	assert.Equal(t, data, decoded)
}
