package simulator

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"alphanonce.com/exchangesimulator/internal/log"

	"github.com/coder/websocket"
)

type Simulator struct {
	config Config
}

func New(config Config) Simulator {
	return Simulator{
		config: config,
	}
}

func (s Simulator) Run() error {
	return http.ListenAndServe(s.config.ServerAddress, http.HandlerFunc(s.requestHandler))
}

func (s Simulator) requestHandler(w http.ResponseWriter, r *http.Request) {
	requestPath := string(r.URL.Path)
	if strings.HasPrefix(requestPath, s.config.HttpBasePath) {
		s.httpRequestHandler(w, r)
	} else if requestPath == s.config.WsEndpoint {
		s.wsRequestHandler(w, r)
	} else {
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}
}

func (s Simulator) httpRequestHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	request, err := convertHttpRequest(r)
	if err != nil {
		logger.Error("Error reading request body", log.Any("error", err))
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	logger.Debug(
		"Received a HTTP request",
		log.Any("request", request),
	)

	response, endTime := s.simulateHttpResponse(request, startTime)
	convertHttpResponse(w, response)
	time.Sleep(time.Until(endTime))

	logger.Debug(
		"Completed a HTTP request",
		log.Any("request", request),
		log.Any("response", response),
	)
}

func (s Simulator) simulateHttpResponse(request HttpRequest, startTime time.Time) (HttpResponse, time.Time) {
	rule, ok := s.config.GetHttpRule(request)
	if !ok {
		response := HttpResponse{
			StatusCode: http.StatusNotFound,
			Body:       []byte("Invalid request"),
		}
		return response, startTime
	}

	return rule.Response(request), startTime.Add(rule.ResponseTime())
}

func convertHttpRequest(r *http.Request) (HttpRequest, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return HttpRequest{}, err
	}

	request := HttpRequest{
		Method:      r.Method,
		Host:        r.Host,
		Path:        r.URL.Path,
		QueryString: r.URL.RawQuery,
		Body:        bodyBytes,
	}
	return request, nil
}

func convertHttpResponse(w http.ResponseWriter, response HttpResponse) {
	w.WriteHeader(response.StatusCode)
	w.Write(response.Body)
}

func (s Simulator) wsRequestHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		logger.Error("Error upgrading to WebSocket", log.Any("error", err))
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "defer close")

	logger.Info("Succeeded upgrading to WebSocket")

	s.handleWsConnection(r.Context(), conn)
}

func (s Simulator) handleWsConnection(ctx context.Context, conn *websocket.Conn) {
	for {
		incomingType, incomingData, err := conn.Read(ctx)
		if err != nil {
			logger.Error("Error reading WebSocket message", log.Any("error", err))
			return
		}
		startTime := time.Now()

		logger.Debug(
			"Received a WebSocket message",
			log.Any("time", startTime),
			log.Group("msg",
				log.Any("data", incomingData),
				log.String("type", incomingType.String()),
			),
		)

		internalRequest := convertWsMessageToInternal(incomingData, incomingType)
		internalResponse, endTime := s.simulateWsResponse(internalRequest, startTime)
		outgoingData, outgoingType := convertWsMessageFromInternal(internalResponse)
		time.Sleep(time.Until(endTime))

		err = conn.Write(ctx, outgoingType, outgoingData)
		if err != nil {
			logger.Error("Error sending WebSocket message", log.Any("error", err))
			return
		}

		logger.Debug(
			"Sent a WebSocket message",
			log.Any("time", endTime),
			log.Any("response_time", endTime.Sub(startTime)),
			log.Group("msg",
				log.Any("data", outgoingData),
				log.String("type", outgoingType.String()),
			),
		)
	}
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

func (s Simulator) simulateWsResponse(message WsMessage, startTime time.Time) (WsMessage, time.Time) {
	rule, ok := s.config.GetWsRule(message)
	if !ok {
		response := WsMessage{
			Type: WsMessageText,
			Data: []byte("Invalid message"),
		}
		return response, startTime
	}

	return rule.Response(message), startTime.Add(rule.ResponseTime())
}
