package ws

type MessageType uint8

const (
	MessageAny    MessageType = 0
	MessageText   MessageType = 1
	MessageBinary MessageType = 2
)

type Message struct {
	Type MessageType `json:"type"`
	Data []byte      `json:"data"`
}
