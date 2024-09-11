package responder

import (
	"time"

	"alphanonce.com/exchangesimulator/internal/rule"
	"alphanonce.com/exchangesimulator/internal/types"
)

// Ensure ResponseFromString implements Responder
var _ rule.Responder = (*ResponseFromString)(nil)

type ResponseFromString struct {
	statusCode   int
	body         string
	responseTime time.Duration
}

func NewResponseFromString(statusCode int, body string, responseTime time.Duration) ResponseFromString {
	return ResponseFromString{
		statusCode:   statusCode,
		body:         body,
		responseTime: responseTime,
	}
}

func (r ResponseFromString) Response(_ types.Request) types.Response {
	return types.Response{
		StatusCode: r.statusCode,
		Body:       []byte(r.body),
	}
}

func (r ResponseFromString) ResponseTime() time.Duration {
	return r.responseTime
}
