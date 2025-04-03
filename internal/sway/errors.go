package sway

import (
	"errors"
	"fmt"
)

var (
	ErrCommandFailed         = errors.New("sway command failed")
	ErrAppLaunchFailed       = errors.New("failed to launch application")
	ErrMarkingFailed         = errors.New("failed to apply mark")
	ErrResizeFailed          = errors.New("failed to resize container")
	ErrFocusFailed           = errors.New("failed to focus container")
	ErrSetLayoutFailed       = errors.New("failed to set container layout")
	ErrInvalidLayout         = errors.New("invalid layout type")
	ErrWorkspaceCreateFailed = errors.New("failed to create workspace")
)

type SwayCommandError struct {
	Command string
	Err     error
	Output  string // Output from stderr
}

func (e *SwayCommandError) Error() string {
	if e.Output != "" {
		return fmt.Sprintf("sway command '%s' failed: %v (output: %s)", e.Command, e.Err, e.Output)
	}
	return fmt.Sprintf("sway command '%s' failed: %v", e.Command, e.Err)
}

func (e *SwayCommandError) Unwrap() error {
	return e.Err
}

func NewSwayCommandError(command string, err error, output string) *SwayCommandError {
	return &SwayCommandError{
		Command: command,
		Err:     err,
		Output:  output,
	}
}

type AppLaunchError struct {
	AppName string
	Command string
	Err     error
}

func (e *AppLaunchError) Error() string {
	return fmt.Sprintf("failed to launch app '%s' with command '%s': %v", e.AppName, e.Command, e.Err)
}

func (e *AppLaunchError) Unwrap() error {
	return e.Err
}

func NewAppLaunchError(appName, command string, err error) *AppLaunchError {
	return &AppLaunchError{
		AppName: appName,
		Command: command,
		Err:     err,
	}
}

type MarkError struct {
	Mark string
	Err  error
}

func (e *MarkError) Error() string {
	return fmt.Sprintf("failed to apply mark '%s': %v", e.Mark, e.Err)
}

func (e *MarkError) Unwrap() error {
	return e.Err
}

func NewMarkError(mark string, err error) *MarkError {
	return &MarkError{
		Mark: mark,
		Err:  err,
	}
}

type ResizeError struct {
	Mark      string
	Size      string
	Dimension string
	Layout    string
	Err       error
}

func (e *ResizeError) Error() string {
	return fmt.Sprintf("failed to resize container with mark '%s' to %s %s (layout: %s): %v",
		e.Mark, e.Dimension, e.Size, e.Layout, e.Err)
}

func (e *ResizeError) Unwrap() error {
	return e.Err
}

func NewResizeError(mark, size, dimension, layout string, err error) *ResizeError {
	return &ResizeError{
		Mark:      mark,
		Size:      size,
		Dimension: dimension,
		Layout:    layout,
		Err:       err,
	}
}
