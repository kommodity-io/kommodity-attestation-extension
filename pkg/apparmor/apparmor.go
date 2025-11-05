// Package apparmor provides utilities to detect the AppArmor (Linux Security Module).
package apparmor

import (
	"fmt"
	"os"
	"strings"
)

const (
	appArmorEnabledPath = "/sys/module/apparmor/parameters/enabled"
)

// IsAppArmorEnabled checks if AppArmor is enabled on the machine.
func IsAppArmorEnabled() (bool, error) {
	data, err := os.ReadFile(appArmorEnabledPath)
	if err != nil {
		return false, fmt.Errorf("failed to read AppArmor enabled file: %w", err)
	}

	return strings.TrimSpace(string(data)) == "Y", nil
}
