// File errors.go: shared error message constants.
//
// Package-level error messages used by dispatchers and routers. They are
// plain strings to keep the package dependency-free; callers often wrap
// or convert them into typed errors as needed.

package cqrs

var (
	// ErrMessageHandlerDoesNotExist reports that no handler is registered
	// for the dispatched command or query.
	ErrMessageHandlerDoesNotExist = "The command or query handler does not exist"
	// ErrMessageHandlerDuplicated reports that a handler is already
	// registered for the command or query.
	ErrMessageHandlerDuplicated = "The command or query handler is duplicated for command or query"
	// ErrMessageEventBusStopped reports that the event bus has been stopped.
	ErrMessageEventBusStopped = "The event bus is stopped"
	// ErrMessageSubscriptionDuplicated reports that the topic already has
	// a subscribed handler.
	ErrMessageSubscriptionDuplicated = "The subscription is duplicated for the topic"
	// ErrMessageSubscriptionDoesNotExist reports that the topic has no
	// subscribed handler.
	ErrMessageSubscriptionDoesNotExist = "The subscription does not exist for the topic"
)
