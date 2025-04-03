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

// Options for executing sway commands
type SwayCommandOptions struct {
	Type           string // command, get_workspaces, get_marks, etc.
	Raw            bool   // Whether to use -r flag for raw output
	ExpectJSON     bool   // Whether the response is expected to be JSON
	ErrorsNonFatal bool   // Whether errors should be treated as non-fatal
}

// Returns standard options for regular sway commands
func DefaultCommandOptions() SwayCommandOptions {
	return SwayCommandOptions{
		Type:           "command",
		Raw:            false,
		ExpectJSON:     true,
		ErrorsNonFatal: false,
	}
}

// Executes a swaymsg command and returns the result
func RunCommand(command string) ([]CommandResponse, error) {
	log.Debug("Executing sway command: %s", command)

	opts := DefaultCommandOptions()
	return executeSwaymsg(command, opts)
}

// Helper for executing swaymsg commands
func executeSwaymsg(command string, opts SwayCommandOptions) ([]CommandResponse, error) {
	args := []string{"-t", opts.Type}

	if opts.Raw {
		args = append(args, "-r")
	}

	// Add -- to prevent swaymsg from interpreting args
	args = append(args, "--", command)

	cmd := exec.Command("swaymsg", args...)

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
		return nil, NewSwayCommandError(command, err, errMsg)
	}

	// If we don't expect JSON, just return empty response
	if !opts.ExpectJSON {
		return []CommandResponse{{Success: true}}, nil
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
		if !resp.Success && !opts.ErrorsNonFatal {
			log.Error("Sway command failed: %s", resp.Error)
			return responses, fmt.Errorf("sway command failed: %s", resp.Error)
		}
		log.Debug("Command response %d: success=%v", i, resp.Success)
	}

	return responses, nil
}

// Helper for sway commands that return JSON data
func executeSwayGetJSON(command string, outputType string, v any) error {
	// Always use raw output for JSON commands
	args := []string{"-t", outputType, "-r"}
	if command != "" {
		args = append(args, "--", command)
	}

	cmd := exec.Command("swaymsg", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		log.Error("Failed to execute sway %s command: %v", outputType, err)
		if errMsg != "" {
			log.Error("Stderr: %s", errMsg)
		}
		return fmt.Errorf("swaymsg error: %w: %s", err, errMsg)
	}

	if err := json.Unmarshal(stdout.Bytes(), v); err != nil {
		log.Error("Failed to parse response: %v", err)
		log.Debug("Raw response: %s", stdout.String())
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

// Retrieves the list of workspaces from sway
func GetWorkspaces() ([]string, error) {
	log.Debug("Getting workspaces from sway")

	type workspace struct {
		Name string `json:"name"`
	}

	var workspaces []workspace
	if err := executeSwayGetJSON("", "get_workspaces", &workspaces); err != nil {
		return nil, err
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
		return fmt.Errorf("%w: failed to switch to workspace '%s': %v",
			ErrWorkspaceCreateFailed, name, err)
	}

	command := fmt.Sprintf("layout %s", layout)
	_, err := RunCommand(command)
	if err != nil {
		return fmt.Errorf("%w: failed to set layout '%s' for workspace '%s': %v",
			ErrWorkspaceCreateFailed, layout, name, err)
	}

	log.Info("Successfully created workspace '%s' with layout '%s'", name, layout)
	return nil
}

// Focuses a container with the specified mark
func FocusByMark(mark string) error {
	command := fmt.Sprintf("[con_mark=\"%s\"] focus", mark)
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
