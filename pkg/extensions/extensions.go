// Package extensions provides utilities for fetching Talos extensions.
package extensions

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Extension represents a Talos extension with its configuration and binary files.
type Extension struct {
	Name   string
	Config File
	Binary File
}

// File represents a file associated with an extension, including its name and hash.
type File struct {
	Name string
	Hash string
}

// ExtensionConfig represents the structure of a Talos extension configuration file.
type ExtensionConfig struct {
	Container struct {
		Entrypoint string `yaml:"entrypoint"`
	} `yaml:"container"`
}

const (
	extensionsDir = "/usr/local/etc/containers"
)

// GetExtensions retrieves the list of Talos extensions from the system.
func GetExtensions() ([]Extension, error) {
	files, err := os.ReadDir(extensionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read extensions directory: %w", err)
	}

	var extensions = make([]Extension, 0)

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		extension, err := parseExtension(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse extension %s: %w", file.Name(), err)
		}

		extensions = append(extensions, *extension)
	}

	return extensions, nil
}

func parseExtension(file os.DirEntry) (*Extension, error) {
	path := fmt.Sprintf("%s/%s", extensionsDir, file.Name())

	//nolint:gosec // Extension file path is controlled
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read extension file %s: %w", path, err)
	}

	var extConfig ExtensionConfig

	err = yaml.Unmarshal(data, &extConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal extension config %s: %w", path, err)
	}

	extensionBinaryData, err := os.ReadFile(extConfig.Container.Entrypoint)
	if err != nil {
		return nil, fmt.Errorf("failed to read extension binary %s: %w", extConfig.Container.Entrypoint, err)
	}

	return &Extension{
		Name: strings.TrimSuffix(file.Name(), ".yaml"),
		Config: File{
			Name: file.Name(),
			Hash: fileHash(data),
		},
		Binary: File{
			Name: extConfig.Container.Entrypoint,
			Hash: fileHash(extensionBinaryData),
		},
	}, nil
}

func fileHash(data []byte) string {
	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:])
}
