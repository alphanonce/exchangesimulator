package simulator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"alphanonce.com/exchangesimulator/internal/simulator/internal/rule/http"
	"alphanonce.com/exchangesimulator/internal/simulator/internal/rule/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	config := Config{ServerAddress: "localhost:8080"}
	sim := New(config)
	assert.Equal(t, config, sim.config)
}

func TestSimulator_simulateHttpResponse(t *testing.T) {
	mockRule := http.NewMockRule(t)
	mockRule.On("MatchRequest", HttpRequest{Method: "GET", Path: "/test"}).Return(true)
	mockRule.On("MatchRequest", mock.Anything).Return(false)
	mockRule.On("Response", mock.Anything).Return(HttpResponse{StatusCode: 200, Body: []byte("OK")})
	mockRule.On("ResponseTime").Return(time.Second)

	config := Config{
		HttpBasePath: "/api",
		HttpRules:    []HttpRule{mockRule},
	}
	sim := New(config)

	tests := []struct {
		name          string
		request       HttpRequest
		expectedResp  HttpResponse
		expectedDelay time.Duration
	}{
		{
			name:          "Matching request",
			request:       HttpRequest{Method: "GET", Path: "/api/test"},
			expectedResp:  HttpResponse{StatusCode: 200, Body: []byte("OK")},
			expectedDelay: time.Second,
		},
		{
			name:          "Non-matching request",
			request:       HttpRequest{Method: "POST", Path: "/api/other"},
			expectedResp:  HttpResponse{StatusCode: 404, Body: []byte("Invalid request")},
			expectedDelay: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			resp, endTime := sim.simulateHttpResponse(tt.request, startTime)
			assert.Equal(t, tt.expectedResp, resp)
			assert.Equal(t, startTime.Add(tt.expectedDelay), endTime)
		})
	}
}

func TestSimulator_simulateWsResponse(t *testing.T) {
	mockPingpongRule := ws.NewMockRule(t)
	mockPingpongRule.On("MatchMessage", WsMessage{Type: WsMessageText, Data: []byte("ping")}).Return(true)
	mockPingpongRule.On("MatchMessage", mock.Anything).Return(false)
	mockPingpongRule.On("Handle", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		connClient := args.Get(2).(WsConnection)
		connClient.Write(ctx, WsMessage{Type: WsMessageText, Data: []byte("pong")})
	}).Return(nil)

	mockRedirectRule := ws.NewMockRule(t)
	mockRedirectRule.On("MatchMessage", WsMessage{Type: WsMessageBinary, Data: []byte("redirect")}).Return(true)
	mockRedirectRule.On("MatchMessage", mock.Anything).Return(false)
	mockRedirectRule.On("Handle", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		message := args.Get(1).(WsMessage)
		connServer := args.Get(3).(WsConnection)
		connServer.Write(ctx, message)
	}).Return(nil)

	config := Config{WsRules: []WsRule{mockPingpongRule, mockRedirectRule}}
	sim := New(config)

	tests := []struct {
		name                  string
		message               WsMessage
		expectedMessageClient WsMessage
		expectedMessageServer WsMessage
	}{
		{
			name:                  "Ping-pong message",
			message:               WsMessage{Type: WsMessageText, Data: []byte("ping")},
			expectedMessageClient: WsMessage{Type: WsMessageText, Data: []byte("pong")},
		},
		{
			name:                  "Redirect message",
			message:               WsMessage{Type: WsMessageBinary, Data: []byte("redirect")},
			expectedMessageServer: WsMessage{Type: WsMessageBinary, Data: []byte("redirect")},
		},
		{
			name:                  "Non-matching message",
			message:               WsMessage{Type: WsMessageText, Data: []byte("hello")},
			expectedMessageClient: WsMessage{Type: WsMessageText, Data: []byte("Invalid message")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockConnClient := ws.NewMockConnection(t)
			mockConnClient.On("Write", ctx, tt.expectedMessageClient).Maybe().Return(nil)
			mockConnServer := ws.NewMockConnection(t)
			mockConnServer.On("Write", ctx, tt.expectedMessageServer).Maybe().Return(nil)

			err := sim.simulateWsResponse(ctx, tt.message, mockConnClient, mockConnServer)
			assert.NoError(t, err)
		})
	}
}

func TestSimulator_saveMessageToFile(t *testing.T) {
	tests := []struct {
		name            string
		message         WsMessage
		expectedContent string
	}{
		{
			name:            "Text message",
			message:         WsMessage{Type: WsMessageText, Data: []byte("Hello, World!")},
			expectedContent: "type: text\ndata: |-\n    Hello, World!\n",
		},
		{
			name:            "Binary message",
			message:         WsMessage{Type: WsMessageBinary, Data: []byte{0x01, 0x02, 0x03, 0x04}},
			expectedContent: "type: binary\ndata: |-\n    01020304\n",
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
			require.Len(t, files, 1)

			content, err := os.ReadFile(filepath.Join(tempDir, files[0].Name()))
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedContent, string(content))
		})
	}
}
