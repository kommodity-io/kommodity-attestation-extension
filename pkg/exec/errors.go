// Package exec provides error definitions for execution operations.
package exec

import "errors"

var (
	// ErrArgMissing is returned when a required command-line argument is missing.
	ErrArgMissing = errors.New("required argument is missing")
)
