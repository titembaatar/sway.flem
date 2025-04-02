package config

import (
	"errors"
	"fmt"
)

var (
	ErrNoWorkspaces              = errors.New("no workspaces defined in configuration")
	ErrInvalidLayoutType         = errors.New("invalid layout type")
	ErrMissingLayout             = errors.New("no layout defined")
	ErrNoContainers              = errors.New("workspace has no containers defined")
	ErrMissingSplit              = errors.New("nested container has no split defined")
	ErrInvalidContainerStructure = errors.New("invalid container structure: must be either an app or have nested containers")
)

type ConfigError struct {
	Err       error
	Workspace string
	Context   string
	Index     int
}

func (e *ConfigError) Error() string {
	if e.Index >= 0 {
		return fmt.Sprintf("%s: workspace '%s', %s at index %d: %v",
			"Configuration error", e.Workspace, e.Context, e.Index, e.Err)
	}

	if e.Context != "" {
		return fmt.Sprintf("%s: workspace '%s', %s: %v",
			"Configuration error", e.Workspace, e.Context, e.Err)
	}

	if e.Workspace != "" {
		return fmt.Sprintf("%s: workspace '%s': %v",
			"Configuration error", e.Workspace, e.Err)
	}

	return fmt.Sprintf("%s: %v", "Configuration error", e.Err)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

func NewConfigError(err error, workspace string, context string, index int) *ConfigError {
	return &ConfigError{
		Err:       err,
		Workspace: workspace,
		Context:   context,
		Index:     index,
	}
}
