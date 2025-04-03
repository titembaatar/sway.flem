package sway

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
)

type AppInfo struct {
	Mark   string
	Size   string
	Layout string
}

// Launches an application and marks it
func LaunchApp(app config.Container, markID string) error {
	log.Info("Launching application: %s", app.App)
	mark := NewMark(markID)

	cmdStr := app.Cmd
	if cmdStr == "" {
		cmdStr = app.App
	}

	// Execute the command to launch the app
	if err := executeCommand(cmdStr); err != nil {
		log.Error("Failed to start application '%s' with command '%s': %v", app.App, cmdStr, err)
		return NewAppLaunchError(app.App, cmdStr, err)
	}

	log.Debug("Application '%s' launched, waiting for it to initialize", app.App)

	// Give the application some time to launch
	if app.Delay != 0 {
		time.Sleep(time.Duration(app.Delay) * time.Second)
	} else {
		time.Sleep(300 * time.Millisecond)
	}

	// Apply mark to the application
	log.Debug("Applying mark '%s' to application", mark.String())
	if err := mark.Apply(); err != nil {
		log.Error("Failed to apply mark '%s' to application '%s': %v", mark.String(), app.App, err)
		return NewMarkError(mark.String(), fmt.Errorf("%w: %v", ErrMarkingFailed, err))
	}

	// Execute post-launch commands if any
	if len(app.Post) > 0 {
		log.Debug("Executing %d post-launch commands for '%s'", len(app.Post), app.App)
		if err := RunPostCmd(app.Post); err != nil {
			log.Warn("Some post-launch commands failed for '%s': %v", app.App, err)
			// Continue execution even if post commands fail
		}
	}

	log.Info("Successfully launched application '%s' with mark '%s'", app.App, mark.String())
	return nil
}

// Executes post-launch commands
func RunPostCmd(commands []string) error {
	if len(commands) == 0 {
		return nil
	}

	log.Info("Executing %d post-launch commands", len(commands))
	var errors []string

	for i, cmdStr := range commands {
		log.Debug("Executing post-launch command %d: %s", i+1, cmdStr)

		err := executeCommand(cmdStr)
		if err != nil {
			log.Error("Failed to execute post-launch command %d: %v", i+1, err)
			errors = append(errors, fmt.Sprintf("command %d: %v", i+1, err))
			continue
		}

		time.Sleep(200 * time.Millisecond)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to execute some post-launch commands: %s", strings.Join(errors, "; "))
	}

	log.Info("All post-launch commands executed successfully")
	return nil
}

// Parses and executes a command string
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

// Checks if a command exists in the PATH
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

// Resizes a list of applications according to their stored information
func ResizeApps(appInfos []AppInfo) {
	log.Info("Resizing %d nodes", len(appInfos))

	for _, app := range appInfos {
		if app.Size == "" {
			log.Debug("Skipping resize for mark '%s' (no size specified)", app.Mark)
			continue
		}

		mark := NewMark(app.Mark)
		log.Debug("Resizing mark '%s' to '%s' with layout '%s'", mark.String(), app.Size, app.Layout)
		if err := ResizeMark(app.Mark, app.Size, app.Layout); err != nil {
			log.Warn("Failed to resize '%s': %v", app.Mark, err)
		}

		time.Sleep(200 * time.Millisecond)
	}
}

// Resizes a node (app or container) with the given mark
func ResizeMark(markID string, size string, layout string) error {
	mark := NewMark(markID)
	dimension := getDimensionForLayout(layout)

	// Focus the container first
	if err := mark.Focus(); err != nil {
		return NewResizeError(markID, size, dimension, layout,
			fmt.Errorf("%w: failed to focus container before resizing", ErrFocusFailed))
	}

	time.Sleep(100 * time.Millisecond)

	// Then resize the focused container
	resizeCmd := mark.ResizeCmd(dimension, size)
	if _, err := RunCommand(resizeCmd); err != nil {
		log.Error("Failed to resize container with mark '%s' to %s %s: %v",
			markID, size, dimension, err)
		return NewResizeError(markID, size, dimension, layout,
			fmt.Errorf("%w: command failed", ErrResizeFailed))
	}

	time.Sleep(100 * time.Millisecond)
	log.Debug("Successfully resized container with mark '%s' to %s %s", markID, size, dimension)

	return nil
}

// Determines the resize dimension based on layout type
func getDimensionForLayout(layout string) string {
	switch layout {
	case "splith", "tabbed", "h", "t", "horizontal":
		return "width"
	case "splitv", "stacking", "v", "s", "vertical":
		return "height"
	default:
		log.Warn("Unknown layout for resizing: %s, defaulting to width", layout)
		return "width"
	}
}
