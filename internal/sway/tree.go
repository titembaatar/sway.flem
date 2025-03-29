package sway

type Node struct {
	ID               int64        `json:"id"`
	Name             string       `json:"name"`
	Type             string       `json:"type"`
	AppID            string       `json:"app_id"`
	Shell            string       `json:"shell"`
	Visible          bool         `json:"visible"`
	Focused          bool         `json:"focused"`
	WorkspaceNum     int          `json:"num,omitempty"`
	Rect             Rect         `json:"rect"`
	Window           *int64       `json:"window"`
	WindowProperties *WindowProps `json:"window_properties"`
	Nodes            []Node       `json:"nodes"`
	FloatingNodes    []Node       `json:"floating_nodes"`
}

type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type WindowProps struct {
	Class    string `json:"class"`
	Instance string `json:"instance"`
	Title    string `json:"title"`
}

type AppNode struct {
	Name     string
	NodeID   int64
	Rect     Rect
	Floating bool
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
			appName := node.AppID
			if appName == "" && node.WindowProperties != nil {
				// X11 application
				appName = "xwayland:" + node.WindowProperties.Class
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

func isAppNode(node *Node) bool {
	return node.AppID != "" || (node.WindowProperties != nil && node.WindowProperties.Class != "")
}
