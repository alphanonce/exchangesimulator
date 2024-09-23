package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequestPredicate(t *testing.T) {
	rp := NewRequestPredicate("GET", "/test")
	assert.Equal(t, "GET", rp.method)
	assert.Equal(t, "/test", rp.path)
}

func TestRequestPredicate_MatchRequest(t *testing.T) {
	tests := []struct {
		name      string
		predicate RequestPredicate
		request   Request
		expected  bool
	}{
		{
			name:      "Exact match",
			predicate: NewRequestPredicate("GET", "/test"),
			request:   Request{Method: "GET", Path: "/test"},
			expected:  true,
		},
		{
			name:      "Method mismatch",
			predicate: NewRequestPredicate("GET", "/test"),
			request:   Request{Method: "POST", Path: "/test"},
			expected:  false,
		},
		{
			name:      "Path mismatch",
			predicate: NewRequestPredicate("GET", "/test"),
			request:   Request{Method: "GET", Path: "/other"},
			expected:  false,
		},
		{
			name:      "Empty predicate matches all",
			predicate: NewRequestPredicate("", ""),
			request:   Request{Method: "POST", Path: "/any"},
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.predicate.MatchRequest(tt.request)
			assert.Equal(t, tt.expected, got)
		})
	}
}
