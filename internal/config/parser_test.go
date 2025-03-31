package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.yaml")

	// Define test cases
	tests := []struct {
		name        string
		configYAML  string
		expectError bool
	}{
		{
			name: "Valid minimal config",
			configYAML: `
workspaces:
  "1":
    layout: "splith"
    apps:
      - name: "app1"
        size: "50ppt"
`,
			expectError: false,
		},
		{
			name: "Valid complex config",
			configYAML: `
workspaces:
  "2":
    layout: "h"
    apps:
      - name: "firefox"
        size: "15ppt"
      - name: "terminal"
        size: "15ppt"
        cmd: "alacritty"
    container:
      split: "v"
      size: "70ppt"
      apps:
        - name: "code"
          size: "15ppt"
        - name: "slack"
          size: "15ppt"
      container:
        split: "splith"
        size: "70ppt"
        apps:
          - name: "spotify"
            size: "30ppt"
          - name: "discord"
            size: "70ppt"
`,
			expectError: false,
		},
		{
			name: "Invalid config - no workspaces",
			configYAML: `
foo: bar
`,
			expectError: true,
		},
		{
			name: "Invalid config - invalid layout",
			configYAML: `
workspaces:
  "1":
    layout: "invalid_layout"
    apps:
      - name: "app1"
        size: "50ppt"
`,
			expectError: true,
		},
		{
			name: "Invalid config - missing app name",
			configYAML: `
workspaces:
  "1":
    layout: "splith"
    apps:
      - size: "50ppt"
`,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Write the test config to a temporary file
			if err := os.WriteFile(configPath, []byte(tc.configYAML), 0644); err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			// Load the config
			config, err := LoadConfig(configPath)

			// Check if the error matches expectation
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// For valid configs, check that the config was parsed correctly
			if !tc.expectError && err == nil {
				if len(config.Workspaces) == 0 {
					t.Errorf("Expected workspaces but found none")
				}
			}
		})
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	// Test loading a non-existent file
	_, err := LoadConfig("non_existent_file.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent file, but got none")
	}
}
