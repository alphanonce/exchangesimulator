package ws

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MessageType uint8

const (
	MessageAny    MessageType = 0
	MessageText   MessageType = 1
	MessageBinary MessageType = 2
)

type Message struct {
	Type MessageType
	Data []byte
}

func (m *Message) MarshalYAML() (any, error) {
	var typeStr string
	switch m.Type {
	case MessageText:
		typeStr = "text"
	case MessageBinary:
		typeStr = "binary"
	default:
		return nil, errors.New("invalid message type")
	}

	var dataValue string
	if m.Type == MessageText {
		dataValue = string(m.Data)
	} else {
		dataValue = hex.EncodeToString(m.Data)
	}

	return &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: "type",
			},
			{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: typeStr,
			},
			{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: "data",
			},
			{
				Kind:  yaml.ScalarNode,
				Style: yaml.LiteralStyle,
				Tag:   "!!str",
				Value: dataValue,
			},
		},
	}, nil
}

func (m *Message) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return errors.New("expected a mapping node")
	}

	var typeStr, dataStr string
	for i := 0; i < len(value.Content); i += 2 {
		key := value.Content[i].Value
		val := value.Content[i+1].Value

		switch key {
		case "type":
			typeStr = val
		case "data":
			dataStr = val
		default:
			return fmt.Errorf("unexpected key: %s", key)
		}
	}

	switch typeStr {
	case "text":
		m.Type = MessageText
		m.Data = []byte(dataStr)
	case "binary":
		m.Type = MessageBinary
		d, err := hex.DecodeString(dataStr)
		if err != nil {
			return fmt.Errorf("failed to decode binary data: %v", err)
		}
		m.Data = d
	default:
		return fmt.Errorf("invalid message type: %s", typeStr)
	}

	return nil
}

func WriteToFile(path string, message Message) error {
	data, err := yaml.Marshal(&message)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadFromFile(path string) (Message, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Message{}, err
	}

	var m Message
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return Message{}, err
	}

	return m, nil
}
