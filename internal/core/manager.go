package core

import (
	"fmt"
	"log"

	"github.com/titembaatar/sway.flem/internal/app"
	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
	"github.com/titembaatar/sway.flem/internal/workspace"
)

func NewManager(config *config.Config, swayClient *sway.Client, verbose bool) *Manager {
	appManager := app.NewManager(swayClient, verbose)

	workspaceManager := workspace.NewManager(swayClient, appManager, verbose)

	return &Manager{
		Config:           config,
		SwayClient:       swayClient,
		WorkspaceManager: workspaceManager,
		AppManager:       appManager,
		Verbose:          verbose,
	}
}

func (m *Manager) Run() error {
	if m.Verbose {
		log.Println("Starting sway.flem workspace management")
	}

	tree, err := m.SwayClient.GetTree()
	if err != nil {
		return fmt.Errorf("getting workspace tree: %w", err)
	}

	workspaces := tree.FindWorkspaces()

	for wsNum, wsConfig := range m.Config.Workspaces {
		if m.Verbose {
			log.Printf("Processing workspace %d", wsNum)
		}

		var currentApps []sway.AppNode
		if ws, exists := workspaces[wsNum]; exists {
			currentApps = ws.FindAllApps()
		}

		if err := m.WorkspaceManager.SetupWorkspace(wsNum, wsConfig, currentApps); err != nil {
			log.Printf("Error setting up workspace %d: %v", wsNum, err)
			// Continue with other workspaces instead of failing entirely
		}
	}

	if m.Config.FocusWorkspace > 0 {
		if m.Verbose {
			log.Printf("Focusing workspace %d", m.Config.FocusWorkspace)
		}

		if err := m.SwayClient.SwitchToWorkspace(m.Config.FocusWorkspace); err != nil {
			return fmt.Errorf("focusing workspace %d: %w", m.Config.FocusWorkspace, err)
		}
	}

	if m.Verbose {
		log.Println("Workspace setup complete!")
	}

	return nil
}

func (m *Manager) Validate() error {
	if !m.isSwayRunning() {
		return fmt.Errorf("Sway window manager is not running")
	}

	if len(m.Config.Workspaces) == 0 {
		return fmt.Errorf("no workspaces defined in configuration")
	}

	return nil
}

func (m *Manager) isSwayRunning() bool {
	_, err := m.SwayClient.GetTree()
	return err == nil
}
