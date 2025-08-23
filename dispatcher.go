// File dispatcher.go: dispatcher contract for commands and queries.
//
// The Dispatcher routes commands and queries to their registered handlers
// and supports handler registration and middleware composition.
package cqrs

import "context"

// Dispatcher is an interface that all dispatcher implementations should
// satisfy. Implementations are responsible for routing commands and
// queries to their handlers and for applying handler middleware.
type Dispatcher interface {
	// Dispatch dispatches a command or query to the appropriate handler.
	Dispatch(context.Context, interface{}) (interface{}, error)
	// AddHandler registers a handler for the given command/query types.
	AddHandler(interface{}, ...interface{}) error
	// Use registers middleware that will be applied to command handlers.
	Use(...CommandHandlerMiddleware)
}
