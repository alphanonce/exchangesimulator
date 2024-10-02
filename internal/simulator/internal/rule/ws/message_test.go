package ws

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteToFile(t *testing.T) {
	tests := []struct {
		name            string
		message         Message
		expectedContent []byte
	}{
		{
			name:            "Text message",
			message:         Message{Type: MessageText, Data: []byte("Hello, World!")},
			expectedContent: []byte(`{"type":1,"data":"Hello, World!"}`),
		},
		{
			name:            "Binary message",
			message:         Message{Type: MessageBinary, Data: []byte{0x01, 0x02, 0x03, 0x04}},
			expectedContent: []byte(`{"type":2,"data":"01020304"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile, err := os.CreateTemp("", "test_message_*.json")
			assert.NoError(t, err)
			defer os.Remove(tempFile.Name())

			err = WriteToFile(tempFile.Name(), tt.message)
			assert.NoError(t, err)

			content, err := os.ReadFile(tempFile.Name())
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedContent, content)
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
			content:         `{"type":1,"data":"Hello, World!"}`,
			expectedMessage: Message{Type: MessageText, Data: []byte("Hello, World!")},
			wantErr:         false,
		},
		{
			name:            "Valid binary message",
			content:         `{"type":2,"data":"01020304"}`,
			expectedMessage: Message{Type: MessageBinary, Data: []byte{0x01, 0x02, 0x03, 0x04}},
			wantErr:         false,
		},
		{
			name:            "Invalid JSON",
			content:         `{"type":1,"data":"Hello, World!"`,
			expectedMessage: Message{},
			wantErr:         true,
		},
		{
			name:            "Invalid hex in binary message",
			content:         `{"type":2,"data":"0102030G"}`,
			expectedMessage: Message{},
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile, err := os.CreateTemp("", "test_message_*.json")
			assert.NoError(t, err)
			defer os.Remove(tempFile.Name())

			err = os.WriteFile(tempFile.Name(), []byte(tt.content), 0644)
			assert.NoError(t, err)

			message, err := ReadFromFile(tempFile.Name())

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
