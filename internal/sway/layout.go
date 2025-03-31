package sway

import (
	"fmt"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
)

// Sets up a workspace with the specified layout
func SetupWorkspace(workspaceName string, workspace config.Workspace) error {
	log.Info("Setting up workspace: %s", workspaceName)

	if err := CreateWorkspace(workspaceName, workspace.Layout); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	workspaceMark := GenerateWorkspaceMark(workspaceName)

	if err := ApplyMark(workspaceMark); err != nil {
		return fmt.Errorf("failed to mark workspace: %w", err)
	}

	// Collect app information for all apps and containers for later resizing
	var allAppInfos []AppInfo

	// Launch top-level apps directly in the workspace
	if len(workspace.Apps) > 0 {
		log.Info("Launching %d top-level applications in workspace %s", len(workspace.Apps), workspaceName)
		appInfos, err := LaunchApps(workspaceMark, workspace.Layout, workspace.Apps)
		if err != nil {
			log.Error("Failed to launch top-level apps: %v", err)
			// Continue to set up the container even if some apps fail
		}
		allAppInfos = append(allAppInfos, appInfos...)
	}

	// Set up the container if present
	if workspace.Container != nil {
		containerMark := GenerateContainerMark(workspaceMark, "main")
		containerAppInfos, err := setupContainer(containerMark, workspace.Container, workspaceMark)
		if err != nil {
			log.Error("Failed to set up container: %v", err)
		}
		allAppInfos = append(allAppInfos, containerAppInfos...)
	}

	ResizeApps(allAppInfos)

	log.Info("Workspace %s setup complete", workspaceName)
	return nil
}

// Recursively sets up a container with the specified layout
// and returns app information for later resizing
func setupContainer(containerMark string, container *config.Container, parentMark string) ([]AppInfo, error) {
	log.Info("Setting up container with split: %s", container.Split)
	var allAppInfos []AppInfo

	if err := FocusMark(parentMark); err != nil {
		return nil, fmt.Errorf("failed to focus parent: %w", err)
	}

	if err := ApplyMark(containerMark); err != nil {
		return nil, fmt.Errorf("failed to mark container: %w", err)
	}

	// Store container info for later resizing
	if container.Size != "" {
		// Get the parent's layout to determine how to resize this container
		parentLayout := getLayoutForContainer(parentMark)
		allAppInfos = append(allAppInfos, AppInfo{
			Mark:   containerMark,
			Size:   container.Size,
			Layout: parentLayout,
		})
	}

	if len(container.Apps) > 0 {
		log.Info("Launching %d applications in container", len(container.Apps))
		appInfos, err := LaunchApps(containerMark, container.Split, container.Apps)
		if err != nil {
			log.Error("Failed to launch apps in container: %v", err)
			// Continue to set up nested container even if some apps fail
		}
		allAppInfos = append(allAppInfos, appInfos...)
	}

	if container.Container != nil {
		nestedMark := GenerateContainerMark(containerMark, "nested")
		nestedAppInfos, err := setupContainer(nestedMark, container.Container, containerMark)
		if err != nil {
			log.Error("Failed to set up nested container: %v", err)
		}
		allAppInfos = append(allAppInfos, nestedAppInfos...)
	}

	return allAppInfos, nil
}

// Determines the layout type of a container to know how to properly resize a child container
func getLayoutForContainer(containerMark string) string {
	defaultLayout := "splith"

	if err := FocusMark(containerMark); err != nil {
		log.Warn("Failed to focus container to determine layout, assuming %s: %v", defaultLayout, err)
		return defaultLayout
	}

	// TODO: In a future version, we can query the container's actual layout
	// using swaymsg -t get_tree and parsing the JSON output

	return defaultLayout
}

// Sets up the entire environment from the configuration
func SetupEnvironment(cfg *config.Config) error {
	log.Info("Setting up environment from configuration")

	for name, workspace := range cfg.Workspaces {
		log.Info("Processing workspace: %s", name)

		if err := SetupWorkspace(name, workspace); err != nil {
			log.Error("Failed to set up workspace %s: %v", name, err)
			// Continue with other workspaces even if one fails
			continue
		}
	}

	log.Info("Environment setup complete")
	return nil
}
