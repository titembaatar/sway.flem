package sway

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
)

// App information for resizing operations
type AppInfo struct {
	Mark   string // Mark identifying this node
	Size   string // Size to set
	Layout string // Parent layout (determines resize dimension)
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

	return nil
}

// Resizes a list of applications according to their stored information
func ResizeApps(appInfos []AppInfo) {
	log.Info("Resizing %d nodes", len(appInfos))

	for _, app := range appInfos {
		if app.Size != "" {
			log.Debug("Resizing mark '%s' to '%s' with layout '%s'", app.Mark, app.Size, app.Layout)
			if err := ResizeMark(app.Mark, app.Size, app.Layout); err != nil {
				log.Warn("Failed to resize '%s': %v", app.Mark, err)
				// Continue with other resizes even if some fail
			}
			time.Sleep(200 * time.Millisecond)
		}
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
