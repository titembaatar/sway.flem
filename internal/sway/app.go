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
func LaunchApp(app config.Container, mark string) error {
	log.Info("Launching application: %s", app.App)

	cmdStr := app.Cmd
	if cmdStr == "" {
		cmdStr = app.App
	}

	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	var cmd *exec.Cmd
	if len(parts) == 1 {
		cmd = exec.Command(parts[0])
	} else {
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	log.Debug("Executing command: %s", cmdStr)
	if err := cmd.Start(); err != nil {
		log.Error("Failed to start application: %v", err)
		return fmt.Errorf("failed to start application: %w", err)
	}

	// Give the application some time to launch
	if app.Delay != 0 {
		time.Sleep(time.Duration(app.Delay) * time.Second)
	} else {
		time.Sleep(300 * time.Millisecond)
	}

	log.Debug("Applying mark to application: %s", mark)
	if err := ApplyMark(mark); err != nil {
		log.Error("Failed to apply mark to application: %v", err)
		return fmt.Errorf("failed to apply mark to application: %w", err)
	}

	// Execute post-launch commands if any
	if len(app.Post) > 0 {
		if err := RunPostCmd(app.Post); err != nil {
			log.Warn("Some post-launch commands failed: %v", err)
			// Continue execution even if post commands fail
		}
	}

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

		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			errors = append(errors, fmt.Sprintf("command %d: empty command", i+1))
			continue
		}

		var cmd *exec.Cmd
		if len(parts) == 1 {
			cmd = exec.Command(parts[0])
		} else {
			cmd = exec.Command(parts[0], parts[1:]...)
		}

		err := cmd.Start()
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

// Resizes a list of applications according to their stored information
func ResizeApps(appInfos []AppInfo) {
	log.Info("Resizing %d nodes", len(appInfos))

	for _, app := range appInfos {
		if app.Size == "" {
			log.Debug("Skipping resize for mark '%s' (no size specified)", app.Mark)
			continue
		}

		log.Debug("Resizing mark '%s' to '%s' with layout '%s'", app.Mark, app.Size, app.Layout)
		if err := ResizeMark(app.Mark, app.Size, app.Layout); err != nil {
			log.Warn("Failed to resize '%s': %v", app.Mark, err)
		}

		time.Sleep(200 * time.Millisecond)
	}
}

// Resizes a node (app or container) with the given mark
func ResizeMark(mark string, size string, layout string) error {
	var dimension string

	// Determine resize dimension based on layout
	switch layout {
	case "splith", "tabbed", "h", "t", "horizontal":
		dimension = "width"
	case "splitv", "stacking", "v", "s", "vertical":
		dimension = "height"
	default:
		// Default to width for unknown layouts
		dimension = "width"
		log.Warn("Unknown layout for resizing: %s, defaulting to width", layout)
	}

	// Use criteria to focus and resize
	focusCmd := fmt.Sprintf("[con_mark=\"%s\"] focus", mark)
	resizeCmd := fmt.Sprintf("resize set %s %s", dimension, size)
	_, errFocus := RunCommand(focusCmd)
	if errFocus != nil {
		return fmt.Errorf("failed to focus node: %w", errFocus)
	}

	time.Sleep(100 * time.Millisecond)

	_, errResize := RunCommand(resizeCmd)
	if errResize != nil {
		return fmt.Errorf("failed to resize node: %w", errResize)
	}

	time.Sleep(100 * time.Millisecond)

	return nil
}
