package sway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/titembaatar/sway.flem/internal/log"
)

const (
	MarkPrefix = "sway_flem"
)

type MarkType string

const (
	MarkTypeWorkspace MarkType = "ws"
	MarkTypeContainer MarkType = "con"
	MarkTypeApp       MarkType = "app"
)

// Generates a mark for a workspace
func GenerateWorkspaceMark(workspaceName string) string {
	return fmt.Sprintf("%s_%s_%s", MarkPrefix, string(MarkTypeWorkspace), workspaceName)
}

// Generates a mark for a container
func GenerateContainerMark(workspaceMark string, containerID string) string {
	return fmt.Sprintf("%s_%s_%s", workspaceMark, string(MarkTypeContainer), containerID)
}

// Generates a mark for an application
func GenerateAppMark(containerMark string, appID string) string {
	return fmt.Sprintf("%s_%s_%s", containerMark, string(MarkTypeApp), appID)
}

// Applies a mark to the currently focused container
func ApplyMark(mark string) error {
	log.Debug("Applying mark '%s' to focused container", mark)
	command := fmt.Sprintf("mark --add %s", mark)
	_, err := RunCommand(command)
	return err
}

// Focuses the container with the specified mark
func FocusMark(mark string) error {
	log.Debug("Focusing container with mark '%s'", mark)
	command := fmt.Sprintf("[con_mark=%s] focus", mark)
	_, err := RunCommand(command)
	return err
}

// Retrieves all nodes with marks
func GetMarkedNodes() (map[string]bool, error) {
	log.Debug("Getting all marked nodes")

	cmd := exec.Command("swaymsg", "-t", "get_marks", "-r")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		log.Error("Failed to get marks: %v", err)
		if errMsg != "" {
			log.Error("Stderr: %s", errMsg)
		}
		return nil, fmt.Errorf("swaymsg error: %w: %s", err, errMsg)
	}

	var marks []string
	if err := json.Unmarshal(stdout.Bytes(), &marks); err != nil {
		log.Error("Failed to parse marks response: %v", err)
		return nil, fmt.Errorf("failed to parse marks response: %w", err)
	}

	// Create a map for quick lookups
	marksMap := make(map[string]bool)
	for _, mark := range marks {
		if strings.HasPrefix(mark, MarkPrefix) {
			marksMap[mark] = true
		}
	}

	log.Debug("Found %d marks with our prefix", len(marksMap))
	return marksMap, nil
}

// Checks if a mark exists
func IsMarkExist(mark string) (bool, error) {
	marks, err := GetMarkedNodes()
	if err != nil {
		return false, err
	}
	return marks[mark], nil
}

// Resizes the container with the specified mark based on layout
func ResizeMark(mark string, size string, layout string) error {
	log.Debug("Resizing container with mark '%s' to '%s' with layout '%s'", mark, size, layout)

	if err := FocusMark(mark); err != nil {
		return err
	}

	var dimension string
	if layout == "splith" || layout == "tabbed" {
		dimension = "width"
	} else if layout == "splitv" || layout == "stacking" {
		dimension = "height"
	} else {
		// Default to width for unknown layouts
		dimension = "width"
		log.Warn("Unknown layout for resizing: %s, defaulting to width", layout)
	}

	command := fmt.Sprintf("resize set %s %s", dimension, size)
	_, err := RunCommand(command)
	return err
}
