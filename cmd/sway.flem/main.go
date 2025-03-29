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
	log.SetPrefix("[sway.flem] ")

	configPath := flag.String("config", "./sway.flem.yaml", "Path to configuration file")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	client := sway.NewClient(*verbose)

	mgr := manager.NewManager(client, cfg)
	if err := mgr.Run(); err != nil {
		log.Fatalf("Error running manager: %v", err)
	}

	log.Println("Workspace setup complete!")
}
