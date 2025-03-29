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

func (c *Client) GetWorkspaceInfo(wsNum int) (*WorkspaceInfo, error) {
	if c.verbose {
		log.Printf("Getting info for workspace %d", wsNum)
	}

	tree, err := c.GetTree()
	if err != nil {
		return nil, err
	}

	workspaces := tree.FindWorkspaces()
	ws, exists := workspaces[wsNum]
	if !exists {
		return nil, fmt.Errorf("workspace %d not found", wsNum)
	}

	info := &WorkspaceInfo{
		Number:         wsNum,
		Layout:         ws.Layout,
		Representation: ws.Representation,
		AppOrder:       ExtractAppOrder(ws.Representation),
	}

	return info, nil
}
