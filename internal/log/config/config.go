package config

import (
	"io"
	"log/slog"
)

type LoggerType uint8

const (
	DefaultLogger LoggerType = iota
	Slog
	Zerolog
)

type FormatType uint8

const (
	DefaultFormat FormatType = iota
	Text
	Json
)

type Level int8

type Config struct {
	Out io.Writer

	Logger LoggerType
	Format FormatType

	AddSource bool
	Level     slog.Leveler
}
