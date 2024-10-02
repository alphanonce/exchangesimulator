package ws

import (
	"encoding/hex"
	"encoding/json"
	"os"
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

type messageOnFile struct {
	Type        MessageType `json:"type"`
	EncodedData string      `json:"data"`
}

func WriteToFile(path string, message Message) error {
	mof := messageOnFile{Type: message.Type}
	if message.Type == MessageText {
		mof.EncodedData = string(message.Data)
	} else {
		mof.EncodedData = hex.EncodeToString(message.Data)
	}

	data, err := json.Marshal(mof)
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

	var mof messageOnFile
	err = json.Unmarshal(data, &mof)
	if err != nil {
		return Message{}, err
	}

	message := Message{Type: mof.Type}
	if message.Type == MessageText {
		message.Data = []byte(mof.EncodedData)
	} else {
		d, err := hex.DecodeString(mof.EncodedData)
		if err != nil {
			return Message{}, err
		}
		message.Data = d
	}

	return message, nil
}
