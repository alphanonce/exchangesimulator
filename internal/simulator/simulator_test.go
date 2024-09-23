package simulator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config := Config{
		ServerAddress: "localhost:8080",
		HttpRules: []HttpRule{
			{
				RequestMatcher: NewHttpRequestPredicate("GET", "/test"),
				Responder:      NewHttpResponseFromString(200, "OK", time.Second),
			},
		},
	}
	sim := New(config)
	assert.Equal(t, config, sim.config)
}

func TestSimulator_simulateHttpResponse(t *testing.T) {
	config := Config{
		ServerAddress: "localhost:8080",
		HttpBasePath:  "/api",
		HttpRules: []HttpRule{
			{
				RequestMatcher: NewHttpRequestPredicate("GET", "/test"),
				Responder:      NewHttpResponseFromString(200, "OK", time.Second),
			},
		},
	}
	sim := New(config)

	tests := []struct {
		name          string
		request       HttpRequest
		expectedCode  int
		expectedBody  string
		expectedDelay time.Duration
	}{
		{
			name:          "Matching request",
			request:       HttpRequest{Method: "GET", Path: "/api/test"},
			expectedCode:  200,
			expectedBody:  "OK",
			expectedDelay: time.Second,
		},
		{
			name:          "Non-matching request",
			request:       HttpRequest{Method: "POST", Path: "/api/other"},
			expectedCode:  404,
			expectedBody:  "Invalid request",
			expectedDelay: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			resp, endTime := sim.simulateHttpResponse(tt.request, startTime)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, string(resp.Body))
			assert.Equal(t, startTime.Add(tt.expectedDelay), endTime)
		})
	}
}

func TestSimulator_simulateWsResponse(t *testing.T) {
	config := Config{
		ServerAddress: "localhost:8080",
		WsEndpoint:    "/ws",
		WsRules: []WsRule{
			{
				MessageMatcher: NewWsMessagePredicate(WsMessageText, []byte("ping")),
				Responder:      NewWsMessageFromString(WsMessageText, "pong", time.Second),
			},
		},
	}
	sim := New(config)

	tests := []struct {
		name          string
		message       WsMessage
		expectedType  WsMessageType
		expectedData  string
		expectedDelay time.Duration
	}{
		{
			name:          "Matching message",
			message:       WsMessage{Type: WsMessageText, Data: []byte("ping")},
			expectedType:  WsMessageText,
			expectedData:  "pong",
			expectedDelay: time.Second,
		},
		{
			name:          "Non-matching message",
			message:       WsMessage{Type: WsMessageText, Data: []byte("hello")},
			expectedType:  WsMessageText,
			expectedData:  "Invalid message",
			expectedDelay: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			resp, endTime := sim.simulateWsResponse(tt.message, startTime)
			assert.Equal(t, tt.expectedType, resp.Type)
			assert.Equal(t, tt.expectedData, string(resp.Data))
			assert.Equal(t, startTime.Add(tt.expectedDelay), endTime)
		})
	}
}
