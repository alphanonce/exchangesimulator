package ws

import (
	"context"
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

func (r MessageFromString) Response(ctx context.Context, _ Message, conn Connection) error {
	message := Message{
		Type: r.messageType,
		Data: r.data,
	}
	time.Sleep(r.responseTime)
	err := conn.Write(ctx, message)
	return err
}
