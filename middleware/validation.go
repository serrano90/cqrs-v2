package middleware

import (
	"context"

	"github.com/serrano90/cqrs-v2"
)

// NewValidationMiddleware
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
