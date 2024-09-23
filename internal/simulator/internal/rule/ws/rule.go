package ws

import "context"

type Rule struct {
	MessageMatcher
	Responder
}

type MessageMatcher interface {
	MatchMessage(Message) bool
}

type Responder interface {
	Response(context.Context, Message, Connection) error
}
