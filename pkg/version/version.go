// Package version provides version information for talos machine.
package version

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	osReleaseFileName = "os-release"
)

// GetTalosVersion retrieves the Talos version from the os-release file.
func GetTalosVersion() (string, error) {
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
