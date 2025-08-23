// File query.go: query type contract.
//
// Queries represent read-only requests for information. Handlers for
// queries typically return a result or an error.
package cqrs

// Query is the interface that all queries should implement.
// The TypeOf string is used to route the query to its handler.
type Query interface {
	TypeOf() string
}
