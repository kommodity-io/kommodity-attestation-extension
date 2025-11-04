// Package squashfs provides utilities to check for squashfs filesystems for talos machine.
package squashfs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	fileExpectedFieldsCount = 4
)

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
