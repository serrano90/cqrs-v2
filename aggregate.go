// File aggregate.go: definitions related to domain aggregates.
//
// This file contains the Aggregate interface used by domain objects to
// collect and expose domain events that should be published by the
// infrastructure layer. E is the event type recorded by the aggregate.
package cqrs

// Aggregate is the interface that domain aggregates should implement.
// It exposes a stable identifier and a small event buffer used by
// dispatchers / event buses to publish changes.
type Aggregate[E Event] interface {
	// GetAggregateID returns the aggregate's identifier.
	GetAggregateID() string
	// TrackEvent records a new domain event for later publication.
	TrackEvent(E)
	// GetEvents returns the slice of recorded domain events.
	GetEvents() []E
	// ClearEvents clears the recorded events after they have been published.
	ClearEvents()
}
