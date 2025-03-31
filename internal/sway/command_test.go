package sway

import (
	"os/exec"
	"testing"
)

func TestRunCommand(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name        string
		command     string
		mockStdout  string
		mockStderr  string
		mockExit    int
		expectError bool
	}{
		{
			name:        "Success",
			command:     "layout tabbed",
			mockStdout:  `[{"success":true}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: false,
		},
		{
			name:        "Command failure",
			command:     "invalid command",
			mockStdout:  `[{"success":false,"error":"Invalid command"}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: true,
		},
		{
			name:        "Command execution error",
			command:     "some command",
			mockStdout:  "",
			mockStderr:  "Error executing command",
			mockExit:    1,
			expectError: true,
		},
		{
			name:        "Invalid JSON response",
			command:     "some command",
			mockStdout:  `Not a JSON response`,
			mockStderr:  "",
			mockExit:    0,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock exec.Command
			execCommand = mockExecCommand(tc.mockStdout, tc.mockStderr, tc.mockExit)

			// Call the function
			_, err := RunCommand(tc.command)

			// Check the result
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestRunCommandWithNoResponse(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name        string
		command     string
		mockStdout  string
		mockStderr  string
		mockExit    int
		expectError bool
	}{
		{
			name:        "Success",
			command:     "layout tabbed",
			mockStdout:  ``,
			mockStderr:  "",
			mockExit:    0,
			expectError: false,
		},
		{
			name:        "Command execution error",
			command:     "some command",
			mockStdout:  "",
			mockStderr:  "Error executing command",
			mockExit:    1,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock exec.Command
			execCommand = mockExecCommand(tc.mockStdout, tc.mockStderr, tc.mockExit)

			// Call the function
			err := RunCommandWithNoResponse(tc.command)

			// Check the result
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetWorkspaces(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name        string
		mockStdout  string
		mockStderr  string
		mockExit    int
		expectError bool
		expectedLen int
	}{
		{
			name:        "Success",
			mockStdout:  `[{"name":"1"},{"name":"2"}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: false,
			expectedLen: 2,
		},
		{
			name:        "Command execution error",
			mockStdout:  "",
			mockStderr:  "Error executing command",
			mockExit:    1,
			expectError: true,
			expectedLen: 0,
		},
		{
			name:        "Invalid JSON response",
			mockStdout:  `Not a JSON response`,
			mockStderr:  "",
			mockExit:    0,
			expectError: true,
			expectedLen: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock exec.Command
			execCommand = mockExecCommand(tc.mockStdout, tc.mockStderr, tc.mockExit)

			// Call the function
			workspaces, err := GetWorkspaces()

			// Check the result
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError && len(workspaces) != tc.expectedLen {
				t.Errorf("Expected %d workspaces, got %d", tc.expectedLen, len(workspaces))
			}
		})
	}
}

func TestSwitchToWorkspace(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name        string
		workspace   string
		mockStdout  string
		mockStderr  string
		mockExit    int
		expectError bool
	}{
		{
			name:        "Success",
			workspace:   "1",
			mockStdout:  `[{"success":true}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: false,
		},
		{
			name:        "Command failure",
			workspace:   "invalid",
			mockStdout:  `[{"success":false,"error":"Invalid workspace"}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock exec.Command
			execCommand = mockExecCommand(tc.mockStdout, tc.mockStderr, tc.mockExit)

			// Call the function
			err := SwitchToWorkspace(tc.workspace)

			// Check the result
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCreateWorkspace(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name        string
		workspace   string
		layout      string
		mockStdout1 string // For SwitchToWorkspace
		mockStderr1 string
		mockExit1   int
		mockStdout2 string // For setting layout
		mockStderr2 string
		mockExit2   int
		expectError bool
	}{
		{
			name:        "Success",
			workspace:   "1",
			layout:      "tabbed",
			mockStdout1: `[{"success":true}]`,
			mockStderr1: "",
			mockExit1:   0,
			mockStdout2: `[{"success":true}]`,
			mockStderr2: "",
			mockExit2:   0,
			expectError: false,
		},
		{
			name:        "Switch workspace failure",
			workspace:   "invalid",
			layout:      "tabbed",
			mockStdout1: `[{"success":false,"error":"Invalid workspace"}]`,
			mockStderr1: "",
			mockExit1:   0,
			mockStdout2: `[{"success":true}]`,
			mockStderr2: "",
			mockExit2:   0,
			expectError: true,
		},
		{
			name:        "Layout setting failure",
			workspace:   "1",
			layout:      "invalid",
			mockStdout1: `[{"success":true}]`,
			mockStderr1: "",
			mockExit1:   0,
			mockStdout2: `[{"success":false,"error":"Invalid layout"}]`,
			mockStderr2: "",
			mockExit2:   0,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// First mock for SwitchToWorkspace
			execCommand = mockExecCommand(tc.mockStdout1, tc.mockStderr1, tc.mockExit1)

			// This is a bit of a hack to test a function that makes two commands
			// We'll replace execCommand after the first command
			oldExecCommand := execCommand
			execCommand = func(command string, args ...string) *exec.Cmd {
				// If this is the layout command, use the second mock
				if args[2] == "layout" {
					return mockExecCommand(tc.mockStdout2, tc.mockStderr2, tc.mockExit2)(command, args...)
				}
				return oldExecCommand(command, args...)
			}

			// Call the function
			err := CreateWorkspace(tc.workspace, tc.layout)

			// Check the result
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
