package sway

import (
	"fmt"
	"slices"

	errs "github.com/titembaatar/sway.flem/internal/errors"
	"github.com/titembaatar/sway.flem/internal/log"
)

type WindowInfo struct {
	ID          int      // Sway internal ID
	Name        string   // Window title
	AppID       string   // Wayland app_id
	Class       string   // X11 window class
	Instance    string   // X11 window instance
	PID         int      // Process ID
	Focused     bool     // Is the window focused
	Workspace   string   // Workspace name
	WorkspaceID int      // Workspace ID
	Marks       []string // Marks applied to this window
}

type Node struct {
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	Type             string            `json:"type"`
	AppID            string            `json:"app_id"`
	PID              int               `json:"pid"`
	Focused          bool              `json:"focused"`
	Marks            []string          `json:"marks"`
	WindowProperties *WindowProperties `json:"window_properties,omitempty"`
	Nodes            []Node            `json:"nodes"`
	FloatingNodes    []Node            `json:"floating_nodes"`
	workspace        string
	workspaceID      int
}

type WindowProperties struct {
	Class string `json:"class"`
}

func GetAllWindows(errorHandler *errs.ErrorHandler) ([]WindowInfo, error) {
	log.Debug("Retrieving all windows from Sway tree")

	var root Node
	swayCmd := NewSwayCmdType("", "get_tree")
	if errorHandler != nil {
		swayCmd.WithErrorHandler(errorHandler)
	}

	if err := swayCmd.GetJSON(&root); err != nil {
		treeErr := errs.Wrap(err, "Failed to get Sway tree information")

		if errorHandler != nil {
			errorHandler.Handle(treeErr)
		}

		return nil, treeErr
	}

	windows := collectWindows(root)
	log.Debug("Found %d windows in Sway tree", len(windows))
	return windows, nil
}

func collectWindows(node Node) []WindowInfo {
	var windows []WindowInfo

	if node.Type == "workspace" && node.Name != "__i3_scratch" {
		node.workspace = node.Name
		node.workspaceID = node.ID
	}

	if node.Type == "con" && isWindow(node) {
		window := WindowInfo{
			ID:          node.ID,
			Name:        node.Name,
			AppID:       node.AppID,
			PID:         node.PID,
			Focused:     node.Focused,
			Workspace:   node.workspace,
			WorkspaceID: node.workspaceID,
			Marks:       node.Marks,
		}

		if node.WindowProperties != nil {
			window.Class = node.WindowProperties.Class
		}

		windows = append(windows, window)
	}

	for _, child := range node.Nodes {
		child.workspace = node.workspace
		child.workspaceID = node.workspaceID
		windows = append(windows, collectWindows(child)...)
	}

	for _, child := range node.FloatingNodes {
		child.workspace = node.workspace
		child.workspaceID = node.workspaceID
		windows = append(windows, collectWindows(child)...)
	}

	return windows
}

func isWindow(node Node) bool {
	return node.AppID != "" || (node.WindowProperties != nil && node.WindowProperties.Class != "")
}

func FindWindowByMark(mark string, windows []WindowInfo) *WindowInfo {
	log.Debug("Searching for window with mark: %s", mark)

	for _, window := range windows {
		if slices.Contains(window.Marks, mark) {
			log.Debug("Found window with mark %s: ID=%d, Name=%s, Workspace=%s",
				mark, window.ID, window.Name, window.Workspace)
			return &window
		}
	}

	log.Debug("No window found with mark: %s", mark)
	return nil
}

func IsAppRunning(mark string) (bool, *WindowInfo, error) {
	windows, err := GetAllWindows(nil)
	if err != nil {
		return false, nil, err
	}

	window := FindWindowByMark(mark, windows)
	return window != nil, window, nil
}

func IsAppRunningWithErrorHandler(mark string, errorHandler *errs.ErrorHandler) (bool, *WindowInfo, error) {
	windows, err := GetAllWindows(errorHandler)
	if err != nil {
		markErr := errs.Wrap(err, fmt.Sprintf("Failed to check if app with mark '%s' is running", mark))

		if errorHandler != nil {
			errorHandler.Handle(markErr)
		}

		return false, nil, markErr
	}

	window := FindWindowByMark(mark, windows)
	return window != nil, window, nil
}
