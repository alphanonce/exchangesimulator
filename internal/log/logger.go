package log

import (
	"log/slog"
	"os"
)

type Logger = slog.Logger

func New(c Config) *Logger {
	var h slog.Handler

	switch c.Logger {
	case DefaultLogger, Slog:
		h = newSlogHandler(c)
	case Zerolog:
		h = newZerologHandler(c)
	default:
		h = newSlogHandler(c)
	}

	return slog.New(h)
}

func NewDefault() *Logger {
	return New(Config{
		Out:       os.Stdout,
		Logger:    Zerolog,
		Format:    Text,
		AddSource: false,
		Level:     LevelDebug,
	})
}
