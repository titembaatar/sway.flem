package config

import (
	"fmt"

	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/pkg/types"
)

func ValidateConfig(config *Config) error {
	if len(config.Workspaces) == 0 {
		return NewConfigError(ErrNoWorkspaces, "", "", -1)
	}

	log.Debug("Validating configuration with %d workspaces", len(config.Workspaces))

	for name, workspace := range config.Workspaces {
		layoutStr := string(workspace.Layout)
		layout, err := types.ParseLayoutType(layoutStr)
		if err != nil {
			return NewConfigError(err, name, "layout", -1)
		}

		workspace.Layout = layout
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
	if !workspace.Layout.IsValid() {
		return NewConfigError(types.ErrInvalidLayoutType, name, "", -1)
	}

	if len(workspace.Containers) == 0 {
		return NewConfigError(ErrNoContainers, name, "", -1)
	}

	for i, container := range workspace.Containers {
		if err := validateContainer(name, container, fmt.Sprintf("container[%d]", i)); err != nil {
			return err
		}
	}

	return nil
}

func validateContainer(workspaceName string, container Container, context string) error {
	isApp := container.App != ""
	isNestedContainer := len(container.Containers) > 0

	if isApp && isNestedContainer {
		return NewConfigError(ErrInvalidContainerStructure, workspaceName, context, -1)
	}

	if !isApp && !isNestedContainer {
		return NewConfigError(ErrInvalidContainerStructure, workspaceName, context, -1)
	}

	if err := validateContainerProperties(workspaceName, container, context); err != nil {
		return err
	}

	if isNestedContainer {
		if err := validateNestedContainer(workspaceName, container, context); err != nil {
			return err
		}
	}

	return nil
}

func validateContainerProperties(workspaceName string, container Container, context string) error {
	if container.Size != "" {
		size, err := types.ParseSize(container.Size)
		if err != nil {
			return NewConfigError(err, workspaceName, fmt.Sprintf("%s.size", context), -1)
		}

		if !size.IsValid() {
			return NewConfigError(types.ErrInvalidSizeFormat, workspaceName, fmt.Sprintf("%s.size", context), -1)
		}
	}

	return nil
}

func validateNestedContainer(workspaceName string, container Container, context string) error {
	if container.Split == "" {
		return NewConfigError(ErrMissingSplit, workspaceName, context, -1)
	}

	splitStr := string(container.Split)
	layout, err := types.ParseLayoutType(splitStr)
	if err != nil {
		return NewConfigError(err, workspaceName, fmt.Sprintf("%s.split", context), -1)
	}

	if !layout.IsValid() {
		return NewConfigError(types.ErrInvalidLayoutType, workspaceName, fmt.Sprintf("%s.split", context), -1)
	}

	// Note: This doesn't actually modify the original container since we're working on a copy
	container.Split = layout

	for i, nestedContainer := range container.Containers {
		nestedContext := fmt.Sprintf("%s.containers[%d]", context, i)
		if err := validateContainer(workspaceName, nestedContainer, nestedContext); err != nil {
			return err
		}
	}

	return nil
}
