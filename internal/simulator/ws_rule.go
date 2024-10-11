package simulator

import (
	"time"

	"alphanonce.com/exchangesimulator/internal/simulator/internal/rule/ws"
)

type WsRule = ws.Rule
type WsConnection = ws.Connection
type WsMessage = ws.Message
type WsMessageType = ws.MessageType

const (
	WsMessageAny    = ws.MessageAny
	WsMessageText   = ws.MessageText
	WsMessageBinary = ws.MessageBinary
)

func NewWsRule(messageMatcher ws.MessageMatcher, messageHandler ws.MessageHandler) ws.RuleImpl {
	return ws.RuleImpl{MessageMatcher: messageMatcher, MessageHandler: messageHandler}
}

func NewWsMessagePredicate(messageType WsMessageType, data []byte) ws.MessagePredicate {
	return ws.NewMessagePredicate(messageType, data)
}

func NewWsJsonMatcher(jsonString string) ws.JsonMessageMatcher {
	return ws.NewJsonMessageMatcher(jsonString)
}

func NewWsMessageFromString(messageType WsMessageType, data string, responseTime time.Duration) ws.MessageFromString {
	return ws.NewMessageFromString(messageType, data, responseTime)
}

func NewWsMessageFromFiles(dirPath string) ws.MessageFromFiles {
	return ws.NewMessageFromFiles(dirPath)
}

func NewWsRedirectHandler() ws.RedirectHandler {
	return ws.NewRedirectHandler()
}
