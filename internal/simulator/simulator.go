package simulator

import (
	"slices"
	"time"

	"alphanonce.com/exchangesimulator/internal/log"
	"alphanonce.com/exchangesimulator/internal/rule"
	"alphanonce.com/exchangesimulator/internal/types"
	"github.com/valyala/fasthttp"
)

type Simulator struct {
	rules []rule.Rule
}

func New(rules []rule.Rule) Simulator {
	return Simulator{
		rules: rules,
	}
}

func (s Simulator) Run(address string) error {
	return fasthttp.ListenAndServe(address, s.requestHandler)
}

func (s Simulator) requestHandler(ctx *fasthttp.RequestCtx) {
	logger.Debug(
		"Received a request",
		log.Any("start_time", ctx.Time()),
		log.String("request", ctx.Request.String()),
	)

	request := getRequest(ctx)
	startTime := ctx.Time()
	response, endTime := s.process(request, startTime)
	setResponse(ctx, response)
	time.Sleep(time.Until(endTime))

	logger.Debug(
		"Completed a request",
		log.Any("start_time", ctx.Time()),
		log.Any("end_time", endTime),
		log.String("request", ctx.Request.String()),
		log.String("response", ctx.Response.String()),
	)
}

func (s Simulator) process(request types.Request, startTime time.Time) (types.Response, time.Time) {
	r, ok := s.findRule(request)
	if !ok {
		return types.Response{Body: []byte("TODO: not implemented")}, startTime
	}

	return r.Response(request), startTime.Add(r.ResponseTime())
}

func (s Simulator) findRule(request types.Request) (rule.Rule, bool) {
	i := slices.IndexFunc(s.rules, func(r rule.Rule) bool { return r.MatchRequest(request) })
	if i == -1 {
		return rule.Rule{}, false
	}
	return s.rules[i], true
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

