package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Ensure RedirectResponder implements Responder
var _ Responder = (*RedirectResponder)(nil)

type RedirectResponder struct {
	targetUrl string
}

func NewRedirectResponder(targetUrl string) RedirectResponder {
	return RedirectResponder{targetUrl: targetUrl}
}

func (r RedirectResponder) Response(request Request) (Response, error) {
	url, err := url.Parse(r.targetUrl)
	if err != nil {
		return Response{}, fmt.Errorf("invalid target URL: %w", err)
	}
	url.Path = request.Path
	url.RawQuery = request.QueryString
	req, err := http.NewRequest(request.Method, url.String(), bytes.NewReader(request.Body))
	if err != nil {
		return Response{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = request.Header

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return Response{}, fmt.Errorf("failed to reach target server: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed to read response data: %w", err)
	}

	return Response{StatusCode: resp.StatusCode, Body: data}, nil
}
