package main

import (
	"time"

	"alphanonce.com/exchangesimulator/internal/log"
	"alphanonce.com/exchangesimulator/internal/rule"
	"alphanonce.com/exchangesimulator/internal/rule/request_matcher"
	"alphanonce.com/exchangesimulator/internal/rule/responder"
	"alphanonce.com/exchangesimulator/internal/server"
	"alphanonce.com/exchangesimulator/internal/simulator"
)

var logger *log.Logger

func init() {
	logger = log.NewDefault().With(log.String("package", "main"))
}

func main() {
	rules := []rule.Rule{
		{
			RequestMatcher: request_matcher.NewRequestPredicate("GET", "/api/v4/public/platform/status"),
			Responder:      responder.NewResponseFromString(200, `{"status":"1"}`, time.Second),
		},
		{
			RequestMatcher: request_matcher.NewRequestPredicate("GET", "/api/v4/public/ping"),
			Responder:      responder.NewResponseFromString(200, `["pong"]`, time.Second),
		},
	}
	s := simulator.NewRuleBasedSimulator(rules)
	sv := server.NewFasthttpServer(s)
	address := "localhost:8080"

	logger.Info("Server is starting", log.String("address", address))
	err := sv.Run(address)
	if err != nil {
		logger.Error("Server encountered an error while running", log.Any("error", err))
		return
	}
}
