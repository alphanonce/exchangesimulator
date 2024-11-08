package ws

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Ensure JsonMessageMatcher implements MessageMatcher
var _ MessageMatcher = (*JsonMessageMatcher)(nil)

type JsonMessageMatcher struct {
	data any
}

func NewJsonMessageMatcher(jsonString string) JsonMessageMatcher {
	m := JsonMessageMatcher{}
	err := json.Unmarshal([]byte(jsonString), &m.data)
	if err != nil {
		panic(fmt.Sprintf("invalid json string `%s`: %s", jsonString, err.Error()))
	}
	return m
}

func (p JsonMessageMatcher) MatchMessage(message Message) bool {
	if message.Type != MessageText {
		return false
	}

	var data any
	err := json.Unmarshal(message.Data, &data)
	return err == nil && reflect.DeepEqual(p.data, data)
}
