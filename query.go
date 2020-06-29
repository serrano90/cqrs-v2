package cqrs

// Query is a interface that all query should implement
type Query interface {
	TypeOf() string
}
