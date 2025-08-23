// File aggregate.go: definitions related to domain aggregates.
//
// This file contains the Aggregate interface used by domain objects to
// collect and expose domain events that should be published by the
// infrastructure layer.
package cqrs

// Aggregate is the interface that domain aggregates should implement.
// It exposes a stable identifier and a small event buffer used by
// dispatchers / event buses to publish changes.
type Aggregate interface {
	// GetAggregateID returns the aggregate's identifier.
	GetAggregateID() string
	// TrackEvent records a new domain event for later publication.
	TrackEvent(Event)
	// GetEvents returns the slice of recorded domain events.
	GetEvents() []Event
	// ClearEvents clears the recorded events after they have been published.
	ClearEvents()
}
