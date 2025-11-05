// Package image provides utilities for fetching the layers of the image on the Talos machine.
package image

import (
	"fmt"
	"os"

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
