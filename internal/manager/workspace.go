package manager

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func NewWorkspaceManager(client *sway.Client) *WorkspaceManager {
	return &WorkspaceManager{
		client: client,
	}
}

func getAppPosition(positions map[string]int, appName string) int {
	pos, found := positions[strings.ToLower(appName)]
	if found {
		return pos
	}

	return -1
}

func (wm *WorkspaceManager) SetupWorkspace(num int, workspace config.Workspace, currentApps []sway.AppNode) error {
	if err := wm.client.SwitchToWorkspace(num); err != nil {
		return fmt.Errorf("switching to workspace: %w", err)
	}

	// Get current workspace info
	wsInfo, err := wm.client.GetWorkspaceInfo(num)
	var currentLayout string
	var appOrder []string

	if err == nil && wsInfo != nil {
		currentLayout = wsInfo.Layout
		appOrder = wsInfo.AppOrder
	}

	// Set the layout if it's different from current or not specified
	if workspace.Layout != "" && workspace.Layout != currentLayout {
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

	// Map app positions by using appOrder
	appPositions := make(map[string]int)
	for i, appName := range appOrder {
		appPositions[strings.ToLower(appName)] = i
	}

	// First pass: Match existing apps with config and determine what needs updates
	for _, configApp := range workspace.Apps {
		found := false

		for _, currentApp := range currentApps {
			if MatchAppName(currentApp.Name, configApp.Name) {
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

	// Order the updates based on the desired app order
	if len(appOrder) > 0 && len(updateApps) > 1 {
		sort.Slice(updateApps, func(i, j int) bool {
			iPos := getAppPosition(appPositions, updateApps[i].Config.Name)
			jPos := getAppPosition(appPositions, updateApps[j].Config.Name)
			return iPos < jPos
		})
	}

	// Update existing apps
	for _, update := range updateApps {
		log.Printf("Updating app: %s", update.Config.Name)
		if err := appManager.UpdateApp(update.NodeID, update.Config, workspace.Layout); err != nil {
			log.Printf("Warning: Failed to update app: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Order the launches based on the desired app order
	if len(launchApps) > 1 {
		// If we have a config-defined order, use it
		sort.Slice(launchApps, func(i, j int) bool {
			// Try to use the position from the existing layout first
			iPos := getAppPosition(appPositions, launchApps[i].Name)
			jPos := getAppPosition(appPositions, launchApps[j].Name)

			// If we can determine both positions, use them
			if iPos >= 0 && jPos >= 0 {
				return iPos < jPos
			}

			// Otherwise, use the order they appear in the config
			return i < j
		})
	}

	// Launch new apps
	for _, app := range launchApps {
		log.Printf("Launching app: %s", app.Name)
		if err := appManager.LaunchApp(app, workspace.Layout); err != nil {
			log.Printf("Warning: Failed to launch app: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
