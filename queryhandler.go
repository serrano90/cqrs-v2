package cqrs

import "context"

// QueryHandler is a interface that all query handler should implement
type QueryHandler interface {
	Handle(context.Context, Query) (interface{}, error)
}
