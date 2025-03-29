package config

type App struct {
	Name     string   `yaml:"name"`               // App name (app_id for Wayland, class for X11)
	Command  string   `yaml:"command,omitempty"`  // Command to launch (default: name)
	Size     string   `yaml:"size,omitempty"`     // Size (width height)
	Position string   `yaml:"position,omitempty"` // Position (for floating windows)
	Floating bool     `yaml:"floating,omitempty"` // Floating state
	Posts    []string `yaml:"post,omitempty"`     // Post-launch commands
	Launcher string   `yaml:"launcher,omitempty"` // Launcher to use (e.g., "tofi-drun")
	Delay    int64    `yaml:"delay,omitempty"`    // Delay in seconds to wait after launching
}

type Workspace struct {
	Layout         string `yaml:"layout,omitempty"`          // Workspace layout
	CloseUnmatched bool   `yaml:"close_unmatched,omitempty"` // Close apps not in config
	Apps           []App  `yaml:"apps"`                      // Apps in this workspace
}

type Config struct {
	FocusWorkspace int               `yaml:"focus_workspace"` // Workspace to focus after setup
	Workspaces     map[int]Workspace `yaml:"workspaces"`      // Workspaces configuration
}
