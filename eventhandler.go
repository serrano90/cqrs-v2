// File eventhandler.go: event handler contract.
//
// Event handlers implement application-specific logic that runs when an
// event is published (e.g., projection updates, notifications, side effects).
package cqrs

import "context"

// EventHandler handles domain events delivered by an EventBus.
// Implementations should ensure their Handle method is safe and resilient
// to failures; the bus implementation controls retry or error handling.
type EventHandler interface {
	// Handle is invoked for each matching event.
	Handle(context.Context, Event)
}
