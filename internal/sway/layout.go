package sway

import (
	"fmt"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/pkg/types"
)

func SetupEnvironment(cfg *config.Config) error {
	log.Info("Setting up environment from configuration")

	for name, workspace := range cfg.Workspaces {
		log.Info("Processing workspace: %s", name)

		if err := SetupWorkspace(name, workspace); err != nil {
			log.Error("Failed to set up workspace %s: %v", name, err)

			continue
		}
	}

	log.Info("Environment setup complete")
	return nil
}

func SetupWorkspace(workspaceName string, workspace config.Workspace) error {
	log.Info("Setting up workspace: %s", workspaceName)

	if err := CreateWorkspace(workspaceName, workspace.Layout.String()); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	var resizeInfo []AppInfo

	info, err := setupContainer(workspaceName, workspace.Containers, workspace.Layout.String(), 0, 0)
	if err != nil {
		log.Error("Failed to setup container in workspace %s: %v", workspaceName, err)
	}
	resizeInfo = append(resizeInfo, info...)

	ResizeApps(resizeInfo)

	log.Info("Workspace %s setup complete", workspaceName)
	return nil
}

func setupContainer(
	workspaceName string,
	containers []config.Container,
	parentLayout string,
	containerID int,
	appID int,
) ([]AppInfo, error) {
	var resizeInfo []AppInfo

	if len(containers) == 0 {
		return resizeInfo, nil
	}

	if containers[0].App == "" && len(containers[0].Containers) > 0 {
		sizeInfo, err := setupContainer(
			workspaceName,
			containers[0].Containers,
			containers[0].Split.String(),
			containerID+1,
			0,
		)
		if err != nil {
			log.Error("Failed to setup nested container %d in workspace %s: %v",
				containerID+1, workspaceName, err)
		}
		resizeInfo = append(resizeInfo, sizeInfo...)
	}

	appMark := NewAppMark(workspaceName, containerID, appID)
	conMark := NewContainerMark(workspaceName, containerID)

	if containers[0].Size != "" {
		appSizeInfo := AppInfo{
			Mark:   appMark.String(),
			Size:   containers[0].Size,
			Layout: parentLayout,
		}
		resizeInfo = append(resizeInfo, appSizeInfo)
	}

	if err := LaunchApp(containers[0], appMark.String()); err != nil {
		log.Error("Failed to launch application %s in workspace %s: %v",
			containers[0].App,
			workspaceName,
			err)
	}

	if err := conMark.Apply(); err != nil {
		log.Error("Failed to apply mark %s in container %d in workspace %s: %v",
			conMark.String(),
			containerID,
			workspaceName,
			err)
	}

	if err := setContainerLayout(parentLayout); err != nil {
		log.Error("Failed to set layout %s for container %d in workspace %s: %v",
			parentLayout,
			containerID,
			workspaceName,
			err)
	}

	if len(containers) > 1 {
		info, err := processContainer(workspaceName, containers[1:], parentLayout, containerID)
		if err != nil {
			log.Error("Failed to process remaining containers in workspace %s: %v",
				workspaceName,
				err)
		}
		resizeInfo = append(resizeInfo, info...)
	}

	return resizeInfo, nil
}

func processContainer(
	workspaceName string,
	containers []config.Container,
	layout string,
	containerID int,
) ([]AppInfo, error) {
	var resizeInfo []AppInfo

	for i, container := range containers {
		if container.App != "" {
			appMark := NewAppMark(workspaceName, containerID, i+1)

			if container.Size != "" {
				appSizeInfo := AppInfo{
					Mark:   appMark.String(),
					Size:   container.Size,
					Layout: layout,
				}
				resizeInfo = append(resizeInfo, appSizeInfo)
			}

			if err := LaunchApp(container, appMark.String()); err != nil {
				log.Error("Failed to launch application %s in workspace %s: %v",
					container.App, workspaceName, err)
			}
		}
	}

	lastContainer := containers[len(containers)-1]
	if lastContainer.Split != "" && len(lastContainer.Containers) > 0 {
		nestedInfo, err := setupContainer(
			workspaceName,
			lastContainer.Containers,
			lastContainer.Split.String(),
			containerID+1,
			0,
		)
		if err != nil {
			log.Error("Failed to setup nested container: %v", err)
		}
		resizeInfo = append(resizeInfo, nestedInfo...)
	}

	return resizeInfo, nil
}

func setContainerLayout(layoutType string) error {
	layout, err := types.ParseLayoutType(layoutType)
	if err != nil {
		return fmt.Errorf("%w: '%s' is not a valid layout type", ErrInvalidLayout, layoutType)
	}

	commands := []string{layout.SplitCommand()}

	switch layout {
	case types.LayoutTabbed:
		commands = append(commands, types.LayoutTabbed.Command())
	case types.LayoutStacking:
		commands = append(commands, types.LayoutStacking.Command())
	default:
	}

	for _, command := range commands {
		if _, err := RunCommand(command); err != nil {
			return err
		}
	}

	return nil
}
