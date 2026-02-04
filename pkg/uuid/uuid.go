package uuidpkg

import (
	"github.com/google/uuid"
)

// New returns a new UUID string.
func New() string {
	return uuid.NewString()
}
