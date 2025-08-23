// File aggregatebase.go: a lightweight base implementation for aggregates.
//
// Provides a simple AggregateBase that generates or accepts an ID and
// keeps an in-memory slice of events recorded by the aggregate.
package cqrs

// Create a new instance of base aggregate
func NewAggregateBase() Aggregate {
	return &AggregateBase{
		id:     NewUUIDString(),
		events: []Event{},
	}
}

// Create a new instance of base aggregate by id
func NewAggregateBaseById(id string) Aggregate {
	return &AggregateBase{
		id:     id,
		events: []Event{},
	}
}

// AggregateBase represents a minimal aggregate implementation used in
// examples and tests. It stores an ID and a slice of tracked events.
type AggregateBase struct {
	id     string
	events []Event
}

// GetAggreagteID returns the id of aggregate
func (a *AggregateBase) GetAggreagteID() string {
	return a.id
}

// TrackEvent appends a new event to the aggregate's event buffer.
func (a *AggregateBase) TrackEvent(e Event) {
	a.events = append(a.events, e)
}

// GetEvents returns a copy of the currently tracked events.
func (a *AggregateBase) GetEvents() []Event {
	return a.events
}

func (a *AggregateBase) ClearEvents() {
	a.events = []Event{}
}
