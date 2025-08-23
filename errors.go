// File errors.go: shared error message constants.
//
// Package-level error messages used by dispatchers and routers. They are
// plain strings to keep the package dependency-free; callers often wrap
// or convert them into typed errors as needed.
package cqrs

var (
	ErrMessageHandlerDoesNotExist = "The command or query handler does not exist"
	ErrMessageHandlerDuplicated   = "The command or query handler is duplicated for command or query"
)
