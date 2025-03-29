package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	return LoadFromBytes(data)
}

func LoadFromBytes(data []byte) (*Config, error) {
	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	ApplyDefaults(&config)

	return &config, nil
}
