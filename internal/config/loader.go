package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoWorkspaces = errors.New("no workspaces defined in configuration")
)

func isValidLayout(layout string) bool {
	validLayouts := map[string]bool{
		"splith":   true,
		"splitv":   true,
		"stacking": true,
		"tabbed":   true,
	}

	return validLayouts[layout]
}

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

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if len(config.Workspaces) == 0 {
		return ErrNoWorkspaces
	}

	for wsNum, workspace := range config.Workspaces {
		if len(workspace.Apps) == 0 {
			return fmt.Errorf("workspace %d has no apps defined", wsNum)
		}

		if workspace.Layout != "" {
			if !isValidLayout(workspace.Layout) {
				return fmt.Errorf("workspace %d: invalid layout '%s'", wsNum, workspace.Layout)
			}
		}
	}

	return nil
}
