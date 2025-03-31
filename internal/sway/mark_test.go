package sway

import (
	"os"
	"os/exec"
	"testing"

	"github.com/titembaatar/sway.flem/internal/log"
)

// mockExecCommand is a variable that will be used to replace exec.Command
var execCommand = exec.Command

func init() {
	// Set log level to none during tests
	log.SetLevel(log.LogLevelNone)
}

// mockExecCommand creates a mock exec.Command that returns predictable output
func mockExecCommand(mockStdout, mockStderr string, mockExitCode int) func(string, ...string) *exec.Cmd {
	return func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestMockExecCommand", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)

		// Set environment variables to control the mock behavior
		cmd.Env = append(os.Environ(),
			"GO_WANT_MOCK_STDOUT="+mockStdout,
			"GO_WANT_MOCK_STDERR="+mockStderr,
			"GO_WANT_MOCK_EXIT_CODE="+string(rune('0'+mockExitCode)),
		)

		return cmd
	}
}

// TestMockExecCommand is not a real test - it's used by mockExecCommand
func TestMockExecCommand(t *testing.T) {
	if os.Getenv("GO_WANT_MOCK_STDOUT") == "" {
		return
	}

	// Write the mock stdout and stderr
	if stdout := os.Getenv("GO_WANT_MOCK_STDOUT"); stdout != "" {
		os.Stdout.WriteString(stdout)
	}
	if stderr := os.Getenv("GO_WANT_MOCK_STDERR"); stderr != "" {
		os.Stderr.WriteString(stderr)
	}

	// Exit with the mock exit code
	exitCode := int(os.Getenv("GO_WANT_MOCK_EXIT_CODE")[0] - '0')
	os.Exit(exitCode)
}

func TestGenerateWorkspaceMark(t *testing.T) {
	tests := []struct {
		workspaceName string
		expected      string
	}{
		{"1", "sway_flem_ws_1"},
		{"main", "sway_flem_ws_main"},
		{"1:Firefox", "sway_flem_ws_1:Firefox"},
	}

	for _, tc := range tests {
		t.Run(tc.workspaceName, func(t *testing.T) {
			result := GenerateWorkspaceMark(tc.workspaceName)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestGenerateContainerMark(t *testing.T) {
	workspaceMark := "sway_flem_ws_1"
	tests := []struct {
		containerID string
		expected    string
	}{
		{"main", "sway_flem_ws_1_con_main"},
		{"nested", "sway_flem_ws_1_con_nested"},
		{"container1", "sway_flem_ws_1_con_container1"},
	}

	for _, tc := range tests {
		t.Run(tc.containerID, func(t *testing.T) {
			result := GenerateContainerMark(workspaceMark, tc.containerID)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestGenerateAppMark(t *testing.T) {
	containerMark := "sway_flem_ws_1_con_main"
	tests := []struct {
		appID    string
		expected string
	}{
		{"firefox", "sway_flem_ws_1_con_main_app_firefox"},
		{"app1", "sway_flem_ws_1_con_main_app_app1"},
		{"terminal", "sway_flem_ws_1_con_main_app_terminal"},
	}

	for _, tc := range tests {
		t.Run(tc.appID, func(t *testing.T) {
			result := GenerateAppMark(containerMark, tc.appID)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestResizeMark(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name        string
		mark        string
		size        string
		layout      string
		mockStdout  string
		mockStderr  string
		mockExit    int
		expectError bool
	}{
		{
			name:        "Success - width resize",
			mark:        "sway_flem_ws_1_app_firefox",
			size:        "50ppt",
			layout:      "splith",
			mockStdout:  `[{"success":true}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: false,
		},
		{
			name:        "Success - height resize",
			mark:        "sway_flem_ws_1_app_firefox",
			size:        "50ppt",
			layout:      "splitv",
			mockStdout:  `[{"success":true}]`,
			mockStderr:  "",
			mockExit:    0,
			expectError: false,
		},
		{
			name:        "Error - focus fails",
			mark:        "sway_flem_ws_1_app_firefox",
			size:        "50ppt",
			layout:      "splith",
			mockStdout:  "",
			mockStderr:  "Failed to focus",
			mockExit:    1,
			expectError: true,
		},
		{
			name:        "Error - resize fails",
			mark:        "sway_flem_ws_1_app_firefox",
			size:        "50ppt",
			layout:      "splith",
			mockStdout:  `[{"success":false,"error":"Failed to resize"}]`,
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
			err := ResizeMark(tc.mark, tc.size, tc.layout)

			// Check the result
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
