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

type CommandResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type SwayCmd struct {
	Command        string             // The command to execute
	Type           string             // Command type: command, get_workspaces, get_marks, etc.
	Raw            bool               // Whether to use -r flag for raw output
	ExpectJSON     bool               // Whether the response is expected to be JSON
	ErrorsNonFatal bool               // Whether errors should be treated as non-fatal
	ErrorHandler   *errs.ErrorHandler // Error handler instance
}

func NewSwayCmd(command string) *SwayCmd {
	return &SwayCmd{
		Command:        command,
		Type:           "command",
		Raw:            false,
		ExpectJSON:     true,
		ErrorsNonFatal: false,
		ErrorHandler:   nil,
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

func (c *SwayCmd) WithErrorHandler(handler *errs.ErrorHandler) *SwayCmd {
	c.ErrorHandler = handler
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

		// Create a Sway command error
		cmdErr := errs.NewSwayCommandError(c.Command, err, errMsg)

		// Set severity based on the non-fatal flag
		if c.ErrorsNonFatal {
			cmdErr.Severity = errs.SeverityWarning
		} else {
			cmdErr.Severity = errs.SeverityError
		}

		// Handle the error if a handler is provided
		if c.ErrorHandler != nil {
			c.ErrorHandler.Handle(cmdErr)
		}

		cmdOp.EndWithError(cmdErr)
		return nil, cmdErr
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

		parseErr := errs.New(err, fmt.Sprintf("Failed to parse sway command response for '%s'", c.Command))
		parseErr.WithCategory("Sway")

		// Handle the error if a handler is provided
		if c.ErrorHandler != nil {
			c.ErrorHandler.Handle(parseErr)
		}

		cmdOp.EndWithError(parseErr)
		return nil, parseErr
	}

	// Check for command success
	for i, resp := range responses {
		if !resp.Success && !c.ErrorsNonFatal {
			log.Error("Sway command '%s' failed: %s", c.Command, resp.Error)

			cmdErr := errs.New(errs.ErrCommandFailed, fmt.Sprintf("Sway command '%s' failed: %s", c.Command, resp.Error))
			cmdErr.WithCategory("Sway")

			// Handle the error if a handler is provided
			if c.ErrorHandler != nil {
				c.ErrorHandler.Handle(cmdErr)
			}

			cmdOp.EndWithError(cmdErr)
			return responses, cmdErr
		}
		log.Debug("Command response %d: success=%v", i, resp.Success)
	}

	log.Debug("Successfully executed sway command '%s'", c.Command)
	cmdOp.End()
	return responses, nil
}

func RunSwayCmd(command string) ([]CommandResponse, error) {
	cmd := NewSwayCmd(command)
	return cmd.Run()
}

func (c *SwayCmd) GetJSON(v any) error {
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

		cmdErr := errs.NewSwayCommandError(c.Command, err, errMsg)

		// Handle the error if a handler is provided
		if c.ErrorHandler != nil {
			c.ErrorHandler.Handle(cmdErr)
		}

		return cmdErr
	}

	if err := json.Unmarshal(stdout.Bytes(), v); err != nil {
		log.Error("Failed to parse response: %v", err)
		log.Debug("Raw response: %s", stdout.String())

		parseErr := errs.New(err, fmt.Sprintf("Failed to parse JSON response for sway command '%s'", c.Command))
		parseErr.WithCategory("Sway")

		// Handle the error if a handler is provided
		if c.ErrorHandler != nil {
			c.ErrorHandler.Handle(parseErr)
		}

		return parseErr
	}

	return nil
}

// Run a sway command with an error handler
func RunSwayCommandWithErrorHandler(command string, errorHandler *errs.ErrorHandler) ([]CommandResponse, error) {
	cmd := NewSwayCmd(command).WithErrorHandler(errorHandler)
	return cmd.Run()
}
