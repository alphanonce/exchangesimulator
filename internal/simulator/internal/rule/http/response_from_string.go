package http

import (
	"time"
)

// Ensure ResponseFromString implements Responder
var _ Responder = (*ResponseFromString)(nil)

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

func (r ResponseFromString) Response(_ Request) (Response, error) {
	resp := Response{
		StatusCode: r.statusCode,
		Body:       []byte(r.body),
	}
	time.Sleep(r.responseTime)
	return resp, nil
}
