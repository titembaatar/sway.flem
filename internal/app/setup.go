package app

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/internal/sway"
)

// Setup initializes and configures the Sway environment based on the configuration
func Setup(cfg *config.Config) error {
	log.Info("Starting environment setup")

	// Check if swaymsg is available
	if err := checkDependencies(); err != nil {
		return err
	}

	// Initialize the environment
	startTime := time.Now()
	if err := sway.SetupEnvironment(cfg); err != nil {
		return fmt.Errorf("failed to setup environment: %w", err)
	}

	elapsed := time.Since(startTime)
	log.Info("Environment setup completed in %.2f seconds", elapsed.Seconds())

	return nil
}

// checkDependencies verifies that all required external dependencies are available
func checkDependencies() error {
	log.Debug("Checking dependencies")

	// Check if swaymsg is available
	if err := checkCommand("swaymsg"); err != nil {
		return fmt.Errorf("swaymsg not found: %w", err)
	}

	log.Debug("All dependencies are available")
	return nil
}

// checkCommand checks if a command is available in the PATH
func checkCommand(command string) error {
	log.Debug("Checking if %s is available", command)

	// Try to run the command with -v flag to check version
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
