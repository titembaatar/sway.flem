package manager

import (
	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

// Manager coordinates workspace and application
type Manager struct {
	Config  *config.Config
	Client  *sway.Client
	Verbose bool
}

// Update to be applied to an existing application
type AppUpdate struct {
	NodeID int64
	Config config.App
}

// Container information
type ContainerInfo struct {
	ID     int64
	Layout string
	Apps   []AppInfo
}

// AppInfo tracks an app that has been launched and needs to be configured
type AppInfo struct {
	App    config.App
	NodeID int64
	Layout string // The layout context for this app
}

// Layout node in a workspace
type LayoutNode struct {
	Type     string
	Children []string
}

// Options for launching an application
type AppLaunchOptions struct {
	Command   string
	Layout    string
	WaitDelay int64
}

// Options for updating an application
type AppUpdateOptions struct {
	Floating bool
	Size     string
	Layout   string
	PostCmds []string
}
