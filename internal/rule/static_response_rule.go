package rule

import (
	"time"

	"alphanonce.com/exchangesimulator/internal/types"
)

type RequestPredicate struct {
	method string
	path   string
}

func (r RequestPredicate) MatchRequest(request types.Request) bool {
	return (r.method == "" || request.Method == r.method) &&
		(r.path == "" || request.Path == r.path)
}

type ResponseFromString struct {
	body         string
	responseTime time.Duration
}

func (r ResponseFromString) Response(_ types.Request) types.Response {
	return types.Response{
		Body: []byte(r.body),
	}
}

func (r ResponseFromString) ResponseTime() time.Duration {
	return r.responseTime
}

type StaticResponseRule struct {
	RequestPredicate
	ResponseFromString
}

func NewStaticResponseRule(method string, path string, responseBody string, responseTime time.Duration) StaticResponseRule {
	return StaticResponseRule{
		RequestPredicate: RequestPredicate{
			method: method,
			path:   path,
		},
		ResponseFromString: ResponseFromString{
			body:         responseBody,
			responseTime: responseTime,
		},
	}
}
