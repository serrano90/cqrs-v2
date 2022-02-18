package cqrs

import "context"

// EventHandler is a interface that all event handler should implement
type EventHandler interface {
	// Handler is the method where orquestated the logic for send event
	Handle(context.Context, Event)
}
