package ws

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscriptionRule_MatchMessage(t *testing.T) {
	tests := []struct {
		name       string
		subMatch   bool
		unsubMatch bool
		expected   bool
	}{
		{
			name:       "subscription message",
			subMatch:   true,
			unsubMatch: false,
			expected:   true,
		},
		{
			name:       "unsubscription message",
			subMatch:   false,
			unsubMatch: true,
			expected:   true,
		},
		{
			name:       "unmatched message",
			subMatch:   false,
			unsubMatch: false,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subMatcher := NewMockMessageMatcher(t)
			unsubMatcher := NewMockMessageMatcher(t)

			rule := NewSubscriptionRule(subMatcher, nil, unsubMatcher, nil, nil)

			subMatcher.On("MatchMessage", mock.Anything).Return(tt.subMatch).Maybe()
			unsubMatcher.On("MatchMessage", mock.Anything).Return(tt.unsubMatch).Maybe()

			assert.Equal(t, tt.expected, rule.MatchMessage(Message{}))
		})
	}
}

func TestSubscriptionRule_Handle_Subscription(t *testing.T) {
	subMatcher := NewMockMessageMatcher(t)
	subHandler := NewMockMessageHandler(t)
	unsubMatcher := NewMockMessageMatcher(t)
	unsubHandler := NewMockMessageHandler(t)
	updateHandler := NewMockMessageHandler(t)
	rule := NewSubscriptionRule(subMatcher, subHandler, unsubMatcher, unsubHandler, updateHandler)

	ctx := context.Background()
	msg := Message{}
	connClient := NewMockConnection(t)
	connServer := NewMockConnection(t)

	// Set up expectations
	subMatcher.On("MatchMessage", msg).Return(true)
	subHandler.On("Handle", ctx, msg, connClient, connServer).Return(nil)
	updateHandler.On("Handle", mock.Anything, msg, connClient, connServer).Return(nil)

	err := rule.Handle(ctx, msg, connClient, connServer)

	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond) // Wait a bit to ensure the goroutine starts
}

func TestSubscriptionRule_Handle_Unsubscription(t *testing.T) {
	subMatcher := NewMockMessageMatcher(t)
	subHandler := NewMockMessageHandler(t)
	unsubMatcher := NewMockMessageMatcher(t)
	unsubHandler := NewMockMessageHandler(t)
	updateHandler := NewMockMessageHandler(t)
	rule := NewSubscriptionRule(subMatcher, subHandler, unsubMatcher, unsubHandler, updateHandler)

	// First, simulate a subscription
	isCancelFuncCalled := false
	rule.updateEnabled = true
	rule.updateCancelFunc = func() {
		isCancelFuncCalled = true
	}

	ctx := context.Background()
	msg := Message{}
	connClient := NewMockConnection(t)
	connServer := NewMockConnection(t)

	// Set up expectations
	subMatcher.On("MatchMessage", msg).Return(false)
	unsubHandler.On("Handle", ctx, msg, connClient, connServer).Return(nil)

	err := rule.Handle(ctx, msg, connClient, connServer)

	assert.NoError(t, err)
	assert.False(t, rule.updateEnabled)
	assert.True(t, isCancelFuncCalled)
}

func TestSubscriptionRule_Handle_MultipleSubscriptions(t *testing.T) {
	subMatcher := NewMockMessageMatcher(t)
	subHandler := NewMockMessageHandler(t)
	unsubMatcher := NewMockMessageMatcher(t)
	unsubHandler := NewMockMessageHandler(t)
	updateHandler := NewMockMessageHandler(t)
	rule := NewSubscriptionRule(subMatcher, subHandler, unsubMatcher, unsubHandler, updateHandler)

	ctx := context.Background()
	msg := Message{}
	connClient := NewMockConnection(t)
	connServer := NewMockConnection(t)

	// Set up expectations
	subMatcher.On("MatchMessage", msg).Return(true)
	subHandler.On("Handle", ctx, msg, connClient, connServer).Return(nil).Twice()
	updateHandler.On("Handle", mock.Anything, msg, connClient, connServer).Return(nil).Once()

	// First subscription
	err := rule.Handle(ctx, msg, connClient, connServer)
	assert.NoError(t, err)

	// Second subscription (should not start another update goroutine)
	err = rule.Handle(ctx, msg, connClient, connServer)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond) // Wait a bit to ensure the goroutine starts
}
