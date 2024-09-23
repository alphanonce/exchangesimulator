package http

// Ensure RequestPredicate implements RequestMatcher
var _ RequestMatcher = (*RequestPredicate)(nil)

type RequestPredicate struct {
	method string
	path   string
}

func NewRequestPredicate(method string, path string) RequestPredicate {
	return RequestPredicate{
		method: method,
		path:   path,
	}
}

func (r RequestPredicate) MatchRequest(request Request) bool {
	return (r.method == "" || request.Method == r.method) &&
		(r.path == "" || request.Path == r.path)
}
