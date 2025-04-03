package sway

import (
	"strings"

	"github.com/titembaatar/sway.flem/internal/log"
)

// Window in the Sway tree
type WindowInfo struct {
	ID      int    // Sway internal ID
	Name    string // Window title
	AppID   string // Wayland app_id
	Class   string // X11 window class
	PID     int    // Process ID
	Focused bool   // Is the window focused
}

// Sway tree JSON structure
type Node struct {
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	Type             string            `json:"type"`
	AppID            string            `json:"app_id"`
	PID              int               `json:"pid"`
	Focused          bool              `json:"focused"`
	WindowProperties *WindowProperties `json:"window_properties,omitempty"`
	Nodes            []Node            `json:"nodes"`
	FloatingNodes    []Node            `json:"floating_nodes"`
}

type WindowProperties struct {
	Class string `json:"class"`
}

func GetAllWindows() ([]WindowInfo, error) {
	log.Debug("Retrieving all windows from Sway tree")

	var root Node
	err := executeSwayGetJSON("", "get_tree", &root)
	if err != nil {
		return nil, err
	}

	windows := collectWindows(root)
	log.Debug("Found %d windows in Sway tree", len(windows))
	return windows, nil
}

func collectWindows(node Node) []WindowInfo {
	var windows []WindowInfo

	if node.Type == "con" {
		if isWindow(node) {
			window := WindowInfo{
				ID:      node.ID,
				Name:    node.Name,
				AppID:   node.AppID,
				PID:     node.PID,
				Focused: node.Focused,
			}

			if node.WindowProperties != nil {
				window.Class = node.WindowProperties.Class
			}

			windows = append(windows, window)
		}
	}

	for _, child := range node.Nodes {
		windows = append(windows, collectWindows(child)...)
	}

	for _, child := range node.FloatingNodes {
		windows = append(windows, collectWindows(child)...)
	}

	return windows
}

func isWindow(node Node) bool {
	return node.AppID != "" || (node.WindowProperties != nil && node.WindowProperties.Class != "")
}

func FindWindowsByApp(appName string, windows []WindowInfo) []WindowInfo {
	log.Debug("Searching for windows with app name: %s", appName)

	appNameLower := strings.ToLower(appName)
	var matches []WindowInfo

	for _, window := range windows {
		if window.AppID != "" && strings.ToLower(window.AppID) == appNameLower {
			log.Debug("Found matching window by app_id: %s (ID: %d, Name: %s)",
				window.AppID, window.ID, window.Name)
			matches = append(matches, window)
			continue
		}

		if window.Class != "" && strings.ToLower(window.Class) == appNameLower {
			log.Debug("Found matching window by class: %s (ID: %d, Name: %s)",
				window.Class, window.ID, window.Name)
			matches = append(matches, window)
		}
	}

	log.Debug("Found %d matching windows for app: %s", len(matches), appName)
	return matches
}

func IsAppRunning(appName string) (bool, []WindowInfo, error) {
	windows, err := GetAllWindows()
	if err != nil {
		return false, nil, err
	}

	matches := FindWindowsByApp(appName, windows)
	return len(matches) > 0, matches, nil
}
