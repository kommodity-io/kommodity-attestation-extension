// Package selinux provides utilities to detect the SELinux (Security-Enhanced Linux).
package selinux

import (
	"fmt"
	"os"
	"strings"
)

const (
	seLinuxEnabledPath = "/sys/fs/selinux/enforce"
)

// GetSELinuxMode checks the current SELinux mode on the machine.
func GetSELinuxMode() (string, error) {
	data, err := os.ReadFile(seLinuxEnabledPath)
	if err != nil {
		return "", fmt.Errorf("failed to read SELinux enforce file: %w", err)
	}

	modes := map[string]string{"1": "Enforcing", "0": "Permissive"}

	return modes[strings.TrimSpace(string(data))], nil
}
