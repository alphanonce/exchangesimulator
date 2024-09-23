package simulator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetHttpRule(t *testing.T) {
	rule1 := HttpRule{
		RequestMatcher: NewHttpRequestPredicate("GET", "/users"),
		Responder:      NewHttpResponseFromString(200, "Users", time.Second),
	}
	rule2 := HttpRule{
		RequestMatcher: NewHttpRequestPredicate("POST", "/users"),
		Responder:      NewHttpResponseFromString(201, "Created", time.Second),
	}
	config := Config{
		HttpBasePath: "/api",
		HttpRules:    []HttpRule{rule1, rule2},
	}

	tests := []struct {
		name         string
		request      HttpRequest
		expectedRule HttpRule
		expectedOk   bool
	}{
		{"Matching GET request", HttpRequest{Method: "GET", Path: "/api/users"}, rule1, true},
		{"Matching POST request", HttpRequest{Method: "POST", Path: "/api/users"}, rule2, true},
		{"Non-matching path", HttpRequest{Method: "GET", Path: "/api/products"}, HttpRule{}, false},
		{"Non-matching method", HttpRequest{Method: "PUT", Path: "/api/users"}, HttpRule{}, false},
		{"Non-matching base path", HttpRequest{Method: "GET", Path: "/users"}, HttpRule{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, ok := config.GetHttpRule(tt.request)
			assert.Equal(t, tt.expectedRule, rule)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}

func TestConfig_GetWsRule(t *testing.T) {
	rule := WsRule{
		MessageMatcher: NewWsMessagePredicate(WsMessageText, []byte("ping")),
		MessageHandler: NewWsMessageFromString(WsMessageText, "pong", time.Second),
	}
	config := Config{
		WsRules: []WsRule{rule},
	}

	tests := []struct {
		name         string
		message      WsMessage
		expectedRule WsRule
		expectedOk   bool
	}{
		{"Matching message", WsMessage{Type: WsMessageText, Data: []byte("ping")}, rule, true},
		{"Non-matching type", WsMessage{Type: WsMessageBinary, Data: []byte("ping")}, WsRule{}, false},
		{"Non-matching data", WsMessage{Type: WsMessageText, Data: []byte("hello")}, WsRule{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, ok := config.GetWsRule(tt.message)
			assert.Equal(t, tt.expectedRule, rule)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}
