package cqrs

// Event is a interface that all event message should implement
type Event interface {
	AggreagateID() string
	TypeOf() string
	Topic() string
	Message() *EventMessage
}

// Message represent a message passed through a message broker
type EventMessage struct {
	Body   []byte
	Header map[string]string
}
