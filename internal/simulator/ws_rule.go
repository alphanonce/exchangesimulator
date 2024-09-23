package simulator

import (
	"time"

	"alphanonce.com/exchangesimulator/internal/simulator/internal/rule/ws"
)

type WsRule = ws.Rule
type WsMessage = ws.Message
type WsMessageType = ws.MessageType

const (
	WsMessageText   = ws.MessageText
	WsMessageBinary = ws.MessageBinary
)

func NewWsMessagePredicate(messageType WsMessageType, data []byte) ws.MessagePredicate {
	return ws.NewMessagePredicate(messageType, data)
}

func NewWsMessageFromString(messageType WsMessageType, data string, responseTime time.Duration) ws.MessageFromString {
	return ws.NewMessageFromString(messageType, data, responseTime)
}
