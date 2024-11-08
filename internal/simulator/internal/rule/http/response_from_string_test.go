package http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewResponseFromString(t *testing.T) {
	statusCode := 200
	body := `{"key": "value"}`
	responseTime := 100 * time.Millisecond

	r := NewResponseFromString(statusCode, body, responseTime)

	assert.Equal(t, statusCode, r.statusCode)
	assert.Equal(t, body, r.body)
	assert.Equal(t, responseTime, r.responseTime)
}

func TestResponseFromString_Response(t *testing.T) {
	r := NewResponseFromString(201, `{"status": "created"}`, 50*time.Millisecond)

	start := time.Now()
	resp, err := r.Response(Request{})
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, Response{StatusCode: 201, Body: []byte(`{"status": "created"}`)}, resp)
	assert.GreaterOrEqual(t, duration, 50*time.Millisecond)
	assert.LessOrEqual(t, duration, 100*time.Millisecond)
}
