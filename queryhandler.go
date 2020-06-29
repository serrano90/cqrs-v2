package cqrs

// QueryHandler is a interface that all query handler should implement
type QueryHandler interface {
	Handle(Query) (interface{}, error)
}
