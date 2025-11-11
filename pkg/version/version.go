// Package version provides version information for talos machine.
package version

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/utils"
)

const (
	osReleaseFileName = "os-release"
)

// Attestable implements the report.Attestable interface for Talos version.
type Attestable struct {
	version   string
	timestamp string
}

// Name returns the name of the attestable component.
func (a *Attestable) Name() string {
	return "talos-version"
}

// Measure returns the measurement of the Talos version.
func (a *Attestable) Measure() (string, error) {
	version, err := getTalosVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get Talos version: %w", err)
	}

	a.timestamp = utils.UnixNowString()
	a.version = version

	return utils.EncodeMeasurement([]byte(version)), nil
}

// Evidence returns metadata about the Talos version.
func (a *Attestable) Evidence() (map[string]string, error) {
	return map[string]string{
		a.Name():    a.version,
		"timestamp": a.timestamp,
	}, nil
}

func getTalosVersion() (string, error) {
	var version string

	osReleaseFilePaths := []string{
		"/etc/",
		"/usr/lib/",
	}

	for _, path := range osReleaseFilePaths {
		releaseFile := path + osReleaseFileName

		_, err := os.Stat(releaseFile)
		if err == nil {
			version, err = getVersionID(releaseFile)
			if err == nil {
				break
			}
		}
	}

	if version == "" {
		return "", ErrVersionNotFound
	}

	return version, nil
}

func getVersionID(path string) (string, error) {
	envMap, err := godotenv.Read(path)
	if err != nil {
		return "", fmt.Errorf("failed to read os-release file: %w", err)
	}

	return envMap["VERSION_ID"], nil
}
