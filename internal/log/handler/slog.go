package handler

import (
	"log/slog"

	"alphanonce.com/exchangesimulator/internal/log/config"
)

func NewSlogHandler(c config.Config) slog.Handler {
	o := slog.HandlerOptions{
		AddSource:   c.AddSource,
		Level:       c.Level,
		ReplaceAttr: nil,
	}

	switch c.Format {
	case config.DefaultFormat, config.Text:
		return slog.NewTextHandler(c.Out, &o)
	case config.Json:
		return slog.NewJSONHandler(c.Out, &o)
	default:
		panic("unknown format")
	}
}
