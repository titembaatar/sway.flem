package sway

import (
	"fmt"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	errs "github.com/titembaatar/sway.flem/internal/errors"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/internal/util"
)

type Environment struct {
	Config       *config.Config
	Workspaces   map[string]*Workspace
	ErrorHandler *errs.ErrorHandler
}

func NewEnvironment(cfg *config.Config, errorHandler *errs.ErrorHandler) *Environment {
	if errorHandler == nil {
		// Create a default error handler if none provided
		errorHandler = errs.NewErrorHandler(false, true, false)
	}

	return &Environment{
		Config:       cfg,
		Workspaces:   make(map[string]*Workspace),
		ErrorHandler: errorHandler,
	}
}

func (e *Environment) Setup() error {
	log.SetComponent(log.ComponentEnv)

	op := log.Operation("Environment setup")
	op.Begin()

	if err := e.validateEnvironment(); err != nil {
		op.EndWithError(err)
		return err
	}

	if err := e.setupWorkspaces(); err != nil {
		op.EndWithError(err)
		return err
	}

	if err := e.focusRequestedWorkspaces(); err != nil {
		// Non-fatal error, just log a warning
		log.Warn("Some workspace focusing operations failed: %v", err)
	}

	op.End()
	return nil
}

func (e *Environment) validateEnvironment() error {
	envOp := log.Operation("Dependency validation")
	envOp.Begin()

	if err := util.CheckCommand("swaymsg"); err != nil {
		envOp.EndWithError(err)
		return err
	}

	log.Debug("All dependencies are available")

	envOp.End()
	return nil
}

func (e *Environment) setupWorkspaces() error {
	setupOp := log.Operation("Workspaces setup")
	setupOp.Begin()

	workspaceCount := len(e.Config.Workspaces)
	log.Info("Configuring %d workspaces", workspaceCount)

	startTime := time.Now()
	errorCount := 0

	for name, workspaceConfig := range e.Config.Workspaces {
		log.Info("Processing workspace: %s", name)

		workspace := NewWorkspace(name, workspaceConfig.Layout.String())
		e.Workspaces[name] = workspace

		if err := workspace.Setup(workspaceConfig, e.ErrorHandler); err != nil {
			errorCount++
			setupErr := errs.Wrap(err, fmt.Sprintf("Failed to set up workspace '%s'", name))
			e.ErrorHandler.Handle(setupErr)

			// Continue with other workspaces even if this one failed
			continue
		}
	}

	elapsed := time.Since(startTime)
	log.Info("Environment setup completed in %.2f seconds", elapsed.Seconds())

	setupOp.End()

	if errorCount > 0 {
		return fmt.Errorf("failed to set up %d workspace(s)", errorCount)
	}

	return nil
}

func (e *Environment) focusRequestedWorkspaces() error {
	if len(e.Config.Focus) == 0 {
		return nil
	}

	focusOp := log.Operation("workspace focusing")
	focusOp.Begin()

	log.Info("Focusing on %d specified workspaces: %v", len(e.Config.Focus), e.Config.Focus)
	return FocusWorkspaces(e.Config.Focus)
}

func SetupEnvironment(cfg *config.Config, errorHandler *errs.ErrorHandler) error {
	env := NewEnvironment(cfg, errorHandler)
	return env.Setup()
}

// No longer needed - using utility function instead
