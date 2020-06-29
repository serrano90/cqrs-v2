package cqrs

// Command is a interface that all command should implement
type Command interface {
	// TypeOf return a type of command
	TypeOf() string
}

// CommandValidate is a interface that command required a validate should implement
type CommandValidate interface {
	// Validate review the command payload.
	// Return error if this is not valid.
	Validate() error
}
