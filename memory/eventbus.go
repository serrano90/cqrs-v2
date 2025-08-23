// File memory/eventbus.go: in-memory EventBus implementation.
//
// A simple in-memory EventBus used for testing and local dispatch. It
// maintains a map from event type to handlers and delivers events
// synchronously to each registered handler.
package memory

import (
	"context"
	"reflect"

	"github.com/serrano90/cqrs-v2"
)

type EventBusInMemory struct {
	handlers map[string]map[string]cqrs.EventHandler
}

func NewEventBusInMemory() cqrs.EventBus {
	return &EventBusInMemory{
		handlers: make(map[string]map[string]cqrs.EventHandler, 0),
	}
}

// Publish synchronously delivers the event to all registered handlers.
func (bus *EventBusInMemory) Publish(ctx context.Context, event cqrs.Event) {
	typeName := event.TypeOf()
	if handlers, ok := bus.handlers[typeName]; ok {
		for _, handle := range handlers {
			handle.Handle(ctx, event)
		}
	}
}

// AddHandler registers the provided handler for each given event type.
func (bus *EventBusInMemory) AddHandler(eventhandler cqrs.EventHandler, events ...cqrs.Event) {
	for _, event := range events {
		typeNameEvent := event.TypeOf()
		if _, ok := bus.handlers[typeNameEvent]; !ok {
			bus.handlers[typeNameEvent] = make(map[string]cqrs.EventHandler, 0)
		}

		typeNameEventHandler := bus.getTypeOf(eventhandler)
		if _, ok := bus.handlers[typeNameEvent][typeNameEventHandler]; !ok {
			bus.handlers[typeNameEvent][typeNameEventHandler] = eventhandler
		}
	}
}

func (bus *EventBusInMemory) getTypeOf(h cqrs.EventHandler) string {
	return reflect.TypeOf(h).Elem().Name()
}
