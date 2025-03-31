package sway

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

// Error constants
var (
	ErrCommandFailed = errors.New("sway command failed")
)

// Sway IPC
type Client struct {
	Verbose bool
}

func NewClient(verbose bool) *Client {
	return &Client{
		Verbose: verbose,
	}
}

// Runs a Sway command via swaymsg
func (c *Client) SwayCmd(cmd string) error {
	if c.Verbose {
		log.Printf("Sway command: %s", cmd)
	}

	out, err := exec.Command("swaymsg", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %v (output: %s)", ErrCommandFailed, err, out)
	}

	return nil
}

// Retrieves the complete Sway tree
func (c *Client) GetTree() (*Node, error) {
	if c.Verbose {
		log.Println("Getting Sway tree")
	}

	cmd := exec.Command("swaymsg", "-t", "get_tree", "--raw")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("getting workspace tree: %w", err)
	}

	var rootNode Node
	if err := json.Unmarshal(output, &rootNode); err != nil {
		return nil, fmt.Errorf("parsing workspace tree: %w", err)
	}

	return &rootNode, nil
}
