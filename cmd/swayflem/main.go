package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/manager"
	"github.com/titembaatar/sway.flem/internal/sway"
)

var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {
	configPath := flag.String("config", "./sway.flem.yaml", "Path to configuration file")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	showVersion := flag.Bool("version", false, "Show version information and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("swayflem %s (built %s)\n", version, buildDate)
		return
	}

	log.SetPrefix("[swayflem] ")
	if *verbose {
		log.Println("Verbose logging enabled")
	}

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	swayClient := sway.NewClient(*verbose)

	mgr := manager.NewManager(cfg, swayClient, *verbose)

	if err := mgr.Validate(); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	setupSignalHandling()

	if err := mgr.Run(); err != nil {
		log.Fatalf("Error running manager: %v", err)
	}

	log.Println("Sway workspace setup completed successfully")
}

func loadConfig(path string) (*config.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", path)
	}

	log.Printf("Loading configuration from %s", path)
	return config.LoadFromFile(path)
}

func setupSignalHandling() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received interrupt signal, exiting gracefully...")
		os.Exit(0)
	}()
}
