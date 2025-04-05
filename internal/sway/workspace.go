package sway

import (
	"fmt"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
)

type Workspace struct {
	Name      string
	Layout    string
	Container *Container
}

func NewWorkspace(name string, layout string) *Workspace {
	return &Workspace{
		Name:   name,
		Layout: layout,
	}
}

func (w *Workspace) Setup(workspace config.Workspace) error {
	log.Info("Setting up workspace: %s", w.Name)

	if err := w.Create(); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	container, err := ProcessContainer(w.Name, workspace.Containers, w.Layout, 0)
	if err != nil {
		log.Error("Failed to process container for workspace %s: %v", w.Name, err)
		return err
	}

	w.Container = container

	if err := container.Setup(); err != nil {
		log.Error("Failed to setup container for workspace %s: %v", w.Name, err)
		return err
	}

	container.ResizeApps()

	log.Info("Workspace %s setup complete", w.Name)
	return nil
}

func (w *Workspace) Switch() error {
	command := fmt.Sprintf("workspace %s", w.Name)
	cmd := NewSwayCmd(command)
	_, err := cmd.Run()
	return err
}

func SwitchWorkspace(workspace string) error {
	ws := NewWorkspace(workspace, "")
	return ws.Switch()
}

func (w *Workspace) Create() error {
	if err := w.Switch(); err != nil {
		return fmt.Errorf("%w: failed to switch to workspace '%s': %v",
			ErrWorkspaceCreateFailed, w.Name, err)
	}

	if w.Layout == "" {
		return nil // No layout to set
	}

	cmd := fmt.Sprintf("layout %s", w.Layout)
	_, err := RunSwayCmd(cmd)
	if err != nil {
		return fmt.Errorf("%w: failed to set layout '%s' for workspace '%s': %v",
			ErrWorkspaceCreateFailed, w.Layout, w.Name, err)
	}

	log.Info("Successfully created workspace '%s' with layout '%s'", w.Name, w.Layout)
	return nil
}

func CreateWorkspace(name string, layout string) error {
	ws := NewWorkspace(name, layout)
	return ws.Create()
}

func FocusWorkspaces(workspaces []string) error {
	log.Info("Focusing on %d workspaces", len(workspaces))
	var errors []string

	for i, workspace := range workspaces {
		ws := NewWorkspace(workspace, "")
		log.Debug("Focusing on workspace: %s", workspace)

		if err := ws.Switch(); err != nil {
			log.Error("Failed to focus on workspace %s: %v", workspace, err)
			errors = append(errors, fmt.Sprintf("workspace %s: %v", workspace, err))
		} else {
			log.Info("Successfully focused on workspace %s (%d of %d)", workspace, i+1, len(workspaces))
			time.Sleep(100 * time.Millisecond)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to focus on some workspaces: %s", strings.Join(errors, "; "))
	}

	return nil
}
