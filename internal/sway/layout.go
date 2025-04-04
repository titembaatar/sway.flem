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
	ws := NewWorkspace(workspaceName, workspace.Layout.String())

	ws.Create()

	var apps []App

	info, err := setupContainer(workspaceName, workspace.Containers, workspace.Layout.String(), 0, 0)
	if err != nil {
		log.Error("Failed to setup container in workspace %s: %v", workspaceName, err)
	}
	apps = append(apps, info...)

	for _, app := range apps {
		app.Resize()
	}

	log.Info("Workspace %s setup complete", workspaceName)
	return nil
}

func setupContainer(
	workspaceName string,
	containers []config.Container,
	parentLayout string,
	containerID int,
	appID int,
) ([]App, error) {
	var apps []App

	if len(containers) == 0 {
		return apps, nil
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
		apps = append(apps, sizeInfo...)
	}

	appMark := NewAppMark(workspaceName, containerID, appID)
	conMark := NewContainerMark(workspaceName, containerID)

	app := App{
		Name:    containers[0].App,
		Command: containers[0].Cmd,
		Mark:    appMark,
		Size:    containers[0].Size,
		Delay:   containers[0].Delay,
		Layout:  parentLayout,
		Post:    containers[0].Post,
	}

	if containers[0].Size != "" {
		apps = append(apps, app)
	}

	if err := app.Launch(); err != nil {
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
		apps = append(apps, info...)
	}

	return apps, nil
}

func processContainer(
	workspaceName string,
	containers []config.Container,
	layout string,
	containerID int,
) ([]App, error) {
	var apps []App

	for i, container := range containers {
		if container.App != "" {
			appMark := NewAppMark(workspaceName, containerID, i+1)

			app := App{
				Name:    containers[0].App,
				Command: containers[0].Cmd,
				Mark:    appMark,
				Size:    containers[0].Size,
				Delay:   containers[0].Delay,
				Layout:  layout,
				Post:    containers[0].Post,
			}

			if container.Size != "" {
				apps = append(apps, app)
			}

			if err := app.Launch(); err != nil {
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
		apps = append(apps, nestedInfo...)
	}

	return apps, nil
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
		cmd := NewSwayCmd(command)
		cmd.Run()
	}

	return nil
}
