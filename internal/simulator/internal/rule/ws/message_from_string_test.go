package ws

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMessageFromString(t *testing.T) {
	messageType := MessageText
	data := "test data"
	responseTime := 100 * time.Millisecond

	message := NewMessageFromString(messageType, data, responseTime)

	assert.Equal(t, messageType, message.messageType)
	assert.Equal(t, []byte(data), message.data)
	assert.Equal(t, responseTime, message.responseTime)
}

func TestMessageFromString_Response(t *testing.T) {
	messageType := MessageBinary
	data := "response data"
	responseTime := 50 * time.Millisecond

	message := NewMessageFromString(messageType, data, responseTime)
	response := message.Response(Message{}) // Input message is ignored

	assert.Equal(t, messageType, response.Type)
	assert.Equal(t, []byte(data), response.Data)
}

func TestMessageFromString_ResponseTime(t *testing.T) {
	messageType := MessageText
	data := "test"
	responseTime := 75 * time.Millisecond

	message := NewMessageFromString(messageType, data, responseTime)

	assert.Equal(t, responseTime, message.ResponseTime())
}
