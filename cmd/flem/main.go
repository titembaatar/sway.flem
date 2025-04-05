package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/titembaatar/sway.flem/internal/app"
	"github.com/titembaatar/sway.flem/internal/config"
	errs "github.com/titembaatar/sway.flem/internal/errors"
	"github.com/titembaatar/sway.flem/internal/log"
)

const (
	version = "0.1.0"
)

type Flags struct {
	ConfigFile  string
	ShowVersion bool
	Verbose     bool
	Debug       bool
	DryRun      bool
	JSONLogs    bool
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmdIndex := 1
	command := os.Args[cmdIndex]

	switch command {
	case "sway":
		runSwayCommand(os.Args[cmdIndex+1:])
	case "-h", "--help":
		printUsage()
		os.Exit(0)
	case "-v", "--version":
		fmt.Printf("flem v%s\n", version)
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

// Handles the 'sway' subcommand
func runSwayCommand(args []string) {
	flags := parseFlags(args)
	configureLogging(flags)

	// Create error handler
	errorHandler := errs.NewErrorHandler(true, flags.Verbose, flags.Debug)

	if flags.ShowVersion {
		fmt.Printf("flem v%s\n", version)
		os.Exit(0)
	}

	if flags.ConfigFile == "" {
		errorHandler.Handle(errs.NewFatal(errs.ErrConfigNotFound, "Config file must be specified with -config flag"))
		fmt.Println("Run 'flem sway -h' for usage information")
		os.Exit(1)
	}

	log.Info("Starting flem sway v%s", version)

	// Load configuration
	cfg, err := config.LoadConfig(flags.ConfigFile)
	if err != nil {
		errorHandler.Handle(errs.Wrap(err, "Failed to load configuration"))
		os.Exit(1)
	}

	for name := range cfg.Workspaces {
		log.Debug("Found workspace configuration: %s", name)
	}

	if flags.DryRun {
		log.Info("Dry run completed successfully. Configuration is valid.")
		os.Exit(0)
	}

	// Setup Sway environment
	if err := app.Setup(cfg); err != nil {
		errorHandler.Handle(errs.Wrap(err, "Failed to setup environment"))
		os.Exit(1)
	}

	log.Info("Sway environment has been successfully configured")

	// Summarize any non-fatal errors or warnings
	errorHandler.SummarizeErrors()
}

// Parses command line flags and returns the parsed values
func parseFlags(args []string) *Flags {
	flags := &Flags{}

	flagSet := flag.NewFlagSet("sway", flag.ExitOnError)

	flagSet.StringVar(&flags.ConfigFile, "config", "", "Path to configuration file")
	flagSet.BoolVar(&flags.ShowVersion, "version", false, "Show version information")
	flagSet.BoolVar(&flags.Verbose, "verbose", false, "Enable verbose logging")
	flagSet.BoolVar(&flags.Debug, "debug", false, "Enable debug mode with extra logging")
	flagSet.BoolVar(&flags.DryRun, "dry-run", false, "Validate configuration without making changes")
	flagSet.BoolVar(&flags.JSONLogs, "json-logs", false, "Output logs in JSON format (useful for log processing)")

	flagSet.Parse(args)

	return flags
}

// Sets up the logging level based on flags
func configureLogging(flags *Flags) {
	if flags.Debug {
		log.SetLevel(log.LogLevelDebug)
	} else if flags.Verbose {
		log.SetLevel(log.LogLevelInfo)
	} else {
		log.SetLevel(log.LogLevelWarn)
	}

	// Check if JSON logging should be enabled via flag
	if flags.JSONLogs {
		log.EnableJSONLogging()
	}
}

func printUsage() {
	fmt.Println("Usage: flem [options] <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  sway                  Configure Sway workspaces")
	fmt.Println("\nGlobal Options:")
	fmt.Println("  -h, --help            Show this help message")
	fmt.Println("  -v, --version         Show version information")
	fmt.Println("\nSway Command Options:")
	fmt.Println("  -config <file>        Path to configuration file (required)")
	fmt.Println("  -version              Show version information")
	fmt.Println("  -verbose              Enable verbose logging")
	fmt.Println("  -debug                Enable debug mode with extra logging")
	fmt.Println("  -dry-run              Validate configuration without making changes")
	fmt.Println("  -json-logs            Output logs in JSON format (useful for log processing)")
	fmt.Println("\nExamples:")
	fmt.Println("  flem sway -config ~/.config/sway/config.yml")
	fmt.Println("  flem sway -config ~/.config/sway/config.yml -verbose")
	fmt.Println("  flem sway -config ~/.config/sway/config.yml -dry-run")
}
