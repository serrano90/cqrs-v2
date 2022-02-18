package cqrs

import "context"

// CommandHandler is a interface that all command handler should implement
type CommandHandler interface {
	Handle(context.Context, Command) (interface{}, error)
}

// CommandHandlerFunc
type CommandHandlerFunc func(context.Context, Command) (interface{}, error)

//CommandHandlerMiddleware is a function that middleware
type CommandHandlerMiddleware func(CommandHandlerFunc) CommandHandlerFunc
