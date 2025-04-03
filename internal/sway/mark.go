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

// Generate a mark for an app
func generateAppMark(workspaceName string, depth, containerID, appIndex int) string {
	if depth == 0 {
		// Top-level app
		return fmt.Sprintf("ws_%s_app_%d", workspaceName, appIndex+1)
	}

	// App in a container
	return fmt.Sprintf("ws_%s_con_%d_app_%d", workspaceName, containerID, appIndex+1)
}

// Generate a mark for a container
func generateContainerMark(workspaceName string, containerID int) string {
	return fmt.Sprintf("ws_%s_con_%d", workspaceName, containerID)
}

// Applies a mark to the currently focused container
func ApplyMark(mark string) error {
	log.Debug("Applying mark '%s' to focused container", mark)
	command := fmt.Sprintf("mark --add %s", mark)

	_, err := RunCommand(command)
	if err != nil {
		return NewMarkError(mark, fmt.Errorf("%w: %v", ErrMarkingFailed, err))
	}

	return nil
}

// Retrieves all nodes with marks (for debugging purposes)
func GetMarkedNodes() ([]string, error) {
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

	// Filter for our marks
	var ourMarks []string
	for _, mark := range marks {
		if strings.HasPrefix(mark, "ws_") {
			ourMarks = append(ourMarks, mark)
		}
	}

	log.Debug("Found %d marks for our applications", len(ourMarks))
	return ourMarks, nil
}
