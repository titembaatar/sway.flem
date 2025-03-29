// Package main provides the entry point for the sway.flem application
package main

import (
	"flag"
	"log"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/manager"
	"github.com/titembaatar/sway.flem/internal/sway"
)

func main() {
	// Set up logging
	log.SetPrefix("[sway.flem] ")

	// Parse command line flags
	configPath := flag.String("config", "./sway.flem.yaml", "Path to configuration file")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create Sway client
	client := sway.NewClient(*verbose)

	// Create and run the workspace manager
	mgr := manager.NewManager(client, cfg)
	if err := mgr.Run(); err != nil {
		log.Fatalf("Error running manager: %v", err)
	}

	log.Println("Workspace setup complete!")
}
