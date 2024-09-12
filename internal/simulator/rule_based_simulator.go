package simulator

import (
	"slices"
	"time"

	"alphanonce.com/exchangesimulator/internal/rule"
	"alphanonce.com/exchangesimulator/internal/types"
)

// Ensure Simulator implements Interface
var _ Simulator = (*RuleBasedSimulator)(nil)

type RuleBasedSimulator struct {
	rules []rule.Rule
}

func NewRuleBasedSimulator(rules []rule.Rule) RuleBasedSimulator {
	return RuleBasedSimulator{
		rules: rules,
	}
}

func (s RuleBasedSimulator) Process(request types.Request, startTime time.Time) (types.Response, time.Time) {
	r, ok := s.findRule(request)
	if !ok {
		return types.Response{Body: []byte("TODO: not implemented")}, startTime
	}

	return r.Response(request), startTime.Add(r.ResponseTime())
}

func (s RuleBasedSimulator) findRule(request types.Request) (rule.Rule, bool) {
	i := slices.IndexFunc(s.rules, func(r rule.Rule) bool { return r.MatchRequest(request) })
	if i == -1 {
		return rule.Rule{}, false
	}
	return s.rules[i], true
}
