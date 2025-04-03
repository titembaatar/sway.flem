package app

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/internal/sway"
)

// Initializes and configures the Sway environment based on the configuration
func Setup(config *config.Config) error {
	log.SetComponent(log.ComponentApp)

	op := log.Operation("environment setup")
	op.Begin()

	if err := validateEnvironment(); err != nil {
		op.EndWithError(err)
		return err
	}

	if err := executeSetup(config); err != nil {
		op.EndWithError(err)
		return err
	}

	if err := focusRequestedWorkspaces(config); err != nil {
		log.Warn("Some workspace focusing operations failed: %v", err)
	}

	op.End()
	return nil
}

// Verifies that all required external dependencies are available
func validateEnvironment() error {
	envOp := log.Operation("dependency validation")
	envOp.Begin()

	log.Debug("Checking for required dependencies")

	if err := checkCommand("swaymsg"); err != nil {
		envOp.EndWithError(err)
		return fmt.Errorf("swaymsg not found: %w", err)
	}

	log.Debug("All dependencies are available")
	envOp.End()
	return nil
}

// Execute the environment setup
func executeSetup(config *config.Config) error {
	setupOp := log.Operation("sway configuration")
	setupOp.Begin()

	workspaceCount := len(config.Workspaces)
	log.Info("Configuring %d workspaces", workspaceCount)

	startTime := time.Now()

	if err := sway.SetupEnvironment(config); err != nil {
		setupOp.EndWithError(err)
		return fmt.Errorf("failed to setup environment: %w", err)
	}

	elapsed := time.Since(startTime)
	log.Info("Environment setup completed in %.2f seconds", elapsed.Seconds())

	setupOp.End()
	return nil
}

// Focus on workspaces specified in the config
func focusRequestedWorkspaces(config *config.Config) error {
	if len(config.Focus) == 0 {
		return nil
	}

	focusOp := log.Operation("workspace focusing")
	focusOp.Begin()

	log.Info("Focusing on %d specified workspaces: %v", len(config.Focus), config.Focus)
	err := sway.FocusWorkspaces(config.Focus)

	if err != nil {
		focusOp.EndWithError(err)
	} else {
		focusOp.End()
	}

	return err
}

// Checks if a command is available in the PATH
func checkCommand(command string) error {
	log.Debug("Checking if %s is available", command)

	cmd := exec.Command(command, "-v")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Required command '%s' is not available: %v", command, err)
		if len(output) > 0 {
			log.Error("Command output: %s", string(output))
		}
		return fmt.Errorf("%s is not available: %w", command, err)
	}

	outputStr := strings.TrimSpace(string(output))
	log.Debug("Command '%s' is available, version: %s", command, outputStr)
	return nil
}
