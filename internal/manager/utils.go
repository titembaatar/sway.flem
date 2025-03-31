package manager

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

// logDebug logs a debug message if verbose mode is enabled
func (m *Manager) logDebug(format string, args ...any) {
	if m.Verbose {
		log.Printf("[DEBUG]: "+format, args...)
	}
}

// logInfo logs an informational message
func (m *Manager) logInfo(format string, args ...any) {
	log.Printf("[INFO ]: "+format, args...)
}

// logWarn logs a warning message
func (m *Manager) logWarn(format string, args ...any) {
	log.Printf("[WARN ]: "+format, args...)
}

// logError logs an error message
func (m *Manager) logError(format string, args ...any) {
	log.Printf("[ERROR]: "+format, args...)
}

// delay pauses execution for the specified number of milliseconds
func (m *Manager) delay(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

// runCommand executes a shell command and waits for it to complete
func (m *Manager) runCommand(cmd string) error {
	m.logDebug("Running command: %s", cmd)
	command := exec.Command("sh", "-c", cmd)
	return command.Run()
}

// runCommandAsync executes a shell command in the background
func (m *Manager) runCommandAsync(cmd string) error {
	m.logDebug("Running command asynchronously: %s", cmd)
	command := exec.Command("sh", "-c", cmd+" &")
	return command.Run()
}

// handleError logs an error and optionally returns it based on severity
// If critical is true, the error is returned for the caller to handle
// If critical is false, the error is just logged as a warning
func (m *Manager) handleError(err error, message string, critical bool) error {
	if err == nil {
		return nil
	}

	fullMessage := fmt.Sprintf("%s: %v", message, err)

	if critical {
		m.logError("%s", fullMessage)
		return fmt.Errorf(fullMessage)
	} else {
		m.logWarn("%s", fullMessage)
		return nil
	}
}
