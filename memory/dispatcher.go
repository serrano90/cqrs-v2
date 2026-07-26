// File memory/dispatcher.go: in-memory generic dispatcher implementation.
//
// A simple dispatcher that maps concrete command/query types to typed
// handlers and applies typed middleware to command handlers. Registration
// and dispatch are exposed as generic package-level functions because Go
// interfaces cannot declare generic methods. Intended for tests and local
// use.
package memory

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/serrano90/cqrs-v2/v3"
)

// DispatcherInMemory routes commands and queries to their typed handlers.
type DispatcherInMemory struct {
	// commandHandlers keyed by concrete command type.
	commandHandlers map[reflect.Type]any
	// queryHandlers keyed by concrete query type.
	queryHandlers map[reflect.Type]any
	// middlewares registered for command handlers; each entry holds a
	// cqrs.CommandHandlerMiddleware[C, R] and is applied when its type
	// parameters match the dispatched command.
	middlewares []any
}

// NewDispatcherInMemory creates an empty in-memory dispatcher.
func NewDispatcherInMemory() *DispatcherInMemory {
	return &DispatcherInMemory{
		commandHandlers: make(map[reflect.Type]any),
		queryHandlers:   make(map[reflect.Type]any),
		middlewares:     make([]any, 0),
	}
}

// AddCommandHandler registers a typed command handler for commands of type C.
// It returns an error when a handler is already registered for C.
func AddCommandHandler[C cqrs.Command, R any](d *DispatcherInMemory, h cqrs.CommandHandler[C, R]) error {
	t, err := typeOf[C]("command")
	if err != nil {
		return err
	}
	if _, ok := d.commandHandlers[t]; ok {
		return fmt.Errorf("%s %s", cqrs.ErrMessageHandlerDuplicated, typeName(t))
	}
	d.commandHandlers[t] = h
	return nil
}

// AddQueryHandler registers a typed query handler for queries of type Q.
// It returns an error when a handler is already registered for Q.
func AddQueryHandler[Q cqrs.Query, R any](d *DispatcherInMemory, h cqrs.QueryHandler[Q, R]) error {
	t, err := typeOf[Q]("query")
	if err != nil {
		return err
	}
	if _, ok := d.queryHandlers[t]; ok {
		return fmt.Errorf("%s %s", cqrs.ErrMessageHandlerDuplicated, typeName(t))
	}
	d.queryHandlers[t] = h
	return nil
}

// Use registers a typed command middleware. The middleware is applied only
// to commands dispatched with matching C and R type parameters.
func Use[C cqrs.Command, R any](d *DispatcherInMemory, mw cqrs.CommandHandlerMiddleware[C, R]) {
	d.middlewares = append(d.middlewares, any(mw))
}

// DispatchCommand dispatches a command to its registered handler, wrapping
// it with every compatible middleware. Middlewares run in registration
// order: the first registered is the outermost.
func DispatchCommand[C cqrs.Command, R any](d *DispatcherInMemory, ctx context.Context, c C) (R, error) {
	var zero R
	t := reflect.TypeOf(c)
	hAny, ok := d.commandHandlers[t]
	if !ok {
		return zero, errors.New(cqrs.ErrMessageHandlerDoesNotExist)
	}
	h, ok := hAny.(cqrs.CommandHandler[C, R])
	if !ok {
		return zero, fmt.Errorf("registered handler has incompatible type for %s", typeName(t))
	}

	fn := cqrs.CommandHandlerFunc[C, R](h.Handle)
	for i := len(d.middlewares) - 1; i >= 0; i-- {
		if mw, ok := d.middlewares[i].(cqrs.CommandHandlerMiddleware[C, R]); ok {
			fn = mw(fn)
		}
	}
	return fn(ctx, c)
}

// DispatchQuery dispatches a query to its registered handler.
func DispatchQuery[Q cqrs.Query, R any](d *DispatcherInMemory, ctx context.Context, q Q) (R, error) {
	var zero R
	t := reflect.TypeOf(q)
	hAny, ok := d.queryHandlers[t]
	if !ok {
		return zero, errors.New(cqrs.ErrMessageHandlerDoesNotExist)
	}
	h, ok := hAny.(cqrs.QueryHandler[Q, R])
	if !ok {
		return zero, fmt.Errorf("registered handler has incompatible type for %s", typeName(t))
	}
	return h.Handle(ctx, q)
}

// typeOf resolves the reflect.Type of the type parameter T. It fails when T
// is an interface type, because the concrete type is needed to route
// messages at dispatch time.
func typeOf[T any](kind string) (reflect.Type, error) {
	var zero T
	t := reflect.TypeOf(zero)
	if t == nil {
		return nil, fmt.Errorf("cannot register handler: %s type must be concrete, not an interface", kind)
	}
	return t, nil
}

// typeName returns a readable name for t, unwrapping pointers so that
// *TestCommand is reported as TestCommand.
func typeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
