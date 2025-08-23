# CQRS V2

This package is a simple implementation in Go of CQRS(Command Query Responsibility Segregation). We have interfaces and implements to use in any project with this solution.

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

## Install

Install cqrs package in your enviroment:
```bash
$ go get github.com/serrano90/cqrs-v2
```

## Examples

Below are two minimal examples that show how to register handlers and dispatch a command and a query using the in-memory dispatcher and the validation middleware included in this repo.

Command example (create a user):

```go
package main

import (
	"context"
	"fmt"

	"github.com/serrano90/cqrs-v2"
	"github.com/serrano90/cqrs-v2/memory"
	"github.com/serrano90/cqrs-v2/middleware"
)

// Command
type CreateUserCommand struct{ Name string }
func (c *CreateUserCommand) TypeOf() string { return "CreateUserCommand" }
func (c *CreateUserCommand) Validate() error {
	if c.Name == "" { return fmt.Errorf("name is required") }
	return nil
}

// Command Handler
type CreateUserHandler struct{}
func (h *CreateUserHandler) Handle(ctx context.Context, cmd cqrs.Command) (interface{}, error) {
	c := cmd.(*CreateUserCommand)
	// pretend we created a user and return an id
	return "user-id-123:" + c.Name, nil
}

func main() {
	d := memory.NewDispatcherInMemory()
	// register the handler for the command type (we pass a prototype value)
	d.AddHandler(&CreateUserHandler{}, &CreateUserCommand{})
	// enable validation middleware (calls Validate() when available)
	d.Use(middleware.NewValidationMiddleware())

	res, err := d.Dispatch(context.Background(), &CreateUserCommand{Name: "Alice"})
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

	"github.com/serrano90/cqrs-v2"
	"github.com/serrano90/cqrs-v2/memory"
)

// Query
type GetUserQuery struct{ ID string }
func (q *GetUserQuery) TypeOf() string { return "GetUserQuery" }

// Query Handler
type GetUserHandler struct{}
func (h *GetUserHandler) Handle(ctx context.Context, q cqrs.Query) (interface{}, error) {
	r := q.(*GetUserQuery)
	// return a simple projection
	return map[string]string{"id": r.ID, "name": "Alice"}, nil
}

func main() {
	d := memory.NewDispatcherInMemory()
	d.AddHandler(&GetUserHandler{}, &GetUserQuery{})

	res, err := d.Dispatch(context.Background(), &GetUserQuery{ID: "user-id-123"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Query result:", res)
}
```

## Material from learning

* [CQRS by Martin Fowler](http://martinfowler.com/bliki/CQRS.html)
* [Domain Event](https://www.martinfowler.com/eaaDev/DomainEvent.html)
* [Domain Driven Design Destilled](https://www.amazon.com/Domain-Driven-Design-Distilled-Vaughn-Vernon-dp-0134434420/dp/0134434420/ref=mt_other?_encoding=UTF8&me=&qid=)

## Based in

* https://github.com/jetbasrawi/go.cqrs