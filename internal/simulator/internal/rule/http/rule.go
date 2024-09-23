package http

import (
	"time"
)

type Rule struct {
	RequestMatcher
	Responder
}

type RequestMatcher interface {
	MatchRequest(Request) bool
}

type Responder interface {
	Response(Request) Response
	ResponseTime() time.Duration
}
