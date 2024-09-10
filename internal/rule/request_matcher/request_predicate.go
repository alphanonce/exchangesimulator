package request_matcher

import "alphanonce.com/exchangesimulator/internal/types"

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

func (r RequestPredicate) MatchRequest(request types.Request) bool {
	return (r.method == "" || request.Method == r.method) &&
		(r.path == "" || request.Path == r.path)
}
