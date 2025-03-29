package app

import (
	"fmt"
	"log"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
	"github.com/titembaatar/sway.flem/internal/util"
)

func NewManager(client *sway.Client, verbose bool) *AppManager {
	return &AppManager{
		Client:  client,
		Verbose: verbose,
	}
}

func (m *AppManager) UpdateApp(nodeID int64, app config.App, layout string) error {
	if m.Verbose {
		log.Printf("Updating app: %s (ID: %d)", app.Name, nodeID)
	}

	if err := m.configureApp(nodeID, app, layout); err != nil {
		return fmt.Errorf("configuring app: %w", err)
	}

	// Run post-launch commands
	for _, cmd := range app.Posts {
		if m.Verbose {
			log.Printf("Running post command: %s", cmd)
		}
		if err := util.RunCommand(cmd); err != nil {
			log.Printf("Warning: post command failed: %v", err)
		}
	}

	return nil
}

func (m *AppManager) configureApp(nodeID int64, app config.App, layout string) error {
	if app.Floating {
		if err := m.Client.SetFloating(nodeID, true); err != nil {
			return fmt.Errorf("setting floating state: %w", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if app.Size != "" {
		if err := m.Client.ResizeWindow(nodeID, app.Size, app.Floating, layout); err != nil {
			return fmt.Errorf("setting size: %w", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if app.Position != "" && app.Floating {
		if err := m.Client.MoveWindow(nodeID, app.Position); err != nil {
			return fmt.Errorf("setting position: %w", err)
		}
	}

	return nil
}

