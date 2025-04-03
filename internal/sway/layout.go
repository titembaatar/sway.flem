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

	var resizeInfo []AppInfo

	if len(workspace.Containers) > 0 {
		containerInfo, err := processContainers(workspaceName, workspace.Layout, workspace.Containers, 0)
		if err != nil {
			log.Error("Failed to process workspace containers: %v", err)
		} else {
			resizeInfo = append(resizeInfo, containerInfo...)
		}
	}

	ResizeApps(resizeInfo)

	log.Info("Workspace %s setup complete", workspaceName)
	return nil
}

// Processes a list of containers at the same level
func processContainers(workspaceName, parentLayout string, containers []config.Container, depth int) ([]AppInfo, error) {
	log.Info("Processing %d containers at depth %d", len(containers), depth)

	var resizeInfo []AppInfo

	containerCounter := depth

	for i, container := range containers {
		isApp := container.App != ""

		if isApp {
			appInfo, err := processAppContainer(workspaceName, parentLayout, container, depth, containerCounter, i)
			if err != nil {
				log.Error("Failed to process app container %s: %v", container.App, err)
				continue
			}

			if appInfo.Size != "" {
				resizeInfo = append(resizeInfo, appInfo)
			}
		} else {
			containerResizeInfo, newContainerID, err := processNestedContainer(
				workspaceName,
				parentLayout,
				container,
				depth,
				containerCounter,
			)

			if err != nil {
				log.Error("Failed to process nested container: %v", err)
				continue
			}

			containerCounter = newContainerID
			resizeInfo = append(resizeInfo, containerResizeInfo...)
		}
	}

	return resizeInfo, nil
}

// Handles a single application container
func processAppContainer(
	workspaceName string,
	parentLayout string,
	container config.Container,
	depth int,
	containerID int,
	index int,
) (AppInfo, error) {
	mark := NewAppMark(workspaceName, depth, containerID, index)

	app := config.Container{
		App:   container.App,
		Cmd:   container.Cmd,
		Size:  container.Size,
		Delay: container.Delay,
		Post:  container.Post,
	}

	if err := LaunchApp(app, mark.String()); err != nil {
		return AppInfo{}, fmt.Errorf("failed to launch app %s: %w", container.App, err)
	}

	return AppInfo{
		Mark:   mark.String(),
		Size:   container.Size,
		Layout: parentLayout,
	}, nil
}

// Handles a container with child containers
func processNestedContainer(
	workspaceName string,
	parentLayout string,
	container config.Container,
	depth int,
	containerID int,
) ([]AppInfo, int, error) {
	var resizeInfo []AppInfo

	if len(container.Containers) == 0 {
		return nil, containerID, fmt.Errorf("container has no child containers")
	}

	firstChild := container.Containers[0]
	containerMark := NewContainerMark(workspaceName, containerID)
	newContainerID := containerID + 1

	if firstChild.App != "" {
		firstChildInfo, err := setupContainerWithApp(
			workspaceName,
			parentLayout,
			container,
			firstChild,
			depth,
			containerID,
			containerMark.String(),
		)

		if err != nil {
			return nil, newContainerID, fmt.Errorf("failed to setup container with app: %w", err)
		}

		resizeInfo = append(resizeInfo, firstChildInfo...)

		if len(container.Containers) > 1 {
			if err := containerMark.Focus(); err != nil {
				return resizeInfo, newContainerID, fmt.Errorf("failed to focus container: %w", err)
			}

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
		log.Warn("First child of container is not an app but another container - this might cause layout issues")

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

	return resizeInfo, newContainerID, nil
}

// Sets up a container by creating its first app and setting the layout
func setupContainerWithApp(
	workspaceName string,
	parentLayout string,
	container config.Container,
	firstChild config.Container,
	depth int,
	containerID int,
	containerMark string,
) ([]AppInfo, error) {
	var resizeInfo []AppInfo

	firstAppMark := NewAppMark(workspaceName, depth+1, containerID, 0).String()

	app := config.Container{
		App:   firstChild.App,
		Cmd:   firstChild.Cmd,
		Size:  firstChild.Size,
		Delay: firstChild.Delay,
		Post:  firstChild.Post,
	}

	if err := LaunchApp(app, firstAppMark); err != nil {
		return nil, fmt.Errorf("failed to launch container's first app: %w", err)
	}

	containerMarkObj := NewMark(containerMark)
	if err := containerMarkObj.Apply(); err != nil {
		log.Warn("Failed to apply container mark: %v", err)
	}

	if err := setContainerLayout(container.Split); err != nil {
		log.Warn("Failed to set container layout: %v", err)
	}

	if container.Size != "" {
		resizeInfo = append(resizeInfo, AppInfo{
			Mark:   containerMark,
			Size:   container.Size,
			Layout: parentLayout,
		})
	}

	if firstChild.Size != "" {
		resizeInfo = append(resizeInfo, AppInfo{
			Mark:   firstAppMark,
			Size:   firstChild.Size,
			Layout: container.Split,
		})
	}

	return resizeInfo, nil
}

// Applies the specified layout to the current container
func setContainerLayout(layoutType string) error {
	splitCmd := fmt.Sprintf("split %s", layoutType)
	_, err := RunCommand(splitCmd)
	return err
}
