// File command.go: domain command types and validation interface.
//
// Defines the Command contract and an optional validation interface that
// commands can implement to validate their payload prior to handling.

package cqrs

// Command represents a domain command.
// Implementations must return a stable, machine-readable type string
// via TypeOf(), which is used to route and handle the command.
type Command interface {
	// TypeOf returns the command type.
	TypeOf() string
}

// CommandValidate is implemented by commands that perform payload validation.
// Implementations should return nil when the command is valid, or a non-nil
// error describing the validation failure.
type CommandValidate interface {
	// Validate reviews the command payload and returns an error on failure.
	Validate() error
}
