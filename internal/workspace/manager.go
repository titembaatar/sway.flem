package workspace

import (
	"fmt"
	"log"
	"time"

	"github.com/titembaatar/sway.flem/internal/app"
	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func NewManager(client *sway.Client, appManager *app.AppManager, verbose bool) *WorkspaceManager {
	return &WorkspaceManager{
		Client:  client,
		Verbose: verbose,
		AppMgr:  appManager,
	}
}

func (wm *WorkspaceManager) SetupWorkspace(wsNum int, workspace config.Workspace, currentApps []sway.AppNode) error {
	if err := wm.Client.SwitchToWorkspace(wsNum); err != nil {
		return fmt.Errorf("switching to workspace: %w", err)
	}

	if workspace.Output != "" {
		if err := wm.moveWorkspaceToOutput(wsNum, workspace.Output); err != nil {
			log.Printf("Warning: Failed to move workspace to output: %v", err)
		}
		time.Sleep(200 * time.Millisecond)
	}

	if workspace.Layout != "" {
		if wm.Verbose {
			log.Printf("Setting workspace %d layout to %s", wsNum, workspace.Layout)
		}
		if err := wm.Client.SetWorkspaceLayout(workspace.Layout); err != nil {
			log.Printf("Warning: Failed to set layout: %v", err)
		}
		time.Sleep(200 * time.Millisecond)
	}

	updateApps, launchApps := wm.categorizeApps(workspace.Apps, currentApps)

	for _, update := range updateApps {
		if wm.Verbose {
			log.Printf("Updating app: %s", update.Config.Name)
		}
		if err := wm.AppMgr.UpdateApp(update.NodeID, update.Config, workspace.Layout); err != nil {
			log.Printf("Warning: Failed to update app: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	for _, app := range launchApps {
		if wm.Verbose {
			log.Printf("Launching app: %s", app.Name)
		}
		if err := wm.AppMgr.LaunchApp(app, workspace.Layout); err != nil {
			log.Printf("Warning: Failed to launch app: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	if workspace.CloseUnmatched {
		wm.closeUnmatchedApps(workspace.Apps, currentApps)
	}

	return nil
}

func (wm *WorkspaceManager) moveWorkspaceToOutput(wsNum int, output string) error {
	wsInfo, err := wm.Client.GetWorkspaceInfo(wsNum)
	if err == nil && wsInfo != nil {
		// Only move if the workspace is on a different output
		if wsInfo.Output != output {
			if wm.Verbose {
				log.Printf("Moving workspace %d to output %s (currently on %s)",
					wsNum, output, wsInfo.Output)
			}
			if err := wm.Client.MoveWorkspaceToOutput(wsNum, output); err != nil {
				return err
			}

			// Switch to workspace again after moving it
			time.Sleep(200 * time.Millisecond)
			if err := wm.Client.SwitchToWorkspace(wsNum); err != nil {
				log.Printf("Warning: Failed to switch to workspace after moving: %v", err)
			}
		}
	}

	return nil
}

func (wm *WorkspaceManager) categorizeApps(configApps []config.App, currentApps []sway.AppNode) ([]AppUpdate, []config.App) {
	var updateApps []AppUpdate
	var launchApps []config.App

	matched := make(map[int64]bool)

	// Match existing apps with config
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

func (wm *WorkspaceManager) closeUnmatchedApps(configApps []config.App, currentApps []sway.AppNode) {
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
			if wm.Verbose {
				log.Printf("Closing app %s (ID: %d) because it's not in workspace config",
					currentApp.Name, currentApp.NodeID)
			}
			if err := wm.Client.KillWindow(currentApp.NodeID); err != nil {
				log.Printf("Warning: Failed to close app: %v", err)
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}
