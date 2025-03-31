package config

import (
	"errors"
	"testing"
)

func TestNormalizeLayoutType(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{"Standard splith", "splith", "splith", false},
		{"Standard splitv", "splitv", "splitv", false},
		{"Standard stacking", "stacking", "stacking", false},
		{"Standard tabbed", "tabbed", "tabbed", false},
		{"Horizontal alias", "horizontal", "splith", false},
		{"Vertical alias", "vertical", "splitv", false},
		{"Stack alias", "stack", "stacking", false},
		{"Tab alias", "tab", "tabbed", false},
		{"H shorthand", "h", "splith", false},
		{"V shorthand", "v", "splitv", false},
		{"S shorthand", "s", "stacking", false},
		{"T shorthand", "t", "tabbed", false},
		{"Invalid layout", "invalid", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NormalizeLayoutType(tc.input)

			// Check error expectation
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check result
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedErr error
	}{
		{
			name: "Valid config",
			config: &Config{
				Workspaces: map[string]Workspace{
					"1": {
						Layout: "splith",
						Apps: []App{
							{Name: "app1", Size: "50ppt"},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Valid config with layout alias",
			config: &Config{
				Workspaces: map[string]Workspace{
					"1": {
						Layout: "h",
						Apps: []App{
							{Name: "app1", Size: "50ppt"},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Valid nested config",
			config: &Config{
				Workspaces: map[string]Workspace{
					"1": {
						Layout: "splith",
						Apps: []App{
							{Name: "app1", Size: "30ppt"},
						},
						Container: &Container{
							Split: "splitv",
							Size:  "70ppt",
							Apps: []App{
								{Name: "app2", Size: "50ppt"},
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Empty workspaces",
			config: &Config{
				Workspaces: map[string]Workspace{},
			},
			expectedErr: ErrNoWorkspaces,
		},
		{
			name: "Invalid layout",
			config: &Config{
				Workspaces: map[string]Workspace{
					"1": {
						Layout: "invalid",
						Apps: []App{
							{Name: "app1", Size: "50ppt"},
						},
					},
				},
			},
			expectedErr: ErrInvalidLayoutType,
		},
		{
			name: "Missing app name",
			config: &Config{
				Workspaces: map[string]Workspace{
					"1": {
						Layout: "splith",
						Apps: []App{
							{Size: "50ppt"}, // Missing name
						},
					},
				},
			},
			expectedErr: ErrMissingAppName,
		},
		{
			name: "Missing container split",
			config: &Config{
				Workspaces: map[string]Workspace{
					"1": {
						Layout: "splith",
						Container: &Container{
							Size: "70ppt",
							Apps: []App{
								{Name: "app1", Size: "50ppt"},
							},
						},
					},
				},
			},
			expectedErr: ErrMissingSplit,
		},
		{
			name: "Missing container size",
			config: &Config{
				Workspaces: map[string]Workspace{
					"1": {
						Layout: "splith",
						Container: &Container{
							Split: "splitv",
							// Size missing
							Apps: []App{
								{Name: "app1", Size: "50ppt"},
							},
						},
					},
				},
			},
			expectedErr: ErrMissingSize,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateConfig(tc.config)

			// Check if we got the expected error type
			if tc.expectedErr == nil && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			} else if tc.expectedErr != nil && err == nil {
				t.Errorf("Expected error %v, got nil", tc.expectedErr)
			} else if tc.expectedErr != nil && err != nil {
				// Unwrap the ConfigError to get the base error
				var configErr *ConfigError
				if errors.As(err, &configErr) {
					if !errors.Is(configErr.Err, tc.expectedErr) {
						t.Errorf("Expected error %v, got %v", tc.expectedErr, configErr.Err)
					}
				} else {
					t.Errorf("Expected ConfigError, got %T", err)
				}
			}
		})
	}
}
