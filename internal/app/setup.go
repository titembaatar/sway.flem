package app

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	errs "github.com/titembaatar/sway.flem/internal/errors"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/internal/sway"
)

// Setup initializes and configures the Sway environment based on the configuration
func Setup(config *config.Config) error {
	log.SetComponent(log.ComponentApp)

	op := log.Operation("environment setup")
	op.Begin()

	// Create an error handler that doesn't exit on fatal errors
	errorHandler := errs.NewErrorHandler(false, true, false)

	if err := validateEnvironment(); err != nil {
		op.EndWithError(err)
		return err
	}

	if err := executeSetup(config, errorHandler); err != nil {
		op.EndWithError(err)
		return err
	}

	if err := focusRequestedWorkspaces(config, errorHandler); err != nil {
		log.Warn("Some workspace focusing operations failed: %v", err)
		// Continue execution - focus failures are not fatal
	}

	// Report any errors or warnings that were encountered
	errorHandler.SummarizeErrors()

	op.End()
	return nil
}

// validateEnvironment verifies that all required external dependencies are available
func validateEnvironment() error {
	envOp := log.Operation("dependency validation")
	envOp.Begin()

	log.Debug("Checking for required dependencies")

	if err := checkCommand("swaymsg"); err != nil {
		fatalErr := errs.NewFatal(errs.ErrCommandNotFound, "Required command 'swaymsg' not found")
		fatalErr.WithSuggestion("Make sure Sway is installed and swaymsg is in your PATH")
		envOp.EndWithError(fatalErr)
		return fatalErr
	}

	log.Debug("All dependencies are available")
	envOp.End()
	return nil
}

// executeSetup applies the environment configuration
func executeSetup(config *config.Config, errorHandler *errs.ErrorHandler) error {
	setupOp := log.Operation("sway configuration")
	setupOp.Begin()

	workspaceCount := len(config.Workspaces)
	log.Info("Configuring %d workspaces", workspaceCount)

	startTime := time.Now()

	if err := sway.SetupEnvironment(config, errorHandler); err != nil {
		setupErr := errs.Wrap(err, "Failed to setup environment")
		setupOp.EndWithError(setupErr)
		return setupErr
	}

	elapsed := time.Since(startTime)
	log.Info("Environment setup completed in %.2f seconds", elapsed.Seconds())

	setupOp.End()
	return nil
}

// focusRequestedWorkspaces focuses on workspaces specified in the config
func focusRequestedWorkspaces(config *config.Config, errorHandler *errs.ErrorHandler) error {
	if len(config.Focus) == 0 {
		return nil
	}

	focusOp := log.Operation("workspace focusing")
	focusOp.Begin()

	log.Info("Focusing on %d specified workspaces: %v", len(config.Focus), config.Focus)
	var lastError error

	for _, focus := range config.Focus {
		workspace := sway.NewWorkspace(focus, "")
		if err := workspace.Switch(); err != nil {
			focusErr := errs.Wrap(err, fmt.Sprintf("Failed to focus on workspace '%s'", focus))
			errorHandler.Handle(focusErr)
			lastError = focusErr
		}
	}

	focusOp.End()
	return lastError
}

// checkCommand checks if a command is available in the PATH
func checkCommand(command string) error {
	log.Debug("Checking if %s is available", command)

	cmd := exec.Command(command, "-v")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Required command '%s' is not available: %v", command, err)
		if len(output) > 0 {
			log.Error("Command output: %s", string(output))
		}
		return errs.NewFatal(errs.ErrCommandNotFound, fmt.Sprintf("Required command '%s' not found", command)).
			WithSuggestion(fmt.Sprintf("Make sure '%s' is installed and in your PATH", command))
	}

	outputStr := strings.TrimSpace(string(output))
	log.Debug("Command '%s' is available, version: %s", command, outputStr)
	return nil
}
