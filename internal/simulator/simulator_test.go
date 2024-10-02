package simulator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"alphanonce.com/exchangesimulator/internal/simulator/internal/rule/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				MessageHandler: NewWsMessageFromString(WsMessageText, "pong", time.Second),
			},
			{
				MessageMatcher: NewWsMessagePredicate(WsMessageBinary, []byte("redirect")),
				MessageHandler: NewWsRedirectHandler(),
			},
		},
	}
	sim := New(config)

	tests := []struct {
		name                  string
		message               WsMessage
		expectedMessageClient WsMessage
		expectedMessageServer WsMessage
		expectedDelay         time.Duration
	}{
		{
			name:                  "Matching message",
			message:               WsMessage{Type: WsMessageText, Data: []byte("ping")},
			expectedMessageClient: WsMessage{Type: WsMessageText, Data: []byte("pong")},
			expectedDelay:         time.Second,
		},
		{
			name:                  "Non-matching message",
			message:               WsMessage{Type: WsMessageText, Data: []byte("hello")},
			expectedMessageClient: WsMessage{Type: WsMessageText, Data: []byte("Invalid message")},
			expectedDelay:         0,
		},
		{
			name:                  "Redirect message",
			message:               WsMessage{Type: WsMessageBinary, Data: []byte("redirect")},
			expectedMessageServer: WsMessage{Type: WsMessageBinary, Data: []byte("redirect")},
			expectedDelay:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockConnClient := new(ws.MockConnection)
			mockConnClient.On("Write", ctx, tt.expectedMessageClient).Maybe().Return(nil)
			mockConnServer := new(ws.MockConnection)
			mockConnServer.On("Write", ctx, tt.expectedMessageServer).Maybe().Return(nil)

			start := time.Now()
			err := sim.simulateWsResponse(ctx, tt.message, mockConnClient, mockConnServer)
			delay := time.Since(start)

			assert.NoError(t, err)
			mockConnClient.AssertExpectations(t)
			mockConnServer.AssertExpectations(t)

			assert.GreaterOrEqual(t, delay, tt.expectedDelay)
		})
	}
}

func TestSimulator_saveMessageToFile(t *testing.T) {
	tests := []struct {
		name            string
		message         WsMessage
		expectedContent []byte
	}{
		{
			name:            "Text message",
			message:         WsMessage{Type: WsMessageText, Data: []byte("Hello, World!")},
			expectedContent: []byte(`{"type":1,"data":"Hello, World!"}`),
		},
		{
			name:            "Binary message",
			message:         WsMessage{Type: WsMessageBinary, Data: []byte{0x01, 0x02, 0x03, 0x04}},
			expectedContent: []byte(`{"type":2,"data":"01020304"}`),
		},
		{
			name:            "Any message",
			message:         WsMessage{Type: WsMessageAny, Data: []byte{0x01, 0x02, 0x03, 0x04}},
			expectedContent: []byte(`{"type":0,"data":"01020304"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "ws_test")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			sim := New(Config{WsRecordDir: tempDir})

			err = sim.saveMessageToFile(tt.message, tempDir)
			assert.NoError(t, err)

			files, err := os.ReadDir(tempDir)
			assert.NoError(t, err)
			assert.Len(t, files, 1)

			content, err := os.ReadFile(filepath.Join(tempDir, files[0].Name()))
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedContent, content)
		})
	}
}
