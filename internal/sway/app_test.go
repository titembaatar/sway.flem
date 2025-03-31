package sway

import (
	"os/exec"
	"testing"

	"github.com/titembaatar/sway.flem/internal/config"
)

func TestLaunchApp(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name        string
		app         config.App
		mark        string
		mockStdout  string
		mockStderr  string
		mockExit    int
		expectError bool
	}{
		{
			name: "Success - launch by name",
			app: config.App{
				Name: "firefox",
				Size: "50ppt",
			},
			mark:        "sway_flem_ws_1_app_firefox",
			mockStdout:  `[{"success":true}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: false,
		},
		{
			name: "Success - launch by command",
			app: config.App{
				Name: "terminal",
				Cmd:  "alacritty",
				Size: "50ppt",
			},
			mark:        "sway_flem_ws_1_app_terminal",
			mockStdout:  `[{"success":true}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: false,
		},
		{
			name: "Error - empty command",
			app: config.App{
				Name: "",
				Cmd:  "",
				Size: "50ppt",
			},
			mark:        "sway_flem_ws_1_app_unknown",
			mockStdout:  `[{"success":true}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: true,
		},
		{
			name: "Error - app launch fails",
			app: config.App{
				Name: "invalid_app",
				Size: "50ppt",
			},
			mark:        "sway_flem_ws_1_app_invalid",
			mockStdout:  ``,
			mockStderr:  "Command not found",
			mockExit:    1,
			expectError: true,
		},
		{
			name: "Error - marking fails",
			app: config.App{
				Name: "firefox",
				Size: "50ppt",
			},
			mark:        "sway_flem_ws_1_app_firefox",
			mockStdout:  `[{"success":false,"error":"Failed to mark"}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock exec.Command for app launch and mark
			execCommand = mockExecCommand(tc.mockStdout, tc.mockStderr, tc.mockExit)

			// Call the function
			err := LaunchApp(tc.app, tc.mark)

			// Check the result
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestLaunchApps(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Create test apps
	apps := []config.App{
		{Name: "firefox", Size: "30ppt"},
		{Name: "terminal", Cmd: "alacritty", Size: "30ppt"},
		{Name: "code", Size: "40ppt"},
	}

	// Mock success response
	execCommand = mockExecCommand(`[{"success":true}]`, "", 0)

	// Call the function
	containerMark := "sway_flem_ws_1_con_main"
	appInfos, err := LaunchApps(containerMark, "splith", apps)

	// Check the results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that the correct number of app infos were returned
	if len(appInfos) != len(apps) {
		t.Errorf("Expected %d app infos, got %d", len(apps), len(appInfos))
	}

	// Check that the app marks were generated correctly
	for i, info := range appInfos {
		expectedMark := GenerateAppMark(containerMark, "app"+string(rune('1'+i)))
		if info.Mark != expectedMark {
			t.Errorf("Expected mark %s, got %s", expectedMark, info.Mark)
		}

		// Check that sizes were preserved
		if info.Size != apps[i].Size {
			t.Errorf("Expected size %s, got %s", apps[i].Size, info.Size)
		}

		// Check that layout was preserved
		if info.Layout != "splith" {
			t.Errorf("Expected layout %s, got %s", "splith", info.Layout)
		}
	}
}

func TestResizeApps(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Create test app infos
	appInfos := []AppInfo{
		{Mark: "sway_flem_ws_1_app_firefox", Size: "30ppt", Layout: "splith"},
		{Mark: "sway_flem_ws_1_app_terminal", Size: "30ppt", Layout: "splith"},
		{Mark: "sway_flem_ws_1_app_code", Size: "40ppt", Layout: "splith"},
		{Mark: "sway_flem_ws_1_app_nosize", Size: "", Layout: "splith"}, // One without size
	}

	// Mock success response
	execCommand = mockExecCommand(`[{"success":true}]`, "", 0)

	// Create a counter to track how many times exec.Command is called
	commandCount := 0
	oldExecCommand := execCommand
	execCommand = func(command string, args ...string) *exec.Cmd {
		commandCount++
		return oldExecCommand(command, args...)
	}

	// Call the function
	ResizeApps(appInfos)

	// Check that resize was called for each app that has a size
	expectedCalls := 0
	for _, info := range appInfos {
		if info.Size != "" {
			expectedCalls++
		}
	}

	// The number of calls should be (2 * expectedCalls) because:
	// - Each resize operation requires a focus command
	// - And then a resize command
	if commandCount != 2*expectedCalls {
		t.Errorf("Expected %d exec.Command calls, got %d", 2*expectedCalls, commandCount)
	}
}
