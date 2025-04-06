package errors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrCommandFailed         = errors.New("sway command failed")
	ErrAppLaunchFailed       = errors.New("failed to launch application")
	ErrMarkingFailed         = errors.New("failed to apply mark")
	ErrResizeFailed          = errors.New("failed to resize container")
	ErrFocusFailed           = errors.New("failed to focus container")
	ErrSetLayoutFailed       = errors.New("failed to set container layout")
	ErrWorkspaceCreateFailed = errors.New("failed to create workspace")
	ErrSwayNotAvailable      = errors.New("sway command not available")
)

type SwayCommandError struct {
	*AppError
	Command string
	Output  string
}

func NewSwayCommandError(command string, err error, output string) *SwayCommandError {
	message := fmt.Sprintf("Sway command '%s' failed", command)

	appErr := New(err, message).WithCategory("Sway")

	if strings.Contains(output, "Unable to connect") {
		appErr.WithSuggestion("Make sure Sway is running and the SWAYSOCK environment variable is set correctly")
	}

	return &SwayCommandError{
		AppError: appErr,
		Command:  command,
		Output:   output,
	}
}

func (e *SwayCommandError) WithOutput(output string) *SwayCommandError {
	e.Output = output
	return e
}

type AppLaunchError struct {
	*AppError
	AppName string
	Command string
}

func NewAppLaunchError(appName, command string, err error) *AppLaunchError {
	message := fmt.Sprintf("Failed to launch application '%s'", appName)

	appErr := New(err, message).WithCategory("Application")
	appErr.WithSuggestion(fmt.Sprintf("Check that '%s' is installed and available in your PATH", command))

	return &AppLaunchError{
		AppError: appErr,
		AppName:  appName,
		Command:  command,
	}
}

type MarkError struct {
	*AppError
	Mark string
}

func NewMarkError(mark string, err error) *MarkError {
	message := fmt.Sprintf("Failed to apply mark '%s'", mark)

	appErr := New(err, message).WithCategory("Sway")

	return &MarkError{
		AppError: appErr,
		Mark:     mark,
	}
}

type ResizeError struct {
	*AppError
	Mark      string
	Size      string
	Dimension string
	Layout    string
}

func NewResizeError(mark, size, dimension, layout string, err error) *ResizeError {
	message := fmt.Sprintf("Failed to resize container '%s' to %s %s", mark, size, dimension)

	appErr := New(err, message).WithCategory("Sway")

	suggestion := fmt.Sprintf("Check that the container with mark '%s' exists and the size '%s' is valid", mark, size)
	appErr.WithSuggestion(suggestion)

	return &ResizeError{
		AppError:  appErr,
		Mark:      mark,
		Size:      size,
		Dimension: dimension,
		Layout:    layout,
	}
}
