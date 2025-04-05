package sway

import (
	"fmt"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	errs "github.com/titembaatar/sway.flem/internal/errors"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/internal/util"
	"github.com/titembaatar/sway.flem/pkg/types"
)

type App struct {
	Name    string
	Command string
	Mark    Mark
	Size    string
	Delay   int64
	Layout  string
	Post    []string
}

func NewApp(container config.Container, markID string) *App {
	cmd := container.Cmd
	if cmd == "" {
		cmd = container.App
	}

	return &App{
		Name:    container.App,
		Command: cmd,
		Mark:    NewMark(markID),
		Size:    container.Size,
		Delay:   container.Delay,
		Layout:  string(container.Split),
		Post:    container.Post,
	}
}

func (a *App) Process(errorHandler *errs.ErrorHandler) error {
	log.Info("Processing application: %s with mark %s", a.Name, a.Mark.String())

	running, _, err := IsAppRunning(a.Mark.String())
	if err != nil {
		// Just log a warning and continue - checking is nice to have but not critical
		appWarn := errs.NewWarning(err, fmt.Sprintf("Failed to check if application '%s' is already running", a.Name))
		appWarn.WithCategory("Application")

		if errorHandler != nil {
			errorHandler.Handle(appWarn)
		} else {
			log.Warn("Failed to check if application is running: %v", err)
		}
	}

	if running {
		return a.Focus(errorHandler)
	} else {
		return a.Launch(errorHandler)
	}
}

func (a *App) Focus(errorHandler *errs.ErrorHandler) error {
	log.Info("Application with mark '%s' is already running, focusing instead of launching", a.Mark.String())

	if err := a.Mark.Focus(errorHandler); err != nil {
		focusErr := errs.Wrap(err, fmt.Sprintf("Failed to focus existing window with mark '%s'", a.Mark.String()))

		if errorHandler != nil {
			errorHandler.Handle(focusErr)
		} else {
			log.Error("Failed to focus existing window with mark '%s': %v", a.Mark.String(), err)
		}

		return focusErr
	}

	if err := a.RunPost(errorHandler); err != nil {
		// Non-fatal error - just log a warning
		postWarn := errs.NewWarning(err, fmt.Sprintf("Some post-focus commands failed for '%s'", a.Name))
		postWarn.WithCategory("Application")

		if errorHandler != nil {
			errorHandler.Handle(postWarn)
		} else {
			log.Warn("Some post-focus commands failed for '%s': %v", a.Name, err)
		}
	}

	log.Info("Successfully focused existing application '%s' with mark '%s'", a.Name, a.Mark.String())
	return nil
}

func (a *App) Launch(errorHandler *errs.ErrorHandler) error {
	log.Info("Launching new instance of application: %s", a.Name)

	if err := a.start(errorHandler); err != nil {
		return err
	}

	if err := a.RunPost(errorHandler); err != nil {
		// Non-fatal error - just log a warning
		postWarn := errs.NewWarning(err, fmt.Sprintf("Some post-launch commands failed for '%s'", a.Name))
		postWarn.WithCategory("Application")

		if errorHandler != nil {
			errorHandler.Handle(postWarn)
		} else {
			log.Warn("Some post-launch commands failed for '%s': %v", a.Name, err)
		}
	}

	log.Info("Successfully launched application '%s' with mark '%s'", a.Name, a.Mark.String())
	return nil
}

func (a *App) start(errorHandler *errs.ErrorHandler) error {
	if err := util.ExecuteCommand(a.Command); err != nil {
		launchErr := errs.NewAppLaunchError(a.Name, a.Command, err)

		if errorHandler != nil {
			errorHandler.Handle(launchErr)
		} else {
			log.Error("Failed to start application '%s' with command '%s': %v", a.Name, a.Command, err)
		}

		return launchErr
	}

	log.Debug("Application '%s' launched, waiting for it to initialize", a.Name)

	if a.Delay != 0 {
		time.Sleep(time.Duration(a.Delay) * time.Second)
	} else {
		time.Sleep(300 * time.Millisecond)
	}

	log.Debug("Applying mark '%s' to application", a.Mark.String())
	if err := a.Mark.Apply(errorHandler); err != nil {
		markErr := errs.Wrap(err, fmt.Sprintf("Failed to apply mark '%s' to application '%s'", a.Mark.String(), a.Name))

		if errorHandler != nil {
			errorHandler.Handle(markErr)
		} else {
			log.Error("Failed to apply mark '%s' to application '%s': %v", a.Mark.String(), a.Name, err)
		}

		return markErr
	}

	return nil
}

func (a *App) RunPost(errorHandler *errs.ErrorHandler) error {
	if len(a.Post) == 0 {
		return nil
	}

	log.Debug("Executing %d post commands for '%s'", len(a.Post), a.Name)
	return RunCommands(a.Post, errorHandler)
}

func (a *App) Resize(errorHandler *errs.ErrorHandler) error {
	if a.Size == "" {
		log.Debug("Skipping resize for mark '%s' (no size specified)", a.Mark.String())
		return nil
	}

	orientation := getOrientation(a.Layout)
	log.Debug("Resizing mark '%s' to '%s' with layout '%s'", a.Mark.String(), a.Size, a.Layout)

	if err := a.Mark.Focus(errorHandler); err != nil {
		resizeErr := errs.NewResizeError(a.Mark.String(), a.Size, orientation, a.Layout,
			errs.ErrFocusFailed)

		if errorHandler != nil {
			errorHandler.Handle(resizeErr)
		}

		return resizeErr
	}

	time.Sleep(100 * time.Millisecond)

	resizeCmd := NewSwayCmd(a.Mark.Resize(orientation, a.Size))
	if errorHandler != nil {
		resizeCmd.WithErrorHandler(errorHandler)
	}

	if _, err := resizeCmd.Run(); err != nil {
		resizeErr := errs.NewResizeError(a.Mark.String(), a.Size, orientation, a.Layout,
			errs.ErrResizeFailed)

		if errorHandler != nil {
			errorHandler.Handle(resizeErr)
		} else {
			log.Error("Failed to resize container with mark '%s' to %s %s: %v",
				a.Mark.String(), a.Size, orientation, err)
		}

		return resizeErr
	}

	time.Sleep(100 * time.Millisecond)

	log.Debug("Successfully resized container with mark '%s' to %s %s", a.Mark.String(), a.Size, orientation)

	return nil
}

func RunCommands(commands []string, errorHandler *errs.ErrorHandler) error {
	if len(commands) == 0 {
		return nil
	}

	log.Info("Executing %d commands", len(commands))
	var errors []string

	for i, cmdStr := range commands {
		log.Debug("Executing command %d: %s", i+1, cmdStr)

		err := executeCommand(cmdStr)
		if err != nil {
			cmdErr := errs.New(err, fmt.Sprintf("Failed to execute command: %s", cmdStr))
			cmdErr.WithCategory("Command")

			if errorHandler != nil {
				errorHandler.Handle(cmdErr)
			} else {
				log.Error("Failed to execute command %d: %v", i+1, err)
			}

			errors = append(errors, fmt.Sprintf("command %d: %v", i+1, err))
			continue
		}

		time.Sleep(400 * time.Millisecond)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to execute some commands: %s", strings.Join(errors, "; "))
	}

	log.Debug("All commands executed successfully")
	return nil
}

func executeCommand(cmdStr string) error {
	return util.ExecuteCommand(cmdStr)
}

func getOrientation(layout string) string {
	layoutType, err := types.ParseLayoutType(layout)
	if err != nil {
		log.Warn("Unknown layout for resizing: %s, defaulting to width", layout)
		return "width"
	}

	return layoutType.Orientation()
}
