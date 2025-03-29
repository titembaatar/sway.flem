// Package manager provides the main orchestration logic for sway.flem
package manager

import (
	"fmt"
	"log"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

type Manager struct {
	client *sway.Client
	config *config.Config
}

func NewManager(client *sway.Client, cfg *config.Config) *Manager {
	return &Manager{
		client: client,
		config: cfg,
	}
}

func (m *Manager) Run() error {
	// Get current state of Sway workspaces
	tree, err := m.client.GetTree()
	if err != nil {
		return fmt.Errorf("getting workspace tree: %w", err)
	}

	workspaces := tree.FindWorkspaces()

	wsManager := NewWorkspaceManager(m.client)

	// Process each workspace in config
	for wsNum, wsConfig := range m.config.Workspaces {
		log.Printf("Processing workspace %d", wsNum)

		// Get current apps in workspace
		var currentApps []sway.AppNode
		if ws, exists := workspaces[wsNum]; exists {
			currentApps = ws.FindAllApps()
		}

		// Setup the workspace
		if err := wsManager.SetupWorkspace(wsNum, wsConfig, currentApps); err != nil {
			log.Printf("Error setting up workspace %d: %v", wsNum, err)
			// Continue with other workspaces instead of failing entirely
		}
	}

	// Focus the final workspace if specified
	if m.config.FocusWorkspace > 0 {
		if err := m.client.SwitchToWorkspace(m.config.FocusWorkspace); err != nil {
			return fmt.Errorf("focusing workspace %d: %w", m.config.FocusWorkspace, err)
		}
	}

	return nil
}
