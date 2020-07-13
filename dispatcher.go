package cqrs

// Dispatcher is a interface that all dispatcher schema should implement
//
// With dispatcher is the form do you destributed the commnad or query to the
// determinated command handler or query handler.
type Dispatcher interface {
	// Manage a command or query for send to determinated handler
	Dispatch(interface{}) (interface{}, error)
	// Add a new to handler for a especified type
	AddHandler(interface{}, ...interface{}) error
	// Use is a function where add middleware
	Use(...CommandHandlerMiddleware)
}
