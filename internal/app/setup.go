package app

import (
	"fmt"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	errs "github.com/titembaatar/sway.flem/internal/errors"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/internal/sway"
	"github.com/titembaatar/sway.flem/internal/util"
)

func Setup(config *config.Config) error {
	log.SetComponent(log.ComponentApp)

	op := log.Operation("environment setup")
	op.Begin()

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
	}

	errorHandler.SummarizeErrors()

	op.End()
	return nil
}

func validateEnvironment() error {
	if err := util.CheckCommand("swaymsg"); err != nil {
		fatalErr := errs.NewFatal(errs.ErrCommandNotFound, "Required command 'swaymsg' not found")
		fatalErr.WithSuggestion("Make sure Sway is installed and swaymsg is in your PATH")
		return fatalErr
	}

	return nil
}

func executeSetup(config *config.Config, errorHandler *errs.ErrorHandler) error {
	setupOp := log.Operation("sway configuration")
	setupOp.Begin()

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

func focusRequestedWorkspaces(config *config.Config, errorHandler *errs.ErrorHandler) error {
	if len(config.Focus) == 0 {
		return nil
	}

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

	return lastError
}
