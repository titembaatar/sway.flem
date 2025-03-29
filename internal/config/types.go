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
	Name     string   `yaml:"name"`
	Command  string   `yaml:"command,omitempty"`
	Size     string   `yaml:"size,omitempty"`
	Position string   `yaml:"position,omitempty"`
	Floating bool     `yaml:"floating,omitempty"`
	Posts    []string `yaml:"post,omitempty"`
	Launcher string   `yaml:"launcher,omitempty"`
	Delay    int64    `yaml:"delay,omitempty"`
}
