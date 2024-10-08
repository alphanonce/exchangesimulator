package ws

import "context"

//go:generate mockery --name=Rule --inpackage --filename=mock_rule.go

type Rule interface {
	MessageMatcher
	MessageHandler
}

type MessageMatcher interface {
	MatchMessage(Message) bool
}

type MessageHandler interface {
	Handle(context.Context, Message, Connection, Connection) error
}

// Ensure RuleImpl implements Rule
var _ Rule = (*RuleImpl)(nil)

type RuleImpl struct {
	MessageMatcher
	MessageHandler
}

func NewRule(messageMatcher MessageMatcher, messageHandler MessageHandler) RuleImpl {
	return RuleImpl{MessageMatcher: messageMatcher, MessageHandler: messageHandler}
}
