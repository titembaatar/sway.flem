package manager

import (
	"log"
	"os/exec"
	"time"
)

func (m *Manager) logVerbose(format string, args ...any) {
	if m.Verbose {
		log.Printf(format, args...)
	}
}

func (m *Manager) delay(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

func (m *Manager) runCommand(cmd string) error {
	command := exec.Command("sh", "-c", cmd)
	return command.Run()
}

func (m *Manager) runCommandAsync(cmd string) error {
	command := exec.Command("sh", "-c", cmd+" &")
	return command.Run()
}
