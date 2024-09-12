package log

import (
	"io"
	"log/slog"
	"time"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

func newZerologHandler(c Config) slog.Handler {
	var out io.Writer
	if c.Format == Text {
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
