package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNoWorkspaces              = errors.New("no workspaces defined in configuration")
	ErrInvalidLayoutType         = errors.New("invalid layout type")
	ErrMissingLayout             = errors.New("missing layout for container")
	ErrNoContainers              = errors.New("workspace has no containers defined")
	ErrMissingSplit              = errors.New("nested container requires a split property")
	ErrInvalidContainerStructure = errors.New("invalid container structure: must be either an app or have nested containers")
	ErrInvalidSizeFormat         = errors.New("invalid size format: must be a number, optionally followed by 'ppt' or 'px' (e.g., '50', '50ppt', '800px')")
)

type ConfigError struct {
	*AppError
	Workspace   string
	Context     string
	ContainerID int
}

func NewConfigError(err error, workspace, context string, containerID int) *ConfigError {
	message := buildConfigErrorMessage(workspace, context, containerID)

	appErr := New(err, message).WithCategory("Config")

	return &ConfigError{
		AppError:    appErr,
		Workspace:   workspace,
		Context:     context,
		ContainerID: containerID,
	}
}

func NewFatalConfigError(err error, workspace, context string, containerID int) *ConfigError {
	message := buildConfigErrorMessage(workspace, context, containerID)

	appErr := NewFatal(err, message).WithCategory("Config")

	return &ConfigError{
		AppError:    appErr,
		Workspace:   workspace,
		Context:     context,
		ContainerID: containerID,
	}
}

func (e *ConfigError) WithConfigFile(file string, line int) *ConfigError {
	e.AppError.WithFile(file, line, 0)
	return e
}

func (e *ConfigError) WithSuggestion(suggestion string) *ConfigError {
	e.AppError.WithSuggestion(suggestion)
	return e
}

func buildConfigErrorMessage(workspace, context string, containerID int) string {
	var locationInfo string

	if workspace != "" {
		locationInfo += fmt.Sprintf("workspace '%s'", workspace)

		if context != "" {
			locationInfo += fmt.Sprintf(", %s", context)
		}

		if containerID >= 0 {
			locationInfo += fmt.Sprintf(", container %d", containerID)
		}
	}

	if locationInfo != "" {
		return fmt.Sprintf("Configuration error in %s", locationInfo)
	}

	return "Configuration error"
}

func GetLayoutSuggestion() string {
	return "Valid layouts are: 'splith'/'h'/'horizontal', 'splitv'/'v'/'vertical', 'tabbed'/'t', 'stacking'/'s'"
}

func GetSizeSuggestion() string {
	return "Size must be specified as a number (e.g., '50') or with units (e.g., '50ppt', '800px')"
}

func GetContainerStructureSuggestion() string {
	return "A container must either have an 'app' property or a 'containers' list with a 'split' layout, but not both"
}
