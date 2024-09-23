package integration

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"alphanonce.com/exchangesimulator/internal/simulator"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
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
				Responder:      simulator.NewWsMessageFromString(simulator.WsMessageText, "pong", 50*time.Millisecond),
			},
		},
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
	ctx := context.Background()
	wsURL := "ws://" + config.ServerAddress + config.WsEndpoint

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	require.NoError(t, err)
	defer conn.Close(websocket.StatusNormalClosure, "")

	t.Run("WebSocket Ping-Pong", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			start := time.Now()

			err := conn.Write(ctx, websocket.MessageText, []byte("ping"))
			require.NoError(t, err)

			msgType, msg, err := conn.Read(ctx)
			require.NoError(t, err)

			duration := time.Since(start)

			assert.Equal(t, websocket.MessageText, msgType)
			assert.Equal(t, "pong", string(msg))
			assert.GreaterOrEqual(t, duration, 50*time.Millisecond)
			assert.Less(t, duration, 100*time.Millisecond) // Allow for some overhead
		}
	})

	t.Run("WebSocket Unmatched Message", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			start := time.Now()

			err := conn.Write(ctx, websocket.MessageText, []byte("hello"))
			require.NoError(t, err)

			msgType, msg, err := conn.Read(ctx)
			require.NoError(t, err)

			duration := time.Since(start)

			assert.Equal(t, websocket.MessageText, msgType)
			assert.Equal(t, "Invalid message", string(msg))
			assert.Less(t, duration, 50*time.Millisecond) // Should be quick as it's not matched
		}
	})
}
