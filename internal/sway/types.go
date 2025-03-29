package sway

import (
	"errors"
)

// Error constants
var (
	ErrCommandFailed = errors.New("sway command failed")
)

// Client for Sway IPC communication
type Client struct {
	Verbose bool
}

// Node in the Sway tree
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
	Layout           string       `json:"layout"`
	Output           string       `json:"output,omitempty"`
}

// Rectangle with position and size
type Rect struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Window metadata
type WindowProps struct {
	Class    string `json:"class"`
	Instance string `json:"instance"`
	Title    string `json:"title"`
}

// Workspace details
type WorkspaceInfo struct {
	Number         int
	Layout         string
	Output         string
	Representation string
	AppOrder       []string
}

// LayoutType constants for different layout types
type LayoutType string

const (
	LayoutTypeVertical   LayoutType = "splitv"
	LayoutTypeHorizontal LayoutType = "splith"
	LayoutTypeStacking   LayoutType = "stacking"
	LayoutTypeTabbed     LayoutType = "tabbed"
)

// AppNode represents an application in the Sway tree
type AppNode struct {
	Name     string
	NodeID   int64
	Rect     Rect
	Floating bool
}
