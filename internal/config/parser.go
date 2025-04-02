package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/titembaatar/sway.flem/internal/log"
	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	log.Info("Loading configuration from %s", absPath)

	file, err := os.Open(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Error("Configuration file not found: %s", absPath)
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		log.Error("Failed to open configuration file: %v", err)
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	if err := decoder.Decode(&config); err != nil {
		log.Error("Failed to parse YAML configuration: %v", err)
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	log.Debug("Successfully parsed configuration, validating...")

	if err := ValidateConfig(&config); err != nil {
		log.Error("Configuration validation failed: %v", err)
		return nil, err
	}

	log.Info("Configuration loaded successfully")
	return &config, nil
}
