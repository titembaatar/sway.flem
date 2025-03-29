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

// App with its position in the layout
type OrderedApp struct {
	App      config.App
	Position int
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
	Position string
	Layout   string
	PostCmds []string
}
