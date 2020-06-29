package cqrs

import (
	"github.com/google/uuid"
)

// NewUUIDString return a new uuid
func NewUUIDString() string {
	newUUID := uuid.New()
	return newUUID.String()
}
