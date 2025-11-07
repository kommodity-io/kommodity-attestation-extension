// Package uuid provides error definitions for UUID operations.
package uuid

import "errors"

var (
	// ErrUUIDNotFound is returned when the machine UUID cannot be found.
	ErrUUIDNotFound = errors.New("could not determine machine UUID")
)
