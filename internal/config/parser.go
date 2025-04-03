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
	log.SetComponent(log.ComponentConfig)

	loadOp := log.Operation("config loading")
	loadOp.Begin()

	absPath, err := filepath.Abs(path)
	if err != nil {
		loadOp.EndWithError(err)
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	log.Info("Loading configuration from %s", absPath)

	file, err := os.Open(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Error("Configuration file not found: %s", absPath)
			loadErr := fmt.Errorf("config file not found: %w", err)
			loadOp.EndWithError(loadErr)
			return nil, loadErr
		}
		log.Error("Failed to open configuration file: %v", err)
		loadErr := fmt.Errorf("failed to open config file: %w", err)
		loadOp.EndWithError(loadErr)
		return nil, loadErr
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	if err := decoder.Decode(&config); err != nil {
		log.Error("Failed to parse YAML configuration: %v", err)
		loadErr := fmt.Errorf("failed to decode config: %w", err)
		loadOp.EndWithError(loadErr)
		return nil, loadErr
	}

	log.Debug("Successfully parsed configuration, validating...")

	validateOp := log.Operation("config validation")
	validateOp.Begin()

	if err := ValidateConfig(&config); err != nil {
		log.Error("Configuration validation failed: %v", err)
		validateOp.EndWithError(err)
		loadOp.EndWithError(err)
		return nil, err
	}

	validateOp.End()

	workspaceCount := len(config.Workspaces)
	log.Info("Configuration loaded successfully with %d workspaces", workspaceCount)

	if log.GetLevel() <= log.LogLevelDebug {
		for name := range config.Workspaces {
			log.Debug("Found workspace configuration: %s", name)
		}
	}

	loadOp.End()
	return &config, nil
}
