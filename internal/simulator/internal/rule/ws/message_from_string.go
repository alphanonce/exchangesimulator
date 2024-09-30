package ws

import (
	"context"
	"time"
)

// Ensure MessageFromString implements MessageHandler
var _ MessageHandler = (*MessageFromString)(nil)

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

func (r MessageFromString) Handle(ctx context.Context, _ Message, connClient Connection, _ Connection) error {
	message := Message{
		Type: r.messageType,
		Data: r.data,
	}
	time.Sleep(r.responseTime)
	err := connClient.Write(ctx, message)
	return err
}
