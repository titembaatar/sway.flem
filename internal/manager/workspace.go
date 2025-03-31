package manager

import (
	"fmt"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

// SetupWorkspace configures a workspace based on the provided configuration
func (m *Manager) SetupWorkspace(wsNum int, workspace config.Workspace, currentApps []sway.AppNode) error {
	m.logDebug("Setting up workspace %d", wsNum)

	if err := m.initializeWorkspace(wsNum, workspace); err != nil {
		return err
	}

	// Use container-aware setup if the workspace has containers
	if workspace.Container != nil {
		m.logDebug("Workspace %d has containers, using container-aware setup", wsNum)
		return m.setupWorkspaceWithContainers(wsNum, workspace, currentApps)
	}

	// Use standard setup for workspaces without containers
	return m.setupStandardWorkspace(wsNum, workspace, currentApps)
}

// initializeWorkspace sets up the basic properties of a workspace
func (m *Manager) initializeWorkspace(wsNum int, workspace config.Workspace) error {
	// Switch to the workspace
	if err := m.Client.SwitchToWorkspace(wsNum); err != nil {
		return fmt.Errorf("switching to workspace: %w", err)
	}

	// Move to specified output if configured
	if workspace.Output != "" {
		if err := m.moveWorkspaceToOutput(wsNum, workspace.Output); err != nil {
			m.logWarn("Failed to move workspace to output: %v", err)
		}
		m.delay(200)
	}

	// Set the workspace layout if configured
	if workspace.Layout != "" {
		if sway.IsValidLayout(workspace.Layout) {
			m.logDebug("Setting workspace %d layout to %s", wsNum, workspace.Layout)
			layoutCmd := m.GetLayoutCommand(workspace.Layout)
			if err := m.Client.ExecuteCommand(layoutCmd); err != nil {
				m.logWarn("Failed to set layout: %v", err)
			}
			m.delay(200)
		} else {
			m.logWarn("Invalid layout '%s' for workspace %d", workspace.Layout, wsNum)
		}
	}

	return nil
}

// setupWorkspaceWithContainers handles setting up a workspace that contains nested containers
func (m *Manager) setupWorkspaceWithContainers(wsNum int, workspace config.Workspace, currentApps []sway.AppNode) error {
	m.logDebug("Setting up workspace %d with containers", wsNum)

	// Tracking list for apps that need to be resized
	var appsToResize []AppInfo

	// Launch apps directly in the workspace (without resizing yet)
	for i, app := range workspace.Apps {
		m.logDebug("Launching workspace app %d: %s", i, app.Name)

		// Launch app
		nodeID, err := m.LaunchApp(app)
		if err != nil {
			m.logWarn("Failed to launch workspace app %s: %v", app.Name, err)
			continue
		}
		m.delay(300)

		// Add to resize list, don't resize yet
		if app.Size != "" {
			appsToResize = append(appsToResize, AppInfo{
				App:    app,
				NodeID: nodeID,
				Layout: workspace.Layout,
			})
		}
	}

	// Set up container structure and resize apps at appropriate times
	if workspace.Container != nil {
		m.logDebug("Setting up container structure in workspace %d", wsNum)
		if err := m.SetupContainerStructure(workspace.Container, workspace.Layout, &appsToResize); err != nil {
			return fmt.Errorf("setting up container: %w", err)
		}
	}

	// Resize any remaining apps that weren't resized during container setup
	if len(appsToResize) > 0 {
		m.logDebug("Resizing remaining apps (%d apps)", len(appsToResize))
		m.resizeAppBatch(appsToResize)
	}

	// Handle CloseUnmatched for workspaces with containers
	if workspace.CloseUnmatched {
		m.logDebug("Closing unmatched applications on workspace %d (including containers)", wsNum)
		allApps := m.CollectAllApps(workspace)
		m.closeUnmatchedApps(allApps, currentApps)
	}

	return nil
}

// setupStandardWorkspace handles setting up a workspace without containers
func (m *Manager) setupStandardWorkspace(wsNum int, workspace config.Workspace, currentApps []sway.AppNode) error {
	m.logDebug("Setting up standard workspace %d without containers", wsNum)

	// First, categorize apps into those to update and those to launch
	updateApps, launchApps := m.categorizeApps(workspace.Apps, currentApps)

	// Launch new apps
	appIDs := m.launchWorkspaceApps(launchApps)

	// Configure all apps (state, sizing, post commands)
	m.configureWorkspaceApps(appIDs, updateApps, workspace.Apps, workspace.Layout)

	// Handle CloseUnmatched if enabled
	if workspace.CloseUnmatched {
		m.logDebug("Closing unmatched applications on workspace %d", wsNum)
		m.closeUnmatchedApps(workspace.Apps, currentApps)
	}

	m.logDebug("Workspace %d setup complete", wsNum)
	return nil
}

// launchWorkspaceApps launches apps specified in the workspace configuration
func (m *Manager) launchWorkspaceApps(apps []config.App) map[string]int64 {
	m.logDebug("Launching %d new applications", len(apps))
	appIDs := make(map[string]int64)

	for _, app := range apps {
		m.logDebug("Launching app: %s", app.Name)
		nodeID, err := m.LaunchApp(app)
		if err != nil {
			m.logWarn("Failed to launch app %s: %v", app.Name, err)
			continue
		}
		appIDs[app.Name] = nodeID
	}

	m.delay(500) // Give apps time to stabilize
	return appIDs
}

// configureWorkspaceApps configures both new and existing apps in the workspace
func (m *Manager) configureWorkspaceApps(
	appIDs map[string]int64,
	updateApps []AppUpdate,
	configApps []config.App,
	layout string,
) {

	// Configure app sizes for all apps
	m.configureAppSizes(appIDs, updateApps, configApps, layout)

	// Run post-launch commands
	m.runWorkspacePostCommands(appIDs, updateApps, configApps)
}

// configureAppSizes resizes applications according to configuration
func (m *Manager) configureAppSizes(
	appIDs map[string]int64,
	updateApps []AppUpdate,
	configApps []config.App,
	layout string,
) {
	m.logDebug("Resizing %d new applications and %d existing applications",
		len(appIDs), len(updateApps))

	// Resize new apps
	for appName, nodeID := range appIDs {
		app := FindAppConfigByName(configApps, appName)
		if app.Size != "" {
			if err := m.ResizeApp(nodeID, app, layout); err != nil {
				m.logWarn("Failed to resize app %s: %v", appName, err)
			}
			m.delay(100)
		}
	}

	// Resize existing apps
	for _, update := range updateApps {
		if update.Config.Size != "" {
			if err := m.ResizeApp(update.NodeID, update.Config, layout); err != nil {
				m.logWarn("Failed to resize existing app %s: %v",
					update.Config.Name, err)
			}
			m.delay(100)
		}
	}
}

// runWorkspacePostCommands executes post-launch commands for apps
func (m *Manager) runWorkspacePostCommands(
	appIDs map[string]int64,
	updateApps []AppUpdate,
	configApps []config.App,
) {
	m.logDebug("Running post-launch commands")

	// Run post commands for new apps
	for appName := range appIDs {
		app := FindAppConfigByName(configApps, appName)
		if len(app.Posts) > 0 {
			if err := m.RunPostCommands(app); err != nil {
				m.logWarn("Failed to run post commands for app %s: %v", appName, err)
			}
		}
	}

	// Run post commands for existing apps
	for _, update := range updateApps {
		if len(update.Config.Posts) > 0 {
			if err := m.RunPostCommands(update.Config); err != nil {
				m.logWarn("Failed to run post commands for existing app %s: %v",
					update.Config.Name, err)
			}
		}
	}
}

// moveWorkspaceToOutput moves a workspace to the specified output
func (m *Manager) moveWorkspaceToOutput(wsNum int, output string) error {
	wsInfo, err := m.Client.GetWorkspaceInfo(wsNum)
	if err == nil && wsInfo != nil {
		if wsInfo.Output != output {
			m.logDebug("Moving workspace %d to output %s (currently on %s)",
				wsNum, output, wsInfo.Output)
			if err := m.Client.MoveWorkspaceToOutput(wsNum, output); err != nil {
				return err
			}

			m.delay(200)
			if err := m.Client.SwitchToWorkspace(wsNum); err != nil {
				m.logWarn("Failed to switch to workspace after moving: %v", err)
			}
		}
	}

	return nil
}

// categorizeApps groups apps into those to update and those to launch
func (m *Manager) categorizeApps(configApps []config.App, currentApps []sway.AppNode) ([]AppUpdate, []config.App) {
	var updateApps []AppUpdate
	var launchApps []config.App

	matched := make(map[int64]bool)

	// Identify apps that already exist and need to be updated
	for _, configApp := range configApps {
		found := false

		for _, currentApp := range currentApps {
			if sway.MatchAppName(currentApp.Name, configApp.Name) {
				found = true
				matched[currentApp.NodeID] = true

				needsUpdate := configApp.Size != ""

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

// closeUnmatchedApps closes apps in the workspace that aren't in the config
func (m *Manager) closeUnmatchedApps(configApps []config.App, currentApps []sway.AppNode) {
	matched := make(map[int64]bool)

	// Mark apps that match the configuration
	for _, configApp := range configApps {
		for _, currentApp := range currentApps {
			if sway.MatchAppName(currentApp.Name, configApp.Name) {
				matched[currentApp.NodeID] = true
			}
		}
	}

	// Close apps that aren't matched
	for _, currentApp := range currentApps {
		if !matched[currentApp.NodeID] {
			m.logDebug("Closing app %s (ID: %d) because it's not in workspace config",
				currentApp.Name, currentApp.NodeID)
			if err := m.Client.KillWindow(currentApp.NodeID); err != nil {
				m.logWarn("Failed to close app %s: %v", currentApp.Name, err)
			}
			m.delay(200)
		}
	}
}
