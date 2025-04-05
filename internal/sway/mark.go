package sway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	errs "github.com/titembaatar/sway.flem/internal/errors"
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

func (m Mark) Focus(errorHandler *errs.ErrorHandler) error {
	log.Debug("Focusing container with mark '%s'", m.ID)

	cmd := fmt.Sprintf("[con_mark=\"%s\"] focus", m.ID)
	swayCmd := NewSwayCmd(cmd)
	if errorHandler != nil {
		swayCmd.WithErrorHandler(errorHandler)
	}

	_, err := swayCmd.Run()

	if err != nil {
		focusErr := errs.New(errs.ErrFocusFailed,
			fmt.Sprintf("Failed to focus container with mark '%s'", m.ID))
		focusErr.WithCategory("Sway")
		focusErr.WithSuggestion(fmt.Sprintf("Check that a container with mark '%s' exists", m.ID))

		if errorHandler != nil {
			errorHandler.Handle(focusErr)
		}

		return focusErr
	}

	return nil
}

func (m Mark) Resize(width string, height string) string {
	return fmt.Sprintf("resize set %s %s", width, height)
}

func (m Mark) Apply(errorHandler *errs.ErrorHandler) error {
	log.Debug("Applying mark '%s' to focused container", m.ID)
	cmd := fmt.Sprintf("mark --add %s", m.ID)
	swayCmd := NewSwayCmd(cmd)
	if errorHandler != nil {
		swayCmd.WithErrorHandler(errorHandler)
	}

	_, err := swayCmd.Run()
	if err != nil {
		markErr := errs.NewMarkError(m.ID, errs.ErrMarkingFailed)
		markErr.WithSuggestion("Make sure a container is focused before trying to mark it")

		if errorHandler != nil {
			errorHandler.Handle(markErr)
		}

		return markErr
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

func GetMarkedNodes(errorHandler *errs.ErrorHandler) ([]Mark, error) {
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

		cmdErr := errs.NewSwayCommandError("get_marks", err, errMsg)

		if errorHandler != nil {
			errorHandler.Handle(cmdErr)
		}

		return nil, cmdErr
	}

	var markIDs []string
	if err := json.Unmarshal(stdout.Bytes(), &markIDs); err != nil {
		log.Error("Failed to parse marks response: %v", err)

		parseErr := errs.New(err, "Failed to parse marks response")
		parseErr.WithCategory("Sway")

		if errorHandler != nil {
			errorHandler.Handle(parseErr)
		}

		return nil, parseErr
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
