package config

import (
	"github.com/titembaatar/sway.flem/pkg/types"
)

type Config struct {
	Workspaces map[string]Workspace `yaml:"workspaces" json:"workspaces"`
	Focus      []string             `yaml:"focus" json:"focus"`
}

type Workspace struct {
	Layout     types.LayoutType `yaml:"layout" json:"layout"`
	Containers []Container      `yaml:"containers" json:"containers"`
}

type Container struct {
	App        string           `yaml:"app" json:"app"`
	Cmd        string           `yaml:"cmd" json:"cmd"`
	Size       string           `yaml:"size" json:"size"`
	Delay      int64            `yaml:"delay" json:"delay"`
	Post       []string         `yaml:"post" json:"post"`
	RerunPost  bool             `yaml:"rerun-post" json:"rerun_post"`
	Split      types.LayoutType `yaml:"split" json:"split"`
	Containers []Container      `yaml:"containers" json:"containers"`
}
