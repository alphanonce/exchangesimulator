package simulator

import (
	"testing"
	"time"

	"alphanonce.com/exchangesimulator/internal/simulator/internal/rule/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	mockRule1 := ws.NewMockRule(t)
	mockRule1.On("MatchMessage", WsMessage{Type: WsMessageText, Data: []byte("ping")}).Return(true)
	mockRule1.On("MatchMessage", mock.Anything).Return(false)

	mockRule2 := ws.NewMockRule(t)
	mockRule2.On("MatchMessage", WsMessage{Type: WsMessageBinary, Data: []byte("pong")}).Return(true)
	mockRule2.On("MatchMessage", mock.Anything).Return(false)

	config := Config{
		WsRules: []WsRule{mockRule1, mockRule2},
	}

	tests := []struct {
		name         string
		message      WsMessage
		expectedRule WsRule
		expectedOk   bool
	}{
		{"Matching message 1", WsMessage{Type: WsMessageText, Data: []byte("ping")}, mockRule1, true},
		{"Matching message 2", WsMessage{Type: WsMessageBinary, Data: []byte("pong")}, mockRule2, true},
		{"Non-matching message 1", WsMessage{Type: WsMessageBinary, Data: []byte("ping")}, nil, false},
		{"Non-matching message 2", WsMessage{Type: WsMessageText, Data: []byte("pong")}, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, ok := config.GetWsRule(tt.message)
			assert.Equal(t, tt.expectedRule, rule)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}
