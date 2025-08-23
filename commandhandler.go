// File commandhandler.go: command handler types and middleware.
//
// Defines the CommandHandler abstraction and the middleware function type
// used to wrap handlers with additional behavior (validation, logging, etc.).
package cqrs

import "context"

// CommandHandler is an interface implemented by concrete command handlers.
type CommandHandler interface {
	Handle(context.Context, Command) (interface{}, error)
}

// CommandHandlerFunc adapts a function to the CommandHandlerFunc shape.
type CommandHandlerFunc func(context.Context, Command) (interface{}, error)

// CommandHandlerMiddleware represents middleware that wraps a
// CommandHandlerFunc and returns a new CommandHandlerFunc.
type CommandHandlerMiddleware func(CommandHandlerFunc) CommandHandlerFunc
