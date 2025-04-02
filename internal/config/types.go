package config

type Config struct {
	Workspaces map[string]Workspace `yaml:"workspaces"`
}

type Workspace struct {
	Layout     string      `yaml:"layout"`
	Containers []Container `yaml:"containers"`
}

type Container struct {
	App        string      `yaml:"app,omitempty"`
	Cmd        string      `yaml:"cmd,omitempty"`
	Size       string      `yaml:"size"`
	Delay      int64       `yaml:"delay"`
	Split      string      `yaml:"split,omitempty"`
	Containers []Container `yaml:"containers,omitempty"`
}
