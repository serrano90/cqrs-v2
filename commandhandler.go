package cqrs

// CommandHandler is a interface that all command handler should implement
type CommandHandler interface {
	Handle(Command) (interface{}, error)
}

// CommandHandlerFunc
type CommandHandlerFunc func(Command) (interface{}, error)

//CommandHandlerMiddleware is a function that middleware
type CommandHandlerMiddleware func(CommandHandlerFunc) CommandHandlerFunc
