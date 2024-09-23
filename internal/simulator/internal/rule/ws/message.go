package ws

type MessageType int

const (
	MessageInvalid MessageType = iota
	MessageText
	MessageBinary
)

type Message struct {
	Type MessageType
	Data []byte
}
