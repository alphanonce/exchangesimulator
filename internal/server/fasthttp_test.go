package server

import (
	"net"
	"testing"
	"time"

	"alphanonce.com/exchangesimulator/internal/mocks"
	"alphanonce.com/exchangesimulator/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestFasthttpServer_serve(t *testing.T) {
	// Create a mock simulator
	mockSim := new(mocks.Simulator)

	// Create a test server
	server := NewFasthttpServer(mockSim)
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := server.serve(ln)
		assert.NoError(t, err)
	}()

	// Create a test client
	client := fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	// Prepare a test request
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("http://" + ln.Addr().String() + "/test")
	req.Header.SetMethod("GET")

	// Set up mock expectations
	expectedResponse := types.Response{
		StatusCode: 200,
		Body:       []byte("Test response"),
	}
	mockSim.On("Process", mock.AnythingOfType("types.Request"), mock.AnythingOfType("time.Time")).
		Return(expectedResponse, time.Now().Add(100*time.Millisecond))

	// Perform the request
	startTime := time.Now()
	err := client.Do(req, resp)
	endTime := time.Now()

	assert.NoError(t, err)

	// Assert the response
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "Test response", string(resp.Body()))

	// Check if the delay was applied
	assert.True(t, endTime.Sub(startTime) >= 100*time.Millisecond)

	// Verify that the mock was called
	mockSim.AssertExpectations(t)
}

func TestGetRequest(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetHost("example.com")
	ctx.Request.SetRequestURI("/test?param=value")
	ctx.Request.SetBody([]byte("request body"))

	request := getRequest(ctx)

	assert.Equal(t, "POST", request.Method)
	assert.Equal(t, "example.com", request.Host)
	assert.Equal(t, "/test", request.Path)
	assert.Equal(t, "param=value", request.QueryString)
	assert.Equal(t, []byte("request body"), request.Body)
}

func TestSetResponse(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	response := types.Response{
		StatusCode: 201,
		Body:       []byte("response body"),
	}

	setResponse(ctx, response)

	assert.Equal(t, 201, ctx.Response.StatusCode())
	assert.Equal(t, []byte("response body"), ctx.Response.Body())
}
