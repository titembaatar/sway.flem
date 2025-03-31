package sway

import (
	"os/exec"
	"testing"

	"github.com/titembaatar/sway.flem/internal/config"
)

func TestSetupWorkspace(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Mock success response
	execCommand = mockExecCommand(`[{"success":true}]`, "", 0)

	// Create a workspace config
	workspace := config.Workspace{
		Layout: "splith",
		Apps: []config.App{
			{Name: "firefox", Size: "30ppt"},
			{Name: "terminal", Size: "30ppt"},
		},
		Container: &config.Container{
			Split: "splitv",
			Size:  "40ppt",
			Apps: []config.App{
				{Name: "code", Size: "50ppt"},
				{Name: "slack", Size: "50ppt"},
			},
		},
	}

	// Call the function
	err := SetupWorkspace("1", workspace)

	// Check the result
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSetupContainer(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Mock success response
	execCommand = mockExecCommand(`[{"success":true}]`, "", 0)

	// Create a container config
	container := &config.Container{
		Split: "splitv",
		Size:  "40ppt",
		Apps: []config.App{
			{Name: "code", Size: "50ppt"},
			{Name: "slack", Size: "50ppt"},
		},
	}

	// Call the function
	containerMark := "sway_flem_ws_1_con_main"
	parentMark := "sway_flem_ws_1"
	_, err := setupContainer(containerMark, container, parentMark)

	// Check the result
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSetupEnvironment(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Mock success response
	execCommand = mockExecCommand(`[{"success":true}]`, "", 0)

	// Create a config with multiple workspaces
	cfg := &config.Config{
		Workspaces: map[string]config.Workspace{
			"1": {
				Layout: "splith",
				Apps: []config.App{
					{Name: "firefox", Size: "50ppt"},
				},
			},
			"2": {
				Layout: "tabbed",
				Apps: []config.App{
					{Name: "terminal", Size: "30ppt"},
					{Name: "code", Size: "70ppt"},
				},
			},
		},
	}

	// Call the function
	err := SetupEnvironment(cfg)

	// Check the result
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSetupWorkspaceErrorHandling(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name                string
		mockStdoutWorkspace string
		mockStderrWorkspace string
		mockExitWorkspace   int
		mockStdoutMark      string
		mockStderrMark      string
		mockExitMark        int
		mockStdoutLayout    string
		mockStderrLayout    string
		mockExitLayout      int
		expectError         bool
	}{
		{
			name:                "Workspace creation fails",
			mockStdoutWorkspace: `[{"success":false,"error":"Failed to create workspace"}]`,
			mockStderrWorkspace: "",
			mockExitWorkspace:   0,
			mockStdoutMark:      `[{"success":true}]`,
			mockStderrMark:      "",
			mockExitMark:        0,
			mockStdoutLayout:    `[{"success":true}]`,
			mockStderrLayout:    "",
			mockExitLayout:      0,
			expectError:         true,
		},
		{
			name:                "Marking fails",
			mockStdoutWorkspace: `[{"success":true}]`,
			mockStderrWorkspace: "",
			mockExitWorkspace:   0,
			mockStdoutMark:      `[{"success":false,"error":"Failed to mark workspace"}]`,
			mockStderrMark:      "",
			mockExitMark:        0,
			mockStdoutLayout:    `[{"success":true}]`,
			mockStderrLayout:    "",
			mockExitLayout:      0,
			expectError:         true,
		},
		{
			name:                "Layout setting fails",
			mockStdoutWorkspace: `[{"success":true}]`,
			mockStderrWorkspace: "",
			mockExitWorkspace:   0,
			mockStdoutMark:      `[{"success":true}]`,
			mockStderrMark:      "",
			mockExitMark:        0,
			mockStdoutLayout:    `[{"success":false,"error":"Failed to set layout"}]`,
			mockStderrLayout:    "",
			mockExitLayout:      0,
			expectError:         true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// This is a bit complex because the function makes multiple exec.Command calls
			// We'll use a counter to track which call we're on
			callCount := 0
			execCommand = func(command string, args ...string) *exec.Cmd {
				callCount++
				if callCount == 1 {
					// First call - workspace
					return mockExecCommand(tc.mockStdoutWorkspace, tc.mockStderrWorkspace, tc.mockExitWorkspace)(command, args...)
				} else if callCount == 2 {
					// Second call - mark
					return mockExecCommand(tc.mockStdoutMark, tc.mockStderrMark, tc.mockExitMark)(command, args...)
				} else {
					// Third call - layout
					return mockExecCommand(tc.mockStdoutLayout, tc.mockStderrLayout, tc.mockExitLayout)(command, args...)
				}
			}

			// Create a minimal workspace config for testing
			workspace := config.Workspace{
				Layout: "splith",
			}

			// Call the function
			err := SetupWorkspace("1", workspace)

			// Check the result
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetLayoutForContainer(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name          string
		containerMark string
		mockStdout    string
		mockStderr    string
		mockExit      int
		expected      string
	}{
		{
			name:          "Focus succeeds",
			containerMark: "sway_flem_ws_1_con_main",
			mockStdout:    `[{"success":true}]`,
			mockStderr:    "",
			mockExit:      0,
			expected:      "splith", // Default layout
		},
		{
			name:          "Focus fails",
			containerMark: "sway_flem_ws_1_con_main",
			mockStdout:    `[{"success":false}]`,
			mockStderr:    "Focus failed",
			mockExit:      1,
			expected:      "splith", // Default layout
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock exec.Command
			execCommand = mockExecCommand(tc.mockStdout, tc.mockStderr, tc.mockExit)

			// Call the function
			layout := getLayoutForContainer(tc.containerMark)

			// Check the result
			if layout != tc.expected {
				t.Errorf("Expected layout %s, got %s", tc.expected, layout)
			}
		})
	}
}
