package cqrs

// EventBus is a interface that all eventbus dispatcher should implement
type EventBus interface {
	// Publish is a methos where send the event to others services
	Publish(Event)
	// AddHandler add the event handler for each event
	AddHandler(EventHandler, ...Event)
}
