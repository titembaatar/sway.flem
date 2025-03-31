package mock

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SwayMock is a mock implementation of the Sway interface
type SwayMock struct {
	// Track what commands were executed
	Commands []string

	// Track marks that have been applied
	Marks map[string]bool

	// Track focuses that have been applied
	Focuses []string

	// Track current workspace
	CurrentWorkspace string

	// Track launched applications
	LaunchedApps []string

	// Define when certain commands should fail
	ShouldFail map[string]bool

	// Control response for GetWorkspaces
	MockWorkspaces []string
}

// NewSwayMock creates a new SwayMock instance
func NewSwayMock() *SwayMock {
	return &SwayMock{
		Commands:       []string{},
		Marks:          make(map[string]bool),
		Focuses:        []string{},
		LaunchedApps:   []string{},
		ShouldFail:     make(map[string]bool),
		MockWorkspaces: []string{"1", "2"},
	}
}

// RunCommand mocks running a swaymsg command
func (m *SwayMock) RunCommand(command string) ([]map[string]any, error) {
	m.Commands = append(m.Commands, command)

	// Check if this command should fail
	for failCmd, shouldFail := range m.ShouldFail {
		if shouldFail && strings.Contains(command, failCmd) {
			return nil, fmt.Errorf("mock failure for command: %s", command)
		}
	}

	// Handle mark command
	if strings.HasPrefix(command, "mark") {
		parts := strings.Fields(command)
		if len(parts) >= 3 {
			markName := parts[2]
			m.Marks[markName] = true
		}
	}

	// Handle focus command
	if strings.Contains(command, "focus") {
		m.Focuses = append(m.Focuses, command)

		// If it's focusing a mark, extract the mark name
		if strings.Contains(command, "con_mark") {
			start := strings.Index(command, "con_mark=")
			if start != -1 {
				// Extract mark between the = and ]
				start += 9 // length of "con_mark="
				end := strings.Index(command[start:], "]")
				if end != -1 {
					markName := command[start : start+end]
					m.Focuses = append(m.Focuses, markName)
				}
			}
		}
	}

	// Handle workspace command
	if strings.HasPrefix(command, "workspace") {
		parts := strings.Fields(command)
		if len(parts) >= 2 {
			m.CurrentWorkspace = parts[1]
		}
	}

	// Create a successful response
	response := []map[string]any{
		{"success": true},
	}

	return response, nil
}

// RunCommandWithNoResponse mocks running a swaymsg command with no response
func (m *SwayMock) RunCommandWithNoResponse(command string) error {
	m.Commands = append(m.Commands, command)

	// Check if this command should fail
	for failCmd, shouldFail := range m.ShouldFail {
		if shouldFail && strings.Contains(command, failCmd) {
			return fmt.Errorf("mock failure for command: %s", command)
		}
	}

	return nil
}

// GetWorkspaces mocks retrieving the list of workspaces
func (m *SwayMock) GetWorkspaces() ([]string, error) {
	return m.MockWorkspaces, nil
}

// LaunchApp mocks launching an application
func (m *SwayMock) LaunchApp(appName string, markName string) error {
	m.LaunchedApps = append(m.LaunchedApps, appName)
	m.Marks[markName] = true

	// Check if this app should fail to launch
	if m.ShouldFail[appName] {
		return fmt.Errorf("mock failure launching app: %s", appName)
	}

	return nil
}

// GetCommandHistory returns all commands executed
func (m *SwayMock) GetCommandHistory() []string {
	return m.Commands
}

// GetMarkHistory returns all marks set
func (m *SwayMock) GetMarkHistory() []string {
	var marks []string
	for mark := range m.Marks {
		marks = append(marks, mark)
	}
	return marks
}

// GetFocusHistory returns all focus commands
func (m *SwayMock) GetFocusHistory() []string {
	return m.Focuses
}

// GetLaunchedApps returns all launched apps
func (m *SwayMock) GetLaunchedApps() []string {
	return m.LaunchedApps
}

// MockWorkspaceJSON generates mock workspace JSON output
func (m *SwayMock) MockWorkspaceJSON() string {
	type workspace struct {
		Name string `json:"name"`
	}

	workspaces := make([]workspace, len(m.MockWorkspaces))
	for i, name := range m.MockWorkspaces {
		workspaces[i] = workspace{Name: name}
	}

	jsonBytes, _ := json.Marshal(workspaces)
	return string(jsonBytes)
}

// MockMarksJSON generates mock marks JSON output
func (m *SwayMock) MockMarksJSON() string {
	marks := make([]string, 0, len(m.Marks))
	for mark := range m.Marks {
		marks = append(marks, mark)
	}

	jsonBytes, _ := json.Marshal(marks)
	return string(jsonBytes)
}

// SetShouldFail configures a command to fail
func (m *SwayMock) SetShouldFail(command string, shouldFail bool) {
	m.ShouldFail[command] = shouldFail
}

// Reset clears all stored state in the mock
func (m *SwayMock) Reset() {
	m.Commands = []string{}
	m.Marks = make(map[string]bool)
	m.Focuses = []string{}
	m.CurrentWorkspace = ""
	m.LaunchedApps = []string{}
	m.ShouldFail = make(map[string]bool)
}
