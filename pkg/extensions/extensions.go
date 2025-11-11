// Package extensions provides utilities for fetching Talos extensions.
package extensions

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kommodity-io/kommodity-attestation-extension/pkg/utils"
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

// Attestable implements the report.Attestable interface for Talos extensions.
type Attestable struct {
	extensions []Extension
	timestamp  string
}

// Name returns the name of the attestable component.
func (a *Attestable) Name() string {
	return "extensions"
}

// Measure returns the measurement of the Talos extensions.
func (a *Attestable) Measure() (string, error) {
	extensions, err := GetExtensions()
	if err != nil {
		return "", fmt.Errorf("failed to get extensions: %w", err)
	}

	data, err := yaml.Marshal(extensions)
	if err != nil {
		return "", fmt.Errorf("failed to marshal extensions: %w", err)
	}

	a.timestamp = utils.UnixNowString()
	a.extensions = extensions

	return utils.EncodeMeasurement(data), nil
}

// Evidence returns metadata about the Talos extensions.
func (a *Attestable) Evidence() (map[string]string, error) {
	evidence := map[string]string{
		"extensions_count": strconv.Itoa(len(a.extensions)),
		"timestamp":        a.timestamp,
	}

	for i, ext := range a.extensions {
		prefix := fmt.Sprintf("extension_%d_", i)
		evidence[prefix+"name"] = ext.Name
		evidence[prefix+"config_name"] = ext.Config.Name
		evidence[prefix+"config_hash"] = ext.Config.Hash
		evidence[prefix+"binary_name"] = ext.Binary.Name
		evidence[prefix+"binary_hash"] = ext.Binary.Hash
	}

	return evidence, nil
}

// GetExtensions retrieves the list of Talos extensions from the system.
func GetExtensions() ([]Extension, error) {
	_, err := os.Stat(extensionsDir)
	if os.IsNotExist(err) {
		return []Extension{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat extensions directory: %w", err)
	}

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
			Hash: utils.EncodeMeasurement(data),
		},
		Binary: File{
			Name: extConfig.Container.Entrypoint,
			Hash: utils.EncodeMeasurement(extensionBinaryData),
		},
	}, nil
}
