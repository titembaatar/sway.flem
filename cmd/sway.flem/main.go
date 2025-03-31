package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	configFile string
	version    bool
)

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&version, "version", false, "Show version information")
	flag.Parse()
}

func main() {
	if version {
		fmt.Println("sway.flem v0.1.0")
		os.Exit(0)
	}

	if configFile == "" {
		fmt.Println("Error: config file must be specified")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Starting sway.flem with config file: %s\n", configFile)
	// TODO: Implement configuration parsing and application launching
	fmt.Println("Not yet implemented")
}
