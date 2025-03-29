package manager

import (
	"sort"
	"strings"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func (m *Manager) GetAppPosition(positions map[string]int, appName string) int {
	pos, found := positions[strings.ToLower(appName)]
	if found {
		return pos
	}

	return -1
}

func (m *Manager) OrderAppsByLayout(apps []config.App, representation string) []config.App {
	if representation == "" {
		return apps
	}

	// Get app order from representation
	appOrder := sway.ExtractAppOrder(representation)
	if len(appOrder) == 0 {
		return apps
	}

	positions := make(map[string]int)
	for i, appName := range appOrder {
		positions[strings.ToLower(appName)] = i
	}

	orderedApps := make([]OrderedApp, 0, len(apps))
	for _, app := range apps {
		pos := m.GetAppPosition(positions, app.Name)
		orderedApps = append(orderedApps, OrderedApp{
			App:      app,
			Position: pos,
		})
	}

	// Sort the apps by their position
	sort.Slice(orderedApps, func(i, j int) bool {
		if orderedApps[i].Position >= 0 && orderedApps[j].Position >= 0 {
			return orderedApps[i].Position < orderedApps[j].Position
		}

		if orderedApps[i].Position >= 0 {
			return true
		}
		if orderedApps[j].Position >= 0 {
			return false
		}

		return i < j
	})

	result := make([]config.App, len(orderedApps))
	for i, ordered := range orderedApps {
		result[i] = ordered.App
	}

	return result
}

func (m *Manager) ParseLayoutRepresentation(repr string) *LayoutNode {
	if repr == "" {
		return nil
	}

	node := &LayoutNode{
		Type:     "",
		Children: []string{},
	}

	if len(repr) >= 2 {
		switch repr[0] {
		case 'V':
			node.Type = "V" // Vertical
		case 'H':
			node.Type = "H" // Horizontal
		case 'T':
			node.Type = "T" // Tabbed
		case 'S':
			node.Type = "S" // Stacking
		}
	}

	startBracket := strings.Index(repr, "[")
	endBracket := strings.LastIndex(repr, "]")

	if startBracket != -1 && endBracket != -1 && startBracket < endBracket {
		content := repr[startBracket+1 : endBracket]

		if strings.Contains(content, "[") {
			// For simplicity, we're not fully recursing into nested structures
			// Just extract the app names
			node.Children = sway.ExtractAppOrder(repr)
		} else {
			// Simple case - just app names
			node.Children = strings.Fields(content)
		}
	}

	return node
}

func (m *Manager) DetermineOptimalLayout(apps []config.App, defaultLayout string) string {
	if len(apps) == 0 {
		return defaultLayout
	}

	if len(apps) == 1 {
		// For a single app, default to splith
		return "splith"
	}

	floatingCount := 0
	for _, app := range apps {
		if app.Floating {
			floatingCount++
		}
	}

	// If all apps are floating, layout doesn't matter
	if floatingCount == len(apps) {
		return defaultLayout
	}

	return defaultLayout
}
