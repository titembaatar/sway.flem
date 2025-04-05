package sway

import (
	"fmt"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/pkg/types"
)

type Container struct {
	Mark   Mark
	Layout string
	Apps   []*App
	Parent *Container
	Nested []*Container
}

func NewContainer(markID string, layout string) *Container {
	return &Container{
		Mark:   NewMark(markID),
		Layout: layout,
		Apps:   make([]*App, 0),
		Nested: make([]*Container, 0),
	}
}

func NewWorkspaceContainer(workspaceName string, containerID int, layout string) *Container {
	markID := fmt.Sprintf("w%s_c%d", workspaceName, containerID)
	return NewContainer(markID, layout)
}

func (c *Container) AddApp(app *App) {
	app.Layout = c.Layout
	c.Apps = append(c.Apps, app)
}

func (c *Container) AddNestedContainer(container *Container) {
	container.Parent = c
	c.Nested = append(c.Nested, container)
}

func (c *Container) SetLayout() error {
	layout, err := types.ParseLayoutType(c.Layout)
	if err != nil {
		return fmt.Errorf("%w: '%s' is not a valid layout type", ErrInvalidLayout, c.Layout)
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
		if _, err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Container) GetAllApps() []*App {
	apps := make([]*App, len(c.Apps))
	copy(apps, c.Apps)

	for _, nested := range c.Nested {
		apps = append(apps, nested.GetAllApps()...)
	}

	return apps
}

func ProcessContainer(
	workspaceName string,
	containers []config.Container,
	parentLayout string,
	containerID int,
) (*Container, error) {
	container := NewWorkspaceContainer(workspaceName, containerID, parentLayout)

	if len(containers) == 0 {
		return container, nil
	}

	firstContainer := containers[0]

	if firstContainer.App == "" && len(firstContainer.Containers) > 0 {
		nestedLayout := firstContainer.Split.String()
		nestedContainer, err := ProcessContainer(
			workspaceName,
			firstContainer.Containers,
			nestedLayout,
			containerID+1,
		)

		if err != nil {
			log.Error("Failed to process nested container: %v", err)
		} else {
			container.AddNestedContainer(nestedContainer)
		}
	} else if firstContainer.App != "" {
		appMark := NewAppMark(workspaceName, containerID, 0)
		app := NewApp(firstContainer, appMark.String())
		container.AddApp(app)
	}

	if len(containers) > 1 {
		for i, ctn := range containers[1:] {
			if ctn.App != "" {
				appMark := NewAppMark(workspaceName, containerID, i+1)
				app := NewApp(ctn, appMark.String())
				container.AddApp(app)
			}
		}

		lastContainer := containers[len(containers)-1]
		if lastContainer.Split != "" && len(lastContainer.Containers) > 0 {
			nestedLayout := lastContainer.Split.String()
			nestedContainer, err := ProcessContainer(
				workspaceName,
				lastContainer.Containers,
				nestedLayout,
				containerID+1,
			)

			if err != nil {
				log.Error("Failed to process nested container: %v", err)
			} else {
				container.AddNestedContainer(nestedContainer)
			}
		}
	}

	return container, nil
}

func (c *Container) Setup() error {
	log.Debug("Setting up container with mark %s and layout %s", c.Mark.String(), c.Layout)

	if len(c.Apps) > 0 {
		firstApp := c.Apps[0]
		if err := firstApp.Process(); err != nil {
			log.Error("Failed to process first app %s: %v", firstApp.Name, err)
			return err
		}

		if err := c.Mark.Apply(); err != nil {
			log.Error("Failed to apply mark %s to container: %v", c.Mark.String(), err)
			return err
		}

		if err := c.SetLayout(); err != nil {
			log.Error("Failed to set layout %s for container: %v", c.Layout, err)
			return err
		}

		for i, app := range c.Apps {
			if i > 0 {
				if err := app.Process(); err != nil {
					log.Error("Failed to process app %s: %v", app.Name, err)
				}
			}
		}
	}

	for _, nested := range c.Nested {
		if err := nested.Setup(); err != nil {
			log.Error("Failed to set up nested container: %v", err)
		}
	}

	return nil
}

func (c *Container) ResizeApps() {
	for _, app := range c.Apps {
		if app.Size != "" {
			if err := app.Resize(); err != nil {
				log.Warn("Failed to resize app %s: %v", app.Name, err)
			}
		}
	}

	for _, nested := range c.Nested {
		nested.ResizeApps()
	}
}
