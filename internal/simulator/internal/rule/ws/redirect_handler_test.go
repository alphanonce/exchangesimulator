package ws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRedirectHandler_Handle(t *testing.T) {
	tests := []struct {
		name            string
		message         Message
		expectedMessage Message
	}{
		{
			name:            "Text message",
			message:         Message{Type: MessageText, Data: []byte("data1")},
			expectedMessage: Message{Type: MessageText, Data: []byte("data1")},
		},
		{
			name:            "Binary message",
			message:         Message{Type: MessageBinary, Data: []byte("data2")},
			expectedMessage: Message{Type: MessageBinary, Data: []byte("data2")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockConnClient := new(MockConnection)
			mockConnServer := new(MockConnection)
			mockConnServer.On("Write", ctx, tt.expectedMessage).Return(nil)

			err := NewRedirectHandler().Handle(ctx, tt.message, mockConnClient, mockConnServer)

			assert.NoError(t, err)
			mockConnClient.AssertExpectations(t)
			mockConnServer.AssertExpectations(t)
			mockConnClient.AssertNotCalled(t, "Write", mock.Anything)
		})
	}
}
