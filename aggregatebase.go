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

// AggregateBase reprecent a new instance of base aggregate
type AggregateBase struct {
	id     string
	events []Event
}

// GetAggreagteID returns the is of aggregate
func (a *AggregateBase) GetAggreagteID() string {
	return a.id
}

// TrackEvent add a new event track for aggregate
func (a *AggregateBase) TrackEvent(e Event) {
	a.events = append(a.events, e)
}

// GetEvents return a array of event tracking
func (a *AggregateBase) GetEvents() []Event {
	return a.events
}

func (a *AggregateBase) ClearEvents() {
	a.events = []Event{}
}
