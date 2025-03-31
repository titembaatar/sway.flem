package manager

import (
	"fmt"

	"github.com/titembaatar/sway.flem/internal/config"
)

// SetupContainerStructure recursively sets up containers and their contents
// It returns a list of apps that need to be resized
func (m *Manager) SetupContainerStructure(container *config.Container, parentLayout string, appsToResize *[]AppInfo) error {
	if container == nil || len(container.Apps) == 0 {
		return nil
	}

	// Launch and configure the first app (container representative)
	firstApp := container.Apps[0]
	nodeID, err := m.launchContainerRepresentative(firstApp)
	if err != nil {
		return err
	}

	// Handle container sizing
	m.handleContainerSizing(container, firstApp, nodeID, parentLayout, appsToResize)

	// Set the container layout
	if container.Layout != "" {
		m.logDebug("Setting container layout to: %s", container.Layout)
		layoutCmd := m.GetLayoutCommand(container.Layout)
		if err := m.Client.ExecuteCommand(layoutCmd); err != nil {
			m.logWarn("Failed to set container layout: %v", err)
		}
		m.delay(200)
	}

	// Launch subsequent apps
	m.launchSubsequentApps(container, appsToResize)

	// Handle nested containers recursively
	if container.Container != nil {
		currentLayout := container.Layout
		if currentLayout == "" {
			currentLayout = parentLayout
		}

		if err := m.SetupContainerStructure(container.Container, currentLayout, appsToResize); err != nil {
			m.logWarn("Issue with nested container: %v", err)
		}
	}

	// Run post-commands for all apps in this container
	m.runContainerPostCommands(container)

	return nil
}

// launchContainerRepresentative launches the first app in a container
func (m *Manager) launchContainerRepresentative(app config.App) (int64, error) {
	m.logDebug("Launching container representative app: %s", app.Name)
	nodeID, err := m.LaunchApp(app)
	if err != nil {
		return 0, fmt.Errorf("launching container representative app: %w", err)
	}
	m.delay(300)
	return nodeID, nil
}

// handleContainerSizing adds resize operations to the resize list
func (m *Manager) handleContainerSizing(
	container *config.Container,
	firstApp config.App,
	nodeID int64,
	parentLayout string,
	appsToResize *[]AppInfo,
) {
	// Add container-level sizing if specified
	if container.Size != "" {
		m.logDebug("Adding container representative to resize list with size %s (parent layout: %s)",
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

	// Perform any pending resize operations
	if len(*appsToResize) > 0 {
		m.logDebug("Resizing %d apps collected so far", len(*appsToResize))
		m.resizeAppBatch(*appsToResize)
		*appsToResize = []AppInfo{} // Clear the list after resizing
	}
}

// launchSubsequentApps launches additional apps in the container
func (m *Manager) launchSubsequentApps(container *config.Container, appsToResize *[]AppInfo) {
	for i := 1; i < len(container.Apps); i++ {
		app := container.Apps[i]
		m.logDebug("Launching subsequent app in container: %s", app.Name)

		nodeID, err := m.LaunchApp(app)
		if err != nil {
			m.logWarn("Failed to launch app %s: %v", app.Name, err)
			continue
		}
		m.delay(300)

		// Add to resize list
		if app.Size != "" {
			*appsToResize = append(*appsToResize, AppInfo{
				App:    app,
				NodeID: nodeID,
				Layout: container.Layout,
			})
		}
	}
}

// runContainerPostCommands executes post-commands for all apps in a container
func (m *Manager) runContainerPostCommands(container *config.Container) {
	for _, app := range container.Apps {
		if len(app.Posts) > 0 {
			if err := m.RunPostCommands(app); err != nil {
				m.logWarn("Failed to run post commands for %s: %v", app.Name, err)
			}
		}
	}
}

// resizeAppBatch applies resize operations to a batch of apps
func (m *Manager) resizeAppBatch(appsToResize []AppInfo) {
	for _, appInfo := range appsToResize {
		m.logDebug("Resizing app %s (ID: %d) with size %s using layout context %s",
			appInfo.App.Name, appInfo.NodeID, appInfo.App.Size, appInfo.Layout)

		if err := m.ResizeApp(appInfo.NodeID, appInfo.App, appInfo.Layout); err != nil {
			m.logWarn("Failed to resize app %s: %v", appInfo.App.Name, err)
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
