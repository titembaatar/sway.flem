package config

import (
	"github.com/titembaatar/sway.flem/internal/log"
)

var validLayoutTypes = map[string]string{
	// layout names
	"splith":   "splith",
	"splitv":   "splitv",
	"stacking": "stacking",
	"tabbed":   "tabbed",

	// aliases
	"horizontal": "splith",
	"h":          "splith",
	"vertical":   "splitv",
	"v":          "splitv",
	"stack":      "stacking",
	"s":          "stacking",
	"tab":        "tabbed",
	"t":          "tabbed",
}

func NormalizeLayoutType(layoutType string) (string, error) {
	if normalized, ok := validLayoutTypes[layoutType]; ok {
		return normalized, nil
	}
	return "", ErrInvalidLayoutType
}

func ValidateConfig(config *Config) error {
	if len(config.Workspaces) == 0 {
		return NewConfigError(ErrNoWorkspaces, "", "", -1)
	}

	log.Debug("Validating configuration with %d workspaces", len(config.Workspaces))

	for name, workspace := range config.Workspaces {
		normalizedLayout, err := NormalizeLayoutType(workspace.Layout)
		if err != nil {
			return NewConfigError(err, name, "layout", -1)
		}
		workspace.Layout = normalizedLayout
		config.Workspaces[name] = workspace

		if err := validateWorkspace(name, workspace); err != nil {
			return err
		}

		log.Debug("Workspace '%s' validated successfully", name)
	}

	log.Info("Configuration validated successfully")
	return nil
}

func validateWorkspace(name string, workspace Workspace) error {
	if workspace.Layout == "" {
		return NewConfigError(ErrMissingLayout, name, "", -1)
	}

	for i, app := range workspace.Apps {
		if app.Name == "" {
			return NewConfigError(ErrMissingAppName, name, "app", i)
		}
	}

	if workspace.Container != nil {
		if err := validateContainer(name, workspace.Container, "container"); err != nil {
			return err
		}
	}

	return nil
}

func validateContainer(workspaceName string, container *Container, context string) error {
	if container.Split == "" {
		return NewConfigError(ErrMissingSplit, workspaceName, context, -1)
	}

	normalizedSplit, err := NormalizeLayoutType(container.Split)
	if err != nil {
		return NewConfigError(err, workspaceName, context, -1)
	}
	container.Split = normalizedSplit

	if container.Size == "" {
		return NewConfigError(ErrMissingSize, workspaceName, context, -1)
	}

	for i, app := range container.Apps {
		if app.Name == "" {
			return NewConfigError(ErrMissingAppName, workspaceName, context+".app", i)
		}
	}

	if container.Container != nil {
		nestedContext := context + ".container"
		if err := validateContainer(workspaceName, container.Container, nestedContext); err != nil {
			return err
		}
	}

	return nil
}
