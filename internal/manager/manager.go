package manager

import (
	"fmt"
	"log"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func NewManager(cfg *config.Config, swayClient *sway.Client, verbose bool) *Manager {
	return &Manager{
		Config:  cfg,
		Client:  swayClient,
		Verbose: verbose,
	}
}

func (m *Manager) Run() error {
	if m.Verbose {
		log.Println("Starting sway.flem workspace management")
	}

	tree, err := m.Client.GetTree()
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

		if err := m.SetupWorkspace(wsNum, wsConfig, currentApps); err != nil {
			log.Printf("Error setting up workspace %d: %v", wsNum, err)
		}
	}

	if m.Config.FocusWorkspace > 0 {
		if m.Verbose {
			log.Printf("Focusing workspace %d", m.Config.FocusWorkspace)
		}

		if err := m.Client.SwitchToWorkspace(m.Config.FocusWorkspace); err != nil {
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
	_, err := m.Client.GetTree()
	return err == nil
}
