package rule

import (
	"time"

	"alphanonce.com/exchangesimulator/internal/types"
)

type Rule interface {
	RequestMatcher
	Responder
}

type RequestMatcher interface {
	MatchRequest(types.Request) bool
}

type Responder interface {
	Response(types.Request) types.Response
	ResponseTime() time.Duration
}
