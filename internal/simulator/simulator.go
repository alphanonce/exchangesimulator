package simulator

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"alphanonce.com/exchangesimulator/internal/log"
	"alphanonce.com/exchangesimulator/internal/simulator/internal/rule/ws"

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
	if s.config.HttpBasePath != "" && strings.HasPrefix(requestPath, s.config.HttpBasePath) {
		s.httpRequestHandler(w, r)
	} else if s.config.WsEndpoint != "" && requestPath == s.config.WsEndpoint {
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
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		logger.Error("Error upgrading to WebSocket", log.Any("error", err))
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")
	logger.Info("Succeeded upgrading to WebSocket")
	connClient := wrapConnection(conn)

	var connServer WsConnection
	if s.config.WsRedirectUrl != "" {
		conn, _, err := websocket.Dial(ctx, s.config.WsRedirectUrl, nil)
		if err != nil {
			logger.Error("Error connecting to WebSocket server", log.String("url", s.config.WsRedirectUrl), log.Any("error", err))
			http.Error(w, "Failed to connect to WebSocket server", http.StatusInternalServerError)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "")
		logger.Info("Succeeded connecting to WebSocket", log.String("url", s.config.WsRedirectUrl))
		connServer = wrapConnection(conn)

		go func() {
			err := s.redirectWsMessageFromServerToClient(ctx, connClient, connServer)
			if err != nil {
				logger.Error("Error redirecting messages from server", log.Any("error", err))
				cancel()
			}
		}()
	}

	err = s.handleWsConnection(ctx, connClient, connServer)
	if err != nil {
		logger.Error("Error handling websocket messages", log.Any("error", err))
		return
	}
}

func (s Simulator) redirectWsMessageFromServerToClient(ctx context.Context, connClient WsConnection, connServer WsConnection) error {
	for {
		message, err := connServer.Read(ctx)
		if err != nil {
			return fmt.Errorf("failed to read from server: %w", err)
		}

		if s.config.WsRecordDir != "" {
			err = s.saveMessageToFile(message, s.config.WsRecordDir)
			if err != nil {
				return fmt.Errorf("failed to save to a file: %w", err)
			}
		}

		err = connClient.Write(ctx, message)
		if err != nil {
			return fmt.Errorf("failed to write to client: %w", err)
		}
	}
}

func (s Simulator) saveMessageToFile(message WsMessage, dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	filename := time.Now().Format(time.RFC3339Nano) + ".yaml"
	path := filepath.Join(dir, filename)

	err = ws.WriteToFile(path, message)
	if err != nil {
		return err
	}

	logger.Info("WebSocket message recorded", log.Any("path", path))
	return nil
}

func (s Simulator) handleWsConnection(ctx context.Context, connClient WsConnection, connServer WsConnection) error {
	for {
		incomingMsg, err := connClient.Read(ctx)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		err = s.simulateWsResponse(ctx, incomingMsg, connClient, connServer)
		if err != nil {
			return fmt.Errorf("failed to handle message: %w", err)
		}
	}
}

func (s Simulator) simulateWsResponse(ctx context.Context, message WsMessage, connClient WsConnection, connServer WsConnection) error {
	rule, ok := s.config.GetWsRule(message)
	if !ok {
		response := WsMessage{
			Type: WsMessageText,
			Data: []byte("Invalid message"),
		}
		return connClient.Write(ctx, response)
	}

	return rule.Handle(ctx, message, connClient, connServer)
}
