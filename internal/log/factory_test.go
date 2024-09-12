package log

import (
	"bytes"
	"testing"

	"alphanonce.com/exchangesimulator/internal/log/config"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config config.Config
	}{
		{
			name: "Default Logger",
			config: config.Config{
				Out:    &bytes.Buffer{},
				Logger: config.DefaultLogger,
				Format: config.DefaultFormat,
			},
		},
		{
			name: "Slog Logger",
			config: config.Config{
				Out:    &bytes.Buffer{},
				Logger: config.Slog,
				Format: config.Json,
			},
		},
		{
			name: "Zerolog Logger",
			config: config.Config{
				Out:    &bytes.Buffer{},
				Logger: config.Zerolog,
				Format: config.Text,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.config)
			assert.NotNil(t, logger)
		})
	}
}

func TestNewDefault(t *testing.T) {
	logger := NewDefault()

	assert.NotNil(t, logger, "Default logger should not be nil")
}
