package ws

import (
	"context"
	"sync"
)

// Ensure SubscriptionRule implements Rule
var _ Rule = (*SubscriptionRule)(nil)

type SubscriptionRule struct {
	subscriptionMessageMatcher   MessageMatcher
	subscriptionResponse         MessageHandler
	unsubscriptionMessageMatcher MessageMatcher
	unsubscriptionResponse       MessageHandler
	updateResponse               MessageHandler
	updateLock                   sync.Mutex
	updateEnabled                bool
	updateCancelFunc             func()
}

func NewSubscriptionRule(
	subscriptionMessageMatcher MessageMatcher,
	subscriptionResponse MessageHandler,
	unsubscriptionMessageMatcher MessageMatcher,
	unsubscriptionResponse MessageHandler,
	updateResponse MessageHandler,
) *SubscriptionRule {
	return &SubscriptionRule{
		subscriptionMessageMatcher:   subscriptionMessageMatcher,
		subscriptionResponse:         subscriptionResponse,
		unsubscriptionMessageMatcher: unsubscriptionMessageMatcher,
		unsubscriptionResponse:       unsubscriptionResponse,
		updateResponse:               updateResponse,
		updateLock:                   sync.Mutex{},
		updateEnabled:                false,
		updateCancelFunc:             nil,
	}
}

func (r *SubscriptionRule) MatchMessage(message Message) bool {
	return r.subscriptionMessageMatcher.MatchMessage(message) ||
		r.unsubscriptionMessageMatcher.MatchMessage(message)
}

func (r *SubscriptionRule) Handle(ctx context.Context, message Message, connClient Connection, connServer Connection) error {
	if r.subscriptionMessageMatcher.MatchMessage(message) {
		return r.handleSubscription(ctx, message, connClient, connServer)
	}
	return r.handleUnsubscription(ctx, message, connClient, connServer)
}

func (r *SubscriptionRule) handleSubscription(ctx context.Context, message Message, connClient Connection, connServer Connection) error {
	r.updateLock.Lock()
	if !r.updateEnabled {
		r.updateEnabled = true
		ctx, cancel := context.WithCancel(ctx)
		r.updateCancelFunc = cancel

		go r.updateResponse.Handle(ctx, message, connClient, connServer)
	}
	r.updateLock.Unlock()

	return r.subscriptionResponse.Handle(ctx, message, connClient, connServer)
}

func (r *SubscriptionRule) handleUnsubscription(ctx context.Context, message Message, connClient Connection, connServer Connection) error {
	r.updateLock.Lock()
	if r.updateEnabled {
		r.updateEnabled = false
		r.updateCancelFunc()
	}
	r.updateLock.Unlock()

	return r.unsubscriptionResponse.Handle(ctx, message, connClient, connServer)
}
