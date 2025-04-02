package config

type Config struct {
	Workspaces map[string]Workspace `yaml:"workspaces"`
}

type Workspace struct {
	Layout     string      `yaml:"layout"`
	Containers []Container `yaml:"containers"`
}

type Container struct {
	App        string      `yaml:"app"`
	Cmd        string      `yaml:"cmd"`
	Size       string      `yaml:"size"`
	Delay      int64       `yaml:"delay"`
	Split      string      `yaml:"split"`
	Containers []Container `yaml:"containers"`
	Post       []string    `yaml:"post"`
}
