package manager

import (
	"fmt"
	"log"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func (m *Manager) SetupWorkspace(wsNum int, workspace config.Workspace, currentApps []sway.AppNode) error {
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

	updateApps, launchApps := m.categorizeApps(workspace.Apps, currentApps)

	for _, update := range updateApps {
		m.logVerbose("Updating app: %s", update.Config.Name)
		if err := m.UpdateApp(update.NodeID, update.Config, workspace.Layout); err != nil {
			log.Printf("Warning: Failed to update app: %v", err)
		}
		m.delay(100)
	}

	for _, app := range launchApps {
		m.logVerbose("Launching app: %s", app.Name)
		if err := m.LaunchApp(app, workspace.Layout); err != nil {
			log.Printf("Warning: Failed to launch app: %v", err)
		}
		m.delay(500)
	}

	if workspace.CloseUnmatched {
		m.closeUnmatchedApps(workspace.Apps, currentApps)
	}

	return nil
}

func (m *Manager) moveWorkspaceToOutput(wsNum int, output string) error {
	wsInfo, err := m.Client.GetWorkspaceInfo(wsNum)
	if err == nil && wsInfo != nil {
		// Only move if the workspace is on a different output
		if wsInfo.Output != output {
			m.logVerbose("Moving workspace %d to output %s (currently on %s)",
				wsNum, output, wsInfo.Output)
			if err := m.Client.MoveWorkspaceToOutput(wsNum, output); err != nil {
				return err
			}

			// Switch to workspace again after moving it
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
