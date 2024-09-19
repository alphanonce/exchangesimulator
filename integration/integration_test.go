package integration

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"alphanonce.com/exchangesimulator/internal/rule"
	"alphanonce.com/exchangesimulator/internal/rule/request_matcher"
	"alphanonce.com/exchangesimulator/internal/rule/responder"
	"alphanonce.com/exchangesimulator/internal/simulator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	addr := "localhost:8081"

	// Define rules for the simulator
	rules := []rule.Rule{
		{
			RequestMatcher: request_matcher.NewRequestPredicate("GET", "/test"),
			Responder:      responder.NewResponseFromString(200, "OK", 100*time.Millisecond),
		},
		{
			RequestMatcher: request_matcher.NewRequestPredicate("POST", "/data"),
			Responder:      responder.NewResponseFromString(201, "Created", 200*time.Millisecond),
		},
	}

	// Create a simulator
	sim := simulator.New(rules)

	// Start the simulator
	go func() {
		err := sim.Run(addr)
		require.NoError(t, err)
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Define test cases
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
			name:           "GET /test",
			method:         "GET",
			path:           "/test",
			expectedStatus: 200,
			expectedBody:   "OK",
			expectedDelay:  100 * time.Millisecond,
		},
		{
			name:           "POST /data",
			method:         "POST",
			path:           "/data",
			body:           "some data",
			expectedStatus: 201,
			expectedBody:   "Created",
			expectedDelay:  200 * time.Millisecond,
		},
		{
			name:           "Unmatched route",
			method:         "GET",
			path:           "/unknown",
			expectedStatus: 200, // Assuming default response for unmatched routes
			expectedBody:   "TODO: not implemented",
			expectedDelay:  0,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			req, err := http.NewRequest(tt.method, "http://"+addr+tt.path, strings.NewReader(tt.body))
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
		})
	}
}
