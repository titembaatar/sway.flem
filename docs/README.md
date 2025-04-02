# sway.flem Documentation

Welcome to the sway.flem documentation! This guide provides detailed information about installing, configuring, and using sway.flem to create complex workspace environments in the Sway window manager.

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Configuration Reference](configuration.md)
4. [Command Line Options](cli.md)
5. [Examples](examples.md)
6. [Troubleshooting](troubleshooting.md)

## Introduction

sway.flem is a tool that automates the creation of complex workspace environments in Sway. It allows you to define your workspaces, containers, and application layouts in a YAML configuration file, and then automatically sets up this environment with a single command.

The tool is designed to be:
- **Simple**: Easy to understand configuration format
- **Flexible**: Support for complex nested layouts
- **Reliable**: Robust error handling and recovery
- **Fast**: Quick setup of even the most complex environments

## Installation

### Prerequisites

- Go 1.16 or higher
- Sway window manager
- swaymsg utility

### Install from Source

```bash
# Clone repository
git clone https://github.com/titembaatar/sway.flem.git
cd sway.flem

# Build
make build

# The binary will be created in the bin directory
./bin/flem -config example.yaml
```

### Install using Go

```bash
# Install directly using go install
go install github.com/titembaatar/sway.flem/cmd/flem@latest

# Run
flem sway -config <path-to-config.yaml>
```

## Next Steps

- Check the [Configuration Reference](configuration.md) for detailed information about the configuration format
- Browse the [Examples](examples.md) for inspiration on how to set up your own workspace
- See the [Command Line Options](cli.md) for all available options
