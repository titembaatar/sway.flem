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
func LaunchApp(app config.App, mark string) error {
	log.Info("Launching application: %s", app.Name)

	cmdStr := app.Cmd
	if cmdStr == "" {
		cmdStr = app.Name
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
	time.Sleep(300 * time.Millisecond)

	log.Debug("Applying mark to application: %s", mark)
	if err := ApplyMark(mark); err != nil {
		log.Error("Failed to apply mark to application: %v", err)
		return fmt.Errorf("failed to apply mark to application: %w", err)
	}

	return nil
}

// Launches applications in a container and returns info for later resizing
func LaunchApps(containerMark string, split string, apps []config.App) ([]AppInfo, error) {
	if len(apps) == 0 {
		return nil, nil
	}

	var appInfos []AppInfo

	// For the first app, we need to set the split layout on the parent container
	if len(apps) > 0 {
		// Focus the parent container first
		if err := FocusMark(containerMark); err != nil {
			return nil, fmt.Errorf("failed to focus parent container: %w", err)
		}

		// Set the split layout for the first app
		splitCmd := fmt.Sprintf("split %s", split)
		if _, err := RunCommand(splitCmd); err != nil {
			return nil, fmt.Errorf("failed to set split layout: %w", err)
		}
	}

	// Launch each app in sequence
	for i, app := range apps {
		appMark := GenerateAppMark(containerMark, fmt.Sprintf("app%d", i+1))

		if err := LaunchApp(app, appMark); err != nil {
			log.Error("Failed to launch app %s: %v", app.Name, err)
			continue // Try to launch remaining apps even if this one fails
		}

		// Store information for later resizing
		appInfos = append(appInfos, AppInfo{
			Mark:   appMark,
			Size:   app.Size,
			Layout: split,
		})

		// After launching the first app, subsequent apps need to focus the parent container
		if i < len(apps)-1 {
			if err := FocusMark(containerMark); err != nil {
				log.Error("Failed to focus parent container: %v", err)
				return appInfos, fmt.Errorf("failed to focus parent container: %w", err)
			}
		}
	}

	return appInfos, nil
}

// Resizes a list of applications according to their stored information
func ResizeApps(appInfos []AppInfo) {
	log.Info("Resizing %d applications", len(appInfos))

	for _, app := range appInfos {
		if app.Size != "" {
			log.Debug("Resizing app with mark '%s' to '%s'", app.Mark, app.Size)
			if err := ResizeMark(app.Mark, app.Size, app.Layout); err != nil {
				log.Warn("Failed to resize app '%s': %v", app.Mark, err)
				// Continue even if resize fails
			}
		}
	}
}
