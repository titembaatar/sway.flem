package sway

import (
	"fmt"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	errs "github.com/titembaatar/sway.flem/internal/errors"
	"github.com/titembaatar/sway.flem/internal/log"
)

type Workspace struct {
	Name         string
	Layout       string
	Container    *Container
	ErrorHandler *errs.ErrorHandler
}

func NewWorkspace(name string, layout string) *Workspace {
	return &Workspace{
		Name:      name,
		Layout:    layout,
		Container: nil,
	}
}

// WithErrorHandler adds an error handler to the workspace
func (w *Workspace) WithErrorHandler(handler *errs.ErrorHandler) *Workspace {
	w.ErrorHandler = handler
	return w
}

func (w *Workspace) Setup(workspace config.Workspace, errorHandler *errs.ErrorHandler) error {
	w.ErrorHandler = errorHandler
	log.Info("Setting up workspace: %s", w.Name)

	if err := w.Create(); err != nil {
		return errs.Wrap(err, fmt.Sprintf("Failed to create workspace '%s'", w.Name))
	}

	container, err := ProcessContainer(w.Name, workspace.Containers, w.Layout, 0, errorHandler)
	if err != nil {
		log.Error("Failed to process container for workspace %s: %v", w.Name, err)
		return err
	}

	w.Container = container

	if err := container.Setup(errorHandler); err != nil {
		log.Error("Failed to setup container for workspace %s: %v", w.Name, err)
		return err
	}

	container.ResizeApps(errorHandler)

	log.Info("Workspace %s setup complete", w.Name)
	return nil
}

func (w *Workspace) Switch() error {
	command := fmt.Sprintf("workspace %s", w.Name)
	cmd := NewSwayCmd(command)
	if w.ErrorHandler != nil {
		cmd.WithErrorHandler(w.ErrorHandler)
	}

	_, err := cmd.Run()
	if err != nil {
		return errs.Wrap(err, fmt.Sprintf("Failed to switch to workspace '%s'", w.Name))
	}
	return nil
}

func SwitchWorkspace(workspace string, errorHandler *errs.ErrorHandler) error {
	ws := NewWorkspace(workspace, "").WithErrorHandler(errorHandler)
	return ws.Switch()
}

func (w *Workspace) Create() error {
	if err := w.Switch(); err != nil {
		createErr := errs.New(errs.ErrWorkspaceCreateFailed,
			fmt.Sprintf("Failed to create workspace '%s'", w.Name))
		createErr.WithCategory("Sway")
		createErr.WithSuggestion("Check that Sway is running and accepting commands")

		if w.ErrorHandler != nil {
			w.ErrorHandler.Handle(createErr)
		}

		return createErr
	}

	if w.Layout == "" {
		return nil // No layout to set
	}

	cmd := fmt.Sprintf("layout %s", w.Layout)
	swayCmd := NewSwayCmd(cmd)
	if w.ErrorHandler != nil {
		swayCmd.WithErrorHandler(w.ErrorHandler)
	}

	_, err := swayCmd.Run()
	if err != nil {
		layoutErr := errs.New(errs.ErrSetLayoutFailed,
			fmt.Sprintf("Failed to set layout '%s' for workspace '%s'", w.Layout, w.Name))
		layoutErr.WithCategory("Sway")

		if w.ErrorHandler != nil {
			w.ErrorHandler.Handle(layoutErr)
		}

		return layoutErr
	}

	log.Info("Successfully created workspace '%s' with layout '%s'", w.Name, w.Layout)
	return nil
}

func CreateWorkspace(name string, layout string, errorHandler *errs.ErrorHandler) error {
	ws := NewWorkspace(name, layout).WithErrorHandler(errorHandler)
	return ws.Create()
}

func FocusWorkspaces(workspaces []string) error {
	log.Info("Focusing on %d workspaces", len(workspaces))
	var errors []string

	for i, workspace := range workspaces {
		ws := NewWorkspace(workspace, "")
		log.Debug("Focusing on workspace: %s", workspace)

		if err := ws.Switch(); err != nil {
			log.Error("Failed to focus on workspace %s: %v", workspace, err)
			errors = append(errors, fmt.Sprintf("workspace %s: %v", workspace, err))
		} else {
			log.Info("Successfully focused on workspace %s (%d of %d)", workspace, i+1, len(workspaces))
			time.Sleep(100 * time.Millisecond)
		}
	}

	if len(errors) > 0 {
		focusErr := errs.New(errs.ErrFocusFailed, "Failed to focus on some workspaces")
		focusErr.WithCategory("Sway")
		focusErr.WithSuggestion("Check that the specified workspaces exist")
		return focusErr
	}

	return nil
}
