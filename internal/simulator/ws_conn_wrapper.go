package simulator

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	"alphanonce.com/exchangesimulator/internal/log"
	"github.com/coder/websocket"
)

// Ensure WsConnWrapper implements WsConnection
var _ WsConnection = (*WsConnWrapper)(nil)

type WsConnWrapper struct {
	conn *websocket.Conn
}

func wrapConnection(conn *websocket.Conn) WsConnWrapper {
	return WsConnWrapper{conn: conn}
}

func (w WsConnWrapper) Read(ctx context.Context) (WsMessage, error) {
	if w.conn == nil {
		return WsMessage{}, errors.New("connection is nil")
	}

	incomingType, incomingData, err := w.conn.Read(ctx)
	if err != nil {
		return WsMessage{}, err
	}

	logger.Debug(
		"Received a WebSocket message",
		log.Any("time", time.Now()),
		log.Group("msg",
			log.String("data", hex.EncodeToString(incomingData)),
			log.String("type", incomingType.String()),
		),
	)

	return convertWsMessageToInternal(incomingData, incomingType), nil
}

func convertWsMessageToInternal(data []byte, messageType websocket.MessageType) WsMessage {
	var t WsMessageType
	switch messageType {
	case websocket.MessageText:
		t = WsMessageText
	case websocket.MessageBinary:
		t = WsMessageBinary
	}

	return WsMessage{
		Type: t,
		Data: data,
	}
}

func (w WsConnWrapper) Write(ctx context.Context, message WsMessage) error {
	if w.conn == nil {
		return errors.New("connection is nil")
	}

	outgoingData, outgoingType := convertWsMessageFromInternal(message)

	err := w.conn.Write(ctx, outgoingType, outgoingData)
	if err != nil {
		return err
	}

	logger.Debug(
		"Sent a WebSocket message",
		log.Group("msg",
			log.String("data", hex.EncodeToString(outgoingData)),
			log.String("type", outgoingType.String()),
		),
	)

	return nil
}

func convertWsMessageFromInternal(message WsMessage) ([]byte, websocket.MessageType) {
	var messageType websocket.MessageType
	switch message.Type {
	case WsMessageText:
		messageType = websocket.MessageText
	case WsMessageBinary:
		messageType = websocket.MessageBinary
	}

	return message.Data, messageType
}
