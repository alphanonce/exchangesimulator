package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "Default Logger",
			config: Config{
				Out:    &bytes.Buffer{},
				Logger: DefaultLogger,
				Format: DefaultFormat,
			},
		},
		{
			name: "Slog Logger",
			config: Config{
				Out:    &bytes.Buffer{},
				Logger: Slog,
				Format: Json,
			},
		},
		{
			name: "Zerolog Logger",
			config: Config{
				Out:    &bytes.Buffer{},
				Logger: Zerolog,
				Format: Text,
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
