package memory_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/serrano90/cqrs-v2/v3"
	"github.com/serrano90/cqrs-v2/v3/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewInstanceOfEventBus(t *testing.T) {
	// The constructor returns the cqrs.EventBus contract
	var b cqrs.EventBus[cqrs.Event]
	if b = memory.NewEventBusInMemory(); b == nil {
		t.Fail()
	}
}

func TestEventBusSubscribeAndPublish(t *testing.T) {
	bus := memory.NewEventBusInMemory()

	h := &TestEventHandler{}
	err := bus.Subscribe("users", h)
	assert.Equal(t, nil, err, "First subscription should succeed")

	// The topic already has a handler
	err = bus.Subscribe("users", &TestEventHandler{})
	assert.Equal(t, errors.New(cqrs.ErrMessageSubscriptionDuplicated), err, "Should return already subscribed error")

	// An event published to the subscribed topic is delivered synchronously
	err = bus.Publish(context.Background(), NewTestEventWithTopic("users"))
	assert.Equal(t, nil, err, "Publish should succeed")
	assert.Equal(t, 1, h.calls, "Handler should receive the published event")

	// An event published to another topic is not delivered
	err = bus.Publish(context.Background(), NewTestEventWithTopic("orders"))
	assert.Equal(t, nil, err, "Publish should succeed")
	assert.Equal(t, 1, h.calls, "Handler should not receive events from other topics")
}

func TestEventBusUnsubscribe(t *testing.T) {
	bus := memory.NewEventBusInMemory()

	h := &TestEventHandler{}
	err := bus.Unsubscribe("users", h)
	assert.Equal(t, errors.New(cqrs.ErrMessageSubscriptionDoesNotExist), err, "Unsubscribing without a subscription should fail")

	err = bus.Subscribe("users", h)
	assert.Equal(t, nil, err, "Subscription should succeed")

	err = bus.Unsubscribe("users", h)
	assert.Equal(t, nil, err, "Unsubscribe should succeed")

	// Events published after unsubscribing are not delivered
	err = bus.Publish(context.Background(), NewTestEventWithTopic("users"))
	assert.Equal(t, nil, err, "Publish should succeed")
	assert.Equal(t, 0, h.calls, "Handler should not receive events after unsubscribing")
}

func TestEventBusUnsubscribeAll(t *testing.T) {
	bus := memory.NewEventBusInMemory()

	h := &TestEventHandler{}
	assert.Equal(t, nil, bus.Subscribe("users", h), "Subscription should succeed")

	err := bus.UnsubscribeAll("users")
	assert.Equal(t, nil, err, "UnsubscribeAll should succeed")

	// Unsubscribing all on a topic without subscriptions is a no-op
	err = bus.UnsubscribeAll("users")
	assert.Equal(t, nil, err, "UnsubscribeAll without subscriptions should succeed")

	// Events published after unsubscribing all are not delivered
	assert.Equal(t, nil, bus.Publish(context.Background(), NewTestEventWithTopic("users")), "Publish should succeed")
	assert.Equal(t, 0, h.calls, "Handler should not receive events after unsubscribing all")
}

func TestEventBusStop(t *testing.T) {
	bus := memory.NewEventBusInMemory()

	h := &TestEventHandler{}
	assert.Equal(t, nil, bus.Subscribe("users", h), "Subscription should succeed")

	err := bus.Stop()
	assert.Equal(t, nil, err, "Stop should succeed")

	// After Stop, every operation returns an error
	err = bus.Stop()
	assert.Equal(t, errors.New(cqrs.ErrMessageEventBusStopped), err, "Stopping again should fail")

	err = bus.Publish(context.Background(), NewTestEventWithTopic("users"))
	assert.Equal(t, errors.New(cqrs.ErrMessageEventBusStopped), err, "Publishing on a stopped bus should fail")

	err = bus.Subscribe("users", &TestEventHandler{})
	assert.Equal(t, errors.New(cqrs.ErrMessageEventBusStopped), err, "Subscribing on a stopped bus should fail")

	err = bus.Unsubscribe("users", h)
	assert.Equal(t, errors.New(cqrs.ErrMessageEventBusStopped), err, "Unsubscribing on a stopped bus should fail")

	err = bus.UnsubscribeAll("users")
	assert.Equal(t, errors.New(cqrs.ErrMessageEventBusStopped), err, "Unsubscribing all on a stopped bus should fail")

	assert.Equal(t, 0, h.calls, "Handler should not receive events after stopping")
}

func TestEventBusPublishWithoutSubscriptions(t *testing.T) {
	bus := memory.NewEventBusInMemory()

	// Publishing an event with no subscriptions should not fail
	err := bus.Publish(context.Background(), NewTestEvent())
	assert.Equal(t, nil, err, "Publish should succeed")
}

type TestEvent struct {
	topic string
}

func NewTestEvent() *TestEvent {
	return &TestEvent{}
}

func NewTestEventWithTopic(topic string) *TestEvent {
	return &TestEvent{topic: topic}
}

func (e *TestEvent) TypeOf() string {
	return reflect.TypeOf(e).Elem().Name()
}

func (e *TestEvent) AggreagateID() string {
	return ""
}

func (e *TestEvent) Topic() string {
	return e.topic
}

func (e *TestEvent) Message() *cqrs.EventMessage {
	return nil
}

// TestEventHandler receives the cqrs.Event interface and counts the
// invocations; delivery is synchronous so no locking is needed.
type TestEventHandler struct {
	calls int
}

func (handler *TestEventHandler) Handle(_ context.Context, _ cqrs.Event) {
	handler.calls++
}
