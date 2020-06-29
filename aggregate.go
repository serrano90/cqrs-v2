package cqrs

import "time"

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
	// GetCreatedAt get a time where the aggreagate was created
	GetCreatedAt() time.Time
	// GetUpdateAt get a time where aggreagate was created
	GetUpdateAt() time.Time
}
