package http

//go:generate mockery --name=Rule --inpackage --filename=mock_rule.go

type Rule interface {
	RequestMatcher
	Responder
}

type RequestMatcher interface {
	MatchRequest(Request) bool
}

type Responder interface {
	Response(Request) (Response, error)
}

// Ensure RuleImpl implements Rule
var _ Rule = (*RuleImpl)(nil)

type RuleImpl struct {
	RequestMatcher
	Responder
}

func NewRule(requestMatcher RequestMatcher, responder Responder) RuleImpl {
	return RuleImpl{RequestMatcher: requestMatcher, Responder: responder}
}
