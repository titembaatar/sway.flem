package config

func ApplyDefaults(config *Config) {
	for num, workspace := range config.Workspaces {
		// default layout if not specified
		if workspace.Layout == "" && config.Defaults.DefaultLayout != "" {
			workspace.Layout = config.Defaults.DefaultLayout
			config.Workspaces[num] = workspace
		}

		// default output if not specified
		if workspace.Output == "" && config.Defaults.DefaultOutput != "" {
			workspace.Output = config.Defaults.DefaultOutput
			config.Workspaces[num] = workspace
		}

		// defaults to each app
		for i, app := range workspace.Apps {
			// command to app name if not specified
			if app.Command == "" {
				app.Command = app.Name
				workspace.Apps[i] = app
			}
		}
	}
}
