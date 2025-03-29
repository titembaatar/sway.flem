package util

import (
	"fmt"
	"os/exec"
	"time"
)

func RunCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	return cmd.Run()
}

func RunCommandWithOutput(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func RunCommandAsync(command string) error {
	cmd := exec.Command("sh", "-c", command+" &")
	return cmd.Run()
}

func RunCommandWithTimeout(command string, timeoutMs int) (string, error) {
	if timeoutMs <= 0 {
		return RunCommandWithOutput(command)
	}

	done := make(chan error, 1)
	var output []byte
	var err error

	// Start the command
	cmd := exec.Command("sh", "-c", command)
	go func() {
		output, err = cmd.CombinedOutput()
		done <- err
	}()

	// Wait for command completion or timeout
	select {
	case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
		// Kill the process if it times out
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("command timed out after %d ms", timeoutMs)
	case err := <-done:
		return string(output), err
	}
}

func DelayExecution(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}
