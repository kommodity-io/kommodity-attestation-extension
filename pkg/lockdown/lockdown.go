// Package lockdown provides utilities for managing Talos kernel lockdown mode.
package lockdown

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	lockdownModePath = "/sys/kernel/security/lockdown"
	lockdownRegex    = `\[(\w+)\]`
)

// GetKernelLockdownMode reads the current kernel lockdown mode from the system.
func GetKernelLockdownMode() (string, error) {
	data, err := os.ReadFile(lockdownModePath)
	if err != nil {
		return "", fmt.Errorf("failed to read lockdown mode: %w", err)
	}

	line := strings.TrimSpace(string(data))

	re := regexp.MustCompile(lockdownRegex)
	mode := re.FindString(line)

	return mode, nil
}
