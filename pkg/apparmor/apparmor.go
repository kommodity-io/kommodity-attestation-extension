// Package apparmor provides utilities to detect the AppArmor (Linux Security Module).
package apparmor

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kommodity-io/kommodity-attestation-extension/pkg/utils"
)

const (
	appArmorEnabledPath = "/sys/module/apparmor/parameters/enabled"
)

// Attestable implements the report.Attestable interface for AppArmor.
type Attestable struct {
	enabled   bool
	timestamp string
}

// Name returns the name of the attestable component.
func (a *Attestable) Name() string {
	return "apparmor"
}

// Measure returns the measurement of the AppArmor status.
func (a *Attestable) Measure() (string, error) {
	enabled, err := isAppArmorEnabled()
	if err != nil {
		return "", fmt.Errorf("failed to determine AppArmor status: %w", err)
	}

	a.timestamp = utils.UnixNowString()
	a.enabled = enabled

	return utils.BoolToMeasurement(enabled), nil
}

// GetPCRs returns the PCR indices relevant to AppArmor.
func (a *Attestable) GetPCRs() (map[int]string, error) {
	return map[int]string{}, nil
}

// Quote returns a dummy quote for AppArmor (WARNING: mock implementation).
func (a *Attestable) Quote(nonce []byte) ([]byte, error) {
	return nonce, nil
}

// Evidence returns metadata about the AppArmor status.
func (a *Attestable) Evidence() (map[string]string, error) {
	return map[string]string{
		a.Name():    strconv.FormatBool(a.enabled),
		"timestamp": a.timestamp,
	}, nil
}

func isAppArmorEnabled() (bool, error) {
	data, err := os.ReadFile(appArmorEnabledPath)
	if err != nil {
		return false, fmt.Errorf("failed to read AppArmor enabled file: %w", err)
	}

	return strings.TrimSpace(string(data)) == "Y", nil
}
