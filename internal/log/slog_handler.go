package log

import (
	"log/slog"
)

func newSlogHandler(c Config) slog.Handler {
	o := slog.HandlerOptions{
		AddSource:   c.AddSource,
		Level:       c.Level,
		ReplaceAttr: nil,
	}

	switch c.Format {
	case DefaultFormat, Text:
		return slog.NewTextHandler(c.Out, &o)
	case Json:
		return slog.NewJSONHandler(c.Out, &o)
	default:
		panic("unknown format")
	}
}
