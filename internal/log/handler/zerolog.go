package handler

import (
	"log/slog"
	"time"

	"alphanonce.com/exchangesimulator/internal/log/config"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

func NewZerologHandler(c config.Config) slog.Handler {
	if c.Format == config.Text {
	}
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        c.Out,
		TimeFormat: time.RFC3339,
	})
	return slogzerolog.Option{
		Level:     c.Level,
		Logger:    &logger,
		AddSource: c.AddSource,
	}.NewZerologHandler()
}
