package ws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessagePredicate(t *testing.T) {
	messageType := MessageText
	data := []byte("test data")

	predicate := NewMessagePredicate(messageType, data)

	assert.Equal(t, messageType, predicate.messageType)
	assert.Equal(t, data, predicate.data)
}

func TestMessagePredicate_MatchMessage(t *testing.T) {
	tests := []struct {
		name      string
		predicate MessagePredicate
		message   Message
		expected  bool
	}{
		{
			name:      "Exact match",
			predicate: NewMessagePredicate(MessageText, []byte("hello")),
			message:   Message{Type: MessageText, Data: []byte("hello")},
			expected:  true,
		},
		{
			name:      "Type match, data mismatch",
			predicate: NewMessagePredicate(MessageText, []byte("hello")),
			message:   Message{Type: MessageText, Data: []byte("world")},
			expected:  false,
		},
		{
			name:      "Type mismatch, data match",
			predicate: NewMessagePredicate(MessageBinary, []byte("hello")),
			message:   Message{Type: MessageText, Data: []byte("hello")},
			expected:  false,
		},
		{
			name:      "Any type, data match",
			predicate: NewMessagePredicate(MessageAny, []byte("hello")),
			message:   Message{Type: MessageText, Data: []byte("hello")},
			expected:  true,
		},
		{
			name:      "Type match, any data",
			predicate: NewMessagePredicate(MessageText, nil),
			message:   Message{Type: MessageText, Data: []byte("any data")},
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.MatchMessage(tt.message)
			assert.Equal(t, tt.expected, result)
		})
	}
}
