package ws

import "context"

//go:generate mockery --name=Connection --inpackage --filename=mock_connection.go

type Connection interface {
	Read(context.Context) (Message, error)
	Write(context.Context, Message) error
}
