// Package utils provides utility functions for the attestation extension.
//
//nolint:revive
package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"time"
)

// BoolToMeasurement converts a boolean value to its hex-encoded string representation.
func BoolToMeasurement(b bool) string {
	var status string
	if b {
		status = "enabled"
	} else {
		status = "disabled"
	}

	return EncodeMeasurement([]byte(status))
}

// EncodeMeasurement encodes the given byte slice to a hex string.
func EncodeMeasurement(data []byte) string {
	hash := sha512.Sum512(data)

	return hex.EncodeToString(hash[:])
}

func UnixNowString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
