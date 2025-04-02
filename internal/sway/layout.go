package sway

import (
	"fmt"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
)

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

// Sets up a workspace with the specified layout
func SetupWorkspace(workspaceName string, workspace config.Workspace) error {
	log.Info("Setting up workspace: %s", workspaceName)

	if err := CreateWorkspace(workspaceName, workspace.Layout); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// Collect all applications that will need to be resized
	var resizeInfo []AppInfo

	// Process all top-level containers
	if len(workspace.Containers) > 0 {
		containerInfo, err := processContainers(workspaceName, workspace.Layout, workspace.Containers, 0)
		if err != nil {
			log.Error("Failed to process workspace containers: %v", err)
		} else {
			resizeInfo = append(resizeInfo, containerInfo...)
		}
	}

	// Perform all resize operations
	ResizeApps(resizeInfo)

	log.Info("Workspace %s setup complete", workspaceName)
	return nil
}

// Process a list of containers at the same level
func processContainers(workspaceName, parentLayout string, containers []config.Container, depth int) ([]AppInfo, error) {
	log.Info("Processing %d containers at depth %d", len(containers), depth)

	var resizeInfo []AppInfo

	// We need to keep track of container IDs at each depth level
	containerCounter := depth

	for i, container := range containers {
		// Determine if this is an app or a nested container
		isApp := container.App != ""

		if isApp {
			// This is an app container
			appMark := generateAppMark(workspaceName, depth, containerCounter, i)

			// Launch the app
			app := config.Container{
				App:   container.App,
				Cmd:   container.Cmd,
				Size:  container.Size,
				Delay: container.Delay,
			}

			if err := LaunchApp(app, appMark); err != nil {
				log.Error("Failed to launch app %s: %v", container.App, err)
				continue
			}

			// Add resize info
			resizeInfo = append(resizeInfo, AppInfo{
				Mark:   appMark,
				Size:   container.Size,
				Layout: parentLayout,
			})
		} else {
			// We need at least one app in the container to create it
			if len(container.Containers) == 0 {
				log.Error("Container has no child containers, skipping")
				continue
			}

			// Process the first child to establish the container
			firstChild := container.Containers[0]
			containerMark := generateContainerMark(workspaceName, containerCounter)
			containerCounter++

			// First child is an app
			if firstChild.App != "" {
				firstAppMark := generateAppMark(workspaceName, depth+1, containerCounter-1, 0)

				// Launch the first app
				app := config.Container{
					App:  firstChild.App,
					Cmd:  firstChild.Cmd,
					Size: firstChild.Size,
				}

				if err := LaunchApp(app, firstAppMark); err != nil {
					log.Error("Failed to launch container's first app: %v", err)
					continue
				}

				// Also mark this app as the container
				if err := ApplyMark(containerMark); err != nil {
					log.Warn("Failed to apply container mark: %v", err)
				}

				// Set the container's layout
				splitCmd := fmt.Sprintf("split %s", container.Split)
				if _, err := RunCommand(splitCmd); err != nil {
					log.Warn("Failed to set split layout: %v", err)
					// Continue anyway as this might be expected
				}

				// Add container resize info (relative to parent)
				resizeInfo = append(resizeInfo, AppInfo{
					Mark:   containerMark,
					Size:   container.Size,
					Layout: parentLayout,
				})

				// Add app resize info (within its container)
				resizeInfo = append(resizeInfo, AppInfo{
					Mark:   firstAppMark,
					Size:   firstChild.Size,
					Layout: container.Split,
				})

				// Process remaining children in this container
				if len(container.Containers) > 1 {
					// Focus the container
					focusCmd := fmt.Sprintf("[con_mark=\"%s\"] focus", containerMark)
					if _, err := RunCommand(focusCmd); err != nil {
						log.Error("Failed to focus container: %v", err)
						continue
					}

					// Process remaining containers (skip the first one)
					childInfo, err := processContainers(
						workspaceName,
						container.Split,
						container.Containers[1:],
						depth+1,
					)
					if err != nil {
						log.Error("Failed to process child containers: %v", err)
					} else {
						resizeInfo = append(resizeInfo, childInfo...)
					}
				}
			} else {
				// First child is another container - recursively process it
				log.Warn("First child of container is not an app but another container - this might cause layout issues")

				// Just process all children and hope for the best
				childInfo, err := processContainers(
					workspaceName,
					container.Split,
					container.Containers,
					depth+1,
				)
				if err != nil {
					log.Error("Failed to process child containers: %v", err)
				} else {
					resizeInfo = append(resizeInfo, childInfo...)
				}
			}
		}
	}

	return resizeInfo, nil
}
