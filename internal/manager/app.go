package manager

import (
	"fmt"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

// LaunchApp starts an application and returns its node ID
func (m *Manager) LaunchApp(app config.App) (int64, error) {
	cmdStr := app.Command
	if cmdStr == "" {
		cmdStr = app.Name
	}

	m.logDebug("Launching: %s", cmdStr)

	if err := m.runCommandAsync(cmdStr); err != nil {
		return 0, fmt.Errorf("launching %s: %w", app.Name, err)
	}

	// Wait for app to appear
	waitTime := 1000 // default 1 second
	if app.Delay > 0 {
		waitTime = int(app.Delay * 1000)
	}
	m.delay(waitTime)

	// Find the app in the Sway tree
	nodeID, err := m.findAppInTree(app.Name)
	if err != nil {
		return 0, fmt.Errorf("could not find newly launched app %s: %w", app.Name, err)
	}

	return nodeID, nil
}

// ResizeApp adjusts the size of an application window
func (m *Manager) ResizeApp(nodeID int64, app config.App, layout string) error {
	if app.Size == "" {
		return nil // Nothing to do
	}

	m.logDebug("Resizing app: %s (ID: %d) with size: %s and layout: %s", app.Name, nodeID, app.Size, layout)

	if err := m.handleSize(nodeID, app, layout); err != nil {
		return err
	}

	m.delay(200) // Small delay to let Sway process the change
	return nil
}

func (m *Manager) handleSize(nodeID int64, app config.App, layout string) error {
	// Check if the size is specified as a percentage
	if strings.HasSuffix(app.Size, "ppt") || strings.HasSuffix(app.Size, "%") {
		parts := strings.Fields(app.Size)
		if len(parts) == 0 {
			return fmt.Errorf("invalid size format: %s", app.Size)
		}

		sizeStr := parts[0]
		sizeStr = strings.TrimSuffix(sizeStr, "ppt")
		sizeStr = strings.TrimSuffix(sizeStr, "%")

		var command string
		if layout == "splitv" {
			m.logDebug("Applying vertical resize: %s", parts[0])
			command = fmt.Sprintf("[id=%d] resize set height %s", nodeID, parts[0])
		} else { // splith
			m.logDebug("Applying horizontal resize: %s", parts[0])
			command = fmt.Sprintf("[id=%d] resize set width %s", nodeID, parts[0])
		}

		if err := m.Client.ExecuteCommand(command); err != nil {
			return fmt.Errorf("resizing in %s layout: %w", layout, err)
		}
	}

	return nil
}

// RunPostCommands executes any commands that should run after an app is launched
func (m *Manager) RunPostCommands(app config.App) error {
	if len(app.Posts) == 0 {
		return nil
	}

	m.logDebug("Running %d post commands for app: %s", len(app.Posts), app.Name)

	for i, cmd := range app.Posts {
		m.logDebug("Running post command %d: %s", i+1, cmd)
		if err := m.runCommand(cmd); err != nil {
			return fmt.Errorf("running post command '%s': %w", cmd, err)
		}
		m.delay(100) // Small delay between commands
	}

	return nil
}

// WaitForApp waits for an app to appear in the tree for up to the specified timeout
func (m *Manager) WaitForApp(appName string, timeoutSec int) (int64, error) {
	m.logDebug("Waiting for app to appear: %s (timeout: %ds)", appName, timeoutSec)

	start := time.Now()
	timeoutDuration := time.Duration(timeoutSec) * time.Second

	for {
		if time.Since(start) > timeoutDuration {
			return 0, fmt.Errorf("timeout waiting for %s to appear", appName)
		}

		nodeID, err := m.findAppInTree(appName)
		if err == nil {
			return nodeID, nil
		}

		m.delay(100) // 100ms delay between checks
	}
}

// findAppInTree looks for an app by name in the Sway tree
// Returns the node ID if found, error otherwise
func (m *Manager) findAppInTree(appName string) (int64, error) {
	tree, err := m.Client.GetTree()
	if err != nil {
		return 0, fmt.Errorf("getting tree: %w", err)
	}

	workspaces := tree.FindWorkspaces()
	for _, ws := range workspaces {
		apps := ws.FindAllApps()
		for _, node := range apps {
			if sway.MatchAppName(node.Name, appName) {
				return node.NodeID, nil
			}
		}
	}

	return 0, fmt.Errorf("app %s not found in tree", appName)
}

// FindAppByID looks for an app in the Sway tree by its node ID
// This is more reliable than searching by name as IDs are unique
func (m *Manager) FindAppByID(nodeID int64) (*sway.AppNode, error) {
	tree, err := m.Client.GetTree()
	if err != nil {
		return nil, fmt.Errorf("getting tree: %w", err)
	}

	workspaces := tree.FindWorkspaces()
	for _, ws := range workspaces {
		apps := ws.FindAllApps()
		for i, node := range apps {
			if node.NodeID == nodeID {
				return &apps[i], nil
			}
		}
	}

	return nil, fmt.Errorf("app with ID %d not found in tree", nodeID)
}

// FindAppConfigByName looks up an app configuration by name
// Note: Multiple apps can have the same name, so this should be used with caution
func FindAppConfigByName(apps []config.App, name string) config.App {
	for _, app := range apps {
		if app.Name == name {
			return app
		}
	}
	return config.App{} // Return empty app if not found
}
