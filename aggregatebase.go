// File aggregatebase.go: a lightweight base implementation for aggregates.
//
// Provides a simple AggregateBase that generates or accepts an ID and
// keeps an in-memory slice of events recorded by the aggregate.

package cqrs

// NewAggregateBase creates a new aggregate base with a random ID.
func NewAggregateBase[E Event]() *AggregateBase[E] {
	return &AggregateBase[E]{
		id:     NewUUIDString(),
		events: []E{},
	}
}

// NewAggregateBaseByID creates a new aggregate base with a fixed ID.
func NewAggregateBaseByID[E Event](id string) *AggregateBase[E] {
	return &AggregateBase[E]{
		id:     id,
		events: []E{},
	}
}

// AggregateBase is a minimal aggregate implementation that stores an ID and
// a slice of tracked events of type E.
type AggregateBase[E Event] struct {
	id     string
	events []E
}

// GetAggregateID returns the id of the aggregate.
func (a *AggregateBase[E]) GetAggregateID() string {
	return a.id
}

// TrackEvent appends a new event to the aggregate's event buffer.
func (a *AggregateBase[E]) TrackEvent(e E) {
	a.events = append(a.events, e)
}

// GetEvents returns the currently tracked events.
func (a *AggregateBase[E]) GetEvents() []E {
	return a.events
}

// ClearEvents clears the recorded events after they have been published.
func (a *AggregateBase[E]) ClearEvents() {
	a.events = []E{}
}
