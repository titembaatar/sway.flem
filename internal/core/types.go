package core

import (
	"github.com/titembaatar/sway.flem/internal/app"
	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
	"github.com/titembaatar/sway.flem/internal/workspace"
)

// Coordinates workspace and app management
type Manager struct {
	Config           *config.Config
	SwayClient       *sway.Client
	WorkspaceManager *workspace.WorkspaceManager
	AppManager       *app.AppManager
	Verbose          bool
}

// Interface for workspace operations
type WorkspaceManager interface {
	SetupWorkspace(wsNum int, wsConfig config.Workspace, currentApps []sway.AppNode) error
}

// Interface for application operations
type AppManager interface {
	LaunchApp(app config.App, layout string) error
	UpdateApp(nodeID int64, app config.App, layout string) error
}
