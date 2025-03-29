package manager

import (
	"fmt"
	"log"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func NewWorkspaceManager(client *sway.Client) *WorkspaceManager {
	return &WorkspaceManager{
		client: client,
	}
}

func (wm *WorkspaceManager) SetupWorkspace(num int, workspace config.Workspace, currentApps []sway.AppNode) error {
	if err := wm.client.SwitchToWorkspace(num); err != nil {
		return fmt.Errorf("switching to workspace: %w", err)
	}

	if workspace.Layout != "" {
		if err := wm.client.SetWorkspaceLayout(workspace.Layout); err != nil {
			log.Printf("Warning: Failed to set layout: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Process actions
	var updateApps []AppUpdate
	var launchApps []config.App

	// Track which current apps are matched with config apps
	matched := make(map[int64]bool)

	for _, configApp := range workspace.Apps {
		found := false

		for _, currentApp := range currentApps {
			if currentApp.Name == configApp.Name {
				// App is already running - might need updating
				found = true
				matched[currentApp.NodeID] = true

				// Check if properties need updating
				needsUpdate := configApp.Floating != currentApp.Floating ||
					configApp.Size != "" ||
					(configApp.Position != "" && configApp.Floating)

				if needsUpdate {
					updateApps = append(updateApps, AppUpdate{
						NodeID: currentApp.NodeID,
						Config: configApp,
					})
				}

				break
			}
		}

		if !found {
			// App is in config but not running
			launchApps = append(launchApps, configApp)
		}
	}

	// Close apps that aren't in config (but only if configured to do so)
	if workspace.CloseUnmatched {
		for _, currentApp := range currentApps {
			if !matched[currentApp.NodeID] {
				log.Printf("Closing app %s (ID: %d) because it's not in workspace config",
					currentApp.Name, currentApp.NodeID)
				if err := wm.client.KillWindow(currentApp.NodeID); err != nil {
					log.Printf("Warning: Failed to close app: %v", err)
				}
				time.Sleep(200 * time.Millisecond)
			}
		}
	}

	appManager := NewAppManager(wm.client)

	for _, update := range updateApps {
		log.Printf("Updating app: %s", update.Config.Name)
		if err := appManager.UpdateApp(update.NodeID, update.Config, workspace.Layout); err != nil {
			log.Printf("Warning: Failed to update app: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	for _, app := range launchApps {
		log.Printf("Launching app: %s", app.Name)
		if err := appManager.LaunchApp(app, workspace.Layout); err != nil {
			log.Printf("Warning: Failed to launch app: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
