package handler

import (
	"io"
	"log/slog"
	"time"

	"alphanonce.com/exchangesimulator/internal/log/config"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

func NewZerologHandler(c config.Config) slog.Handler {
	var out io.Writer
	if c.Format == config.Text {
		out = zerolog.ConsoleWriter{
			Out:        c.Out,
			TimeFormat: time.RFC3339,
		}
	} else {
		out = c.Out
	}
	logger := zerolog.New(out)
	return slogzerolog.Option{
		Level:     c.Level,
		Logger:    &logger,
		AddSource: c.AddSource,
	}.NewZerologHandler()
}
