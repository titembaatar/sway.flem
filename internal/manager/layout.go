package manager

import (
	"strings"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

// ParseLayoutRepresentation converts a Sway layout representation string to a LayoutNode
func (m *Manager) ParseLayoutRepresentation(repr string) *LayoutNode {
	if repr == "" {
		return nil
	}

	node := &LayoutNode{
		Type:     "",
		Children: []string{},
	}

	// Determine layout type from first character
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

	// Extract content within brackets
	startBracket := strings.Index(repr, "[")
	endBracket := strings.LastIndex(repr, "]")

	if startBracket != -1 && endBracket != -1 && startBracket < endBracket {
		content := repr[startBracket+1 : endBracket]

		if strings.Contains(content, "[") {
			// Complex nested structure - extract app names only
			node.Children = sway.ExtractAppOrder(repr)
		} else {
			// Simple case - just app names
			node.Children = strings.Fields(content)
		}
	}

	return node
}

// DetermineOptimalLayout selects the best layout based on the apps and default
func (m *Manager) DetermineOptimalLayout(apps []config.App, defaultLayout string) string {
	if len(apps) == 0 {
		return defaultLayout
	}

	// Single app doesn't need a complex layout
	if len(apps) == 1 {
		return "splith"
	}

	// Otherwise use the specified default
	return defaultLayout
}

// GetLayoutFromName converts a layout name to its Sway command form
// Supports both full names and shorthand versions
func (m *Manager) GetLayoutFromName(name string) string {
	switch strings.ToLower(name) {
	case "vertical", "splitv", "v":
		return "splitv"
	case "horizontal", "splith", "h":
		return "splith"
	case "tabbed", "t":
		return "tabbed"
	case "stacking", "stack", "s":
		return "stacking"
	default:
		return "splith" // Default to horizontal split
	}
}

// IsTabOrStack checks if a layout is tabbed or stacking
func (m *Manager) IsTabOrStack(layout string) bool {
	normalizedLayout := m.GetLayoutFromName(layout)
	return normalizedLayout == "tabbed" || normalizedLayout == "stacking"
}

// GetLayoutCommand returns the Sway command to set a specific layout
func (m *Manager) GetLayoutCommand(layout string) string {
	normalizedLayout := m.GetLayoutFromName(layout)
	return "layout " + normalizedLayout
}
