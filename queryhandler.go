// File queryhandler.go: generic query handler contract.
//
// QueryHandler is implemented by components that serve queries and return
// results. Q is the concrete query type and R is the result type, so
// callers get a typed response instead of interface{}.
package cqrs

import "context"

// QueryHandler handles read-only queries of type Q and returns a typed
// response R or an error.
type QueryHandler[Q Query, R any] interface {
	Handle(context.Context, Q) (R, error)
}
