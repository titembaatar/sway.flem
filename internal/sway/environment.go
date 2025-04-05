package sway

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
)

type Environment struct {
	Config     *config.Config
	Workspaces map[string]*Workspace
}

func NewEnvironment(cfg *config.Config) *Environment {
	return &Environment{
		Config:     cfg,
		Workspaces: make(map[string]*Workspace),
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
		log.Warn("Some workspace focusing operations failed: %v", err)
	}

	op.End()
	return nil
}

func (e *Environment) validateEnvironment() error {
	envOp := log.Operation("Dependency validation")
	envOp.Begin()

	if err := checkCommand("swaymsg"); err != nil {
		envOp.EndWithError(err)
		return fmt.Errorf("swaymsg not found: %w", err)
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

	for name, workspaceConfig := range e.Config.Workspaces {
		log.Info("Processing workspace: %s", name)

		workspace := NewWorkspace(name, workspaceConfig.Layout.String())
		e.Workspaces[name] = workspace

		if err := workspace.Setup(workspaceConfig); err != nil {
			log.Error("Failed to set up workspace %s: %v", name, err)
			continue
		}
	}

	elapsed := time.Since(startTime)
	log.Info("Environment setup completed in %.2f seconds", elapsed.Seconds())

	setupOp.End()
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

func SetupEnvironment(cfg *config.Config) error {
	env := NewEnvironment(cfg)
	return env.Setup()
}

func checkCommand(command string) error {
	log.Debug("Checking if %s is available", command)

	cmd := exec.Command(command, "-v")
	output, err := cmd.CombinedOutput()

	if err != nil {
		if len(output) > 0 {
			log.Error("Command output: %s", string(output))
		}
		return fmt.Errorf("%s is not available: %w", command, err)
	}

	outputStr := strings.TrimSpace(string(output))
	log.Debug("Command '%s' is available, version: %s", command, outputStr)

	return nil
}
