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
	Name      string
	Command   string
	Mark      Mark
	Size      string
	Delay     int64
	Layout    string
	Post      []string
	RerunPost bool
}

func NewApp(container config.Container, markID string) *App {
	cmd := container.Cmd
	if cmd == "" {
		cmd = container.App
	}

	return &App{
		Name:      container.App,
		Command:   cmd,
		Mark:      NewMark(markID),
		Size:      container.Size,
		Delay:     container.Delay,
		Layout:    string(container.Split),
		Post:      container.Post,
		RerunPost: container.RerunPost,
	}
}

func (a *App) Process(errorHandler *errs.ErrorHandler) error {
	log.Info("Processing application: %s with mark %s", a.Name, a.Mark.String())

	running, _, err := IsAppRunning(a.Mark.String())
	if err != nil {
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
			log.Error("Failed to focus existing window: %v", err)
		}

		return focusErr
	}

	if a.RerunPost && len(a.Post) > 0 {
		log.Debug("RerunPost is true, executing post-focus commands for '%s'", a.Name)
		if err := a.RunPost(errorHandler); err != nil {
			postWarn := errs.NewWarning(err, fmt.Sprintf("Some post-focus commands failed for '%s'", a.Name))
			postWarn.WithCategory("Application")

			if errorHandler != nil {
				errorHandler.Handle(postWarn)
			} else {
				log.Warn("Some post-focus commands failed: %v", err)
			}
		}
	} else if len(a.Post) > 0 {
		log.Debug("Skipping post-focus commands for '%s' (RerunPost is false)", a.Name)
	}

	log.Info("Successfully focused existing application '%s'", a.Name)
	return nil
}

func (a *App) Launch(errorHandler *errs.ErrorHandler) error {
	log.Info("Launching new instance of application: %s", a.Name)

	if err := a.start(errorHandler); err != nil {
		return err
	}

	if err := a.RunPost(errorHandler); err != nil {
		postWarn := errs.NewWarning(err, fmt.Sprintf("Some post-launch commands failed for '%s'", a.Name))
		postWarn.WithCategory("Application")

		if errorHandler != nil {
			errorHandler.Handle(postWarn)
		} else {
			log.Warn("Some post-launch commands failed: %v", err)
		}
	}

	log.Info("Successfully launched application '%s'", a.Name)
	return nil
}

func (a *App) start(errorHandler *errs.ErrorHandler) error {
	if err := util.ExecuteCommand(a.Command); err != nil {
		launchErr := errs.NewAppLaunchError(a.Name, a.Command, err)

		if errorHandler != nil {
			errorHandler.Handle(launchErr)
		} else {
			log.Error("Failed to start application '%s': %v", a.Name, err)
		}

		return launchErr
	}

	log.Debug("Application '%s' launched, waiting for it to initialize", a.Name)

	if a.Delay != 0 {
		time.Sleep(time.Duration(a.Delay) * time.Second)
	} else {
		time.Sleep(300 * time.Millisecond)
	}

	log.Debug("Applying mark to application")
	if err := a.Mark.Apply(errorHandler); err != nil {
		markErr := errs.Wrap(err, fmt.Sprintf("Failed to apply mark to application '%s'", a.Name))

		if errorHandler != nil {
			errorHandler.Handle(markErr)
		} else {
			log.Error("Failed to apply mark: %v", err)
		}

		return markErr
	}

	return nil
}

func (a *App) RunPost(errorHandler *errs.ErrorHandler) error {
	postCmdCount := len(a.Post)
	if postCmdCount == 0 {
		return nil
	}

	log.Debug("Executing %d post commands for '%s'", postCmdCount, a.Name)
	return RunCommands(a.Post, errorHandler)
}

func (a *App) Resize(errorHandler *errs.ErrorHandler) error {
	if a.Size == "" {
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
			log.Error("Failed to resize container: %v", err)
		}

		return resizeErr
	}

	time.Sleep(100 * time.Millisecond)
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
