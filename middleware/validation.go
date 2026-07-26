// File validation.go: generic command validation middleware.
//
// Provides a middleware that inspects a command for the CommandValidate
// interface and runs its Validate method before passing the command on.
package middleware

import (
	"context"

	"github.com/serrano90/cqrs-v2/v3"
)

// NewValidationMiddleware returns a CommandHandlerMiddleware for commands
// of type C that runs Validate when C implements cqrs.CommandValidate and
// rejects the command when validation fails.
func NewValidationMiddleware[C cqrs.Command, R any]() cqrs.CommandHandlerMiddleware[C, R] {
	return func(next cqrs.CommandHandlerFunc[C, R]) cqrs.CommandHandlerFunc[C, R] {
		return func(ctx context.Context, cmd C) (R, error) {
			if v, ok := any(cmd).(cqrs.CommandValidate); ok {
				if err := v.Validate(); err != nil {
					var zero R
					return zero, err
				}
			}
			return next(ctx, cmd)
		}
	}
}
