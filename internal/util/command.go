package util

import (
	"fmt"
	"os/exec"
	"strings"

	errs "github.com/titembaatar/sway.flem/internal/errors"
	"github.com/titembaatar/sway.flem/internal/log"
)

func CheckCommand(command string) error {
	log.Debug("Checking if %s is available", command)

	cmd := exec.Command(command, "-v")
	output, err := cmd.CombinedOutput()

	if err != nil {
		if len(output) > 0 {
			log.Error("Command output: %s", string(output))
		}

		return errs.NewFatal(errs.ErrCommandNotFound,
			fmt.Sprintf("Command '%s' is not available", command)).
			WithSuggestion(fmt.Sprintf("Make sure %s is installed and available in your PATH", command))
	}

	outputStr := strings.TrimSpace(string(output))
	log.Debug("Command '%s' is available, version: %s", command, outputStr)

	return nil
}

func ExecuteCommand(cmdStr string) error {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	if err := ValidateCommand(parts[0]); err != nil {
		return err
	}

	var cmd *exec.Cmd
	if len(parts) == 1 {
		cmd = exec.Command(parts[0])
	} else {
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	log.Debug("Executing command: %s", cmdStr)
	return cmd.Start()
}

// ValidateCommand checks if a command is available in PATH
// Returns nil if the command is available or contains path separators
func ValidateCommand(command string) error {
	// If command path is absolute or relative, don't check PATH
	if strings.ContainsAny(command, "/\\") {
		return nil
	}

	_, err := exec.LookPath(command)
	if err != nil {
		return errs.New(errs.ErrCommandNotFound,
			fmt.Sprintf("Command '%s' not found in PATH", command)).
			WithSuggestion(fmt.Sprintf("Make sure '%s' is installed and available in your PATH", command))
	}

	return nil
}
