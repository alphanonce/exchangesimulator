package log

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

type Leveler = slog.Leveler
type Level = slog.Level

const (
	LevelDebug Level = slog.LevelDebug
	LevelInfo  Level = slog.LevelInfo
	LevelWarn  Level = slog.LevelWarn
	LevelError Level = slog.LevelError
)

type Config struct {
	Out io.Writer

	Logger LoggerType
	Format FormatType

	AddSource bool
	Level     Leveler
}
