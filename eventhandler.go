// File eventhandler.go: generic event handler contract.
//
// Event handlers implement application-specific logic that runs when an
// event is published (e.g., projection updates, notifications, side effects).
// E is the concrete event type the handler subscribes to.
package cqrs

import "context"

// EventHandler handles domain events of type E delivered by an event bus.
// Implementations should ensure their Handle method is safe and resilient
// to failures; the bus implementation controls retry or error handling.
type EventHandler[E Event] interface {
	// Handle is invoked for each matching event.
	Handle(context.Context, E)
}
