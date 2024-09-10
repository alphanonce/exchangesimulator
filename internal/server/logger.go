package server

import (
	"log/slog"
	"os"

	"alphanonce.com/exchangesimulator/internal/log"
	"alphanonce.com/exchangesimulator/internal/log/config"
)

var logger *slog.Logger

func init() {
	logger = log.New(config.Config{
		Out:       os.Stdout,
		Logger:    config.Zerolog,
		Format:    config.Json,
		AddSource: false,
		Level:     slog.LevelDebug,
	}).With(slog.String("package", "server"))
}
