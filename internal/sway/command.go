package sway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/titembaatar/sway.flem/internal/log"
)

type CommandResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// Executes a swaymsg command and returns the result
func RunCommand(command string) ([]CommandResponse, error) {
	log.Debug("Executing sway command: %s", command)

	cmd := exec.Command("swaymsg", "-t", "command", "--", command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		log.Error("Failed to execute sway command: %v", err)
		if errMsg != "" {
			log.Error("Stderr: %s", errMsg)
		}
		return nil, fmt.Errorf("swaymsg error: %w: %s", err, errMsg)
	}

	// Parse the JSON response
	var responses []CommandResponse
	if err := json.Unmarshal(stdout.Bytes(), &responses); err != nil {
		log.Error("Failed to parse sway command response: %v", err)
		log.Debug("Raw response: %s", stdout.String())
		return nil, fmt.Errorf("failed to parse sway command response: %w", err)
	}

	// Check for command success
	for i, resp := range responses {
		if !resp.Success {
			log.Error("Sway command failed: %s", resp.Error)
			return responses, fmt.Errorf("sway command failed: %s", resp.Error)
		}
		log.Debug("Command response %d: success=%v", i, resp.Success)
	}

	return responses, nil
}

// Retrieves the list of workspaces from sway
func GetWorkspaces() ([]string, error) {
	log.Debug("Getting workspaces from sway")

	cmd := exec.Command("swaymsg", "-t", "get_workspaces", "-r")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		log.Error("Failed to get workspaces: %v", err)
		if errMsg != "" {
			log.Error("Stderr: %s", errMsg)
		}
		return nil, fmt.Errorf("swaymsg error: %w: %s", err, errMsg)
	}

	type workspace struct {
		Name string `json:"name"`
	}

	var workspaces []workspace
	if err := json.Unmarshal(stdout.Bytes(), &workspaces); err != nil {
		log.Error("Failed to parse workspace response: %v", err)
		log.Debug("Raw response: %s", stdout.String())
		return nil, fmt.Errorf("failed to parse workspace response: %w", err)
	}

	names := make([]string, len(workspaces))
	for i, ws := range workspaces {
		names[i] = ws.Name
	}

	log.Debug("Found %d workspaces: %s", len(names), strings.Join(names, ", "))
	return names, nil
}

// Switches to the specified workspace
func SwitchToWorkspace(workspace string) error {
	command := fmt.Sprintf("workspace %s", workspace)
	_, err := RunCommand(command)
	return err
}

// Creates a new workspace with the specified name and layout
func CreateWorkspace(name string, layout string) error {
	if err := SwitchToWorkspace(name); err != nil {
		return err
	}

	command := fmt.Sprintf("layout %s", layout)
	_, err := RunCommand(command)
	return err
}

// Switches focus to each of the specified workspaces in order.
func FocusWorkspaces(workspaces []string) error {
	log.Info("Focusing on %d workspaces", len(workspaces))
	var errors []string

	for i, workspace := range workspaces {
		log.Debug("Focusing on workspace: %s", workspace)
		if err := SwitchToWorkspace(workspace); err != nil {
			log.Error("Failed to focus on workspace %s: %v", workspace, err)
			errors = append(errors, fmt.Sprintf("workspace %s: %v", workspace, err))
		} else {
			log.Info("Successfully focused on workspace %s (%d of %d)", workspace, i+1, len(workspaces))
			time.Sleep(100 * time.Millisecond)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to focus on some workspaces: %s", strings.Join(errors, "; "))
	}

	return nil
}
