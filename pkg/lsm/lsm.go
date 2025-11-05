// Package lsm provides utilities to detect the LSM (Linux Security Module).
package lsm

import (
	"fmt"
	"os"
	"strings"
)

const (
	seLinuxEnabledPath  = "/sys/fs/selinux/enforce"
	appArmorEnabledPath = "/sys/module/apparmor/parameters/enabled"
)

// IsLSMEnabled checks if any Linux Security Module (SELinux or AppArmor) is enabled on the machine.
func IsLSMEnabled() (bool, error) {
	selinuxEnabled, err := isSELinuxEnabled()
	if err != nil {
		return false, fmt.Errorf("failed to check SELinux status: %w", err)
	}

	if selinuxEnabled {
		return true, nil
	}

	appArmorEnabled, err := isAppArmorEnabled()
	if err != nil {
		return false, fmt.Errorf("failed to check AppArmor status: %w", err)
	}

	if appArmorEnabled {
		return true, nil
	}

	return false, nil
}

func isSELinuxEnabled() (bool, error) {
	data, err := os.ReadFile(seLinuxEnabledPath)
	if err != nil {
		return false, fmt.Errorf("failed to read SELinux enforce file: %w", err)
	}

	return strings.TrimSpace(string(data)) == "1", nil
}

func isAppArmorEnabled() (bool, error) {
	data, err := os.ReadFile(appArmorEnabledPath)
	if err != nil {
		return false, fmt.Errorf("failed to read AppArmor enabled file: %w", err)
	}

	return strings.TrimSpace(string(data)) == "Y", nil
}
