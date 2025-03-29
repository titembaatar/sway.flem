package app

import (
	"github.com/titembaatar/sway.flem/internal/sway"
)

// Manages application operations
type AppManager struct {
	Client  *sway.Client
	Verbose bool
}

// Launch applications
type LaunchOptions struct {
	Command   string
	Layout    string
	WaitDelay int64
}

// Update application properties
type UpdateOptions struct {
	Floating bool
	Size     string
	Position string
	Layout   string
	PostCmds []string
}
