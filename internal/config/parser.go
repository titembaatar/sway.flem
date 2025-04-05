package config

import (
	"fmt"
	"os"
	"path/filepath"

	errs "github.com/titembaatar/sway.flem/internal/errors"
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
		return nil, errs.Wrap(err, "Failed to get absolute path for config file")
	}

	log.Info("Loading configuration from %s", absPath)

	file, err := os.Open(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Error("Configuration file not found: %s", absPath)
			loadErr := errs.NewFatalConfigError(errs.ErrConfigNotFound, "", "", -1)
			loadErr.WithConfigFile(absPath, 0)
			loadErr.WithSuggestion(fmt.Sprintf("Check that the file exists at '%s'", absPath))
			loadOp.EndWithError(loadErr)
			return nil, loadErr
		}
		log.Error("Failed to open configuration file: %v", err)
		loadErr := errs.Wrap(err, "Failed to open config file")
		loadOp.EndWithError(loadErr)
		return nil, loadErr
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	if err := decoder.Decode(&config); err != nil {
		log.Error("Failed to parse YAML configuration: %v", err)
		loadErr := errs.NewFatalConfigError(errs.ErrConfigParse, "", "", -1)
		loadErr.WithConfigFile(absPath, 0)
		loadErr.WithSuggestion("Check your YAML syntax for errors")
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
