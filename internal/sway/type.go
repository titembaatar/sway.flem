package sway

type Client struct {
	verbose bool
}

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
