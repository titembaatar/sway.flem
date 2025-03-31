package config

// Top-level configuration
type Config struct {
	Workspaces map[string]Workspace `yaml:"workspaces"`
}

// Sway workspace configuration
type Workspace struct {
	Layout    string     `yaml:"layout"`
	Apps      []App      `yaml:"apps,omitempty"`
	Container *Container `yaml:"container,omitempty"`
}

// Container within a workspace
type Container struct {
	Split     string     `yaml:"split"`
	Size      string     `yaml:"size"`
	Apps      []App      `yaml:"apps,omitempty"`
	Container *Container `yaml:"container,omitempty"`
}

// Application to be launched
type App struct {
	Name string `yaml:"name"`
	Size string `yaml:"size"`
	Cmd  string `yaml:"cmd,omitempty"` // Optional command to execute instead of name
}
