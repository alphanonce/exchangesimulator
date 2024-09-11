package simulator

import (
	"testing"
	"time"

	"alphanonce.com/exchangesimulator/internal/rule"
	"alphanonce.com/exchangesimulator/internal/rule/request_matcher"
	"alphanonce.com/exchangesimulator/internal/rule/responder"
	"alphanonce.com/exchangesimulator/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestNewRuleBasedSimulator(t *testing.T) {
	rules := []rule.Rule{
		{
			RequestMatcher: request_matcher.NewRequestPredicate("GET", "/test"),
			Responder:      responder.NewResponseFromString(200, "OK", time.Second),
		},
	}
	sim := NewRuleBasedSimulator(rules)
	assert.Equal(t, rules, sim.rules)
}

func TestRuleBasedSimulator_Process(t *testing.T) {
	rules := []rule.Rule{
		{
			RequestMatcher: request_matcher.NewRequestPredicate("GET", "/test"),
			Responder:      responder.NewResponseFromString(200, "OK", time.Second),
		},
	}
	sim := NewRuleBasedSimulator(rules)

	tests := []struct {
		name      string
		request   types.Request
		wantCode  int
		wantBody  string
		wantDelay time.Duration
	}{
		{
			name:      "Matching request",
			request:   types.Request{Method: "GET", Path: "/test"},
			wantCode:  200,
			wantBody:  "OK",
			wantDelay: time.Second,
		},
		{
			name:      "Non-matching request",
			request:   types.Request{Method: "POST", Path: "/other"},
			wantCode:  0,
			wantBody:  "TODO: not implemented",
			wantDelay: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			resp, endTime := sim.Process(tt.request, startTime)
			assert.Equal(t, tt.wantCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody, string(resp.Body))
			assert.Equal(t, startTime.Add(tt.wantDelay), endTime)
		})
	}
}
