package manager

import (
	"fmt"
	"log"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

// setupWorkspaceWithContainers handles setting up a workspace that contains nested containers
func (m *Manager) setupWorkspaceWithContainers(wsNum int, workspace config.Workspace, currentApps []sway.AppNode) error {
	m.logVerbose("Setting up workspace %d with containers", wsNum)

	// Tracking list for apps that need to be resized
	var appsToResize []AppInfo

	// STEP 1: Launch apps directly in the workspace (without resizing yet)
	for i, app := range workspace.Apps {
		m.logVerbose("Launching workspace app %d: %s", i, app.Name)

		// Launch app
		nodeID, err := m.LaunchApp(app)
		if err != nil {
			log.Printf("Warning: Failed to launch workspace app %s: %v", app.Name, err)
			continue
		}
		m.delay(300)

		// Set floating state immediately if needed
		if app.Floating {
			if err := m.SetFloatingState(nodeID, app); err != nil {
				log.Printf("Warning: Failed to set floating state for workspace app %s: %v", app.Name, err)
			}
			m.delay(100)
		}

		// Add to resize list, don't resize yet
		if app.Size != "" {
			appsToResize = append(appsToResize, AppInfo{
				App:    app,
				NodeID: nodeID,
				Layout: workspace.Layout,
			})
		}
	}

	// STEP 2: Set up container structure and resize apps at appropriate times
	if workspace.Container != nil {
		m.logVerbose("Setting up container structure in workspace %d", wsNum)
		if err := m.SetupContainerStructure(workspace.Container, workspace.Layout, &appsToResize); err != nil {
			return fmt.Errorf("setting up container: %w", err)
		}
	}

	// STEP 3: Resize any remaining apps that weren't resized during container setup
	if len(appsToResize) > 0 {
		m.logVerbose("Resizing remaining apps (%d apps)", len(appsToResize))
		m.resizeAppBatch(appsToResize)
	}

	// STEP 4: Handle CloseUnmatched for workspaces with containers
	if workspace.CloseUnmatched {
		m.logVerbose("Closing unmatched applications on workspace %d (including containers)", wsNum)
		allApps := m.CollectAllApps(workspace)
		m.closeUnmatchedApps(allApps, currentApps)
	}

	return nil
}

func (m *Manager) SetupWorkspace(wsNum int, workspace config.Workspace, currentApps []sway.AppNode) error {
	m.logVerbose("Setting up workspace %d", wsNum)

	if err := m.initializeWorkspace(wsNum, workspace); err != nil {
		return err
	}

	// If workspace has containers, use container-aware setup
	if workspace.Container != nil {
		m.logVerbose("Workspace %d has containers, using container-aware setup", wsNum)
		return m.setupWorkspaceWithContainers(wsNum, workspace, currentApps)
	}

	// Standard non-container setup
	var appsToResize []AppInfo

	// Launch all apps first (without resizing)
	for i, app := range workspace.Apps {
		m.logVerbose("Launching workspace app %d: %s", i, app.Name)

		// Launch app
		nodeID, err := m.LaunchApp(app)
		if err != nil {
			log.Printf("Warning: Failed to launch workspace app %s: %v", app.Name, err)
			continue
		}
		m.delay(300)

		// Set floating state immediately
		if app.Floating {
			if err := m.SetFloatingState(nodeID, app); err != nil {
				log.Printf("Warning: Failed to set floating state for workspace app %s: %v", app.Name, err)
			}
			m.delay(100)
		}

		// Add to resize list
		if app.Size != "" {
			appsToResize = append(appsToResize, AppInfo{
				App:    app,
				NodeID: nodeID,
				Layout: workspace.Layout,
			})
		}

		// Run post commands
		if len(app.Posts) > 0 {
			if err := m.RunPostCommands(app); err != nil {
				log.Printf("Warning: Failed to run post commands for workspace app %s: %v", app.Name, err)
			}
		}
	}

	// Resize all apps at once
	if len(appsToResize) > 0 {
		m.logVerbose("Resizing all workspace apps (%d apps)", len(appsToResize))
		m.resizeAppBatch(appsToResize)
	}

	if workspace.CloseUnmatched {
		m.logVerbose("Closing unmatched applications on workspace %d", wsNum)
		m.closeUnmatchedApps(workspace.Apps, currentApps)
	}

	m.logVerbose("Workspace %d setup complete", wsNum)
	return nil
}

func (m *Manager) initializeWorkspace(wsNum int, workspace config.Workspace) error {
	if err := m.Client.SwitchToWorkspace(wsNum); err != nil {
		return fmt.Errorf("switching to workspace: %w", err)
	}

	if workspace.Output != "" {
		if err := m.moveWorkspaceToOutput(wsNum, workspace.Output); err != nil {
			log.Printf("Warning: Failed to move workspace to output: %v", err)
		}
		m.delay(200)
	}

	if workspace.Layout != "" {
		if sway.IsValidLayout(workspace.Layout) {
			m.logVerbose("Setting workspace %d layout to %s", wsNum, workspace.Layout)
			if err := m.Client.SetWorkspaceLayout(workspace.Layout); err != nil {
				log.Printf("Warning: Failed to set layout: %v", err)
			}
			m.delay(200)
		} else {
			log.Printf("Warning: Invalid layout '%s' for workspace %d", workspace.Layout, wsNum)
		}
	}

	return nil
}

func (m *Manager) launchNewApps(apps []config.App) map[string]int64 {
	m.logVerbose("Launching new applications")
	appIDs := make(map[string]int64)

	for _, app := range apps {
		m.logVerbose("Launching app: %s", app.Name)
		nodeID, err := m.LaunchApp(app)
		if err != nil {
			log.Printf("Warning: Failed to launch app %s: %v", app.Name, err)
			continue
		}
		appIDs[app.Name] = nodeID
	}

	m.delay(500)

	return appIDs
}

func (m *Manager) configureApps(appIDs map[string]int64, updateApps []AppUpdate, configApps []config.App, layout string) {
	m.configureFloatingState(appIDs, updateApps, configApps)

	m.configureAppSizes(appIDs, updateApps, configApps, layout)

	m.runPostCommands(appIDs, updateApps, configApps)
}

func (m *Manager) configureFloatingState(appIDs map[string]int64, updateApps []AppUpdate, configApps []config.App) {
	m.logVerbose("Setting floating states for applications")

	for appName, nodeID := range appIDs {
		app := FindAppByName(configApps, appName)
		if err := m.SetFloatingState(nodeID, app); err != nil {
			log.Printf("Warning: Failed to set floating state for app %s: %v", appName, err)
		}
		m.delay(100)
	}

	for _, update := range updateApps {
		if err := m.SetFloatingState(update.NodeID, update.Config); err != nil {
			log.Printf("Warning: Failed to set floating state for existing app: %v", err)
		}
		m.delay(100)
	}
}

func (m *Manager) configureAppSizes(appIDs map[string]int64, updateApps []AppUpdate, configApps []config.App, layout string) {
	m.logVerbose("Resizing applications")

	for appName, nodeID := range appIDs {
		app := FindAppByName(configApps, appName)
		if app.Size != "" {
			if err := m.ResizeApp(nodeID, app, layout); err != nil {
				log.Printf("Warning: Failed to resize app %s: %v", appName, err)
			}
			m.delay(100)
		}
	}

	for _, update := range updateApps {
		if update.Config.Size != "" {
			if err := m.ResizeApp(update.NodeID, update.Config, layout); err != nil {
				log.Printf("Warning: Failed to resize existing app: %v", err)
			}
			m.delay(100)
		}
	}
}

func (m *Manager) runPostCommands(appIDs map[string]int64, updateApps []AppUpdate, configApps []config.App) {
	m.logVerbose("Running post-launch commands")

	for appName := range appIDs {
		app := FindAppByName(configApps, appName)
		if len(app.Posts) > 0 {
			if err := m.RunPostCommands(app); err != nil {
				log.Printf("Warning: Failed to run post commands for app %s: %v", appName, err)
			}
		}
	}

	for _, update := range updateApps {
		if len(update.Config.Posts) > 0 {
			if err := m.RunPostCommands(update.Config); err != nil {
				log.Printf("Warning: Failed to run post commands for existing app: %v", err)
			}
		}
	}
}

func (m *Manager) moveWorkspaceToOutput(wsNum int, output string) error {
	wsInfo, err := m.Client.GetWorkspaceInfo(wsNum)
	if err == nil && wsInfo != nil {
		if wsInfo.Output != output {
			m.logVerbose("Moving workspace %d to output %s (currently on %s)",
				wsNum, output, wsInfo.Output)
			if err := m.Client.MoveWorkspaceToOutput(wsNum, output); err != nil {
				return err
			}

			m.delay(200)
			if err := m.Client.SwitchToWorkspace(wsNum); err != nil {
				log.Printf("Warning: Failed to switch to workspace after moving: %v", err)
			}
		}
	}

	return nil
}

func (m *Manager) categorizeApps(configApps []config.App, currentApps []sway.AppNode) ([]AppUpdate, []config.App) {
	var updateApps []AppUpdate
	var launchApps []config.App

	matched := make(map[int64]bool)

	for _, configApp := range configApps {
		found := false

		for _, currentApp := range currentApps {
			if sway.MatchAppName(currentApp.Name, configApp.Name) {
				found = true
				matched[currentApp.NodeID] = true

				needsUpdate := configApp.Floating != currentApp.Floating || configApp.Size != ""

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
			launchApps = append(launchApps, configApp)
		}
	}

	return updateApps, launchApps
}

func (m *Manager) closeUnmatchedApps(configApps []config.App, currentApps []sway.AppNode) {
	matched := make(map[int64]bool)

	for _, configApp := range configApps {
		for _, currentApp := range currentApps {
			if sway.MatchAppName(currentApp.Name, configApp.Name) {
				matched[currentApp.NodeID] = true
			}
		}
	}

	for _, currentApp := range currentApps {
		if !matched[currentApp.NodeID] {
			m.logVerbose("Closing app %s (ID: %d) because it's not in workspace config",
				currentApp.Name, currentApp.NodeID)
			if err := m.Client.KillWindow(currentApp.NodeID); err != nil {
				log.Printf("Warning: Failed to close app: %v", err)
			}
			m.delay(200)
		}
	}
}
