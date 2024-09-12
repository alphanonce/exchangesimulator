package log

import (
	"log/slog"
	"time"
)

func Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

func Bool(key string, v bool) slog.Attr {
	return slog.Bool(key, v)
}

func Duration(key string, v time.Duration) slog.Attr {
	return slog.Duration(key, v)
}

func Float64(key string, v float64) slog.Attr {
	return slog.Float64(key, v)
}

func Group(key string, args ...any) slog.Attr {
	return slog.Group(key, args...)
}

func Int(key string, v int) slog.Attr {
	return slog.Int(key, v)
}

func Int64(key string, v int64) slog.Attr {
	return slog.Int64(key, v)
}

func String(key string, v string) slog.Attr {
	return slog.String(key, v)
}

func Time(key string, v time.Time) slog.Attr {
	return slog.Time(key, v)
}

func Uint64(key string, v uint64) slog.Attr {
	return slog.Uint64(key, v)
}
