package sway

import "strings"

func isAppNode(node *Node) bool {
	return node.AppID != "" || (node.WindowProperties != nil && node.WindowProperties.Class != "")
}

func MatchAppName(runningApp string, configApp string) bool {
	runningLower := strings.ToLower(runningApp)
	configLower := strings.ToLower(configApp)

	return runningLower == configLower
}

func (n *Node) FindWorkspaces() map[int]*Node {
	workspaces := make(map[int]*Node)

	for i := range n.Nodes {
		output := &n.Nodes[i]
		if output.Type == "output" {
			for j := range output.Nodes {
				workspace := &output.Nodes[j]
				if workspace.Type == "workspace" && workspace.WorkspaceNum > 0 {
					workspace.Output = output.Name
					workspaces[workspace.WorkspaceNum] = workspace
				}
			}
		}
	}

	return workspaces
}

func (n *Node) FindAllApps() []AppNode {
	var apps []AppNode

	var processNodes func(node *Node, isFloating bool)

	processNodes = func(node *Node, isFloating bool) {
		if isAppNode(node) {
			var appName string
			if node.AppID != "" {
				// Wayland application
				appName = node.AppID
			} else if node.WindowProperties != nil && node.WindowProperties.Class != "" {
				// X11 application
				appName = node.WindowProperties.Class
			}

			apps = append(apps, AppNode{
				Name:     appName,
				NodeID:   node.ID,
				Rect:     node.Rect,
				Floating: isFloating,
			})
		}

		for i := range node.Nodes {
			childNode := &node.Nodes[i]
			processNodes(childNode, false)
		}

		for i := range node.FloatingNodes {
			floatingNode := &node.FloatingNodes[i]
			processNodes(floatingNode, true)
		}
	}

	processNodes(n, false)

	return apps
}

func ExtractAppOrder(repr string) []string {
	// Remove layout indicators
	simplified := strings.ReplaceAll(repr, "V[", "")
	simplified = strings.ReplaceAll(simplified, "H[", "")
	simplified = strings.ReplaceAll(simplified, "T[", "")
	simplified = strings.ReplaceAll(simplified, "S[", "")
	simplified = strings.ReplaceAll(simplified, "]", "")

	// Split the remaining string by spaces
	parts := strings.Fields(simplified)
	return parts
}

func ParseRepresentation(repr string) *LayoutNode {
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

	// Extract the content inside the brackets
	startBracket := strings.Index(repr, "[")
	endBracket := strings.LastIndex(repr, "]")

	if startBracket != -1 && endBracket != -1 && startBracket < endBracket {
		content := repr[startBracket+1 : endBracket]

		// Handle nested layouts
		if strings.Contains(content, "[") {
			// For simplicity, we're not fully recursing into nested structures
			// Just extract the app names
			node.Children = ExtractAppOrder(repr)
		} else {
			// Simple case - just app names
			node.Children = strings.Fields(content)
		}
	}

	return node
}

func GetLayoutTypeFromRepresentation(repr string) string {
	if len(repr) == 0 {
		return ""
	}

	switch repr[0] {
	case 'V':
		return string(LayoutTypeVertical)
	case 'H':
		return string(LayoutTypeHorizontal)
	case 'S':
		return string(LayoutTypeStacking)
	case 'T':
		return string(LayoutTypeTabbed)
	default:
		return ""
	}
}
