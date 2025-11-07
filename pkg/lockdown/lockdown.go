// Package lockdown provides utilities for managing Talos kernel lockdown mode.
package lockdown

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kommodity-io/kommodity-attestation-extension/pkg/utils"
)

const (
	lockdownModePath = "/sys/kernel/security/lockdown"
	lockdownRegex    = `\[(\w+)\]`
)

// Attestable implements the report.Attestable interface for kernel lockdown mode.
type Attestable struct {
	mode      string
	timestamp string
}

// Name returns the name of the attestable component.
func (a *Attestable) Name() string {
	return "lockdown"
}

// Measure returns the measurement of the kernel lockdown mode.
func (a *Attestable) Measure() (string, error) {
	mode, err := GetKernelLockdownMode()
	if err != nil {
		return "", fmt.Errorf("failed to get kernel lockdown mode: %w", err)
	}

	a.timestamp = utils.UnixNowString()
	a.mode = mode

	return utils.EncodeMeasurement([]byte(mode)), nil
}

// GetPCRs returns the PCR indices relevant to kernel lockdown mode.
func (a *Attestable) GetPCRs() (map[int]string, error) {
	return map[int]string{}, nil
}

// Quote returns a dummy quote for kernel lockdown mode (not implemented).
func (a *Attestable) Quote(nonce []byte) ([]byte, error) {
	return nonce, nil
}

// Evidence returns metadata about the kernel lockdown mode.
func (a *Attestable) Evidence() (map[string]string, error) {
	return map[string]string{
		a.Name():    a.mode,
		"timestamp": a.timestamp,
	}, nil
}

// GetKernelLockdownMode reads the current kernel lockdown mode from the system.
func GetKernelLockdownMode() (string, error) {
	data, err := os.ReadFile(lockdownModePath)
	if err != nil {
		return "", fmt.Errorf("failed to read lockdown mode: %w", err)
	}

	line := strings.TrimSpace(string(data))

	re := regexp.MustCompile(lockdownRegex)
	mode := re.FindString(line)

	return strings.ToLower(mode), nil
}
