package main

import (
	"time"

	"alphanonce.com/exchangesimulator/internal/log"
	"alphanonce.com/exchangesimulator/internal/simulator"
)

var logger *log.Logger

func init() {
	logger = log.NewDefault().With(log.String("package", "main"))
}

func main() {
	config := simulator.Config{
		ServerAddress: "localhost:8080",
		HttpBasePath:  "/api",
		HttpRules: []simulator.HttpRule{
			{
				RequestMatcher: simulator.NewHttpRequestPredicate("GET", "/v4/public/platform/status"),
				Responder:      simulator.NewHttpResponseFromString(200, `{"status":"1"}`, time.Second),
			},
			{
				RequestMatcher: simulator.NewHttpRequestPredicate("GET", "/v4/public/ping"),
				Responder:      simulator.NewHttpResponseFromString(200, `["pong"]`, time.Second),
			},
		},
		WsEndpoint: "/ws",
		WsRules: []simulator.WsRule{
			{
				MessageMatcher: simulator.NewWsMessagePredicate(simulator.WsMessageText, []byte("ping\n")),
				MessageHandler: simulator.NewWsMessageFromString(simulator.WsMessageText, "pong", time.Second),
			},
			{
				MessageMatcher: simulator.NewWsMessagePredicate(simulator.WsMessageText, []byte("pong\n")),
				MessageHandler: simulator.NewWsMessageFromString(simulator.WsMessageText, "ping", time.Second),
			},
		},
	}
	sim := simulator.New(config)

	logger.Info("Server is starting", log.String("address", config.ServerAddress))
	err := sim.Run()
	if err != nil {
		logger.Error("Server encountered an error while running", log.Any("error", err))
		return
	}
}
