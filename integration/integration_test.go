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

	// Create a temporary directory for test files and for saving messages
	tempDir := t.TempDir()

	// Create test YAML files
	err := os.Mkdir(filepath.Join(tempDir, "http"), 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Join(tempDir, "ws", "single_file"), 0755)
	assert.NoError(t, err)
	err = os.Mkdir(filepath.Join(tempDir, "ws", "multiple_files"), 0755)
	assert.NoError(t, err)
	testFiles := []struct {
		Path string
		Data []byte
	}{
		{
			Path: filepath.Join(tempDir, "http", "response.yaml"),
			Data: []byte("status: 200\nbody: 'test content'\n"),
		},
		{
			Path: filepath.Join(tempDir, "ws", "single_file", "2000-01-23T12:34:56.000000+09:00.yaml"),
			Data: []byte("type: text\ndata: 'test content 1'\n"),
		},
		{
			Path: filepath.Join(tempDir, "ws", "multiple_files", "2000-01-23T12:34:56.000000+09:00.yaml"),
			Data: []byte("type: text\ndata: 'test content 1'\n"),
		},
		{
			Path: filepath.Join(tempDir, "ws", "multiple_files", "2000-01-23T12:34:56.010000+09:00.yaml"),
			Data: []byte("type: text\ndata: 'test content 2'\n"),
		},
		{
			Path: filepath.Join(tempDir, "ws", "multiple_files", "2000-01-23T12:34:56.020000+09:00.yaml"),
			Data: []byte("type: text\ndata: 'test content 3'\n"),
		},
	}
	for _, f := range testFiles {
		err := os.WriteFile(f.Path, f.Data, 0644)
		assert.NoError(t, err)
	}

	// Define the configuration for the simulator
	err = os.Mkdir(filepath.Join(tempDir, "redirect"), 0755)
	assert.NoError(t, err)

	config := simulator.Config{
		ServerAddress: "localhost:8081",
		HttpBasePath:  "/http",
		HttpRules: []simulator.HttpRule{
			simulator.NewHttpRule(
				simulator.NewHttpRequestPredicate("GET", "/test"),
				simulator.NewHttpResponseFromString(200, "OK", 100*time.Millisecond),
			),
			simulator.NewHttpRule(
				simulator.NewHttpRequestPredicate("POST", "/data"),
				simulator.NewHttpResponseFromString(201, "Created", 200*time.Millisecond),
			),
			simulator.NewHttpRule(
				simulator.NewHttpRequestPredicate("GET", "/file"),
				simulator.NewHttpResponseFromFile(filepath.Join(tempDir, "http", "response.yaml"), 150*time.Millisecond),
			),
			// https://developers.binance.com/docs/binance-spot-api-docs/rest-api#test-connectivity
			simulator.NewHttpRule(
				simulator.NewHttpRequestPredicate("GET", "/api/v3/ping"),
				simulator.NewHttpRedirectResponder("https://api.binance.com", filepath.Join(tempDir, "record", "http")),
			),
		},
		WsEndpoint: "/ws",
		WsRules: []simulator.WsRule{
			simulator.NewWsRule(
				simulator.NewWsMessagePredicate(simulator.WsMessageText, []byte("ping")),
				simulator.NewWsMessageFromString(simulator.WsMessageText, "pong", 10*time.Millisecond),
			),
			simulator.NewWsRule(
				simulator.NewWsJsonMatcher(`{ "id": 1, "method": "time", "params": [] }`),
				simulator.NewWsMessageFromString(simulator.WsMessageText, `{ "id": 1, "result": 1493285895, "error": null }`, 5*time.Millisecond),
			),
			simulator.NewWsRule(
				simulator.NewWsMessagePredicate(simulator.WsMessageText, []byte("single_file")),
				simulator.NewWsMessageFromFiles(filepath.Join(tempDir, "ws", "single_file")),
			),
			simulator.NewWsRule(
				simulator.NewWsMessagePredicate(simulator.WsMessageText, []byte("multiple_files")),
				simulator.NewWsMessageFromFiles(filepath.Join(tempDir, "ws", "multiple_files")),
			),
			simulator.NewWsSubscriptionRule(
				simulator.NewWsMessagePredicate(simulator.WsMessageText, []byte("subscribe")),
				simulator.NewWsMessageFromString(simulator.WsMessageText, "subscribe_success", 5*time.Millisecond),
				simulator.NewWsMessagePredicate(simulator.WsMessageText, []byte("unsubscribe")),
				simulator.NewWsMessageFromString(simulator.WsMessageText, "unsubscribe_success", 5*time.Millisecond),
				simulator.NewWsMessageFromString(simulator.WsMessageText, "subscription_update", 10*time.Millisecond),
			),
			simulator.NewWsRule(
				simulator.NewWsMessagePredicate(simulator.WsMessageText, []byte("redirect")),
				simulator.NewWsRedirectHandler(),
			),
		},
		WsRedirectUrl: mockServerURL,
		WsRecordDir:   filepath.Join(tempDir, "redirect", "ws"),
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
		testHttp(t, config, filepath.Join(tempDir, "record", "http"))
	})

	// Run WebSocket tests
	t.Run("WebSocket Tests", func(t *testing.T) {
		testWsBasicTest(t, config)
		testWsMessageMatchers(t, config)
		testWsMessageHandlers(t, config)
		testWsSubscription(t, config)
		testWsRedirection(t, config)
	})
}

func testHttp(t *testing.T, config simulator.Config, recordDir string) {
	tests := []struct {
		name                string
		method              string
		path                string
		body                string
		expectedStatus      int
		expectedBody        string
		expectedDelay       time.Duration
		expectedFileContent string
	}{
		{
			name:           "GET /http/test",
			method:         "GET",
			path:           "/http/test",
			expectedStatus: 200,
			expectedBody:   "OK",
			expectedDelay:  100 * time.Millisecond,
		},
		{
			name:           "POST /http/data",
			method:         "POST",
			path:           "/http/data",
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
		},
		{
			name:           "Matched endpoint, unmatched path",
			method:         "GET",
			path:           "/http/unknown",
			expectedStatus: 404,
			expectedBody:   "Invalid request",
		},
		{
			name:           "Response from file",
			method:         "GET",
			path:           "/http/file",
			expectedStatus: 200,
			expectedBody:   "test content",
		},
		{
			// https://developers.binance.com/docs/binance-spot-api-docs/rest-api#test-connectivity
			name:                "Redirect to Binance",
			method:              "GET",
			path:                "/http/api/v3/ping",
			body:                "",
			expectedStatus:      200,
			expectedBody:        "{}",
			expectedFileContent: "status: 200\nbody: |-\n    {}\n",
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
			if tt.expectedDelay > 0 {
				assert.GreaterOrEqual(t, duration, tt.expectedDelay)
				assert.Less(t, duration, tt.expectedDelay+50*time.Millisecond) // Allow for some overhead
			}
			if tt.expectedFileContent != "" {
				files, err := os.ReadDir(recordDir)
				require.NoError(t, err)
				assert.Len(t, files, 1)

				content, err := os.ReadFile(filepath.Join(recordDir, files[0].Name()))
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFileContent, string(content))

				err = os.Remove(filepath.Join(recordDir, files[0].Name()))
				require.NoError(t, err)
			}
		})
	}
}

func testWsBasicTest(t *testing.T, config simulator.Config) {
	tests := []struct {
		name            string
		msgType         websocket.MessageType
		msgData         []byte
		expectedMsgType websocket.MessageType
		expectedMsgData []byte
		expectedDelay   time.Duration
	}{
		{
			name:            "WebSocket basic test",
			msgType:         websocket.MessageText,
			msgData:         []byte("ping"),
			expectedMsgType: websocket.MessageText,
			expectedMsgData: []byte("pong"),
			expectedDelay:   10 * time.Millisecond,
		},
		{
			name:            "WebSocket unmatched message",
			msgType:         websocket.MessageText,
			msgData:         []byte("hello"),
			expectedMsgType: websocket.MessageText,
			expectedMsgData: []byte("Invalid message"),
			expectedDelay:   0,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

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
			}
		})
	}
}

func testWsMessageMatchers(t *testing.T, config simulator.Config) {
	tests := []struct {
		name                string
		msgType             websocket.MessageType
		msgData             []byte
		expectedMsgType     websocket.MessageType
		expectedMsgData     []byte
		expectedDelay       time.Duration
		expectedFileContent string
	}{
		{
			name:            "WebSocket JSON request",
			msgType:         websocket.MessageText,
			msgData:         []byte(`{ "method": "time", "params": [], "id": 1 }`),
			expectedMsgType: websocket.MessageText,
			expectedMsgData: []byte(`{ "id": 1, "result": 1493285895, "error": null }`),
			expectedDelay:   5 * time.Millisecond,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

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
			}
		})
	}
}

func testWsMessageHandlers(t *testing.T, config simulator.Config) {
	tests := []struct {
		name            string
		msgType         websocket.MessageType
		msgData         []byte
		expectedMsgType websocket.MessageType
		expectedMsgData [][]byte
		expectedDelay   time.Duration
	}{
		{
			name:            "WebSocket message from a single file",
			msgType:         websocket.MessageText,
			msgData:         []byte("single_file"),
			expectedMsgType: websocket.MessageText,
			expectedMsgData: [][]byte{
				[]byte("test content 1"),
			},
			expectedDelay: 0,
		},
		{
			name:            "WebSocket message from multiple files",
			msgType:         websocket.MessageText,
			msgData:         []byte("multiple_files"),
			expectedMsgType: websocket.MessageText,
			expectedMsgData: [][]byte{
				[]byte("test content 1"),
				[]byte("test content 2"),
				[]byte("test content 3"),
			},
			expectedDelay: 20 * time.Millisecond,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()

			wsURL := "ws://" + config.ServerAddress + config.WsEndpoint
			conn, _, err := websocket.Dial(ctx, wsURL, nil)
			require.NoError(t, err)
			defer conn.Close(websocket.StatusNormalClosure, "")

			for i := 0; i < 3; i++ {
				start := time.Now()

				err := conn.Write(ctx, tt.msgType, tt.msgData)
				require.NoError(t, err)

				for _, expected := range tt.expectedMsgData {
					msgType, msgData, err := conn.Read(ctx)
					require.NoError(t, err)

					assert.Equal(t, tt.expectedMsgType, msgType)
					assert.Equal(t, expected, msgData)
				}

				duration := time.Since(start)

				assert.GreaterOrEqual(t, duration, tt.expectedDelay)
				assert.Less(t, duration, 2*tt.expectedDelay+10*time.Millisecond)
			}
		})
	}
}

func testWsSubscription(t *testing.T, config simulator.Config) {
	tests := []struct {
		name                  string
		subMsgType            websocket.MessageType
		subMsgData            []byte
		unsubMsgType          websocket.MessageType
		unsubMsgData          []byte
		expectedSubMsgType    websocket.MessageType
		expectedSubMsgData    []byte
		expectedUpdateMsgType websocket.MessageType
		expectedUpdateMsgData []byte
		expectedUnsubMsgType  websocket.MessageType
		expectedUnsubMsgData  []byte
	}{
		{
			name:                  "WebSocket subscription test",
			subMsgType:            websocket.MessageText,
			subMsgData:            []byte("subscribe"),
			unsubMsgType:          websocket.MessageText,
			unsubMsgData:          []byte("unsubscribe"),
			expectedSubMsgType:    websocket.MessageText,
			expectedSubMsgData:    []byte("subscribe_success"),
			expectedUpdateMsgType: websocket.MessageText,
			expectedUpdateMsgData: []byte("subscription_update"),
			expectedUnsubMsgType:  websocket.MessageText,
			expectedUnsubMsgData:  []byte("unsubscribe_success"),
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			wsURL := "ws://" + config.ServerAddress + config.WsEndpoint
			conn, _, err := websocket.Dial(ctx, wsURL, nil)
			require.NoError(t, err)
			defer conn.Close(websocket.StatusNormalClosure, "")

			for i := 0; i < 3; i++ {
				err := conn.Write(ctx, tt.subMsgType, tt.subMsgData)
				require.NoError(t, err)

				subMsgType, subMsgData, err := conn.Read(ctx)
				require.NoError(t, err)

				updateMsgType, updateMsgData, err := conn.Read(ctx)
				require.NoError(t, err)

				err = conn.Write(ctx, tt.unsubMsgType, tt.unsubMsgData)
				require.NoError(t, err)

				unsubMsgType, unsubMsgData, err := conn.Read(ctx)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedSubMsgType, subMsgType)
				assert.Equal(t, tt.expectedSubMsgData, subMsgData)
				assert.Equal(t, tt.expectedUpdateMsgType, updateMsgType)
				assert.Equal(t, tt.expectedUpdateMsgData, updateMsgData)
				assert.Equal(t, tt.expectedUnsubMsgType, unsubMsgType)
				assert.Equal(t, tt.expectedUnsubMsgData, unsubMsgData)
			}
		})
	}
}

func testWsRedirection(t *testing.T, config simulator.Config) {
	tests := []struct {
		name                string
		msgType             websocket.MessageType
		msgData             []byte
		expectedMsgType     websocket.MessageType
		expectedMsgData     []byte
		expectedDelay       time.Duration
		expectedFileContent string
	}{
		{
			name:                "WebSocket Redirect Message",
			msgType:             websocket.MessageText,
			msgData:             []byte("redirect"),
			expectedMsgType:     websocket.MessageText,
			expectedMsgData:     []byte("echoed: redirect"),
			expectedDelay:       0,
			expectedFileContent: "type: text\ndata: |-\n    echoed: redirect\n",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

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

				files, err := os.ReadDir(config.WsRecordDir)
				require.NoError(t, err)
				assert.Len(t, files, 1)

				content, err := os.ReadFile(filepath.Join(config.WsRecordDir, files[0].Name()))
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFileContent, string(content))

				err = os.Remove(filepath.Join(config.WsRecordDir, files[0].Name()))
				require.NoError(t, err)
			}
		})
	}
}
