package sway

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

var (
	ErrCommandFailed = errors.New("sway command failed")
)

type Client struct {
	verbose bool
}

func NewClient(verbose bool) *Client {
	return &Client{
		verbose: verbose,
	}
}

func (c *Client) ExecuteCommand(cmd string) error {
	if c.verbose {
		log.Printf("Executing: swaymsg %s\n", cmd)
	}

	out, err := exec.Command("swaymsg", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %v (output: %s)", ErrCommandFailed, err, out)
	}

	return nil
}

func (c *Client) GetTree() (*Node, error) {
	if c.verbose {
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
