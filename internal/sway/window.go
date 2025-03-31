package sway

import (
	"fmt"
	"time"
)

type WindowInfo struct {
	ID int64
}

// Focuses a window by its ID
func (c *Client) FocusWindow(id int64) error {
	if c.Verbose {
		fmt.Printf("Focusing window ID: %d\n", id)
	}
	return c.SwayCmd(fmt.Sprintf("[id=%d] focus", id))
}

// Resize a window by its ID
func (c *Client) ResizeWindow(id int64, size string, layout string) error {
	if size == "" {
		return nil
	}

	if c.Verbose {
		fmt.Printf("Resizing tiled window ID %d to %s using layout %s\n", id, size, layout)
	}

	switch layout {
	case "splith", "tabbed":
		return c.SwayCmd(fmt.Sprintf("[id=%d] resize set width %s", id, size))
	case "splitv", "stacking":
		return c.SwayCmd(fmt.Sprintf("[id=%d] resize set height %s", id, size))
	default:
		return c.SwayCmd(fmt.Sprintf("[id=%d] resize set %s %s", id, size, size))
	}
}

// Closes a window by its ID
func (c *Client) KillWindow(id int64) error {
	if c.Verbose {
		fmt.Printf("Killing window ID: %d\n", id)
	}
	return c.SwayCmd(fmt.Sprintf("[id=%d] kill", id))
}

// Searches for a window by id
func (c *Client) FindWindow(id int64) (*WindowInfo, error) {
	tree, err := c.GetTree()
	if err != nil {
		return nil, err
	}

	apps := FindApps(tree)
	for _, app := range apps {
		if MatchApp(app.ID, id) {
			return &app, nil
		}
	}

	return nil, fmt.Errorf("window not found: %d", id)
}

// Wait until a window appears
func (c *Client) WaitForWindow(name string, timeoutSec int) (*WindowInfo, error) {
	start := time.Now()
	timeout := time.Duration(timeoutSec) * time.Second

	if c.Verbose {
		fmt.Printf("Waiting for window '%s' to appear (timeout: %ds)\n", name, timeoutSec)
	}

	for time.Since(start) < timeout {
		app, err := c.FindWindow(name)
		if err == nil {
			return app, nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil, fmt.Errorf("timeout waiting for window: %s", name)
}
