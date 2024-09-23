package ws

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewMessageFromString(t *testing.T) {
	messageType := MessageText
	data := "test data"
	responseTime := 100 * time.Millisecond

	h := NewMessageFromString(messageType, data, responseTime)

	assert.Equal(t, messageType, h.messageType)
	assert.Equal(t, []byte(data), h.data)
	assert.Equal(t, responseTime, h.responseTime)
}

func TestMessageFromString_Response(t *testing.T) {
	messageType := MessageBinary
	data := "response data"
	responseTime := 50 * time.Millisecond

	h := NewMessageFromString(messageType, data, responseTime)

	ctx := context.Background()
	mockConn := new(MockConnection)
	mockConn.On("Write", ctx, Message{Type: messageType, Data: []byte(data)}).Return(nil)

	err := h.Response(ctx, Message{}, mockConn)

	assert.NoError(t, err)
	mockConn.AssertExpectations(t)
}

func TestMessageFromString_ResponseTime(t *testing.T) {
	messageType := MessageText
	data := "test"
	responseTime := 75 * time.Millisecond

	h := NewMessageFromString(messageType, data, responseTime)

	start := time.Now()
	ctx := context.Background()
	mockConn := new(MockConnection)
	mockConn.On("Write", ctx, mock.Anything).Return(nil)

	err := h.Response(ctx, Message{}, mockConn)

	elapsed := time.Since(start)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, elapsed, responseTime)
}
