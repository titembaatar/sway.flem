package sway

import (
	"fmt"

	"github.com/titembaatar/sway.flem/internal/config"
	errs "github.com/titembaatar/sway.flem/internal/errors"
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

func (c *Container) SetLayout(errorHandler *errs.ErrorHandler) error {
	layout, err := types.ParseLayoutType(c.Layout)
	if err != nil {
		layoutErr := errs.New(errs.ErrInvalidLayoutType,
			fmt.Sprintf("'%s' is not a valid layout type for container %s", c.Layout, c.Mark.String()))
		layoutErr.WithCategory("Sway")
		layoutErr.WithSuggestion(errs.GetLayoutSuggestion())

		if errorHandler != nil {
			errorHandler.Handle(layoutErr)
		}

		return layoutErr
	}

	log.Debug("Setting layout '%s' for container %s", layout, c.Mark.String())

	commands := []string{layout.SplitCommand()}
	switch layout {
	case types.LayoutTabbed:
		commands = append(commands, types.LayoutTabbed.Command())
	case types.LayoutStacking:
		commands = append(commands, types.LayoutStacking.Command())
	}

	for _, command := range commands {
		cmd := NewSwayCmd(command)
		if errorHandler != nil {
			cmd.WithErrorHandler(errorHandler)
		}

		_, err := cmd.Run()
		if err != nil {
			cmdErr := errs.New(errs.ErrSetLayoutFailed,
				fmt.Sprintf("Failed to set layout '%s' for container %s", layout, c.Mark.String()))
			cmdErr.WithCategory("Sway")

			if errorHandler != nil {
				errorHandler.Handle(cmdErr)
			}

			return cmdErr
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
	errorHandler *errs.ErrorHandler,
) (*Container, error) {
	nextID := containerID + 1
	return buildContainerTree(workspaceName, containers, parentLayout, containerID, &nextID, errorHandler)
}

func buildContainerTree(
	workspaceName string,
	configContainers []config.Container,
	parentLayout string,
	currentContainerID int,
	nextContainerID *int,
	errorHandler *errs.ErrorHandler,
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
				warnMsg := fmt.Sprintf("Nested container definition for container ID %d is missing 'split' layout, inheriting parent layout '%s'",
					currentContainerID, parentLayout)

				if errorHandler != nil {
					layoutWarn := errs.NewWarning(errs.ErrMissingSplit, warnMsg)
					layoutWarn.WithCategory("Config")
					layoutWarn.WithSuggestion("Add a 'split' property to specify the layout for nested containers")
					errorHandler.Handle(layoutWarn)
				}

				log.Warn("Nested container definition for container ID %d is missing 'split' layout, inheriting parent layout '%s'",
					currentContainerID, parentLayout)
			}

			nestedID := *nextContainerID
			*nextContainerID++

			log.Debug("Recursing for nested container. CurrentID: %d, Assigned NestedID: %d",
				currentContainerID, nestedID)

			nestedContainer, err := buildContainerTree(
				workspaceName,
				cfgContainer.Containers,
				nestedLayout,
				nestedID,
				nextContainerID,
				errorHandler,
			)
			if err != nil {
				nestErr := errs.Wrap(err, fmt.Sprintf("Error building nested container (id %d) within container %d",
					nestedID, currentContainerID))

				if errorHandler != nil {
					errorHandler.Handle(nestErr)
				}

				return nil, nestErr
			}

			container.AddNestedContainer(nestedContainer)
			log.Debug("Added nested container with mark %s to container %d", nestedContainer.Mark.String(), currentContainerID)

		} else {
			invalidErr := errs.New(errs.ErrInvalidContainerStructure,
				fmt.Sprintf("Invalid container definition inside container ID %d (neither app nor nested structure)",
					currentContainerID))
			invalidErr.WithCategory("Config")
			invalidErr.WithSuggestion(errs.GetContainerStructureSuggestion())

			if errorHandler != nil {
				errorHandler.Handle(invalidErr)
			}

			log.Warn("Skipping invalid container definition inside container ID %d", currentContainerID)
		}
	}

	log.Debug("Finished processing container level ID %d. Apps: %d, Nested: %d",
		currentContainerID, len(container.Apps), len(container.Nested))
	return container, nil
}

func (c *Container) Setup(errorHandler *errs.ErrorHandler) error {
	log.Debug("Setting up container with mark %s and layout %s", c.Mark.String(), c.Layout)

	if len(c.Apps) > 0 {
		firstApp := c.Apps[0]
		log.Debug("Processing first app '%s' for container %s", firstApp.Name, c.Mark.String())
		if err := firstApp.Process(errorHandler); err != nil {
			appErr := errs.Wrap(err, fmt.Sprintf("Failed to process first app '%s' for container %s",
				firstApp.Name, c.Mark.String()))

			if errorHandler != nil {
				errorHandler.Handle(appErr)
			}

			return appErr
		}
	} else if len(c.Nested) == 0 {
		emptyWarn := errs.NewWarning(nil,
			fmt.Sprintf("Container %s has no apps and no nested containers to set up", c.Mark.String()))
		emptyWarn.WithCategory("Container")

		if errorHandler != nil {
			errorHandler.Handle(emptyWarn)
		}

		log.Warn("Container %s has no apps and no nested containers to set up", c.Mark.String())
		return nil
	}

	if len(c.Apps) > 0 || len(c.Nested) > 0 {
		if err := c.Mark.Apply(errorHandler); err != nil {
			markErr := errs.Wrap(err, fmt.Sprintf("Failed to apply mark %s", c.Mark.String()))

			if errorHandler != nil {
				errorHandler.Handle(markErr)
			}

			return markErr
		}

		if err := c.SetLayout(errorHandler); err != nil {
			layoutErr := errs.Wrap(err, fmt.Sprintf("Failed to set layout for container %s", c.Mark.String()))

			if errorHandler != nil {
				errorHandler.Handle(layoutErr)
			}

			return layoutErr
		}
	}

	for i, app := range c.Apps {
		if i > 0 {
			log.Debug("Processing subsequent app '%s' (%d/%d) for container %s",
				app.Name, i+1, len(c.Apps), c.Mark.String())
			if err := app.Process(errorHandler); err != nil {
				appErr := errs.Wrap(err, fmt.Sprintf("Failed to process app '%s' in container %s",
					app.Name, c.Mark.String()))

				if errorHandler != nil {
					errorHandler.Handle(appErr)
				} else {
					log.Error("Failed to process subsequent app '%s': %v", app.Name, err)
				}
			}
		}
	}

	for _, nested := range c.Nested {
		log.Debug("Recursively setting up nested container %s within %s", nested.Mark.String(), c.Mark.String())
		if err := nested.Setup(errorHandler); err != nil {
			nestedErr := errs.Wrap(err, fmt.Sprintf("Failed to set up nested container %s within %s",
				nested.Mark.String(), c.Mark.String()))

			if errorHandler != nil {
				errorHandler.Handle(nestedErr)
			} else {
				log.Error("Failed to set up nested container %s: %v", nested.Mark.String(), err)
			}
		}
	}

	log.Debug("Finished setting up container %s", c.Mark.String())
	return nil
}

func (c *Container) ResizeApps(errorHandler *errs.ErrorHandler) {
	log.Debug("Resizing apps in container %s", c.Mark.String())

	for _, app := range c.Apps {
		if app.Size != "" {
			if err := app.Resize(errorHandler); err != nil {
				if errorHandler != nil {
					resizeErr := errs.Wrap(err, fmt.Sprintf("Failed to resize app %s in container %s",
						app.Name, c.Mark.String()))
					errorHandler.Handle(resizeErr)
				} else {
					log.Warn("Failed to resize app %s: %v", app.Name, err)
				}
			}
		}
	}

	for _, nested := range c.Nested {
		nested.ResizeApps(errorHandler)
	}
}
