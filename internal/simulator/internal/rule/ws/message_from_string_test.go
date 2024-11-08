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

func TestMessageFromString_Handle(t *testing.T) {
	tests := []struct {
		name            string
		handler         MessageFromString
		expectedMessage Message
		expectedDelay   time.Duration
	}{
		{
			name:            "Text message",
			handler:         NewMessageFromString(MessageText, "data1", 5*time.Millisecond),
			expectedMessage: Message{Type: MessageText, Data: []byte("data1")},
			expectedDelay:   5 * time.Millisecond,
		},
		{
			name:            "Binary message",
			handler:         NewMessageFromString(MessageBinary, "data2", 5*time.Millisecond),
			expectedMessage: Message{Type: MessageBinary, Data: []byte("data2")},
			expectedDelay:   5 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockConnClient := NewMockConnection(t)
			mockConnClient.On("Write", ctx, tt.expectedMessage).Return(nil)
			mockConnServer := NewMockConnection(t)

			start := time.Now()
			err := tt.handler.Handle(ctx, Message{}, mockConnClient, mockConnServer)
			duration := time.Since(start)

			assert.NoError(t, err)
			assert.GreaterOrEqual(t, duration, tt.expectedDelay)
			assert.LessOrEqual(t, duration, 2*tt.expectedDelay)
			mockConnServer.AssertNotCalled(t, "Write", mock.Anything)
		})
	}
}
