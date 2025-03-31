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
	Layout         string     `yaml:"layout,omitempty"`
	Output         string     `yaml:"output,omitempty"`
	CloseUnmatched bool       `yaml:"close_unmatched,omitempty"`
	Apps           []App      `yaml:"apps"`
	Container      *Container `yaml:"container,omitempty"`
}

type Container struct {
	Layout    string     `yaml:"layout,omitempty"`
	Size      string     `yaml:"size,omitempty"`
	Apps      []App      `yaml:"apps"`
	Container *Container `yaml:"container,omitempty"` // Nested container
}

type App struct {
	Name     string   `yaml:"name"`
	Command  string   `yaml:"command,omitempty"` // Command to launch the app (defaults to Name if empty)
	Size     string   `yaml:"size,omitempty"`
	Floating bool     `yaml:"floating,omitempty"`
	Posts    []string `yaml:"post,omitempty"`  // Commands to run after launching the app
	Delay    int64    `yaml:"delay,omitempty"` // Delay in seconds before configuring the app
}
