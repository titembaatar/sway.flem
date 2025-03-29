package sway

import "strings"

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
	collectApps(n, false, &apps)
	return apps
}

func collectApps(node *Node, isFloating bool, apps *[]AppNode) {
	if isAppNode(node) {
		var appName string
		if node.AppID != "" {
			// Wayland application
			appName = node.AppID
		} else if node.WindowProperties != nil && node.WindowProperties.Class != "" {
			// X11 application
			appName = node.WindowProperties.Class
		}

		*apps = append(*apps, AppNode{
			Name:     appName,
			NodeID:   node.ID,
			Rect:     node.Rect,
			Floating: isFloating,
		})
	}

	// Process regular nodes
	for i := range node.Nodes {
		childNode := &node.Nodes[i]
		collectApps(childNode, false, apps)
	}

	// Process floating nodes
	for i := range node.FloatingNodes {
		floatingNode := &node.FloatingNodes[i]
		collectApps(floatingNode, true, apps)
	}
}

func isAppNode(node *Node) bool {
	return node.AppID != "" || (node.WindowProperties != nil && node.WindowProperties.Class != "")
}

func MatchAppName(runningApp string, configApp string) bool {
	runningLower := strings.ToLower(runningApp)
	configLower := strings.ToLower(configApp)

	return runningLower == configLower
}

func ExtractAppOrder(repr string) []string {
	simplified := strings.ReplaceAll(repr, "V[", "")
	simplified = strings.ReplaceAll(simplified, "H[", "")
	simplified = strings.ReplaceAll(simplified, "T[", "")
	simplified = strings.ReplaceAll(simplified, "S[", "")
	simplified = strings.ReplaceAll(simplified, "]", "")

	return strings.Fields(simplified)
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
