package config

import (
	"fmt"

	errs "github.com/titembaatar/sway.flem/internal/errors"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/pkg/types"
)

func ValidateConfig(config *Config) error {
	if len(config.Workspaces) == 0 {
		return errs.NewFatalConfigError(errs.ErrNoWorkspaces, "", "", -1).
			WithSuggestion("Add at least one workspace to your configuration file")
	}

	log.Debug("Validating configuration with %d workspaces", len(config.Workspaces))

	for name, workspace := range config.Workspaces {
		layoutStr := string(workspace.Layout)
		layout, err := types.ParseLayoutType(layoutStr)
		if err != nil {
			configErr := errs.NewConfigError(err, name, "layout", -1).
				WithSuggestion(errs.GetLayoutSuggestion())
			return configErr
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
		return errs.NewConfigError(errs.ErrInvalidLayoutType, name, "", -1).
			WithSuggestion(errs.GetLayoutSuggestion())
	}

	if len(workspace.Containers) == 0 {
		return errs.NewConfigError(errs.ErrNoContainers, name, "", -1).
			WithSuggestion("Add at least one container to the workspace")
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
		return errs.NewConfigError(errs.ErrInvalidContainerStructure, workspaceName, context, -1).
			WithSuggestion(errs.GetContainerStructureSuggestion())
	}

	if !isApp && !isNestedContainer {
		return errs.NewConfigError(errs.ErrInvalidContainerStructure, workspaceName, context, -1).
			WithSuggestion("Specify either an 'app' property or a list of nested 'containers'")
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
			return errs.NewConfigError(err, workspaceName, fmt.Sprintf("%s.size", context), -1).
				WithSuggestion(errs.GetSizeSuggestion())
		}

		if !size.IsValid() {
			return errs.NewConfigError(errs.ErrInvalidSizeFormat, workspaceName, fmt.Sprintf("%s.size", context), -1).
				WithSuggestion(errs.GetSizeSuggestion())
		}
	}

	return nil
}

func validateNestedContainer(workspaceName string, container Container, context string) error {
	if container.Split == "" {
		return errs.NewConfigError(errs.ErrMissingSplit, workspaceName, context, -1).
			WithSuggestion("Add a 'split' property to specify the layout for nested containers")
	}

	splitStr := string(container.Split)
	layout, err := types.ParseLayoutType(splitStr)
	if err != nil {
		return errs.NewConfigError(err, workspaceName, fmt.Sprintf("%s.split", context), -1).
			WithSuggestion(errs.GetLayoutSuggestion())
	}

	if !layout.IsValid() {
		return errs.NewConfigError(errs.ErrInvalidLayoutType, workspaceName, fmt.Sprintf("%s.split", context), -1).
			WithSuggestion(errs.GetLayoutSuggestion())
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
