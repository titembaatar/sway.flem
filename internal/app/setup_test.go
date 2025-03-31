package app

import (
	"os"
	"os/exec"
	"testing"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
)

func init() {
	// Set log level to none during tests
	log.SetLevel(log.LogLevelNone)
}

// mockExecCommand is a variable that will be used to replace exec.Command
var execCommand = exec.Command

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

func TestSetup(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Mock success response
	execCommand = mockExecCommand("1.0.0\n", "", 0)

	// Create a minimal config for testing
	cfg := &config.Config{
		Workspaces: map[string]config.Workspace{
			"1": {
				Layout: "splith",
				Apps: []config.App{
					{Name: "firefox", Size: "50ppt"},
				},
			},
		},
	}

	// Call the function
	err := Setup(cfg)

	// We expect this to fail in tests without mocking all the sway functions
	// We're mostly checking that the dependencies check works
	if err == nil {
		t.Error("Expected error due to missing setup implementation in test, but got none")
	}
}

func TestCheckDependencies(t *testing.T) {
	// Save the original exec.Command and restore after test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name        string
		mockStdout  string
		mockStderr  string
		mockExit    int
		expectError bool
	}{
		{
			name:        "Swaymsg available",
			mockStdout:  "swaymsg 1.0.0\n",
			mockStderr:  "",
			mockExit:    0,
			expectError: false,
		},
		{
			name:        "Swaymsg not available",
			mockStdout:  "",
			mockStderr:  "Command not found",
			mockExit:    1,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock exec.Command
			execCommand = mockExecCommand(tc.mockStdout, tc.mockStderr, tc.mockExit)

			// Call the function
			err := checkDependencies()

			// Check the result
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCheckCommand(t *testing.T) {
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
			name:        "Command available",
			command:     "swaymsg",
			mockStdout:  "swaymsg 1.0.0\n",
			mockStderr:  "",
			mockExit:    0,
			expectError: false,
		},
		{
			name:        "Command not available",
			command:     "swaymsg",
			mockStdout:  "",
			mockStderr:  "Command not found",
			mockExit:    1,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock exec.Command
			execCommand = mockExecCommand(tc.mockStdout, tc.mockStderr, tc.mockExit)

			// Call the function
			err := checkCommand(tc.command)

			// Check the result
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
