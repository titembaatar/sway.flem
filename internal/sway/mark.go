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

func (m Mark) String() string {
	return m.ID
}

func NewMark(id string) Mark {
	return Mark{ID: id}
}

func NewAppMark(workspaceName string, containerID, appIndex int) Mark {
	id := fmt.Sprintf("w%s_c%d_a%d", workspaceName, containerID, appIndex+1)
	return Mark{ID: id}
}

func NewContainerMark(workspaceName string, containerID int) Mark {
	return Mark{ID: fmt.Sprintf("w%s_c%d", workspaceName, containerID)}
}

func (m Mark) Focus() error {
	log.Debug("Focusing container with mark '%s'", m.ID)

	cmd := fmt.Sprintf("[con_mark=\"%s\"] focus", m.ID)
	swayCmd := NewSwayCmd(cmd)

	_, err := swayCmd.Run()

	if err != nil {
		return fmt.Errorf("failed to focus container with mark '%s': %w", m.ID, err)
	}

	return nil
}

func (m Mark) Resize(width string, height string) string {
	return fmt.Sprintf("resize set %s %s", width, height)
}

func (m Mark) Apply() error {
	log.Debug("Applying mark '%s' to focused container", m.ID)
	cmd := fmt.Sprintf("mark --add %s", m.ID)
	swayCmd := NewSwayCmd(cmd)

	_, err := swayCmd.Run()
	if err != nil {
		return NewMarkError(m.ID, fmt.Errorf("%w: %v", ErrMarkingFailed, err))
	}

	return nil
}

func (m Mark) IsApp() bool {
	return strings.Contains(m.ID, "_a") && strings.Contains(m.ID, "_c")
}

func (m Mark) IsContainer() bool {
	return strings.Contains(m.ID, "_c") && !strings.Contains(m.ID, "_a")
}

func (m Mark) GetWorkspace() string {
	if len(m.ID) < 2 || !strings.HasPrefix(m.ID, "w") {
		return ""
	}

	parts := strings.Split(m.ID, "_")
	if len(parts) == 0 {
		return ""
	}

	return strings.TrimPrefix(parts[0], "w")
}

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

	var marks []Mark
	for _, id := range markIDs {
		if len(id) > 2 && strings.HasPrefix(id, "w") && strings.Contains(id, "_c") {
			marks = append(marks, Mark{ID: id})
		}
	}

	log.Debug("Found %d marks for our applications", len(marks))
	return marks, nil
}
