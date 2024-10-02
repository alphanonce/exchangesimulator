package ws

import (
	"encoding/json"
	"testing"
)

func TestNewJsonMessageMatcher(t *testing.T) {
	tests := []struct {
		name        string
		jsonString  string
		shouldPanic bool
	}{
		{
			name:        "Valid JSON",
			jsonString:  `{"key": "value"}`,
			shouldPanic: false,
		},
		{
			name:        "Invalid JSON",
			jsonString:  `{"key": "value"`,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewJsonMessageMatcher should have panicked")
					}
				}()
			}
			matcher := NewJsonMessageMatcher(tt.jsonString)
			if !tt.shouldPanic && matcher.data == nil {
				t.Errorf("Expected non-nil data, got nil")
			}
		})
	}
}

func TestJsonMessageMatcher_MatchMessage(t *testing.T) {
	tests := []struct {
		name          string
		matcherJSON   string
		messageType   MessageType
		messageData   string
		expectedMatch bool
	}{
		{
			name:          "Matching JSON - object",
			matcherJSON:   `{"number1": 123, "number2": 456.789, "string": "abc"}`,
			messageType:   MessageText,
			messageData:   `{"number2": 456.789, "string": "abc", "number1": 123}`,
			expectedMatch: true,
		},
		{
			name:          "Matching JSON - array",
			matcherJSON:   `[1, "abc", 3.5]`,
			messageType:   MessageText,
			messageData:   `  [  1   ,  "abc"  ,3.5] `,
			expectedMatch: true,
		},
		{
			name:          "Non-matching JSON",
			matcherJSON:   `{"key": "value"}`,
			messageType:   MessageText,
			messageData:   `{"key": "different"}`,
			expectedMatch: false,
		},
		{
			name:          "Binary message type",
			matcherJSON:   `{"key": "value"}`,
			messageType:   MessageBinary,
			messageData:   `{"key": "value"}`,
			expectedMatch: false,
		},
		{
			name:          "Invalid JSON in message",
			matcherJSON:   `{"key": "value"}`,
			messageType:   MessageText,
			messageData:   `{"key": "value"`,
			expectedMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewJsonMessageMatcher(tt.matcherJSON)
			message := Message{
				Type: tt.messageType,
				Data: json.RawMessage(tt.messageData),
			}
			result := matcher.MatchMessage(message)
			if result != tt.expectedMatch {
				t.Errorf("Expected match: %v, got: %v", tt.expectedMatch, result)
			}
		})
	}
}
