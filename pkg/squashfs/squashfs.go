// Package squashfs provides utilities to check for squashfs filesystems for talos machine.
package squashfs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kommodity-io/kommodity-attestation-extension/pkg/utils"
)

const (
	fileExpectedFieldsCount = 4
)

// Attestable implements the report.Attestable interface for squashfs.
type Attestable struct {
	enabled   bool
	timestamp string
}

// Name returns the name of the attestable component.
func (a *Attestable) Name() string {
	return "squashfs"
}

// Measure checks if the root filesystem is a read-only squashfs and returns its status.
func (a *Attestable) Measure() (string, error) {
	enabled, err := IsRootSquashfsReadOnly()
	if err != nil {
		return "", fmt.Errorf("failed to measure squashfs: %w", err)
	}

	a.timestamp = utils.UnixNowString()
	a.enabled = enabled

	return utils.BoolToMeasurement(enabled), nil
}

// GetPCRs returns the PCR indices relevant to squashfs.
func (a *Attestable) GetPCRs() (map[int]string, error) {
	return map[int]string{}, nil
}

// Quote returns a dummy quote for squashfs (WARNING: mock implementation).
func (a *Attestable) Quote(nonce []byte) ([]byte, error) {
	return nonce, nil
}

// Evidence returns metadata about the squashfs status.
func (a *Attestable) Evidence() (map[string]string, error) {
	return map[string]string{
		a.Name():    strconv.FormatBool(a.enabled),
		"timestamp": a.timestamp,
	}, nil
}

// IsRootSquashfsReadOnly checks if the root filesystem is a read-only squashfs.
func IsRootSquashfsReadOnly() (bool, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return false, fmt.Errorf("failed to open /proc/mounts: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < fileExpectedFieldsCount {
			continue
		}

		mountPoint := fields[1]
		fsType := fields[2]
		options := fields[3]

		if mountPoint == "/" && fsType == "squashfs" && strings.Contains(options, "ro") {
			return true, nil
		}
	}

	return false, fmt.Errorf("root is not a read-only squashfs filesystem: %w", scanner.Err())
}
