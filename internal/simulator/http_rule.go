package simulator

import (
	"time"

	"alphanonce.com/exchangesimulator/internal/simulator/internal/rule/http"
)

type HttpRule = http.Rule
type HttpRequest = http.Request
type HttpResponse = http.Response

func NewHttpRule(requestMatcher http.RequestMatcher, responder http.Responder) http.RuleImpl {
	return http.NewRule(requestMatcher, responder)
}

func NewHttpRequestPredicate(method string, path string) http.RequestPredicate {
	return http.NewRequestPredicate(method, path)
}

func NewHttpResponseFromString(statusCode int, body string, responseTime time.Duration) http.ResponseFromString {
	return http.NewResponseFromString(statusCode, body, responseTime)
}
