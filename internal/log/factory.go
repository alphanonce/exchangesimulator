package log

import (
	"log/slog"
	"os"

	"alphanonce.com/exchangesimulator/internal/log/config"
	"alphanonce.com/exchangesimulator/internal/log/handler"
)

func New(c config.Config) *slog.Logger {
	var h slog.Handler

	switch c.Logger {
	case config.DefaultLogger, config.Slog:
		h = handler.NewSlogHandler(c)
	case config.Zerolog:
		h = handler.NewZerologHandler(c)
	default:
		h = handler.NewSlogHandler(c)
	}

	return slog.New(h)
}

func NewDefault() *slog.Logger {
	return New(config.Config{
		Out:       os.Stdout,
		Logger:    config.Zerolog,
		Format:    config.Text,
		AddSource: false,
		Level:     slog.LevelDebug,
	})
}
