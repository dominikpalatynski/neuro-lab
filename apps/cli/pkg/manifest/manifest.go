package manifest

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       interface{} `yaml:"spec"`
}

type Metadata struct {
	Name string `yaml:"name"`
}

func ParseManifest(filename string) (*Manifest, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}
