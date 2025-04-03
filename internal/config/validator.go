package config

import (
	"fmt"
	"regexp"
	"strconv"

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

// Regular expression for valid size formats:
// - Digits only (e.g., "50") - defaults to ppt in Sway
// - Digits followed by "ppt" (e.g., "50ppt")
// - Digits followed by "px" (e.g., "800px")
var sizeRegex = regexp.MustCompile(`^(\d+)(ppt|px)?$`)

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
		if err := ValidateSize(container.Size); err != nil {
			return NewConfigError(err, workspaceName, fmt.Sprintf("%s.size", context), -1)
		}
	}

	return nil
}

func validateNestedContainer(workspaceName string, container Container, context string) error {
	if container.Split == "" {
		return NewConfigError(ErrMissingSplit, workspaceName, context, -1)
	}

	normalizedSplit, err := NormalizeLayoutType(container.Split)
	if err != nil {
		return NewConfigError(err, workspaceName, fmt.Sprintf("%s.split", context), -1)
	}
	container.Split = normalizedSplit

	for i, nestedContainer := range container.Containers {
		nestedContext := fmt.Sprintf("%s.containers[%d]", context, i)
		if err := validateContainer(workspaceName, nestedContainer, nestedContext); err != nil {
			return err
		}
	}

	return nil
}

func NormalizeLayoutType(layoutType string) (string, error) {
	if normalized, ok := validLayoutTypes[layoutType]; ok {
		return normalized, nil
	}
	return "", ErrInvalidLayoutType
}

func ValidateSize(size string) error {
	if size == "" {
		return nil
	}

	if !sizeRegex.MatchString(size) {
		return ErrInvalidSizeFormat
	}

	matches := sizeRegex.FindStringSubmatch(size)
	if len(matches) < 2 {
		return ErrInvalidSizeFormat
	}

	numStr := matches[1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return ErrInvalidSizeFormat
	}

	if num < 0 {
		return ErrInvalidSizeFormat
	}

	return nil
}
