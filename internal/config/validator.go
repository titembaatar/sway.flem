package config

import (
	"fmt"

	"github.com/titembaatar/sway.flem/internal/log"
)

var validLayoutTypes = map[string]string{
	// layout names
	"splith":   "splith",
	"splitv":   "splitv",
	"stacking": "stacking",
	"tabbed":   "tabbed",

	// aliases
	"horizontal": "splith",
	"h":          "splith",
	"vertical":   "splitv",
	"v":          "splitv",
	"stack":      "stacking",
	"s":          "stacking",
	"tab":        "tabbed",
	"t":          "tabbed",
}

func NormalizeLayoutType(layoutType string) (string, error) {
	if normalized, ok := validLayoutTypes[layoutType]; ok {
		return normalized, nil
	}
	return "", ErrInvalidLayoutType
}

func ValidateConfig(config *Config) error {
	if len(config.Workspaces) == 0 {
		return NewConfigError(ErrNoWorkspaces, "", "", -1)
	}

	log.Debug("Validating configuration with %d workspaces", len(config.Workspaces))

	for name, workspace := range config.Workspaces {
		normalizedLayout, err := NormalizeLayoutType(workspace.Layout)
		if err != nil {
			return NewConfigError(err, name, "layout", -1)
		}
		workspace.Layout = normalizedLayout
		config.Workspaces[name] = workspace

		if err := validateWorkspace(name, workspace); err != nil {
			return err
		}

		log.Debug("Workspace '%s' validated successfully", name)
	}

	log.Info("Configuration validated successfully")
	return nil
}

func validateWorkspace(name string, workspace Workspace) error {
	if workspace.Layout == "" {
		return NewConfigError(ErrMissingLayout, name, "", -1)
	}

	if len(workspace.Containers) == 0 {
		return NewConfigError(ErrNoContainers, name, "", -1)
	}

	// Validate each container in the workspace
	for i, container := range workspace.Containers {
		if err := validateContainer(name, container, fmt.Sprintf("container[%d]", i)); err != nil {
			return err
		}
	}

	return nil
}

func validateContainer(workspaceName string, container Container, context string) error {
	// Check if this is an app container or a nested container
	isApp := container.App != ""
	isNestedContainer := len(container.Containers) > 0

	// A container must be either an app or have nested containers, not both
	if isApp && isNestedContainer {
		return NewConfigError(ErrInvalidContainerStructure, workspaceName, context, -1)
	}

	// A container must be either an app or have nested containers, not neither
	if !isApp && !isNestedContainer {
		return NewConfigError(ErrInvalidContainerStructure, workspaceName, context, -1)
	}

	// If this is a nested container, validate split and nested containers
	if isNestedContainer {
		if container.Split == "" {
			return NewConfigError(ErrMissingSplit, workspaceName, context, -1)
		}

		// Normalize split
		normalizedSplit, err := NormalizeLayoutType(container.Split)
		if err != nil {
			return NewConfigError(err, workspaceName, fmt.Sprintf("%s.split", context), -1)
		}
		container.Split = normalizedSplit

		// Validate nested containers
		for i, nestedContainer := range container.Containers {
			nestedContext := fmt.Sprintf("%s.containers[%d]", context, i)
			if err := validateContainer(workspaceName, nestedContainer, nestedContext); err != nil {
				return err
			}
		}
	}

	return nil
}
