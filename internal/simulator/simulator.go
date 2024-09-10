package simulator

import (
	"slices"
	"time"

	"alphanonce.com/exchangesimulator/internal/rule"
	"alphanonce.com/exchangesimulator/internal/types"
)

type Simulator struct {
	rules []rule.Rule
}

func NewSimulator(rules []rule.Rule) Simulator {
	return Simulator{
		rules: rules,
	}
}

func (s Simulator) Process(request types.Request, startTime time.Time) (types.Response, time.Time) {
	r, ok := s.findRule(request)
	if !ok {
		return types.Response{Body: []byte("TODO: not implemented")}, startTime
	}

	return r.Response(request), startTime.Add(r.ResponseTime())
}

func (s Simulator) findRule(request types.Request) (rule.Rule, bool) {
	i := slices.IndexFunc(s.rules, func(r rule.Rule) bool { return r.MatchRequest(request) })
	if i == -1 {
		return nil, false
	}
	return s.rules[i], true
}
