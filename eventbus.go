// File eventbus.go: publish/subscribe event bus contract.
//
// The EventBus transports domain events between the code that produces
// them and the EventHandlers that consume them. It follows a
// publish/subscribe model: events are published to the topic they carry,
// and every handler subscribed to that topic consumes them. Subscriptions
// are managed on the bus itself: Subscribe binds a handler to a topic,
// Unsubscribe and UnsubscribeAll release topic subscriptions, and Stop
// shuts down the whole bus. Implementations may be in-memory,
// broker-backed, or remote and should be safe for concurrent use.
//
// Go interfaces cannot declare generic methods, so the contract is
// parameterized by the event type E. Instantiate it with a concrete event
// type for a single-event bus, or with the Event interface itself to
// transport every event type on one bus.

package cqrs

import "context"

// EventBus is the interface that event bus implementations should satisfy
// for events of type E. Publish publishes an event and Subscribe listens
// on a topic and dispatches the published events to an EventHandler.
type EventBus[E Event] interface {
	// Publish publishes an event to the subscriptions on its topic. It
	// returns an error when the bus cannot accept the event, for example
	// after it has been stopped.
	Publish(context.Context, E) error
	// Subscribe registers a consumer on the given topic that dispatches
	// published events of type E to the handler. It returns an error when
	// the handler is already subscribed to the topic or the bus has been
	// stopped.
	Subscribe(string, EventHandler[E]) error
	// Unsubscribe stops consuming events on the given topic for the
	// given handler. It returns an error when no such subscription exists.
	// Events already published but not yet handled may be discarded.
	Unsubscribe(string, EventHandler[E]) error
	// UnsubscribeAll stops consuming events on the given topic for all
	// handlers. Events already published but not yet handled may be
	// discarded.
	UnsubscribeAll(string) error
	// Stop stops the bus and all its subscriptions. After Stop, publishing
	// and subscribing return an error. Events already published but not
	// yet handled may be discarded.
	Stop() error
}
