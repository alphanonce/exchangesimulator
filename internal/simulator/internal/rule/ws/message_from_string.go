package ws

import (
	"time"
)

// Ensure MessageFromString implements Responder
var _ Responder = (*MessageFromString)(nil)

type MessageFromString struct {
	messageType  MessageType
	data         []byte
	responseTime time.Duration
}

func NewMessageFromString(messageType MessageType, data string, responseTime time.Duration) MessageFromString {
	return MessageFromString{
		messageType:  messageType,
		data:         []byte(data),
		responseTime: responseTime,
	}
}

func (r MessageFromString) Response(_ Message) Message {
	return Message{
		Type: r.messageType,
		Data: r.data,
	}
}

func (r MessageFromString) ResponseTime() time.Duration {
	return r.responseTime
}
