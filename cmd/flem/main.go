package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/titembaatar/sway.flem/internal/app"
	"github.com/titembaatar/sway.flem/internal/config"
	"github.com/titembaatar/sway.flem/internal/log"
)

const (
	version = "0.1.0"
)

var (
	configFile  string
	showVersion bool
	verbose     bool
	debug       bool
	dryRun      bool
)

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode with extra logging")
	flag.BoolVar(&dryRun, "dry-run", false, "Validate configuration without making changes")
	flag.Parse()
}

func main() {
	// Check if we have a command ("sway")
	args := flag.Args()
	if len(args) == 0 {
		log.Error("Command required (e.g., 'flem sway')")
		fmt.Println("Usage: flem [options] <command>")
		fmt.Println("Commands:")
		fmt.Println("  sway    Configure Sway workspaces")
		flag.Usage()
		os.Exit(1)
	}

	command := args[0]
	if command != "sway" {
		log.Error("Unknown command: %s", command)
		fmt.Println("Usage: flem [options] <command>")
		fmt.Println("Commands:")
		fmt.Println("  sway    Configure Sway workspaces")
		flag.Usage()
		os.Exit(1)
	}

	if debug {
		log.SetLevel(log.LogLevelDebug)
	} else if verbose {
		log.SetLevel(log.LogLevelInfo)
	} else {
		log.SetLevel(log.LogLevelWarn)
	}

	if showVersion {
		fmt.Printf("flem v%s\n", version)
		os.Exit(0)
	}

	if configFile == "" {
		log.Error("Config file must be specified")
		flag.Usage()
		os.Exit(1)
	}

	log.Info("Starting flem sway v%s", version)

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatal("Failed to load configuration: %v", err)
	}

	for name := range cfg.Workspaces {
		log.Debug("Found workspace configuration: %s", name)
	}

	if dryRun {
		log.Info("Dry run completed successfully. Configuration is valid.")
		os.Exit(0)
	}

	if err := app.Setup(cfg); err != nil {
		log.Fatal("Failed to setup environment: %v", err)
	}

	log.Info("Sway environment has been successfully configured")
}
