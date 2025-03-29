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

	// Set default values if needed
	applyDefaults(&config)

	return &config, nil
}

func applyDefaults(config *Config) {
	for num, workspace := range config.Workspaces {
		if workspace.Layout == "" && config.Defaults.DefaultLayout != "" {
			workspace.Layout = config.Defaults.DefaultLayout
			config.Workspaces[num] = workspace
		}

		if workspace.Output == "" && config.Defaults.DefaultOutput != "" {
			workspace.Output = config.Defaults.DefaultOutput
			config.Workspaces[num] = workspace
		}

		for i, app := range workspace.Apps {
			if app.Command == "" {
				app.Command = app.Name
				workspace.Apps[i] = app
			}

			if !app.Floating && config.Defaults.DefaultFloating {
				app.Floating = true
				workspace.Apps[i] = app
			}
		}
	}
}
