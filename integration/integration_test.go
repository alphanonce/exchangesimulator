package integration

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"alphanonce.com/exchangesimulator/internal/simulator"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func echo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	for {
		incomingType, incomingData, err := conn.Read(ctx)
		if err != nil {
			return
		}

		outgoingType := incomingType
		outgoingData := append([]byte("echoed: "), incomingData...)

		err = conn.Write(ctx, outgoingType, outgoingData)
		if err != nil {
			return
		}
	}
}

func TestIntegration(t *testing.T) {
	// Start a mock WebSocket server
	mockServer := httptest.NewServer(http.HandlerFunc(echo))
	defer mockServer.Close()
	mockServerURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")

	// Create a temporary directory for saving messages
	tempDir := t.TempDir()

	// Define the configuration for the simulator
	config := simulator.Config{
		ServerAddress: "localhost:8081",
		HttpBasePath:  "/api",
		HttpRules: []simulator.HttpRule{
			{
				RequestMatcher: simulator.NewHttpRequestPredicate("GET", "/test"),
				Responder:      simulator.NewHttpResponseFromString(200, "OK", 100*time.Millisecond),
			},
			{
				RequestMatcher: simulator.NewHttpRequestPredicate("POST", "/data"),
				Responder:      simulator.NewHttpResponseFromString(201, "Created", 200*time.Millisecond),
			},
		},
		WsEndpoint: "/ws",
		WsRules: []simulator.WsRule{
			{
				MessageMatcher: simulator.NewWsMessagePredicate(simulator.WsMessageText, []byte("ping")),
				MessageHandler: simulator.NewWsMessageFromString(simulator.WsMessageText, "pong", 50*time.Millisecond),
			},
			// TODO: add a WsRule with JsonMessageMatcher and MessageFromFiles
			{
				MessageMatcher: simulator.NewWsMessagePredicate(simulator.WsMessageText, []byte("redirect")),
				MessageHandler: simulator.NewWsRedirectHandler(),
			},
		},
		WsRedirectUrl: mockServerURL,
		WsRecordDir:   tempDir,
	}

	// Create a simulator
	sim := simulator.New(config)

	// Start the simulator
	go func() {
		err := sim.Run()
		require.NoError(t, err)
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Run HTTP tests
	t.Run("HTTP Tests", func(t *testing.T) {
		runHttpTests(t, config)
	})

	// Run WebSocket tests
	t.Run("WebSocket Tests", func(t *testing.T) {
		runWsTests(t, config)
	})
}

func runHttpTests(t *testing.T, config simulator.Config) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
		expectedBody   string
		expectedDelay  time.Duration
	}{
		{
			name:           "GET /api/test",
			method:         "GET",
			path:           "/api/test",
			expectedStatus: 200,
			expectedBody:   "OK",
			expectedDelay:  100 * time.Millisecond,
		},
		{
			name:           "POST /api/data",
			method:         "POST",
			path:           "/api/data",
			body:           "some data",
			expectedStatus: 201,
			expectedBody:   "Created",
			expectedDelay:  200 * time.Millisecond,
		},
		{
			name:           "Unmatched route",
			method:         "GET",
			path:           "/unknown",
			expectedStatus: 404,
			expectedBody:   "Invalid endpoint\n",
			expectedDelay:  0,
		},
		{
			name:           "Matched endpoint, unmatched path",
			method:         "GET",
			path:           "/api/unknown",
			expectedStatus: 404,
			expectedBody:   "Invalid request",
			expectedDelay:  0,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			req, err := http.NewRequest(tt.method, "http://"+config.ServerAddress+tt.path, strings.NewReader(tt.body))
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			duration := time.Since(start)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, string(body))
			assert.GreaterOrEqual(t, duration, tt.expectedDelay)
			assert.Less(t, duration, tt.expectedDelay+50*time.Millisecond) // Allow for some overhead
		})
	}
}

func runWsTests(t *testing.T, config simulator.Config) {
	tests := []struct {
		name                string
		msgType             websocket.MessageType
		msgData             []byte
		expectedMsgType     websocket.MessageType
		expectedMsgData     []byte
		expectedDelay       time.Duration
		expectedFileContent []byte
	}{
		{
			name:                "WebSocket Ping-Pong",
			msgType:             websocket.MessageText,
			msgData:             []byte("ping"),
			expectedMsgType:     websocket.MessageText,
			expectedMsgData:     []byte("pong"),
			expectedDelay:       50 * time.Millisecond,
			expectedFileContent: nil,
		},
		{
			name:                "WebSocket Unmatched Message",
			msgType:             websocket.MessageText,
			msgData:             []byte("hello"),
			expectedMsgType:     websocket.MessageText,
			expectedMsgData:     []byte("Invalid message"),
			expectedDelay:       0,
			expectedFileContent: nil,
		},
		{
			name:                "WebSocket Redirect Message",
			msgType:             websocket.MessageText,
			msgData:             []byte("redirect"),
			expectedMsgType:     websocket.MessageText,
			expectedMsgData:     []byte("echoed: redirect"),
			expectedDelay:       0,
			expectedFileContent: []byte(`{"type":1,"data":"echoed: redirect"}`),
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			wsURL := "ws://" + config.ServerAddress + config.WsEndpoint

			conn, _, err := websocket.Dial(ctx, wsURL, nil)
			require.NoError(t, err)
			defer conn.Close(websocket.StatusNormalClosure, "")

			for i := 0; i < 3; i++ {
				start := time.Now()

				err := conn.Write(ctx, tt.msgType, tt.msgData)
				require.NoError(t, err)

				msgType, msgData, err := conn.Read(ctx)
				require.NoError(t, err)

				duration := time.Since(start)

				assert.Equal(t, tt.expectedMsgType, msgType)
				assert.Equal(t, tt.expectedMsgData, msgData)
				assert.GreaterOrEqual(t, duration, tt.expectedDelay)
				assert.Less(t, duration, 2*tt.expectedDelay+10*time.Millisecond)

				if tt.expectedFileContent != nil {
					files, err := os.ReadDir(config.WsRecordDir)
					require.NoError(t, err)
					assert.Len(t, files, 1)

					content, err := os.ReadFile(filepath.Join(config.WsRecordDir, files[0].Name()))
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedFileContent, content)

					err = os.Remove(filepath.Join(config.WsRecordDir, files[0].Name()))
					require.NoError(t, err)
				}
			}
		})
	}
}
