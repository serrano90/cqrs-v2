// File eventbus.go: event publishing interface.
//
// The EventBus is responsible for publishing events to registered
// EventHandlers. Implementations may be in-memory, broker-backed, or
// remote and should be safe for concurrent use when required.
package cqrs

import "context"

// EventBus is the interface that event bus implementations should satisfy.
type EventBus interface {
	// Publish publishes an event to all registered handlers.
	Publish(context.Context, Event)
	// AddHandler registers an EventHandler to receive the specified events.
	AddHandler(EventHandler, ...Event)
}
