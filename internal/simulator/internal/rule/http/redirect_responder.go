package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"alphanonce.com/exchangesimulator/internal/log"
)

// Ensure RedirectResponder implements Responder
var _ Responder = (*RedirectResponder)(nil)

type RedirectResponder struct {
	targetUrl string
	recordDir string
}

func NewRedirectResponder(targetUrl string, recordDir string) RedirectResponder {
	return RedirectResponder{
		targetUrl: targetUrl,
		recordDir: recordDir,
	}
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

	response := Response{StatusCode: resp.StatusCode, Body: data}

	if r.recordDir != "" {
		err = r.saveResponseToFile(response)
		if err != nil {
			return Response{}, fmt.Errorf("failed to save to a file: %w", err)
		}
	}
	return response, nil
}

func (r RedirectResponder) saveResponseToFile(response Response) error {
	err := os.MkdirAll(r.recordDir, 0755)
	if err != nil {
		return err
	}

	filename := time.Now().Format(time.RFC3339Nano) + ".yaml"
	path := filepath.Join(r.recordDir, filename)

	err = WriteToFile(path, response)
	if err != nil {
		return err
	}

	logger.Info("Http response recorded", log.Any("path", path))
	return nil
}
