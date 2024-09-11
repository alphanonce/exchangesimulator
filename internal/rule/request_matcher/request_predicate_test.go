package request_matcher

import (
	"testing"

	"alphanonce.com/exchangesimulator/internal/types"

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
		request   types.Request
		want      bool
	}{
		{
			name:      "Exact match",
			predicate: NewRequestPredicate("GET", "/test"),
			request:   types.Request{Method: "GET", Path: "/test"},
			want:      true,
		},
		{
			name:      "Method mismatch",
			predicate: NewRequestPredicate("GET", "/test"),
			request:   types.Request{Method: "POST", Path: "/test"},
			want:      false,
		},
		{
			name:      "Path mismatch",
			predicate: NewRequestPredicate("GET", "/test"),
			request:   types.Request{Method: "GET", Path: "/other"},
			want:      false,
		},
		{
			name:      "Empty predicate matches all",
			predicate: NewRequestPredicate("", ""),
			request:   types.Request{Method: "POST", Path: "/any"},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.predicate.MatchRequest(tt.request)
			assert.Equal(t, tt.want, got)
		})
	}
}
