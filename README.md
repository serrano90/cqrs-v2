# CQRS

This package is a simple implementation in Go of CQRS (Command Query Responsibility Segregation). It provides generic (type-safe) interfaces and in-memory implementations to use in any project with this solution.

The API is built on Go generics: commands, queries, events, and their handlers are all parameterized by concrete types, so handlers receive and return typed values instead of `interface{}` and no type assertions are needed.

The generics-based API is a breaking change released as major version `v3` (module path `github.com/serrano90/cqrs-v2/v3`).

## Definitions

### CQRS

CQRS (Command Query Responsibility Segregation) is a pattern that separates operations that change state (commands) from operations that read state (queries).
The goal is to optimize, scale and reason about reads and writes independently. In this repository we provide small interfaces and in-memory implementations to help structure an application using CQRS.

#### Command

A Command represents an intent to change the system state. Commands are typically imperatively named (e.g. `CreateUser`, `UpdateOrderStatus`) and are handled by a command handler which executes business logic and may emit events.

Properties:
- Mutates state
- Handled by a single handler
- Can be validated before execution

#### Query

A Query represents a request for information and must not change system state. Queries are handled by query handlers which return data to the caller. Queries should be side-effect free and optimized for reading.

Properties:
- Read-only
- Can have multiple handlers or projections optimized for specific read models
- Should be cheap and fast to execute

#### Event

An Event represents a fact that already happened in the domain (e.g. `UserCreated`). Events are recorded by aggregates and published through an event bus to any number of registered event handlers (projections, notifications, side effects).

## Requirements

Go 1.26 or later is required.

## Install

Install the cqrs package in your environment:
```bash
$ go get github.com/serrano90/cqrs-v2/v3
```

## API overview

Go interfaces cannot declare generic methods, so the dispatcher exposes registration and dispatch as generic package-level functions. The in-memory event bus implements `cqrs.EventBus[cqrs.Event]`, so its operations are ordinary methods:

| Operation | Function |
|---|---|
| Register a command handler | `memory.AddCommandHandler[C, R](d, handler)` |
| Register a query handler | `memory.AddQueryHandler[Q, R](d, handler)` |
| Register command middleware | `memory.Use[C, R](d, middleware)` |
| Dispatch a command | `memory.DispatchCommand[C, R](d, ctx, cmd)` |
| Dispatch a query | `memory.DispatchQuery[Q, R](d, ctx, query)` |
| Create an event bus | `memory.NewEventBusInMemory()` |
| Publish an event | `bus.Publish(ctx, event)` |
| Subscribe to a topic | `bus.Subscribe(topic, handler)` |
| Unsubscribe a handler from a topic | `bus.Unsubscribe(topic, handler)` |
| Unsubscribe all handlers from a topic | `bus.UnsubscribeAll(topic)` |
| Stop the bus | `bus.Stop()` |

Notes:
- `C`/`Q`/`E` must be concrete types (usually pointers, e.g. `*CreateUserCommand`), not interfaces; the dispatcher routes messages by their concrete type.
- The type parameters used at registration must match the ones used at dispatch. Registering with `[*CreateUserCommand, string]` and dispatching with a different result type returns an error.
- Middleware registered with `memory.Use[C, R]` is applied only to commands dispatched with the same `[C, R]` pair. Middlewares run in registration order: the first registered is the outermost.
- The root package defines the `cqrs.EventBus[E]` contract (publish, subscribe, unsubscribe, stop). `memory.NewEventBusInMemory()` returns an implementation instantiated with the `cqrs.Event` interface, so one bus transports every event type and handlers assert the concrete event type they expect.
- The event bus follows a publish/subscribe model. `bus.Subscribe(topic, handler)` registers the handler for a topic — one handler per topic — and `bus.Publish` delivers the event synchronously to the handler subscribed to the event's `Topic()`. Subscriptions are released with `bus.Unsubscribe(topic, handler)`, `bus.UnsubscribeAll(topic)`, or `bus.Stop()`, which shuts down the whole bus.

## Examples

Below are minimal examples that show how to register handlers and dispatch a command and a query using the in-memory dispatcher and the validation middleware included in this repo.

Command example (create a user):

```go
package main

import (
	"context"
	"fmt"

	"github.com/serrano90/cqrs-v2/v3/memory"
	"github.com/serrano90/cqrs-v2/v3/middleware"
)

// Command
type CreateUserCommand struct{ Name string }
func (c *CreateUserCommand) TypeOf() string { return "CreateUserCommand" }
func (c *CreateUserCommand) Validate() error {
	if c.Name == "" { return fmt.Errorf("name is required") }
	return nil
}

// Command Handler: receives the concrete command type and returns a typed result.
type CreateUserHandler struct{}
func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) (string, error) {
	// pretend we created a user and return an id
	return "user-id-123:" + cmd.Name, nil
}

func main() {
	d := memory.NewDispatcherInMemory()
	// register the handler for the command type
	memory.AddCommandHandler[*CreateUserCommand, string](d, &CreateUserHandler{})
	// enable validation middleware (calls Validate() when available)
	memory.Use(d, middleware.NewValidationMiddleware[*CreateUserCommand, string]())

	// res is a string: no type assertion needed
	res, err := memory.DispatchCommand[*CreateUserCommand, string](d, context.Background(), &CreateUserCommand{Name: "Alice"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Create command result:", res)
}
```

Query example (read a user):

```go
package main

import (
	"context"
	"fmt"

	"github.com/serrano90/cqrs-v2/v3/memory"
)

// Query
type GetUserQuery struct{ ID string }
func (q *GetUserQuery) TypeOf() string { return "GetUserQuery" }

// Query Handler: receives the concrete query type and returns a typed projection.
type GetUserHandler struct{}
func (h *GetUserHandler) Handle(ctx context.Context, q *GetUserQuery) (map[string]string, error) {
	// return a simple projection
	return map[string]string{"id": q.ID, "name": "Alice"}, nil
}

func main() {
	d := memory.NewDispatcherInMemory()
	memory.AddQueryHandler[*GetUserQuery, map[string]string](d, &GetUserHandler{})

	res, err := memory.DispatchQuery[*GetUserQuery, map[string]string](d, context.Background(), &GetUserQuery{ID: "user-id-123"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Query result:", res)
}
```

Event example (publish a domain event):

```go
package main

import (
	"context"
	"fmt"

	"github.com/serrano90/cqrs-v2/v3"
	"github.com/serrano90/cqrs-v2/v3/memory"
)

// Event
type UserCreatedEvent struct{ ID string }
func (e *UserCreatedEvent) AggreagateID() string { return e.ID }
func (e *UserCreatedEvent) TypeOf() string { return "UserCreatedEvent" }
func (e *UserCreatedEvent) Topic() string { return "users" }
func (e *UserCreatedEvent) Message() *cqrs.EventMessage { return nil }

// Event Handler: receives the cqrs.Event interface and asserts the
// concrete event type it expects.
type SendWelcomeEmailHandler struct{}
func (h *SendWelcomeEmailHandler) Handle(ctx context.Context, e cqrs.Event) {
	if event, ok := e.(*UserCreatedEvent); ok {
		fmt.Println("sending welcome email to user", event.ID)
	}
}

func main() {
	bus := memory.NewEventBusInMemory()
	defer bus.Stop()

	// subscribe the handler to the topic
	if err := bus.Subscribe("users", &SendWelcomeEmailHandler{}); err != nil {
		panic(err)
	}

	// deliver the event synchronously to the handler subscribed to "users"
	if err := bus.Publish(context.Background(), &UserCreatedEvent{ID: "user-id-123"}); err != nil {
		panic(err)
	}
}
```

Aggregates: embed `cqrs.AggregateBase[E]` to record events of type `E` and publish them after handling a command:

```go
type User struct {
	*cqrs.AggregateBase[*UserCreatedEvent]
	Name string
}

func NewUser(name string) *User {
	u := &User{
		AggregateBase: cqrs.NewAggregateBase[*UserCreatedEvent](),
		Name:          name,
	}
	u.TrackEvent(&UserCreatedEvent{ID: u.GetAggregateID()})
	return u
}

// after handling the command:
// for _, e := range u.GetEvents() { bus.Publish(ctx, e) }
// u.ClearEvents()
```

## Material from learning

* [CQRS by Martin Fowler](http://martinfowler.com/bliki/CQRS.html)
* [Domain Event](https://www.martinfowler.com/eaaDev/DomainEvent.html)
* [Domain Driven Design Destilled](https://www.amazon.com/Domain-Driven-Design-Distilled-Vaughn-Vernon-dp-0134434420/dp/0134434420/ref=mt_other?_encoding=UTF8&me=&qid=)

## Based in

* https://github.com/jetbasrawi/go.cqrs
