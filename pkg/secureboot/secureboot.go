// Package secureboot provides secure boot information for talos machine.
package secureboot

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kommodity-io/kommodity-attestation-extension/pkg/utils"
)

const (
	efiVarsDir        = "/sys/firmware/efi/efivars"
	efiGlobalVarGUID  = "8be4df61-93ca-11d2-aa0d-00e098032b8c" // EFI_GLOBAL_VARIABLE
	secureBootVarPath = efiVarsDir + "/SecureBoot-" + efiGlobalVarGUID
)

// Attestable implements the report.Attestable interface for Secure Boot.
type Attestable struct {
	enabled   bool
	timestamp string
}

// Name returns the name of the attestable component.
func (a *Attestable) Name() string {
	return "secure-boot"
}

// Measure returns the measurement of the Secure Boot status.
func (a *Attestable) Measure() (string, error) {
	enabled, err := IsSecureBootEnabled()
	if err != nil {
		return "", fmt.Errorf("failed to determine Secure Boot status: %w", err)
	}

	a.timestamp = utils.UnixNowString()
	a.enabled = enabled

	return utils.BoolToMeasurement(enabled), nil
}

// Evidence returns metadata about the Secure Boot status.
func (a *Attestable) Evidence() (map[string]string, error) {
	return map[string]string{
		a.Name():    strconv.FormatBool(a.enabled),
		"timestamp": a.timestamp,
	}, nil
}

// IsSecureBootEnabled checks if Secure Boot is enabled on the machine.
func IsSecureBootEnabled() (bool, error) {
	data, err := os.ReadFile(secureBootVarPath)
	if err != nil {
		return false, fmt.Errorf("failed to read SecureBoot variable: %w", err)
	}

	// SecureBoot value is the last byte; 1 means enabled, 0 means disabled
	return data[len(data)-1] == 1, nil
}
