<p align="center">
  <a href="https://github.com/titembaatar/sway.flem">
    <img src="https://img.shields.io/badge/version-v0.1.0-b93d4d?style=flat-square&labelColor=446f5e">
  </a>
  <img src="https://img.shields.io/badge/go-1.24.1-12adad?style=flat-square&logoColor=dceae4&labelColor=446f5e">
  <a href="https://github.com/titembaatar/sway.flem/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-a353c6?style=flat-square&logoColor=dceae4&labelColor=446f5e">
  </a>
</p>

# sway.flem

**Workspace Automation for Sway Window Manager**

Effortlessly set up and organize your working environment with a single command.

## ✨ Features

- Rapid workspace configuration
- Flexible nested container layouts
- Control over window sizes and positions

## 🚀 Quick Start

```bash
# Install via Go
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

## 📦 Installation Options

### Go Install
```bash
go install github.com/titembaatar/sway.flem/cmd/flem@latest
```

### From Source
```bash
# Clone the repository
git clone https://github.com/titembaatar/sway.flem.git
cd sway.flem

# Build the project
make build

# (Optional) Install to $GOPATH/bin
go install ./cmd/flem
```

## 🛠️ Usage

```
flem sway -config <config-file>
```

### Options

- `-config`: Path to configuration file (required)
- `-version`: Show version information
- `-verbose`: Enable verbose logging
- `-debug`: Enable debug mode
- `-dry-run`: Validate configuration without changes

## 📝 Configuration Example

```yaml
workspaces:
  1:
    layout: h
    containers:
      - app: firefox
        size: 50
        post:
          - "firefox --new-tab https://example.com"
        rerun-post: false  # Only run post commands on initial launch (default behavior)
      - app: code
        size: 50  # technically optional because of the way Sway works
        post:
          - "code --goto README.md:10"
        rerun-post: true  # Run post commands every time this app is focused
```

## 🤝 Contributing

Contributions are welcome! Feel free to submit a Pull Request.

## 📄 License

MIT License
