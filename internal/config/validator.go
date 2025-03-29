package config

import (
	"errors"
	"fmt"
)

var (
	ErrNoWorkspaces = errors.New("no workspaces defined in configuration")
)

func ValidateConfig(config *Config) error {
	if len(config.Workspaces) == 0 {
		return ErrNoWorkspaces
	}

	// Validate each workspace
	for wsNum, workspace := range config.Workspaces {
		if err := validateWorkspace(wsNum, workspace); err != nil {
			return err
		}
	}

	return nil
}

func validateWorkspace(wsNum int, workspace Workspace) error {
	if len(workspace.Apps) == 0 {
		return fmt.Errorf("workspace %d has no apps defined", wsNum)
	}

	if workspace.Layout != "" && !isValidLayout(workspace.Layout) {
		return fmt.Errorf("workspace %d: invalid layout '%s'", wsNum, workspace.Layout)
	}

	for i, app := range workspace.Apps {
		if err := validateApp(wsNum, i, app); err != nil {
			return err
		}
	}

	return nil
}

func validateApp(wsNum, appIndex int, app App) error {
	if app.Name == "" {
		return fmt.Errorf("workspace %d, app %d: name is required", wsNum, appIndex)
	}

	return nil
}

func isValidLayout(layout string) bool {
	validLayouts := map[string]bool{
		"splith":   true,
		"splitv":   true,
		"stacking": true,
		"tabbed":   true,
	}

	return validLayouts[layout]
}
