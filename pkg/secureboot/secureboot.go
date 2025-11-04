// Package secureboot provides secure boot information for talos machine.
package secureboot

import (
	"fmt"
	"os"
)

const (
	efiVarsDir        = "/sys/firmware/efi/efivars"
	efiGlobalVarGUID  = "8be4df61-93ca-11d2-aa0d-00e098032b8c" // EFI_GLOBAL_VARIABLE
	secureBootVarPath = efiVarsDir + "/SecureBoot-" + efiGlobalVarGUID
)

// IsSecureBootEnabled checks if Secure Boot is enabled on the machine.
func IsSecureBootEnabled() (bool, error) {
	data, err := os.ReadFile(secureBootVarPath)
	if err != nil {
		return false, fmt.Errorf("failed to read SecureBoot variable: %w", err)
	}

	// SecureBoot value is the last byte; 1 means enabled, 0 means disabled
	return data[len(data)-1] == 1, nil
}
