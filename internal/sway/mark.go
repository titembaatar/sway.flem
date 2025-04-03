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
	MarkPrefix = "flem"
)

type Mark struct {
	ID string
}

func NewMark(id string) Mark {
	return Mark{ID: id}
}

// Creates a mark for an application in a workspace
func NewAppMark(workspaceName string, depth, containerID, appIndex int) Mark {
	var id string
	if depth == 0 {
		// Top-level app
		id = fmt.Sprintf("ws_%s_app_%d", workspaceName, appIndex+1)
	} else {
		// App in a container
		id = fmt.Sprintf("ws_%s_con_%d_app_%d", workspaceName, containerID, appIndex+1)
	}
	return Mark{ID: id}
}

// Creates a mark for a container in a workspace
func NewContainerMark(workspaceName string, containerID int) Mark {
	return Mark{ID: fmt.Sprintf("ws_%s_con_%d", workspaceName, containerID)}
}

// String representation of the mark
func (m Mark) String() string {
	return m.ID
}

// Focus a container with this mark
func (m Mark) FocusCmd() string {
	return fmt.Sprintf("[con_mark=\"%s\"] focus", m.ID)
}

// Resize a container with this mark
func (m Mark) ResizeCmd(dimension string, size string) string {
	return fmt.Sprintf("resize set %s %s", dimension, size)
}

// Applies the mark to the currently focused container
func (m Mark) Apply() error {
	log.Debug("Applying mark '%s' to focused container", m.ID)
	command := fmt.Sprintf("mark --add %s", m.ID)

	_, err := RunCommand(command)
	if err != nil {
		return NewMarkError(m.ID, fmt.Errorf("%w: %v", ErrMarkingFailed, err))
	}

	return nil
}

// Focuses the container with this mark
func (m Mark) Focus() error {
	log.Debug("Focusing container with mark '%s'", m.ID)
	_, err := RunCommand(m.FocusCmd())
	if err != nil {
		return fmt.Errorf("failed to focus container with mark '%s': %w", m.ID, err)
	}
	return nil
}

// Does the mark represents a top-level workspace app
func (m Mark) IsWorkspaceApp() bool {
	return strings.Contains(m.ID, "_app_") && !strings.Contains(m.ID, "_con_")
}

// Does the mark represents an app inside a container
func (m Mark) IsContainerApp() bool {
	return strings.Contains(m.ID, "_app_") && strings.Contains(m.ID, "_con_")
}

// Does the mark represents a container
func (m Mark) IsContainer() bool {
	return strings.Contains(m.ID, "_con_") && !strings.Contains(m.ID, "_app_")
}

// Extracts the workspace name from the mark
func (m Mark) GetWorkspace() string {
	parts := strings.Split(m.ID, "_")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// Retrieves all nodes with marks (for debugging purposes)
func GetMarkedNodes() ([]Mark, error) {
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

	var markIDs []string
	if err := json.Unmarshal(stdout.Bytes(), &markIDs); err != nil {
		log.Error("Failed to parse marks response: %v", err)
		return nil, fmt.Errorf("failed to parse marks response: %w", err)
	}

	// Filter for our marks
	var marks []Mark
	for _, id := range markIDs {
		if strings.HasPrefix(id, "ws_") {
			marks = append(marks, Mark{ID: id})
		}
	}

	log.Debug("Found %d marks for our applications", len(marks))
	return marks, nil
}
