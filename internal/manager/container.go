package manager

import (
	"fmt"
	"log"

	"github.com/titembaatar/sway.flem/internal/config"
)

// SetupContainerStructure recursively sets up containers and their contents
// It returns a list of apps that need to be resized
func (m *Manager) SetupContainerStructure(container *config.Container, parentLayout string, appsToResize *[]AppInfo) error {
	if container == nil || len(container.Apps) == 0 {
		return nil
	}

	// STEP 1: Launch the first app in the container (container representative)
	firstApp := container.Apps[0]
	m.logVerbose("Launching first app in container (container representative): %s", firstApp.Name)
	nodeID, err := m.LaunchApp(firstApp)
	if err != nil {
		return fmt.Errorf("launching container representative app: %w", err)
	}
	m.delay(300)

	// STEP 2: Add first app to resize list with container size (using parent layout context)
	if container.Size != "" {
		m.logVerbose("Adding container representative to resize list with size %s (parent layout: %s)",
			container.Size, parentLayout)
		// Create a temporary app with container's size for the resize list
		containerApp := config.App{
			Name: firstApp.Name,
			Size: container.Size,
		}
		*appsToResize = append(*appsToResize, AppInfo{
			App:    containerApp,
			NodeID: nodeID,
			Layout: parentLayout, // Use parent layout for container sizing
		})
	}

	// Also add the app itself with its own size (using container layout)
	if firstApp.Size != "" {
		// Regular app resize with app's own size
		*appsToResize = append(*appsToResize, AppInfo{
			App:    firstApp,
			NodeID: nodeID,
			Layout: container.Layout,
		})
	}

	// STEP 3: Resize all apps collected so far
	m.logVerbose("Resizing all apps collected so far (%d apps)", len(*appsToResize))
	m.resizeAppBatch(*appsToResize)

	// Clear the list of apps to resize since we've resized them
	*appsToResize = []AppInfo{}

	// STEP 4: Set the container layout
	if container.Layout != "" {
		m.logVerbose("Setting container layout to: %s", container.Layout)
		if err := m.Client.ExecuteCommand(container.Layout); err != nil {
			log.Printf("Warning: Failed to set container layout: %v", err)
		}
		m.delay(200)
	}

	// STEP 5: Launch and set floating state for remaining apps in this container
	for i := 1; i < len(container.Apps); i++ {
		app := container.Apps[i]
		m.logVerbose("Launching subsequent app in container: %s", app.Name)
		nodeID, err := m.LaunchApp(app)
		if err != nil {
			log.Printf("Warning: Failed to launch app %s: %v", app.Name, err)
			continue
		}
		m.delay(300)

		// Set floating state immediately if needed
		if app.Floating {
			if err := m.SetFloatingState(nodeID, app); err != nil {
				log.Printf("Warning: Failed to set floating state for %s: %v", app.Name, err)
			}
			m.delay(100)
		}

		// Add to resize list
		if app.Size != "" {
			*appsToResize = append(*appsToResize, AppInfo{
				App:    app,
				NodeID: nodeID,
				Layout: container.Layout,
			})
		}
	}

	// STEP 6: Recursively handle nested container
	if container.Container != nil {
		currentLayout := container.Layout
		if currentLayout == "" {
			currentLayout = parentLayout
		}

		err := m.SetupContainerStructure(container.Container, currentLayout, appsToResize)
		if err != nil {
			log.Printf("Warning: Issue with nested container: %v", err)
		}
	}

	// STEP 7: Run post-commands for all apps in this container
	for _, app := range container.Apps {
		if len(app.Posts) > 0 {
			if err := m.RunPostCommands(app); err != nil {
				log.Printf("Warning: Failed to run post commands for %s: %v", app.Name, err)
			}
		}
	}

	return nil
}

// resizeAppBatch applies resize operations to a batch of apps
func (m *Manager) resizeAppBatch(appsToResize []AppInfo) {
	for _, appInfo := range appsToResize {
		m.logVerbose("Resizing app %s (ID: %d) with size %s using layout context %s",
			appInfo.App.Name, appInfo.NodeID, appInfo.App.Size, appInfo.Layout)

		if err := m.ResizeApp(appInfo.NodeID, appInfo.App, appInfo.Layout); err != nil {
			log.Printf("Warning: Failed to resize app %s: %v", appInfo.App.Name, err)
		}
		m.delay(200)
	}
}

// CollectAllApps gathers all apps from a workspace and its containers
func (m *Manager) CollectAllApps(workspace config.Workspace) []config.App {
	var allApps []config.App

	// Add apps directly in the workspace
	allApps = append(allApps, workspace.Apps...)

	// Add apps from containers recursively
	m.collectAppsFromContainer(workspace.Container, &allApps)

	return allApps
}

// collectAppsFromContainer recursively collects apps from a container
func (m *Manager) collectAppsFromContainer(container *config.Container, allApps *[]config.App) {
	if container == nil {
		return
	}

	// Add apps from this container
	*allApps = append(*allApps, container.Apps...)

	// Recursively collect from nested container
	m.collectAppsFromContainer(container.Container, allApps)
}
