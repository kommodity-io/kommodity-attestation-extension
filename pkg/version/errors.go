// Package version provides error definitions for version operations.
package version

import "errors"

var (
	// ErrVersionNotFound is returned when the version cannot be determined.
	ErrVersionNotFound = errors.New("could not determine talos version")
)
