package sway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/titembaatar/sway.flem/internal/log"
)

type CommandResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type SwayCmd struct {
	Command        string // The command to execute
	Type           string // Command type: command, get_workspaces, get_marks, etc.
	Raw            bool   // Whether to use -r flag for raw output
	ExpectJSON     bool   // Whether the response is expected to be JSON
	ErrorsNonFatal bool   // Whether errors should be treated as non-fatal
}

func NewSwayCmd(command string) *SwayCmd {
	return &SwayCmd{
		Command:        command,
		Type:           "command",
		Raw:            false,
		ExpectJSON:     true,
		ErrorsNonFatal: false,
	}
}

func NewSwayCmdType(command string, cmdType string) *SwayCmd {
	cmd := NewSwayCmd(command)
	cmd.Type = cmdType
	return cmd
}

func NewRawSwayCmd(command string, cmdType string) *SwayCmd {
	cmd := NewSwayCmdType(command, cmdType)
	cmd.Raw = true
	return cmd
}

func (c *SwayCmd) WithType(cmdType string) *SwayCmd {
	c.Type = cmdType
	return c
}

func (c *SwayCmd) WithRaw(raw bool) *SwayCmd {
	c.Raw = raw
	return c
}

func (c *SwayCmd) WithExpectJSON(expectJSON bool) *SwayCmd {
	c.ExpectJSON = expectJSON
	return c
}

func (c *SwayCmd) WithErrorsNonFatal(nonFatal bool) *SwayCmd {
	c.ErrorsNonFatal = nonFatal
	return c
}

func (c *SwayCmd) Run() ([]CommandResponse, error) {
	log.SetComponent(log.ComponentSway)
	log.Debug("Executing sway command: %s", c.Command)

	cmdOp := log.Operation(fmt.Sprintf("sway command '%s'", c.Command))
	cmdOp.Begin()

	args := []string{"-t", c.Type}

	if c.Raw {
		args = append(args, "-r")
	}

	// Add -- to prevent swaymsg from interpreting args
	args = append(args, "--", c.Command)

	log.Debug("Full swaymsg command: swaymsg %s", strings.Join(args, " "))
	cmd := exec.Command("swaymsg", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		log.Error("Failed to execute sway command '%s': %v", c.Command, err)
		if errMsg != "" {
			log.Error("Command stderr: %s", errMsg)
		}
		cmdOp.EndWithError(err)
		return nil, NewSwayCommandError(c.Command, err, errMsg)
	}

	// If we don't expect JSON, just return empty response
	if !c.ExpectJSON {
		cmdOp.End()
		return []CommandResponse{{Success: true}}, nil
	}

	// Parse the JSON response
	var responses []CommandResponse
	if err := json.Unmarshal(stdout.Bytes(), &responses); err != nil {
		log.Error("Failed to parse response for command '%s': %v", c.Command, err)
		log.Debug("Raw response: %s", stdout.String())
		cmdOp.EndWithError(err)
		return nil, fmt.Errorf("failed to parse sway command response: %w", err)
	}

	// Check for command success
	for i, resp := range responses {
		if !resp.Success && !c.ErrorsNonFatal {
			log.Error("Sway command '%s' failed: %s", c.Command, resp.Error)
			cmdErr := fmt.Errorf("sway command failed: %s", resp.Error)
			cmdOp.EndWithError(cmdErr)
			return responses, cmdErr
		}
		log.Debug("Command response %d: success=%v", i, resp.Success)
	}

	log.Debug("Successfully executed sway command '%s'", c.Command)
	cmdOp.End()
	return responses, nil
}

func (c *SwayCmd) GetJSON(v interface{}) error {
	c.Raw = true

	cmd := exec.Command("swaymsg", "-t", c.Type, "-r", "--", c.Command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		log.Error("Failed to execute sway %s command: %v", c.Type, err)
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
