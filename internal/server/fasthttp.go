package server

import (
	"log/slog"
	"time"

	"alphanonce.com/exchangesimulator/internal/simulator"
	"alphanonce.com/exchangesimulator/internal/types"

	"github.com/valyala/fasthttp"
)

type FasthttpServer struct {
	simulator simulator.Simulator
}

func NewFasthttpServer(s simulator.Simulator) FasthttpServer {
	return FasthttpServer{
		simulator: s,
	}
}

func (s FasthttpServer) Run(address string) error {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		logger.Debug(
			"Received a request",
			slog.Any("start_time", ctx.Time()),
			slog.String("request", ctx.Request.String()),
		)

		request := getRequest(ctx)
		response, endTime := s.simulator.Process(request, ctx.Time())
		setResponse(ctx, response)
		time.Sleep(time.Until(endTime))

		logger.Debug(
			"Completed a request",
			slog.Any("start_time", ctx.Time()),
			slog.Any("end_time", endTime),
			slog.String("request", ctx.Request.String()),
			slog.String("response", ctx.Response.String()),
		)
	}

	err := fasthttp.ListenAndServe(address, requestHandler)
	if err != nil {
		return err
	}

	return nil
}

func getRequest(ctx *fasthttp.RequestCtx) types.Request {
	return types.Request{
		Method:      string(ctx.Request.Header.Method()),
		Host:        string(ctx.Request.Header.Host()),
		Path:        string(ctx.Request.URI().Path()),
		QueryString: string(ctx.Request.URI().QueryString()),
		Body:        ctx.Request.Body(),
	}
}

func setResponse(ctx *fasthttp.RequestCtx, response types.Response) {
	ctx.Response.SetStatusCode(response.StatusCode)
	ctx.Response.SetBody(response.Body)
}
