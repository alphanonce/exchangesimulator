package http

import (
	"time"
)

// Ensure ResponseFromFile implements Responder
var _ Responder = (*ResponseFromFile)(nil)

type ResponseFromFile struct {
	filePath     string
	responseTime time.Duration
}

func NewResponseFromFile(filePath string, responseTime time.Duration) ResponseFromFile {
	return ResponseFromFile{
		filePath:     filePath,
		responseTime: responseTime,
	}
}

func (r ResponseFromFile) Response(_ Request) (Response, error) {
	startTime := time.Now()

	response, err := ReadFromFile(r.filePath)
	if err != nil {
		return response, err
	}

	time.Sleep(time.Until(startTime.Add(r.responseTime)))
	return response, nil
}
