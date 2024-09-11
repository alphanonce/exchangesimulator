package responder

import (
	"testing"
	"time"

	"alphanonce.com/exchangesimulator/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestNewResponseFromString(t *testing.T) {
	statusCode := 200
	body := `{"key": "value"}`
	responseTime := 100 * time.Millisecond

	r := NewResponseFromString(statusCode, body, responseTime)

	assert.Equal(t, statusCode, r.Response(types.Request{}).StatusCode)
	assert.Equal(t, []byte(body), r.Response(types.Request{}).Body)
	assert.Equal(t, responseTime, r.ResponseTime())
}

func TestResponseFromString_Response(t *testing.T) {
	r := NewResponseFromString(201, `{"status": "created"}`, 50*time.Millisecond)

	response := r.Response(types.Request{})

	assert.Equal(t, 201, response.StatusCode)
	assert.Equal(t, []byte(`{"status": "created"}`), response.Body)
}

func TestResponseFromString_ResponseTime(t *testing.T) {
	expectedResponseTime := 75 * time.Millisecond
	r := NewResponseFromString(200, `{"status": "ok"}`, expectedResponseTime)

	actualResponseTime := r.ResponseTime()

	assert.Equal(t, expectedResponseTime, actualResponseTime)
}
