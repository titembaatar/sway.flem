package app

import (
	"fmt"
	"log"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
	"github.com/titembaatar/sway.flem/internal/util"
)

func (m *AppManager) LaunchApp(app config.App, layout string) error {
	cmdStr := app.Command
	if cmdStr == "" {
		cmdStr = app.Name
	}

	if m.Verbose {
		log.Printf("Launching: %s", cmdStr)
	}

	if err := util.RunCommandAsync(cmdStr); err != nil {
		return fmt.Errorf("launching %s: %w", app.Name, err)
	}

	nodeID, err := m.waitForAppToAppear(app)
	if err != nil {
		return err
	}

	return m.configureApp(nodeID, app, layout)
}

func (m *AppManager) waitForAppToAppear(app config.App) (int64, error) {
	waitTime := 1000

	if app.Delay > 0 {
		waitTime = int(app.Delay * 1000)
	}

	time.Sleep(time.Duration(waitTime) * time.Millisecond)

	tree, err := m.Client.GetTree()
	if err != nil {
		return 0, fmt.Errorf("getting tree after launch: %w", err)
	}

	workspaces := tree.FindWorkspaces()
	for _, ws := range workspaces {
		apps := ws.FindAllApps()
		for _, node := range apps {
			if sway.MatchAppName(node.Name, app.Name) {
				return node.NodeID, nil
			}
		}
	}

	return 0, fmt.Errorf("could not find newly launched app %s in the Sway tree", app.Name)
}

