// File queryhandler.go: query handler contract.
//
// QueryHandler is implemented by components that serve queries and return
// results. They may be used alongside Dispatchers that route queries.
package cqrs

import "context"

// QueryHandler handles read-only queries and returns a response or error.
type QueryHandler interface {
	Handle(context.Context, Query) (interface{}, error)
}
