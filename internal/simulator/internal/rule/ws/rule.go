package ws

import (
	"time"
)

type Rule struct {
	MessageMatcher
	Responder
}

type MessageMatcher interface {
	MatchMessage(Message) bool
}

type Responder interface {
	Response(Message) Message
	ResponseTime() time.Duration
}
