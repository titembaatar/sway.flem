package sway

import (
	"fmt"
	"regexp"
)

// WorkspaceInfo contains workspace details
type WorkspaceInfo struct {
	Number         int
	Name           string
	Output         string
	Layout         string
	Representation string
	AppOrder       []string
}

// SwitchToWorkspace focuses the specified workspace
func (c *Client) SwitchToWorkspace(num int) error {
	if c.Verbose {
		fmt.Printf("Switching to workspace %d\n", num)
	}
	return c.SwayCmd(fmt.Sprintf("workspace number %d", num))
}

// MoveWindowToWorkspace moves a window to specified workspace
func (c *Client) MoveWindowToWorkspace(id int64, num int) error {
	if c.Verbose {
		fmt.Printf("Moving window %d to workspace %d\n", id, num)
	}
	return c.SwayCmd(fmt.Sprintf("[id=%d] move container to workspace number %d", id, num))
}

// MoveWorkspaceToOutput moves a workspace to specified output
func (c *Client) MoveWorkspaceToOutput(num int, output string) error {
	if c.Verbose {
		fmt.Printf("Moving workspace %d to output %s\n", num, output)
	}
	return c.SwayCmd(fmt.Sprintf("workspace number %d, move workspace to output %s", num, output))
}

// GetWorkspaceInfo gets information about a workspace
func (c *Client) GetWorkspaceInfo(num int) (*WorkspaceInfo, error) {
	tree, err := c.GetTree()
	if err != nil {
		return nil, err
	}

	workspaces := tree.FindWorkspaces()
	ws, exists := workspaces[num]
	if !exists {
		return nil, fmt.Errorf("workspace %d not found", num)
	}

	info := &WorkspaceInfo{
		Number:         num,
		Name:           ws.Name,
		Output:         ws.Output,
		Layout:         ws.Layout,
		Representation: ws.Representation,
		AppOrder:       ExtractAppOrder(ws.Representation),
	}

	return info, nil
}

// SetWorkspaceLayout sets the layout for the current workspace
func (c *Client) SetWorkspaceLayout(layout string) error {
	if !IsValidLayout(layout) {
		return fmt.Errorf("invalid layout: %s", layout)
	}

	if c.Verbose {
		fmt.Printf("Setting workspace layout to %s\n", layout)
	}

	return c.SetLayout(layout)
}

// IsNumericWorkspace checks if a workspace name is numeric
func IsNumericWorkspace(name string) bool {
	return regexp.MustCompile(`^\d+$`).MatchString(name)
}

