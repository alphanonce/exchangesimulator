package server

import (
	"log/slog"

	"alphanonce.com/exchangesimulator/internal/log"
)

var logger *slog.Logger

func init() {
	logger = log.NewDefault().With(slog.String("package", "server"))
}
