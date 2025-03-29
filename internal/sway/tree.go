package sway

func isAppNode(node *Node) bool {
	return node.AppID != "" || (node.WindowProperties != nil && node.WindowProperties.Class != "")
}

func (n *Node) FindWorkspaces() map[int]*Node {
	workspaces := make(map[int]*Node)

	for i := range n.Nodes {
		output := &n.Nodes[i]
		if output.Type == "output" {
			for j := range output.Nodes {
				workspace := &output.Nodes[j]
				if workspace.Type == "workspace" && workspace.WorkspaceNum > 0 {
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
