// Package selinux provides utilities to detect the SELinux (Security-Enhanced Linux).
package selinux

import (
	"fmt"
	"os"
	"strings"

	"github.com/kommodity-io/kommodity-attestation-extension/pkg/utils"
)

const (
	seLinuxEnabledPath = "/sys/fs/selinux/enforce"
)

// Attestable implements the report.Attestable interface for SELinux.
type Attestable struct {
	mode      string
	timestamp string
}

// Name returns the name of the attestable component.
func (a *Attestable) Name() string {
	return "selinux"
}

// Measure returns the measurement of the SELinux status.
func (a *Attestable) Measure() (string, error) {
	mode, err := GetSELinuxMode()
	if err != nil {
		return "", fmt.Errorf("failed to determine SELinux status: %w", err)
	}

	a.timestamp = utils.UnixNowString()
	a.mode = mode

	return utils.EncodeMeasurement([]byte(mode)), nil
}

// GetPCRs returns the PCR indices relevant to SELinux.
func (a *Attestable) GetPCRs() (map[int]string, error) {
	return map[int]string{}, nil
}

// Quote returns a dummy quote for SELinux (WARNING: mock implementation).
func (a *Attestable) Quote(nonce []byte) ([]byte, error) {
	return nonce, nil
}

// Evidence returns metadata about the SELinux status.
func (a *Attestable) Evidence() (map[string]string, error) {
	return map[string]string{
		a.Name():    a.mode,
		"timestamp": a.timestamp,
	}, nil
}

// GetSELinuxMode checks the current SELinux mode on the machine.
func GetSELinuxMode() (string, error) {
	data, err := os.ReadFile(seLinuxEnabledPath)
	if err != nil {
		return "", fmt.Errorf("failed to read SELinux enforce file: %w", err)
	}

	modes := map[string]string{"1": "Enforcing", "0": "Permissive"}

	return modes[strings.TrimSpace(string(data))], nil
}
