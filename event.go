// File event.go: domain event contract and broker message payload.
//
// Events describe facts that already happened in the domain. They expose
// the aggregate they belong to, a routing type, a topic, and an optional
// broker message payload.

package cqrs

// Event is the interface that all event messages should implement.
type Event interface {
	AggreagateID() string
	TypeOf() string
	Topic() string
	Message() *EventMessage
}

// EventMessage represents a message passed through a message broker.
type EventMessage struct {
	Body   []byte
	Header map[string]string
}
