// File validation.go: command validation middleware.
//
// Provides a middleware that inspects a command for the CommandValidate
// interface and runs its Validate method before passing the command on.
package middleware

import (
	"context"

	"github.com/serrano90/cqrs-v2"
)

// NewValidationMiddleware returns a CommandHandlerMiddleware that runs
// Validate on commands implementing cqrs.CommandValidate and rejects the
// command when validation fails.
func NewValidationMiddleware() cqrs.CommandHandlerMiddleware {
	return func(next cqrs.CommandHandlerFunc) cqrs.CommandHandlerFunc {
		return func(ctx context.Context, cmd cqrs.Command) (interface{}, error) {
			if c, ok := cmd.(cqrs.CommandValidate); ok {
				err := c.Validate()
				if err != nil {
					return nil, err
				}
			}
			return next(ctx, cmd)
		}
	}
}
