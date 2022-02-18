package cqrs

import "context"

// EventBus is a interface that all eventbus dispatcher should implement
type EventBus interface {
	// Publish publishes events to all registered event handlers
	Publish(context.Context, Event)
	// AddHandler registers and event handler for all events specified
	AddHandler(EventHandler, ...Event)
}
