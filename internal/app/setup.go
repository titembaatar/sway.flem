package app

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/internal/sway"
)

// Initializes and configures the Sway environment based on the configuration
func Setup(config *config.Config) error {
	log.Info("Starting environment setup")

	if err := validateEnvironment(); err != nil {
		return err
	}

	if err := executeSetup(config); err != nil {
		return err
	}

	if err := focusRequestedWorkspaces(config); err != nil {
		log.Warn("Some workspace focusing operations failed: %v", err)
	}

	return nil
}

// Verifies that all required external dependencies are available
func validateEnvironment() error {
	log.Debug("Validating environment and dependencies")

	if err := checkCommand("swaymsg"); err != nil {
		return fmt.Errorf("swaymsg not found: %w", err)
	}

	log.Debug("All dependencies are available")
	return nil
}

// Execute the environment setup
func executeSetup(config *config.Config) error {
	log.Info("Configuring Sway environment")

	startTime := time.Now()

	if err := sway.SetupEnvironment(config); err != nil {
		return fmt.Errorf("failed to setup environment: %w", err)
	}

	elapsed := time.Since(startTime)
	log.Info("Environment setup completed in %.2f seconds", elapsed.Seconds())

	return nil
}

// Focus on workspaces specified in the config
func focusRequestedWorkspaces(config *config.Config) error {
	if len(config.Focus) == 0 {
		return nil
	}

	log.Info("Focusing on specified workspaces: %v", config.Focus)
	return sway.FocusWorkspaces(config.Focus)
}

// Checks if a command is available in the PATH
func checkCommand(command string) error {
	log.Debug("Checking if %s is available", command)

	cmd := exec.Command(command, "-v")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("%s is not available: %v", command, err)
		if len(output) > 0 {
			log.Error("Command output: %s", string(output))
		}
		return fmt.Errorf("%s is not available: %w", command, err)
	}

	log.Debug("%s is available, version: %s", command, string(output))
	return nil
}
