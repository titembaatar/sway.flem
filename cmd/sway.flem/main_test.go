package main

import (
	"flag"
	"os"
	"testing"

	"github.com/titembaatar/sway.flem/internal/log"
)

func TestMain(m *testing.M) {
	// Save original args
	oldArgs := os.Args

	// Set log level to none for all tests
	log.SetLevel(log.LogLevelNone)

	// Run tests
	code := m.Run()

	// Restore original args
	os.Args = oldArgs

	os.Exit(code)
}

func TestFlagsDefault(t *testing.T) {
	// Create a new flag set for testing
	testFlags := flag.NewFlagSet("test", flag.ContinueOnError)

	// Create local variables for flags
	var (
		testConfigFile  string
		testShowVersion bool
		testVerbose     bool
		testDebug       bool
		testDryRun      bool
	)

	// Define flags just like in init()
	testFlags.StringVar(&testConfigFile, "config", "", "Path to configuration file")
	testFlags.BoolVar(&testShowVersion, "version", false, "Show version information")
	testFlags.BoolVar(&testVerbose, "verbose", false, "Enable verbose logging")
	testFlags.BoolVar(&testDebug, "debug", false, "Enable debug mode with extra logging")
	testFlags.BoolVar(&testDryRun, "dry-run", false, "Validate configuration without making changes")

	// Set up minimal args
	args := []string{}

	// Parse flags
	err := testFlags.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Check that flags have their default values
	if testConfigFile != "" {
		t.Errorf("Expected empty configFile, got %s", testConfigFile)
	}
	if testShowVersion != false {
		t.Errorf("Expected showVersion to be false, got %t", testShowVersion)
	}
	if testVerbose != false {
		t.Errorf("Expected verbose to be false, got %t", testVerbose)
	}
	if testDebug != false {
		t.Errorf("Expected debug to be false, got %t", testDebug)
	}
	if testDryRun != false {
		t.Errorf("Expected dryRun to be false, got %t", testDryRun)
	}
}

func TestFlagsWithValues(t *testing.T) {
	// Create a new flag set for testing
	testFlags := flag.NewFlagSet("test", flag.ContinueOnError)

	// Create local variables for flags
	var (
		testConfigFile  string
		testShowVersion bool
		testVerbose     bool
		testDebug       bool
		testDryRun      bool
	)

	// Define flags just like in init()
	testFlags.StringVar(&testConfigFile, "config", "", "Path to configuration file")
	testFlags.BoolVar(&testShowVersion, "version", false, "Show version information")
	testFlags.BoolVar(&testVerbose, "verbose", false, "Enable verbose logging")
	testFlags.BoolVar(&testDebug, "debug", false, "Enable debug mode with extra logging")
	testFlags.BoolVar(&testDryRun, "dry-run", false, "Validate configuration without making changes")

	// Set up args with all flags
	args := []string{
		"-config", "test_config.yaml",
		"-version",
		"-verbose",
		"-debug",
		"-dry-run",
	}

	// Parse flags
	err := testFlags.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Check that flags have the expected values
	if testConfigFile != "test_config.yaml" {
		t.Errorf("Expected configFile to be 'test_config.yaml', got %s", testConfigFile)
	}
	if testShowVersion != true {
		t.Errorf("Expected showVersion to be true, got %t", testShowVersion)
	}
	if testVerbose != true {
		t.Errorf("Expected verbose to be true, got %t", testVerbose)
	}
	if testDebug != true {
		t.Errorf("Expected debug to be true, got %t", testDebug)
	}
	if testDryRun != true {
		t.Errorf("Expected dryRun to be true, got %t", testDryRun)
	}
}

