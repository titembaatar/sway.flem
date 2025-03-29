package sway

import (
	"fmt"
	"strings"
)

func (c *Client) MoveWorkspaceToOutput(num int, output string) error {
	return c.ExecuteCommand(fmt.Sprintf("workspace number %d, move workspace to output %s", num, output))
}

func (c *Client) SwitchToWorkspace(num int) error {
	return c.ExecuteCommand(fmt.Sprintf("workspace number %d", num))
}

func (c *Client) SetWorkspaceLayout(layout string) error {
	return c.ExecuteCommand(fmt.Sprintf("layout %s", layout))
}

func (c *Client) KillWindow(id int64) error {
	return c.ExecuteCommand(fmt.Sprintf("[con_id=%d] kill", id))
}

func (c *Client) SetFloating(id int64, floating bool) error {
	action := "enable"
	if !floating {
		action = "disable"
	}
	return c.ExecuteCommand(fmt.Sprintf("[con_id=%d] floating %s", id, action))
}

func (c *Client) MoveWindow(id int64, position string) error {
	var cmd string

	switch position {
	case "top":
		cmd = fmt.Sprintf("[con_id=%d] move position 0 0", id)
	case "bottom":
		cmd = fmt.Sprintf("[con_id=%d] move position 0 999999", id)
	case "left":
		cmd = fmt.Sprintf("[con_id=%d] move position 0 center", id)
	case "right":
		cmd = fmt.Sprintf("[con_id=%d] move position 999999 center", id)
	case "center", "middle":
		cmd = fmt.Sprintf("[con_id=%d] move position center", id)
	default:
		// Assume coordinates
		cmd = fmt.Sprintf("[con_id=%d] move position %s", id, position)
	}

	return c.ExecuteCommand(cmd)
}

func (c *Client) ResizeWindow(id int64, size string, isFloating bool, layout string) error {
	parts := strings.Fields(size)

	if isFloating {
		if len(parts) == 1 {
			// If only one value is provided, use it for both width and height
			return c.ExecuteCommand(fmt.Sprintf("[con_id=%d] resize set %s %s", id, parts[0], parts[0]))
		}
		return c.ExecuteCommand(fmt.Sprintf("[con_id=%d] resize set %s", id, size))
	}

	switch layout {
	case "tabbed", "stacking":
		return nil // These layouts don't support resizing individual windows
	case "splitv":
		// In vertical split, only set height
		if len(parts) == 1 {
			return c.ExecuteCommand(fmt.Sprintf("[con_id=%d] resize set height %s", id, parts[0]))
		} else if len(parts) >= 2 {
			return c.ExecuteCommand(fmt.Sprintf("[con_id=%d] resize set height %s", id, parts[1]))
		}
	case "splith", "":
		// In horizontal split, only set width
		if len(parts) >= 1 {
			return c.ExecuteCommand(fmt.Sprintf("[con_id=%d] resize set width %s", id, parts[0]))
		}
	}

	return fmt.Errorf("cannot apply size: invalid layout or size format")
}
