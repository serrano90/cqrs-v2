package memory_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/serrano90/cqrs-v2"
	"github.com/serrano90/cqrs-v2/memory"
)

func TestNewInstanceOfEventBus(t *testing.T) {
	if d := memory.NewEventBusInMemory(); d == nil {
		t.Fail()
	}
}

func TestEventBusAddHandler(t *testing.T) {
	tests := map[string]struct {
		events       []cqrs.Event
		eventHandler cqrs.EventHandler
	}{
		"when does not exit any event or event handler regitered": {
			events: []cqrs.Event{
				NewTestEvent(),
			},
			eventHandler: NewTestEventHandler(),
		},
		"when the event exist and need add a new handler": {
			events: []cqrs.Event{
				NewTestEvent(),
			},
			eventHandler: NewTestEventHandler2(),
		},
	}

	bus := memory.NewEventBusInMemory()
	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		bus.AddHandler(test.eventHandler, test.events...)
	}
}

func TestEventBusPublish(t *testing.T) {
	tests := map[string]struct {
		event        cqrs.Event
		eventHandler []cqrs.EventHandler
		eventPublish cqrs.Event
	}{
		"success": {
			event: NewTestEvent(),
			eventHandler: []cqrs.EventHandler{
				NewTestEventHandler(),
				NewTestEventHandler2(),
			},
			eventPublish: NewTestEvent(),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		bus := memory.NewEventBusInMemory()
		for _, h := range test.eventHandler {
			bus.AddHandler(h, test.event)
		}

		bus.Publish(context.Background(), test.eventPublish)
	}
}

type TestEvent struct{}

func NewTestEvent() cqrs.Event {
	return &TestEvent{}
}

func (e *TestEvent) TypeOf() string {
	return reflect.TypeOf(e).Elem().Name()
}

func (e *TestEvent) AggreagateID() string {
	return ""
}

func (e *TestEvent) Topic() string {
	return ""
}

func (e *TestEvent) Message() *cqrs.EventMessage {
	return nil
}

type TestEventHandler struct{}

func NewTestEventHandler() cqrs.EventHandler {
	return &TestEventHandler{}
}

func (handler *TestEventHandler) Handle(ctx context.Context, e cqrs.Event) {}

type TestEventHandler2 struct{}

func NewTestEventHandler2() cqrs.EventHandler {
	return &TestEventHandler2{}
}

func (handler *TestEventHandler2) Handle(ctx context.Context, e cqrs.Event) {}
