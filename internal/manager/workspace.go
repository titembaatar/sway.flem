package manager

import (
	"fmt"
	"log"
	"strings"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func (m *Manager) SetupWorkspace(wsNum int, workspace config.Workspace, currentApps []sway.AppNode) error {
	m.logVerbose("Setting up workspace %d", wsNum)

	if err := m.initializeWorkspace(wsNum, workspace); err != nil {
		return err
	}

	updateApps, launchApps := m.categorizeApps(workspace.Apps, currentApps)

	appIDs := m.launchNewApps(launchApps)

	m.configureApps(appIDs, updateApps, workspace.Apps, workspace.Layout)

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

	m.configureAppPositions(appIDs, updateApps, configApps, layout)

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

func (m *Manager) configureAppPositions(appIDs map[string]int64, updateApps []AppUpdate, configApps []config.App, layout string) {
	m.logVerbose("Positioning windows")

	m.positionFloatingApps(appIDs, updateApps, configApps)

	m.positionTiledApps(appIDs, updateApps, configApps, layout)
}

func (m *Manager) positionFloatingApps(appIDs map[string]int64, updateApps []AppUpdate, configApps []config.App) {
	for appName, nodeID := range appIDs {
		app := FindAppByName(configApps, appName)
		if app.Floating && app.Position != "" {
			if err := m.PositionApp(nodeID, app); err != nil {
				log.Printf("Warning: Failed to position floating app %s: %v", appName, err)
			}
			m.delay(100)
		}
	}

	for _, update := range updateApps {
		if update.Config.Floating && update.Config.Position != "" {
			if err := m.PositionApp(update.NodeID, update.Config); err != nil {
				log.Printf("Warning: Failed to position existing floating app: %v", err)
			}
			m.delay(100)
		}
	}
}

func (m *Manager) positionTiledApps(appIDs map[string]int64, updateApps []AppUpdate, configApps []config.App, layout string) {
	if layout == "splitv" {
		topApps := make([]AppPosition, 0)
		middleApps := make([]AppPosition, 0)
		bottomApps := make([]AppPosition, 0)

		for appName, nodeID := range appIDs {
			app := FindAppByName(configApps, appName)
			if !app.Floating && app.Position != "" {
				pos := getPositionType(app.Position)
				appPos := AppPosition{NodeID: nodeID, App: app}

				switch pos {
				case "top":
					topApps = append(topApps, appPos)
				case "middle", "center":
					middleApps = append(middleApps, appPos)
				case "bottom":
					bottomApps = append(bottomApps, appPos)
				}
			}
		}

		for _, update := range updateApps {
			if !update.Config.Floating && update.Config.Position != "" {
				pos := getPositionType(update.Config.Position)
				appPos := AppPosition{NodeID: update.NodeID, App: update.Config}

				switch pos {
				case "top":
					topApps = append(topApps, appPos)
				case "middle", "center":
					middleApps = append(middleApps, appPos)
				case "bottom":
					bottomApps = append(bottomApps, appPos)
				}
			}
		}

		m.applyVerticalOrdering(topApps, middleApps, bottomApps)
	} else if layout == "splith" {
		leftApps := make([]AppPosition, 0)
		middleApps := make([]AppPosition, 0)
		rightApps := make([]AppPosition, 0)

		for appName, nodeID := range appIDs {
			app := FindAppByName(configApps, appName)
			if !app.Floating && app.Position != "" {
				pos := getPositionType(app.Position)
				appPos := AppPosition{NodeID: nodeID, App: app}

				switch pos {
				case "left":
					leftApps = append(leftApps, appPos)
				case "middle", "center":
					middleApps = append(middleApps, appPos)
				case "right":
					rightApps = append(rightApps, appPos)
				}
			}
		}

		for _, update := range updateApps {
			if !update.Config.Floating && update.Config.Position != "" {
				pos := getPositionType(update.Config.Position)
				appPos := AppPosition{NodeID: update.NodeID, App: update.Config}

				switch pos {
				case "left":
					leftApps = append(leftApps, appPos)
				case "middle", "center":
					middleApps = append(middleApps, appPos)
				case "right":
					rightApps = append(rightApps, appPos)
				}
			}
		}

		m.applyHorizontalOrdering(leftApps, middleApps, rightApps)
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

type AppPosition struct {
	NodeID int64
	App    config.App
}

func getPositionType(position string) string {
	position = strings.ToLower(position)

	if position == "top" || position == "up" {
		return "top"
	} else if position == "bottom" || position == "down" {
		return "bottom"
	} else if position == "left" {
		return "left"
	} else if position == "right" {
		return "right"
	} else if position == "middle" || position == "center" {
		return "middle"
	}

	return position
}

func (m *Manager) applyVerticalOrdering(topApps, middleApps, bottomApps []AppPosition) {
	for _, app := range bottomApps {
		m.logVerbose("Moving app to bottom position: %s (ID: %d)", app.App.Name, app.NodeID)
		if err := m.Client.ExecuteCommand(fmt.Sprintf("[con_id=%d] move position 0 999999", app.NodeID)); err != nil {
			log.Printf("Warning: Failed to position app at bottom: %v", err)
		}
		m.delay(100)
	}

	for _, app := range middleApps {
		m.logVerbose("Moving app to middle position: %s (ID: %d)", app.App.Name, app.NodeID)
		if err := m.Client.ExecuteCommand(fmt.Sprintf("[con_id=%d] move position center", app.NodeID)); err != nil {
			log.Printf("Warning: Failed to position app at middle: %v", err)
		}
		m.delay(100)
	}

	for _, app := range topApps {
		m.logVerbose("Moving app to top position: %s (ID: %d)", app.App.Name, app.NodeID)
		if err := m.Client.ExecuteCommand(fmt.Sprintf("[con_id=%d] move position 0 0", app.NodeID)); err != nil {
			log.Printf("Warning: Failed to position app at top: %v", err)
		}
		m.delay(100)
	}
}

func (m *Manager) applyHorizontalOrdering(leftApps, middleApps, rightApps []AppPosition) {
	for _, app := range rightApps {
		m.logVerbose("Moving app to right position: %s (ID: %d)", app.App.Name, app.NodeID)
		if err := m.Client.ExecuteCommand(fmt.Sprintf("[con_id=%d] move position 999999 center", app.NodeID)); err != nil {
			log.Printf("Warning: Failed to position app to right: %v", err)
		}
		m.delay(100)
	}

	for _, app := range middleApps {
		m.logVerbose("Moving app to middle position: %s (ID: %d)", app.App.Name, app.NodeID)
		if err := m.Client.ExecuteCommand(fmt.Sprintf("[con_id=%d] move position center", app.NodeID)); err != nil {
			log.Printf("Warning: Failed to position app to middle: %v", err)
		}
		m.delay(100)
	}

	for _, app := range leftApps {
		m.logVerbose("Moving app to left position: %s (ID: %d)", app.App.Name, app.NodeID)
		if err := m.Client.ExecuteCommand(fmt.Sprintf("[con_id=%d] move position 0 center", app.NodeID)); err != nil {
			log.Printf("Warning: Failed to position app to left: %v", err)
		}
		m.delay(100)
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

