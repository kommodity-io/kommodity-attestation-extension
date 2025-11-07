// Package uuid provides utilities for fetching the Talos machine UUID.
package uuid

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

// GetMachineUUID retrieves and returns the machine UUID from the system's DMI information.
func GetMachineUUID() (string, error) {
	candidates := []string{
		"/sys/class/dmi/id/product_uuid",
		"/sys/devices/virtual/dmi/id/product_uuid",
	}

	var raw string

	for _, p := range candidates {
		//nolint:gosec // Secure and known file
		b, err := os.ReadFile(p)
		if err == nil && len(b) > 0 {
			raw = strings.TrimSpace(string(b))

			break
		}
	}

	if raw == "" {
		return "", ErrUUIDNotFound
	}

	// Canonicalize to RFC 4122 lowercase 8-4-4-4-12 form
	u, err := uuid.Parse(strings.ToLower(raw))
	if err != nil {
		return "", fmt.Errorf("invalid UUID in SMBIOS (%q): %w", raw, err)
	}

	return u.String(), nil
}
