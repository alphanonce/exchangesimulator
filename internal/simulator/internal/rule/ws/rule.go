package ws

import "context"

type Rule struct {
	MessageMatcher
	MessageHandler
}

type MessageMatcher interface {
	MatchMessage(Message) bool
}

type MessageHandler interface {
	Handle(context.Context, Message, Connection, Connection) error
}
