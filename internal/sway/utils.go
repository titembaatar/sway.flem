package sway

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func IsNumericWorkspace(name string) bool {
	return regexp.MustCompile(`^\d+$`).MatchString(name)
}

func FormatPosition(position string) string {
	specialPositions := map[string]string{
		"center":  "center",
		"middle":  "center",
		"top":     "0 0",
		"bottom":  "0 999999",
		"left":    "0 center",
		"right":   "999999 center",
		"pointer": "cursor",
		"cursor":  "cursor",
		"mouse":   "cursor",
	}

	if formatted, ok := specialPositions[strings.ToLower(position)]; ok {
		return formatted
	}

	return position
}

func (c *Client) WaitForWindowAppearance(appName string, timeoutSec int) (int64, error) {
	start := time.Now()
	timeoutDuration := time.Duration(timeoutSec) * time.Second

	// Poll for the window
	for {
		if time.Since(start) > timeoutDuration {
			return 0, fmt.Errorf("timeout waiting for %s to appear", appName)
		}

		tree, err := c.GetTree()
		if err != nil {
			return 0, err
		}

		// Search for the app in all workspaces
		workspaces := tree.FindWorkspaces()
		for _, ws := range workspaces {
			apps := ws.FindAllApps()
			for _, app := range apps {
				if MatchAppName(app.Name, appName) {
					return app.NodeID, nil
				}
			}
		}

		// Wait a bit before checking again
		time.Sleep(100 * time.Millisecond)
	}
}

func IsValidLayout(layout string) bool {
	validLayouts := map[string]bool{
		"splith":   true,
		"splitv":   true,
		"stacking": true,
		"tabbed":   true,
	}

	return validLayouts[layout]
}

func ParseSize(size string) (width, height string, err error) {
	parts := strings.Fields(size)

	if len(parts) == 0 {
		return "", "", fmt.Errorf("empty size string")
	}

	if len(parts) == 1 {
		// Use same value for both dimensions
		return parts[0], parts[0], nil
	}

	if len(parts) >= 2 {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("invalid size format: %s", size)
}
