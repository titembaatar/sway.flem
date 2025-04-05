package sway

import (
	"fmt"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
	"github.com/titembaatar/sway.flem/pkg/types"
)

type Container struct {
	Mark   Mark         // Unique mark applied to this container in Sway.
	Layout string       // Layout type for this container (e.g., "splith", "tabbed").
	Apps   []*App       // List of applications directly within this container.
	Parent *Container   // Reference to the parent container, if any.
	Nested []*Container // List of nested containers within this container.
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
		return fmt.Errorf("%w: '%s' is not a valid layout type for container mark %s", ErrInvalidLayout, c.Layout, c.Mark.String())
	}

	log.Debug("Setting layout '%s' for container mark '%s'", layout, c.Mark.String())

	commands := []string{layout.SplitCommand()}
	switch layout {
	case types.LayoutTabbed:
		commands = append(commands, types.LayoutTabbed.Command())
	case types.LayoutStacking:
		commands = append(commands, types.LayoutStacking.Command())
	}

	for _, command := range commands {
		cmd := NewSwayCmd(command)
		_, err := cmd.Run()
		if err != nil {
			return fmt.Errorf("%w: failed to execute command '%s' for container mark %s: %v", ErrSetLayoutFailed, command, c.Mark.String(), err)
		}
	}

	log.Debug("Successfully set layout '%s' for container mark '%s'", layout, c.Mark.String())
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
	nextID := containerID + 1
	return buildContainerTree(workspaceName, containers, parentLayout, containerID, &nextID)
}

func buildContainerTree(
	workspaceName string,
	configContainers []config.Container,
	parentLayout string,
	currentContainerID int,
	nextContainerID *int,
) (*Container, error) {

	container := NewWorkspaceContainer(workspaceName, currentContainerID, parentLayout)
	log.Debug("Processing container level with ID %d and layout %s", currentContainerID, parentLayout)

	appIndexCounter := 0

	for _, cfgContainer := range configContainers {
		isApp := cfgContainer.App != ""
		isNested := !isApp && len(cfgContainer.Containers) > 0

		if isApp {
			appMark := NewAppMark(workspaceName, currentContainerID, appIndexCounter)
			app := NewApp(cfgContainer, appMark.String())
			container.AddApp(app)
			log.Debug("Added app '%s' with mark %s to container %d", app.Name, app.Mark.String(), currentContainerID)
			appIndexCounter++ // Increment index for the next app in *this* container.

		} else if isNested {
			nestedLayout := cfgContainer.Split.String()
			if nestedLayout == "" {
				nestedLayout = parentLayout
				log.Warn("Nested container definition for container ID %d is missing 'split' layout, inheriting parent layout '%s'", currentContainerID, parentLayout)
			}

			nestedID := *nextContainerID
			*nextContainerID++

			log.Debug("Recursing for nested container. CurrentID: %d, Assigned NestedID: %d, Next Available ID: %d", currentContainerID, nestedID, *nextContainerID)

			nestedContainer, err := buildContainerTree(
				workspaceName,
				cfgContainer.Containers,
				nestedLayout,
				nestedID,
				nextContainerID,
			)
			if err != nil {
				return nil, fmt.Errorf("error building nested container (id %d) within container %d: %w", nestedID, currentContainerID, err)
			}

			container.AddNestedContainer(nestedContainer)
			log.Debug("Added nested container with mark %s to container %d", nestedContainer.Mark.String(), currentContainerID)

		} else {
			log.Warn("Skipping invalid container definition inside container ID %d (neither app nor nested structure)", currentContainerID)
		}
	}

	log.Debug("Finished processing container level ID %d. Apps: %d, Nested: %d", currentContainerID, len(container.Apps), len(container.Nested))
	return container, nil
}

func (c *Container) Setup() error {
	log.Debug("Setting up container with mark %s and layout %s", c.Mark.String(), c.Layout)

	if len(c.Apps) > 0 {
		firstApp := c.Apps[0]
		log.Debug("Processing first app '%s' for container %s", firstApp.Name, c.Mark.String())
		if err := firstApp.Process(); err != nil {
			log.Error("Failed to process first app '%s' in container %s: %v", firstApp.Name, c.Mark.String(), err)
			return fmt.Errorf("failed to process first app '%s' for container %s: %w", firstApp.Name, c.Mark.String(), err)
		}
	} else if len(c.Nested) == 0 {
		log.Warn("Container %s has no apps and no nested containers to set up.", c.Mark.String())
		return nil
	}

	if len(c.Apps) > 0 || len(c.Nested) > 0 {
		log.Debug("Applying mark %s to container", c.Mark.String())
		if err := c.Mark.Apply(); err != nil {
			return fmt.Errorf("failed to apply mark %s: %w", c.Mark.String(), err)
		}

		log.Debug("Setting layout %s for container %s", c.Layout, c.Mark.String())
		if err := c.SetLayout(); err != nil {
			return fmt.Errorf("failed to set layout for container %s: %w", c.Mark.String(), err)
		}
	}

	for i, app := range c.Apps {
		if i > 0 {
			log.Debug("Processing subsequent app '%s' (%d/%d) for container %s", app.Name, i+1, len(c.Apps), c.Mark.String())
			if err := app.Process(); err != nil {
				log.Error("Failed to process subsequent app '%s' in container %s: %v", app.Name, c.Mark.String(), err)
			}
		}
	}

	for _, nested := range c.Nested {
		log.Debug("Recursively setting up nested container %s within %s", nested.Mark.String(), c.Mark.String())
		if err := nested.Setup(); err != nil {
			log.Error("Failed to set up nested container %s within %s: %v", nested.Mark.String(), c.Mark.String(), err)
		}
	}

	log.Debug("Finished setting up container %s", c.Mark.String())
	return nil
}

func (c *Container) ResizeApps() {
	log.Debug("Resizing apps in container %s", c.Mark.String())
	for _, app := range c.Apps {
		if app.Size != "" {
			if err := app.Resize(); err != nil {
				log.Warn("Failed to resize app %s (mark %s) in container %s: %v", app.Name, app.Mark.String(), c.Mark.String(), err)
			}
		}
	}

	for _, nested := range c.Nested {
		nested.ResizeApps()
	}
}
