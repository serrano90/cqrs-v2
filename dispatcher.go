package cqrs

import "context"

// Dispatcher is an interface that all dispatcher schema should implement.
//
// The dispatcher is the form do you distribute the command or query to the determined handler
type Dispatcher interface {
	// Dispatch dispatches a command or query to the specified controller to be served
	Dispatch(context.Context, interface{}) (interface{}, error)
	// AddHandler registers and command handler for a specified command type
	AddHandler(interface{}, ...interface{}) error
	// Use registers the handler middlewares where will pass the command before being dispatched to command handler
	Use(...CommandHandlerMiddleware)
}
