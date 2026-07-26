// File commandhandler.go: generic command handler types and middleware.
//
// Defines the CommandHandler abstraction and the middleware function type
// used to wrap handlers with additional behavior (validation, logging, etc.).
// The type parameters provide compile-time safety: C is the concrete command
// type and R is the result type produced by the handler.
package cqrs

import "context"

// CommandHandler is implemented by concrete command handlers.
// C is the command type the handler accepts and R is the result it returns.
type CommandHandler[C Command, R any] interface {
	Handle(context.Context, C) (R, error)
}

// CommandHandlerFunc adapts a function to the CommandHandler shape.
type CommandHandlerFunc[C Command, R any] func(context.Context, C) (R, error)

// CommandHandlerMiddleware represents middleware that wraps a
// CommandHandlerFunc and returns a new CommandHandlerFunc.
type CommandHandlerMiddleware[C Command, R any] func(CommandHandlerFunc[C, R]) CommandHandlerFunc[C, R]
