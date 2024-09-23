package ws

import "slices"

// Ensure MessagePredicate implements MessageMatcher
var _ MessageMatcher = (*MessagePredicate)(nil)

type MessagePredicate struct {
	messageType MessageType
	data        []byte
}

func NewMessagePredicate(messageType MessageType, data []byte) MessagePredicate {
	return MessagePredicate{
		messageType: messageType,
		data:        data,
	}
}

func (p MessagePredicate) MatchMessage(message Message) bool {
	return (p.messageType == MessageInvalid || message.Type == p.messageType) &&
		(p.data == nil || slices.Equal(message.Data, p.data))
}
