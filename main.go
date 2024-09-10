package main

import (
	"log/slog"
	"os"
	"time"

	"alphanonce.com/exchangesimulator/internal/log"
	"alphanonce.com/exchangesimulator/internal/log/config"
	"alphanonce.com/exchangesimulator/internal/rule"
	"alphanonce.com/exchangesimulator/internal/server"
	"alphanonce.com/exchangesimulator/internal/simulator"
)

var logger *slog.Logger

func init() {
	logger = log.New(config.Config{
		Out:       os.Stdout,
		Logger:    config.Zerolog,
		Format:    config.Json,
		AddSource: false,
		Level:     slog.LevelDebug,
	}).With(slog.String("package", "main"))
}

func main() {
	rules := []rule.Rule{
		rule.NewStaticResponseRule(
			"GET",
			"/api/v4/public/platform/status",
			`{"status":"1"}`,
			time.Second,
		),
		rule.NewStaticResponseRule(
			"GET",
			"/api/v4/public/ping",
			`["pong"]`,
			time.Second,
		),
	}
	s := simulator.NewSimulator(rules)
	sv := server.NewFasthttpServer(s)

	err := sv.Run("localhost:8080")
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
