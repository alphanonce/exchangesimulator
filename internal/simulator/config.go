package simulator

import (
	"slices"
	"strings"
)

type Config struct {
	// HttpBasePath and WsEndpoint must not be a prefix of each other
	ServerAddress string
	HttpBasePath  string
	HttpRules     []HttpRule
	WsEndpoint    string
	WsRules       []WsRule
	WsRedirectUrl string
	WsRecordDir   string
}

func (c *Config) GetHttpRule(request HttpRequest) (HttpRule, bool) {
	suffix, found := strings.CutPrefix(request.Path, c.HttpBasePath)
	if !found {
		return HttpRule{}, false
	}

	tmpRequest := request
	tmpRequest.Path = suffix
	i := slices.IndexFunc(c.HttpRules, func(r HttpRule) bool { return r.MatchRequest(tmpRequest) })
	if i == -1 {
		return HttpRule{}, false
	}

	return c.HttpRules[i], true
}

func (c *Config) GetWsRule(message WsMessage) (WsRule, bool) {
	i := slices.IndexFunc(c.WsRules, func(r WsRule) bool { return r.MatchMessage(message) })
	if i == -1 {
		return nil, false
	}

	return c.WsRules[i], true
}
