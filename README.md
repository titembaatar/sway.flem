# sway.flem

![Version](https://img.shields.io/badge/version-0.1.0-blue)
![Go](https://img.shields.io/badge/go-%3E%3D1.16-blue)
![License](https://img.shields.io/badge/license-MIT-green)

**Auto-launcher for complex Sway workspace environments**

sway.flem helps you instantly set up your perfect working environment in Sway by automatically launching applications and organizing them into precise layouts based on a YAML configuration file.

## Features

- Instantly launch multiple applications in the right places
- Configure complex nested container layouts
- Precise control over window sizes and positions
- Support for all Sway layout types (horizontal, vertical, tabbed, stacking)
- Easy to understand configuration format

## Quick Start

```bash
# Install
go install github.com/titembaatar/sway.flem/cmd/sway.flem@latest

# Create a basic configuration
cat > ~/.config/sway.flem.yaml << EOF
workspaces:
  "1":
    layout: "h"
    apps:
      - name: "firefox"
        size: "60ppt"
      - name: "terminal"
        size: "40ppt"
EOF

# Launch your workspace setup
sway.flem -config ~/.config/sway.flem.yaml
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
go install ./cmd/sway.flem
```

## Basic Usage

```
sway.flem -config <config-file>
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

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details
