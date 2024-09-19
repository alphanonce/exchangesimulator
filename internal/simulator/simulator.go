package simulator

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"alphanonce.com/exchangesimulator/internal/log"
	"alphanonce.com/exchangesimulator/internal/rule"
	"alphanonce.com/exchangesimulator/internal/types"
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
	return http.ListenAndServe(address, http.HandlerFunc(s.requestHandler))
}

func (s Simulator) requestHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	request, err := getRequest(r)
	if err != nil {
		logger.Error("Error reading request body", log.Any("error", err))
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	logger.Debug(
		"Received a request",
		log.Any("start_time", startTime),
		log.String("request", fmt.Sprintf("%+v", request)),
	)

	response, endTime := s.process(request, startTime)
	setResponse(w, response)
	time.Sleep(time.Until(endTime))

	logger.Debug(
		"Completed a request",
		log.Any("start_time", startTime),
		log.Any("end_time", endTime),
		log.String("request", fmt.Sprintf("%+v", request)),
		log.String("response", fmt.Sprintf("%+v", response)),
	)
}

func (s Simulator) process(request types.Request, startTime time.Time) (types.Response, time.Time) {
	r, ok := s.findRule(request)
	if !ok {
		return types.Response{StatusCode: 404, Body: []byte("TODO: not implemented")}, startTime
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

func getRequest(r *http.Request) (types.Request, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return types.Request{}, err
	}

	request := types.Request{
		Method:      r.Method,
		Host:        r.Host,
		Path:        r.URL.Path,
		QueryString: "", // TODO
		Body:        bodyBytes,
	}
	return request, nil
}

func setResponse(w http.ResponseWriter, response types.Response) {
	w.WriteHeader(response.StatusCode)
	w.Write(response.Body)
}

