// Package cmdline provides functionality to parse the /proc/cmdline file.
package cmdline

import (
	"fmt"
	"os"
	"strings"
)

const (
	procCmdlinePath    = "/proc/cmdline"
	expectedParamCount = 2
)

// ParseProcCmdline reads and parses the /proc/cmdline file into a map of key-value pairs.
func ParseProcCmdline() (map[string]string, error) {
	data, err := os.ReadFile(procCmdlinePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", procCmdlinePath, err)
	}

	result := make(map[string]string)

	params := strings.Fields(string(data))
	for _, param := range params {
		parts := strings.Split(param, "=")

		if len(parts) == expectedParamCount {
			result[parts[0]] = parts[1]
		} else {
			result[parts[0]] = ""
		}
	}

	return result, nil
}
