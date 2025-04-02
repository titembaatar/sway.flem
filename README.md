# sway.flem

![Version](https://img.shields.io/badge/version-0.1.0-blue)
![Go](https://img.shields.io/badge/go-%3E%3D1.16-blue)
![License](https://img.shields.io/badge/license-MIT-green)

**Auto-launcher for complex Sway workspace environments**

sway.flem helps you instantly set up your perfect working environment in Sway by automatically launching applications and organizing them into precise layouts based on a YAML configuration file.

## Features

- Launch multiple applications in the right places
- Configure nested container layouts
- Control over window sizes and positions (within Sway limitation)
- Support for all Sway layout types (horizontal, vertical, tabbed, stacking)
- Easy to understand configuration format

## Quick Start

```bash
# Install
go install github.com/titembaatar/sway.flem/cmd/flem@latest

# Create a basic configuration
cat > ~/.config/sway/config.yml << EOF
workspaces:
  "1":
    layout: "h"
    containers:
      - app: "app1"
        size: "60ppt"
      - app: "app2"
        size: "40ppt"
EOF

# Launch your workspace setup
flem sway -config ~/.config/sway/config.yml
```

## Installation

### From Source

```bash
# Clone repository
git clone https://github.com/titembaatar/sway.flem.git
cd sway.flem

# Build
make build

# (Optional) Install to $GOPATH/bin
go install ./cmd/flem
```

## Basic Usage

```
flem sway -config <config-file>
```

### Options

```
  -config string   Path to configuration file
  -version         Show version information
  -verbose         Enable verbose logging
  -debug           Enable debug mode with extra logging
  -dry-run         Validate configuration without making changes
```

## Simple Configuration Example

```yaml
workspaces:
  "2":
    layout: "h" # horizontal layout
    apps:
      - name: "firefox"
        size: "50ppt"
      - name: "code"
        size: "50ppt"
```

## Advanced Features

- Support for multiple workspaces
- Nested container layouts
- Custom commands for launching applications
- Detailed logging options
- Automatic error recovery

See the [documentation](docs/README.md) for more advanced configuration examples.

## Contributing

Contributions are welcome! Feel free to submit a Pull Request.
