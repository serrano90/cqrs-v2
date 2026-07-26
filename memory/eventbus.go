// File memory/eventbus.go: in-memory publish/subscribe event bus.
//
// A simple in-memory event bus used for testing and local dispatch. It
// implements the cqrs.EventBus contract for the cqrs.Event interface:
// events are routed by their topic and delivered synchronously to the
// handler subscribed to that topic. Each topic holds at most one handler,
// and handlers receive the cqrs.Event interface and assert the concrete
// event type they expect.
package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/serrano90/cqrs-v2/v3"
)

// EventBusInMemory delivers published events to topic subscriptions.
type EventBusInMemory struct {
	// mu guards stopped and subscriptions.
	mu      sync.RWMutex
	stopped bool
	// subscriptions keyed by topic.
	subscriptions map[string]cqrs.EventHandler[cqrs.Event]
}

// NewEventBusInMemory creates an empty in-memory event bus that satisfies
// the cqrs.EventBus contract.
func NewEventBusInMemory() cqrs.EventBus[cqrs.Event] {
	return &EventBusInMemory{
		subscriptions: make(map[string]cqrs.EventHandler[cqrs.Event]),
	}
}

// Publish delivers the event synchronously to the handler subscribed to
// the event's topic. The handler is invoked after the lock is released so
// handlers can safely call back into the bus.
func (b *EventBusInMemory) Publish(ctx context.Context, e cqrs.Event) error {
	b.mu.RLock()
	if b.stopped {
		b.mu.RUnlock()
		return errors.New(cqrs.ErrMessageEventBusStopped)
	}
	h, ok := b.subscriptions[e.Topic()]
	b.mu.RUnlock()

	if ok {
		h.Handle(ctx, e)
	}
	return nil
}

// Subscribe registers the handler for the given topic. It returns an
// error when the topic already has a handler or the bus has been stopped.
func (b *EventBusInMemory) Subscribe(topic string, h cqrs.EventHandler[cqrs.Event]) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.stopped {
		return errors.New(cqrs.ErrMessageEventBusStopped)
	}
	if _, ok := b.subscriptions[topic]; ok {
		return errors.New(cqrs.ErrMessageSubscriptionDuplicated)
	}
	b.subscriptions[topic] = h
	return nil
}

// Unsubscribe removes the handler subscribed to the given topic. It
// returns an error when the topic has no subscription or the bus has been
// stopped.
func (b *EventBusInMemory) Unsubscribe(topic string, h cqrs.EventHandler[cqrs.Event]) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.stopped {
		return errors.New(cqrs.ErrMessageEventBusStopped)
	}
	if _, ok := b.subscriptions[topic]; !ok {
		return errors.New(cqrs.ErrMessageSubscriptionDoesNotExist)
	}
	delete(b.subscriptions, topic)
	return nil
}

// UnsubscribeAll removes every subscription on the given topic. It
// returns an error when the bus has been stopped.
func (b *EventBusInMemory) UnsubscribeAll(topic string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.stopped {
		return errors.New(cqrs.ErrMessageEventBusStopped)
	}
	delete(b.subscriptions, topic)
	return nil
}

// Stop stops the bus and releases every subscription. After Stop, all bus
// operations return an error, including stopping it again.
func (b *EventBusInMemory) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.stopped {
		return errors.New(cqrs.ErrMessageEventBusStopped)
	}
	b.stopped = true
	b.subscriptions = make(map[string]cqrs.EventHandler[cqrs.Event])
	return nil
}
