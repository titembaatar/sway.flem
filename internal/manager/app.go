package manager

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func NewAppManager(client *sway.Client) *AppManager {
	return &AppManager{
		client: client,
	}
}

func MatchAppName(runningApp string, configApp string) bool {
	runningLower := strings.ToLower(runningApp)
	configLower := strings.ToLower(configApp)

	return runningLower == configLower
}

func (am *AppManager) LaunchApp(app config.App, layout string) error {
	cmdStr := app.Command
	if cmdStr == "" {
		cmdStr = app.Name
	}

	log.Printf("Launching: %s", cmdStr)
	cmd := exec.Command("sh", "-c", cmdStr+" &")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("launching %s: %w", app.Name, err)
	}

	if app.Delay == 0 {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Duration(app.Delay) * time.Second)

	tree, err := am.client.GetTree()
	if err != nil {
		return fmt.Errorf("getting tree after launch: %w", err)
	}

	workspaces := tree.FindWorkspaces()
	for _, ws := range workspaces {
		apps := ws.FindAllApps()
		for _, node := range apps {
			if MatchAppName(node.Name, app.Name) {
				return am.configureApp(node.NodeID, app, layout)
			}
		}
	}

	log.Printf("Warning: Could not find newly launched app %s in the Sway tree", app.Name)

	for _, post := range app.Posts {
		log.Printf("Running post command: %s", post)
		cmd := exec.Command("sh", "-c", post)
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: post command failed: %v", err)
		}
	}

	return nil
}

func (am *AppManager) UpdateApp(nodeID int64, app config.App, layout string) error {
	if err := am.configureApp(nodeID, app, layout); err != nil {
		return err
	}

	for _, post := range app.Posts {
		log.Printf("Running post command: %s", post)
		cmd := exec.Command("sh", "-c", post)
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: post command failed: %v", err)
		}
	}

	return nil
}

func (am *AppManager) configureApp(nodeID int64, app config.App, layout string) error {
	tree, err := am.client.GetTree()
	if err == nil {
		// Try to find which workspace this app is in
		for _, output := range tree.Nodes {
			if output.Type != "output" || output.Name == "__i3" {
				continue
			}

			for _, ws := range output.Nodes {
				if ws.Type != "workspace" {
					continue
				}

				// Check if this workspace contains our app
				found := false
				for _, node := range ws.FindAllApps() {
					if node.NodeID == nodeID {
						found = true
						break
					}
				}

				if found && ws.Representation != "" {
					log.Printf("App %s is in workspace with representation: %s",
						app.Name, ws.Representation)
				}
			}
		}
	}

	// Standard configuration logic
	if app.Floating {
		if err := am.client.SetFloating(nodeID, true); err != nil {
			return fmt.Errorf("setting floating state: %w", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if app.Size != "" {
		if err := am.client.ResizeWindow(nodeID, app.Size, app.Floating, layout); err != nil {
			return fmt.Errorf("setting size: %w", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if app.Position != "" && app.Floating {
		if err := am.client.MoveWindow(nodeID, app.Position); err != nil {
			return fmt.Errorf("setting position: %w", err)
		}
	}

	return nil
}
