package sway

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

func NewClient(verbose bool) *Client {
	return &Client{
		Verbose: verbose,
	}
}

func (c *Client) ExecuteCommand(cmd string) error {
	if c.Verbose {
		log.Printf("Executing: swaymsg %s\n", cmd)
	}

	out, err := exec.Command("swaymsg", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %v (output: %s)", ErrCommandFailed, err, out)
	}

	return nil
}

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

func (c *Client) GetWorkspaceInfo(wsNum int) (*WorkspaceInfo, error) {
	if c.Verbose {
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
		Output:         ws.Output,
		Representation: ws.Representation,
		AppOrder:       ExtractAppOrder(ws.Representation),
	}

	return info, nil
}

func (c *Client) GetOutputs() ([]string, error) {
	cmd := exec.Command("swaymsg", "-t", "get_outputs", "--raw")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("getting outputs: %w", err)
	}

	var outputs []struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(output, &outputs); err != nil {
		return nil, fmt.Errorf("parsing outputs: %w", err)
	}

	var outputNames []string
	for _, output := range outputs {
		// Skip special outputs like __i3
		if output.Name != "__i3" {
			outputNames = append(outputNames, output.Name)
		}
	}

	return outputNames, nil
}
