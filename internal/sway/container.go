package sway

import (
	"fmt"
)

// SetSplit applies a split to the currently focused container
func (c *Client) SetSplit(splitType string) error {
	validSplits := map[string]string{
		"h":          "horizontal",
		"horizontal": "horizontal",
		"v":          "vertical",
		"vertical":   "vertical",
		"t":          "toggle",
		"toggle":     "toggle",
		"n":          "none",
		"none":       "none",
	}

	split, valid := validSplits[splitType]
	if !valid {
		return fmt.Errorf("invalid split type: %s", splitType)
	}

	if c.Verbose {
		fmt.Printf("Setting split: %s\n", split)
	}

	return c.ExecuteCommand(fmt.Sprintf("split %s", split))
}

// SetLayout sets the layout of the currently focused container
func (c *Client) SetLayout(layoutType string) error {
	validLayouts := map[string]bool{
		"splith":   true,
		"splitv":   true,
		"stacking": true,
		"tabbed":   true,
		"default":  true,
	}

	if !validLayouts[layoutType] {
		return fmt.Errorf("invalid layout type: %s", layoutType)
	}

	if c.Verbose {
		fmt.Printf("Setting layout: %s\n", layoutType)
	}

	return c.ExecuteCommand(fmt.Sprintf("layout %s", layoutType))
}

// IsValidLayout checks if a layout name is valid
func IsValidLayout(layout string) bool {
	validLayouts := map[string]bool{
		"splith":   true,
		"splitv":   true,
		"stacking": true,
		"tabbed":   true,
		"default":  true,
	}

	return validLayouts[layout]
}

// SetNodeLayout sets a layout for a specific node by ID
func (c *Client) SetNodeLayout(id int64, layout string) error {
	if !IsValidLayout(layout) {
		return fmt.Errorf("invalid layout: %s", layout)
	}

	if c.Verbose {
		fmt.Printf("Setting node %d layout to %s\n", id, layout)
	}

	// First focus the node, then set its layout
	if err := c.FocusWindow(id); err != nil {
		return fmt.Errorf("focusing node for layout change: %w", err)
	}

	return c.SetLayout(layout)
}

// SetNodeSplit sets a split for a specific node by ID
func (c *Client) SetNodeSplit(id int64, split string) error {
	if c.Verbose {
		fmt.Printf("Setting node %d split to %s\n", id, split)
	}

	// First focus the node, then set its split
	if err := c.FocusWindow(id); err != nil {
		return fmt.Errorf("focusing node for split change: %w", err)
	}

	return c.SetSplit(split)
}
