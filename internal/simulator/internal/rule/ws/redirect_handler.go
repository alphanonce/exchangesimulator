package ws

import (
	"context"
)

// Ensure RedirectHandler implements MessageHandler
var _ MessageHandler = (*RedirectHandler)(nil)

type RedirectHandler struct{}

func NewRedirectHandler() RedirectHandler {
	return RedirectHandler{}
}

func (r RedirectHandler) Handle(ctx context.Context, message Message, _ Connection, connServer Connection) error {
	err := connServer.Write(ctx, message)
	return err
}
