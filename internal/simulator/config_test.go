package simulator

import (
	"testing"

	"alphanonce.com/exchangesimulator/internal/simulator/internal/rule/http"
	"alphanonce.com/exchangesimulator/internal/simulator/internal/rule/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfig_GetHttpRule(t *testing.T) {
	mockRule1 := http.NewMockRule(t)
	mockRule1.On("MatchRequest", HttpRequest{Method: "GET", Path: "/users"}).Return(true)
	mockRule1.On("MatchRequest", mock.Anything).Return(false)

	mockRule2 := http.NewMockRule(t)
	mockRule2.On("MatchRequest", HttpRequest{Method: "POST", Path: "/users"}).Return(true)
	mockRule2.On("MatchRequest", mock.Anything).Return(false)

	config := Config{
		HttpBasePath: "/api",
		HttpRules:    []HttpRule{mockRule1, mockRule2},
	}

	tests := []struct {
		name         string
		request      HttpRequest
		expectedRule HttpRule
		expectedOk   bool
	}{
		{"Matching GET request", HttpRequest{Method: "GET", Path: "/api/users"}, mockRule1, true},
		{"Matching POST request", HttpRequest{Method: "POST", Path: "/api/users"}, mockRule2, true},
		{"Non-matching path", HttpRequest{Method: "GET", Path: "/api/products"}, nil, false},
		{"Non-matching method", HttpRequest{Method: "PUT", Path: "/api/users"}, nil, false},
		{"Non-matching base path", HttpRequest{Method: "GET", Path: "/users"}, nil, false},
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
