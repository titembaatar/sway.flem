package sway

// Node in the Sway tree
type Node struct {
	ID               int64        `json:"id"`
	Name             string       `json:"name"`
	Type             string       `json:"type"`
	AppID            string       `json:"app_id"`
	WorkspaceNum     int          `json:"num,omitempty"`
	Layout           string       `json:"layout"`
	WindowProperties *WindowProps `json:"window_properties"`
	Nodes            []Node       `json:"nodes"`
	FloatingNodes    []Node       `json:"floating_nodes"`
	Output           string       `json:"output,omitempty"`
}

// Window metadata
type WindowProps struct {
	Class string `json:"class"`
}

// Returns all workspaces in the tree
func (n *Node) FindWorkspaces() map[int]*Node {
	workspaces := make(map[int]*Node)

	if n.Type == "root" {
		// Search in outputs
		for i := range n.Nodes {
			output := &n.Nodes[i]
			if output.Type == "output" && output.Name != "__i3" {
				// Search for workspaces in this output
				for j := range output.Nodes {
					workspace := &output.Nodes[j]
					if workspace.Type == "workspace" && workspace.WorkspaceNum > 0 {
						workspace.Output = output.Name
						workspaces[workspace.WorkspaceNum] = workspace
					}
				}
			}
		}
	}

	return workspaces
}

func isAppNode(node *Node) bool {
	return (node.AppID != "" ||
		(node.WindowProperties != nil && node.WindowProperties.Class != "")) &&
		node.Type == "con"
}

// Collects app nodes
func FindApps(node *Node) []WindowInfo {
	var apps []WindowInfo

	for i := range node.Nodes {
		if isAppNode(&node.Nodes[i]) {
			appName := node.AppID
			if appName == "" && node.WindowProperties != nil {
				appName = node.WindowProperties.Class
			}

			apps = append(apps, WindowInfo{
				ID: node.ID,
			})
		}
	}

	return apps
}
