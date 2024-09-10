package responder

import (
	"time"

	"alphanonce.com/exchangesimulator/internal/types"
)

type ResponseFromString struct {
	body         string
	responseTime time.Duration
}

func NewResponseFromString(body string, responseTime time.Duration) ResponseFromString {
	return ResponseFromString{
		body:         body,
		responseTime: responseTime,
	}
}

func (r ResponseFromString) Response(_ types.Request) types.Response {
	return types.Response{
		Body: []byte(r.body),
	}
}

func (r ResponseFromString) ResponseTime() time.Duration {
	return r.responseTime
}
