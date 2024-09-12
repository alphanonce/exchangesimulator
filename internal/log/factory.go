package log

import (
	"log/slog"
	"os"
)

func New(c Config) *slog.Logger {
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

func NewDefault() *slog.Logger {
	return New(Config{
		Out:       os.Stdout,
		Logger:    Zerolog,
		Format:    Text,
		AddSource: false,
		Level:     slog.LevelDebug,
	})
}
