package sway

import (
	"fmt"
	"strings"

	"github.com/titembaatar/sway.flem/internal/log"
)

type Workspace struct {
	Name   string
	Layout string
}

func NewWorkspace(name string, layout string) *Workspace {
	return &Workspace{
		Name:   name,
		Layout: layout,
	}
}

func (w *Workspace) List() ([]string, error) {
	log.Debug("Getting workspaces from sway")

	type workspaceInfo struct {
		Name string `json:"name"`
	}

	var workspaces []workspaceInfo

	cmd := NewRawSwayCmd("", "get_workspaces")
	if err := cmd.GetJSON(&workspaces); err != nil {
		return nil, err
	}

	names := make([]string, len(workspaces))
	for i, ws := range workspaces {
		names[i] = ws.Name
	}

	log.Debug("Found %d workspaces: %s", len(names), strings.Join(names, ", "))
	return names, nil
}

func (w *Workspace) Switch() error {
	command := fmt.Sprintf("workspace %s", w.Name)
	cmd := NewSwayCmd(command)
	_, err := cmd.Run()
	return err
}

func (w *Workspace) Create() error {
	if err := w.Switch(); err != nil {
		return fmt.Errorf("%w: failed to switch to workspace '%s': %v",
			ErrWorkspaceCreateFailed, w.Name, err)
	}

	if w.Layout == "" {
		return nil // No layout to set
	}

	command := fmt.Sprintf("layout %s", w.Layout)
	cmd := NewSwayCmd(command)
	_, err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%w: failed to set layout '%s' for workspace '%s': %v",
			ErrWorkspaceCreateFailed, w.Layout, w.Name, err)
	}

	log.Info("Successfully created workspace '%s' with layout '%s'", w.Name, w.Layout)
	return nil
}
