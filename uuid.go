// File uuid.go: small UUID helper.
//
// Convenience helper to produce UUID strings used as default aggregate IDs
// in examples and tests.
package cqrs

import (
	"github.com/google/uuid"
)

// NewUUIDString returns a new UUID string.
func NewUUIDString() string {
	newUUID := uuid.New()
	return newUUID.String()
}
