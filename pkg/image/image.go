// Package image provides utilities for fetching the layers of the image on the Talos machine.
package image

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kommodity-io/kommodity-attestation-extension/pkg/utils"
	"go.yaml.in/yaml/v3"
)

// LayerList represents the structure of the image layers configuration file.
type LayerList struct {
	Layers []Layer `yaml:"layers"`
}

// Layer represents a single layer in the image with its metadata.
type Layer struct {
	Image    string   `yaml:"image"`
	Metadata Metadata `yaml:"metadata"`
}

// Metadata represents the metadata associated with an image layer.
type Metadata struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Author  string `yaml:"author"`
}

const (
	layersPath = "/etc/extensions.yaml"
)

// Attestable implements the report.Attestable interface for image layers.
type Attestable struct {
	layers    []Layer
	timestamp string
}

// Name returns the name of the attestable component.
func (a *Attestable) Name() string {
	return "image"
}

// Measure returns the measurement of the image layers.
func (a *Attestable) Measure() (string, error) {
	layers, err := GetImageLayers()
	if err != nil {
		return "", fmt.Errorf("failed to get image layers: %w", err)
	}

	data, err := yaml.Marshal(layers)
	if err != nil {
		return "", fmt.Errorf("failed to marshal image layers: %w", err)
	}

	return utils.EncodeMeasurement(data), nil
}

// Evidence returns metadata about the image layers.
func (a *Attestable) Evidence() (map[string]string, error) {
	evidence := map[string]string{
		"layers_count": strconv.Itoa(len(a.layers)),
		"timestamp":    a.timestamp,
	}

	for i, layer := range a.layers {
		prefix := fmt.Sprintf("layer_%d_", i)
		evidence[prefix+"name"] = layer.Metadata.Name
		evidence[prefix+"version"] = layer.Metadata.Version
		evidence[prefix+"author"] = layer.Metadata.Author
		evidence[prefix+"image"] = layer.Image
	}

	return evidence, nil
}

// GetImageLayers retrieves the list of image layers from the system.
func GetImageLayers() (*LayerList, error) {
	return parseLayers(layersPath)
}

func parseLayers(path string) (*LayerList, error) {
	//nolint:gosec // Reading well-known file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read layers file %s: %w", path, err)
	}

	var layers LayerList

	err = yaml.Unmarshal(data, &layers)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal layers file %s: %w", path, err)
	}

	return &layers, nil
}
