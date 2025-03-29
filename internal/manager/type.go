package manager

import (
	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/sway"
)

type Manager struct {
	client *sway.Client
	config *config.Config
}

type WorkspaceManager struct {
	client *sway.Client
}

type AppManager struct {
	client *sway.Client
}

type AppUpdate struct {
	NodeID int64
	Config config.App
}
