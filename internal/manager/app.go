package manager

import (
	"fmt"
	"log"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func (m *Manager) LaunchApp(app config.App, layout string) error {
	cmdStr := app.Command
	if cmdStr == "" {
		cmdStr = app.Name
	}

	m.logVerbose("Launching: %s", cmdStr)

	if err := m.runCommandAsync(cmdStr); err != nil {
		return fmt.Errorf("launching %s: %w", app.Name, err)
	}

	nodeID, err := m.waitForAppToAppear(app)
	if err != nil {
		return err
	}

	return m.configureApp(nodeID, app, layout)
}

func (m *Manager) UpdateApp(nodeID int64, app config.App, layout string) error {
	m.logVerbose("Updating app: %s (ID: %d)", app.Name, nodeID)

	if err := m.configureApp(nodeID, app, layout); err != nil {
		return fmt.Errorf("configuring app: %w", err)
	}

	// Run post-launch commands
	for _, cmd := range app.Posts {
		m.logVerbose("Running post command: %s", cmd)
		if err := m.runCommand(cmd); err != nil {
			log.Printf("Warning: post command failed: %v", err)
		}
	}

	return nil
}

func (m *Manager) waitForAppToAppear(app config.App) (int64, error) {
	waitTime := 1000 // default 1 second

	if app.Delay > 0 {
		waitTime = int(app.Delay * 1000)
	}

	m.delay(waitTime)

	tree, err := m.Client.GetTree()
	if err != nil {
		return 0, fmt.Errorf("getting tree after launch: %w", err)
	}

	workspaces := tree.FindWorkspaces()
	for _, ws := range workspaces {
		apps := ws.FindAllApps()
		for _, node := range apps {
			if sway.MatchAppName(node.Name, app.Name) {
				return node.NodeID, nil
			}
		}
	}

	return 0, fmt.Errorf("could not find newly launched app %s in the Sway tree", app.Name)
}

func (m *Manager) configureApp(nodeID int64, app config.App, layout string) error {
	if app.Floating {
		if err := m.Client.SetFloating(nodeID, true); err != nil {
			return fmt.Errorf("setting floating state: %w", err)
		}
		m.delay(100)
	}

	if app.Size != "" {
		if err := m.Client.ResizeWindow(nodeID, app.Size, app.Floating, layout); err != nil {
			return fmt.Errorf("setting size: %w", err)
		}
		m.delay(100)
	}

	if app.Position != "" && app.Floating {
		if err := m.Client.MoveWindow(nodeID, app.Position); err != nil {
			return fmt.Errorf("setting position: %w", err)
		}
		m.delay(100)
	}

	m.delay(100)

	return nil
}
