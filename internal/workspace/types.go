package workspace

import (
	"github.com/titembaatar/sway.flem/internal/app"
	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

// Manages workspace operations
type WorkspaceManager struct {
	Client  *sway.Client
	AppMgr  *app.AppManager
	Verbose bool
}

// For updating application properties
type AppUpdate struct {
	NodeID int64
	Config config.App
}

// Layout structure
type LayoutNode struct {
	Type     string
	Children []string
}

// App with position information
type OrderedApp struct {
	App      config.App
	Position int
}
