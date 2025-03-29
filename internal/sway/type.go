package sway

type Client struct {
	verbose bool
}

// swaymsg -t get_tree --raw response
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
	Representation   string       `json:"representation,omitempty"`
	Layout           string       `json:"layout"` // Make sure this field exists
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

// Workspaces

type WorkspaceInfo struct {
	Number         int
	Layout         string
	Representation string
	AppOrder       []string
}

type LayoutNode struct {
	Type     string   // "V", "H", "S", "T", or app name
	Children []string // App names in this container
}

type LayoutType string

const (
	LayoutTypeVertical   LayoutType = "splitv"
	LayoutTypeHorizontal LayoutType = "splith"
	LayoutTypeStacking   LayoutType = "stacking"
	LayoutTypeTabbed     LayoutType = "tabbed"
)

// Apps

type AppNode struct {
	Name     string
	NodeID   int64
	Rect     Rect
	Floating bool
}
