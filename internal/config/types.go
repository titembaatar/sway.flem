package config

type Config struct {
	FocusWorkspace int               `yaml:"focus_workspace"`
	Workspaces     map[int]Workspace `yaml:"workspaces"`
	Defaults       DefaultsConfig    `yaml:"defaults,omitempty"`
}

// Fallback settings
type DefaultsConfig struct {
	DefaultLayout   string `yaml:"default_layout,omitempty"`
	DefaultOutput   string `yaml:"default_output,omitempty"`
	DefaultFloating bool   `yaml:"default_floating,omitempty"`
}

type Workspace struct {
	Layout         string `yaml:"layout,omitempty"`
	Output         string `yaml:"output,omitempty"`
	CloseUnmatched bool   `yaml:"close_unmatched,omitempty"`
	Apps           []App  `yaml:"apps"`
}

type App struct {
	Name     string   `yaml:"name"`               // Name of the app (required)
	Command  string   `yaml:"command,omitempty"`  // Command to launch the app (defaults to Name if empty)
	Size     string   `yaml:"size,omitempty"`     // Size of the app window, e.g. "800x600" or "50ppt 70ppt"
	Position string   `yaml:"position,omitempty"` // Position of the app window, e.g. "center", "top", "0 0"
	Floating bool     `yaml:"floating,omitempty"` // Whether the app should be floating or tiled
	Posts    []string `yaml:"post,omitempty"`     // Commands to run after launching the app
	Launcher string   `yaml:"launcher,omitempty"` // Custom launcher to use (if any)
	Delay    int64    `yaml:"delay,omitempty"`    // Delay in seconds before configuring the app
}
