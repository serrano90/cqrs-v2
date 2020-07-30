package cqrs

// Aggreagate is a interface that all aggregator should implement
type Aggregate interface {
	// GetAggreagteID returns a id of the resource
	GetAggreagteID() string
	// TrackEvent add a new event that should applied
	TrackEvent(Event)
	// GetEvents returns all change that should applied in others applications
	GetEvents() []Event
	// ClearEvents remove all events from the aggregator
	ClearEvents()
}
