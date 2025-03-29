package manager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func (m *Manager) LaunchApp(app config.App) (int64, error) {
	cmdStr := app.Command
	if cmdStr == "" {
		cmdStr = app.Name
	}

	m.logVerbose("Launching: %s", cmdStr)

	if err := m.runCommandAsync(cmdStr); err != nil {
		return 0, fmt.Errorf("launching %s: %w", app.Name, err)
	}

	waitTime := 1000 // default 1 second
	if app.Delay > 0 {
		waitTime = int(app.Delay * 1000)
	}
	m.delay(waitTime)

	nodeID, err := m.findAppInTree(app.Name)
	if err != nil {
		return 0, fmt.Errorf("could not find newly launched app %s: %w", app.Name, err)
	}

	return nodeID, nil
}

func (m *Manager) ResizeApp(nodeID int64, app config.App, layout string) error {
	if app.Size == "" {
		return nil // Nothing to do
	}

	m.logVerbose("Resizing app: %s (ID: %d) with size: %s", app.Name, nodeID, app.Size)

	if !app.Floating && (layout == "splitv" || layout == "splith") {
		if strings.HasSuffix(app.Size, "ppt") || strings.HasSuffix(app.Size, "%") {
			sizeStr := app.Size
			sizeStr = strings.TrimSuffix(sizeStr, "ppt")
			sizeStr = strings.TrimSuffix(sizeStr, "%")

			size, err := strconv.ParseFloat(sizeStr, 64)
			if err != nil {
				return fmt.Errorf("parsing size value: %w", err)
			}

			if strings.HasSuffix(app.Size, "%") {
				size = size / 100.0
			}

			if layout == "splitv" {
				m.logVerbose("Applying vertical resize: %s", app.Size)
				command := fmt.Sprintf("[con_id=%d] resize set height %s", nodeID, app.Size)
				if err := m.Client.ExecuteCommand(command); err != nil {
					return fmt.Errorf("vertical resize: %w", err)
				}
			} else { // splith
				m.logVerbose("Applying horizontal resize: %s", app.Size)
				command := fmt.Sprintf("[con_id=%d] resize set width %s", nodeID, app.Size)
				if err := m.Client.ExecuteCommand(command); err != nil {
					return fmt.Errorf("horizontal resize: %w", err)
				}
			}

			return nil
		}
	}

	if err := m.Client.ResizeWindow(nodeID, app.Size, app.Floating, layout); err != nil {
		return fmt.Errorf("setting size: %w", err)
	}

	return nil
}

func (m *Manager) SetFloatingState(nodeID int64, app config.App) error {
	m.logVerbose("Setting floating state for app: %s (ID: %d) to: %v", app.Name, nodeID, app.Floating)

	if err := m.Client.SetFloating(nodeID, app.Floating); err != nil {
		return fmt.Errorf("setting floating state: %w", err)
	}

	return nil
}

func (m *Manager) RunPostCommands(app config.App) error {
	if len(app.Posts) == 0 {
		return nil
	}

	m.logVerbose("Running post commands for app: %s", app.Name)

	for _, cmd := range app.Posts {
		m.logVerbose("Running post command: %s", cmd)
		if err := m.runCommand(cmd); err != nil {
			return fmt.Errorf("running post command: %w", err)
		}
	}

	return nil
}

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

	return 0, fmt.Errorf("app not found in tree")
}

func FindAppByName(apps []config.App, name string) config.App {
	for _, app := range apps {
		if app.Name == name {
			return app
		}
	}
	return config.App{} // Should not happen if appIDs is built correctly
}
