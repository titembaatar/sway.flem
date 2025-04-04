package sway

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
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

func (a *App) Process() error {
	log.Info("Processing application: %s with mark %s", a.Name, a.Mark.String())

	running, _, err := IsAppRunning(a.Mark.String())
	if err != nil {
		log.Warn("Failed to check if application is running: %v", err)
	}

	if running {
		return a.Focus()
	} else {
		return a.Launch()
	}
}

func (a *App) Focus() error {
	log.Info("Application with mark '%s' is already running, focusing instead of launching", a.Mark.String())

	if err := a.Mark.Focus(); err != nil {
		log.Error("Failed to focus existing window with mark '%s': %v", a.Mark.String(), err)
		return err
	}

	if err := a.RunPost(); err != nil {
		log.Warn("Some post-focus commands failed for '%s': %v", a.Name, err)
	}

	log.Info("Successfully focused existing application '%s' with mark '%s'", a.Name, a.Mark.String())
	return nil
}

func (a *App) Launch() error {
	log.Info("Launching new instance of application: %s", a.Name)

	if err := a.start(); err != nil {
		return err
	}

	if err := a.RunPost(); err != nil {
		log.Warn("Some post-launch commands failed for '%s': %v", a.Name, err)
	}

	log.Info("Successfully launched application '%s' with mark '%s'", a.Name, a.Mark.String())
	return nil
}

func (a *App) start() error {
	if err := executeCommand(a.Command); err != nil {
		log.Error("Failed to start application '%s' with command '%s': %v", a.Name, a.Command, err)
		return NewAppLaunchError(a.Name, a.Command, err)
	}

	log.Debug("Application '%s' launched, waiting for it to initialize", a.Name)

	if a.Delay != 0 {
		time.Sleep(time.Duration(a.Delay) * time.Second)
	} else {
		time.Sleep(300 * time.Millisecond)
	}

	log.Debug("Applying mark '%s' to application", a.Mark.String())
	if err := a.Mark.Apply(); err != nil {
		log.Error("Failed to apply mark '%s' to application '%s': %v", a.Mark.String(), a.Name, err)
		return NewMarkError(a.Mark.String(), fmt.Errorf("%w: %v", ErrMarkingFailed, err))
	}

	return nil
}

func (a *App) RunPost() error {
	if len(a.Post) == 0 {
		return nil
	}

	log.Debug("Executing %d post commands for '%s'", len(a.Post), a.Name)
	return RunCommands(a.Post)
}

func (a *App) Resize() error {
	if a.Size == "" {
		log.Debug("Skipping resize for mark '%s' (no size specified)", a.Mark.String())

		return nil
	}

	orientation := getOrientation(a.Layout)
	log.Debug("Resizing mark '%s' to '%s' with layout '%s'", a.Mark.String(), a.Size, a.Layout)

	if err := a.Mark.Focus(); err != nil {
		return NewResizeError(a.Mark.String(), a.Size, orientation, a.Layout,
			fmt.Errorf("%w: failed to focus container before resizing", ErrFocusFailed))
	}

	time.Sleep(100 * time.Millisecond)

	resizeCmd := NewSwayCmd(a.Mark.Resize(orientation, a.Size))

	if _, err := resizeCmd.Run(); err != nil {
		log.Error("Failed to resize container with mark '%s' to %s %s: %v",
			a.Mark.String(), a.Size, orientation, err)

		return NewResizeError(a.Mark.String(), a.Size, orientation, a.Layout,
			fmt.Errorf("%w: command failed", ErrResizeFailed))
	}

	time.Sleep(100 * time.Millisecond)

	log.Debug("Successfully resized container with mark '%s' to %s %s", a.Mark.String(), a.Size, orientation)

	return nil
}

func RunCommands(commands []string) error {
	if len(commands) == 0 {
		return nil
	}

	log.Info("Executing %d commands", len(commands))
	var errors []string

	for i, cmdStr := range commands {
		log.Debug("Executing command %d: %s", i+1, cmdStr)

		err := executeCommand(cmdStr)
		if err != nil {
			log.Error("Failed to execute command %d: %v", i+1, err)
			errors = append(errors, fmt.Sprintf("command %d: %v", i+1, err))
			continue
		}

		time.Sleep(200 * time.Millisecond)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to execute some commands: %s", strings.Join(errors, "; "))
	}

	log.Debug("All commands executed successfully")
	return nil
}

func executeCommand(cmdStr string) error {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	if err := validateCommand(parts[0]); err != nil {
		return err
	}

	var cmd *exec.Cmd
	if len(parts) == 1 {
		cmd = exec.Command(parts[0])
	} else {
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	log.Debug("Executing command: %s", cmdStr)
	return cmd.Start()
}

func validateCommand(command string) error {
	if strings.ContainsAny(command, "/\\") {
		return nil
	}

	_, err := exec.LookPath(command)
	if err != nil {
		log.Error("Command '%s' not found in PATH: %v", command, err)
		return fmt.Errorf("command '%s' not found in PATH: %w", command, err)
	}

	return nil
}

func getOrientation(layout string) string {
	layoutType, err := types.ParseLayoutType(layout)
	if err != nil {
		log.Warn("Unknown layout for resizing: %s, defaulting to width", layout)
		return "width"
	}

	return layoutType.Orientation()
}
