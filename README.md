# 🪟 sway.flem

A Go application to manage Sway workspaces based on YAML configuration files.
This tool allows you to automatically set up, restore, and maintain your Sway workspace layout
when launching or reloading SwayWM.

## ✨ Features

- Define multiple workspaces with different layouts
- Launch applications in specific workspaces
- Configure window properties:
  - Floating state
  - Size
- Run post-launch commands
- Support for both Wayland and X11 applications
- Gracefully handle existing applications (close, resize)

## 🔧 Installation

### Prerequisites

- Go 1.24.1 or later
- Sway window manager
- YAML support

### Building

```bash
# Clone the repository
git clone https://github.com/titembaatar/sway.flem.git
cd sway.flem/

# Build the project
make build

# Install to your system
sudo make install
```

## 🚀 Usage

```bash
swayflem [--config path/to/config.yml] [--verbose] [--version]
```

Options:
- `--config`: Path to the configuration file (default: ./sway.flem.yaml)
- `--verbose`: Enable verbose logging
- `--version`: Show version information and exit
- `--help`: Show help message

### Automatic Sway Integration

To run sway.flem automatically when Sway loads/reloads, add the following to your Sway config file:

```
# Execute sway.flem on reload
exec_always "swayflem --config /path/to/your/config.yml"
```

## 📝 Configuration

The configuration uses YAML format and follows this structure:

```yaml
# Optional: Specify which workspace to focus after setup
focus_workspace: 1

# Default settings (optional)
defaults:
  default_layout: "splith"  # Default layout for workspaces
  default_output: "DP-1"    # Default output for workspaces
  default_floating: false   # Whether apps should be floating by default

# Workspace configurations
workspaces:
  1:  # Workspace number
    layout: "tabbed"  # Optional: splith, splitv, stacking, tabbed
    output: "DP-1"    # Optional: Output name
    close_unmatched: false  # Optional: Close windows not in config
    apps:
      - name: "app_name_1"  # app_id for Wayland or class for X11
        command: "cmd_to_launch_app"  # Optional if same as name
        size: "20ppt 100ppt"  # width height (optional)
        floating: true  # Set window to floating mode (optional)
        posts:  # Commands to run after app launches
          - "cmd 1"
          - "cmd 2"
        delay: 1  # Delay in seconds before configuring (optional)
```

### 🧩 Configuration Properties

#### Workspace Properties

- `layout`: The layout for the workspace (`splith`, `splitv`, `stacking`, `tabbed`)
- `output`: The output to display the workspace on
- `close_unmatched`: Whether to close windows in the workspace that aren't defined in the config

#### App Properties

- `name`: The app name (app_id for Wayland, class for X11 apps)
- `command`: The command to launch the app (defaults to `name` if omitted)
- `size`: The size of the app window
  - Can be specified in pixels (`800 600`) or percentage points (`50ppt 100ppt`)
  - In `splitv` layouts, only height is applied
  - In `splith` layouts, only width is applied
  - In `tabbed`/`stacking` layouts, resize has no effect
- `floating`: Whether the window should be floating
- `posts`: Commands to run after launching the app
- `delay`: Delay in seconds before configuring the app (useful for slow-starting apps)

### 📋 Example Configurations

Simple example:

```yaml
focus_workspace: 1
workspaces:
  1:
    layout: "splith"
    apps:
      - name: "kitty"
        size: "50ppt"
      - name: "firefox"
        size: "50ppt"
```

Check the `examples/` directory for more comprehensive configuration examples.

## 🔍 Troubleshooting

- If windows don't appear or have incorrect properties, use the `--verbose` flag to see more details
- Use `swaymsg -t get_tree | grep -E 'app_id|class'` to find the correct app_id or class for your applications
- Check the Sway log for errors: `journalctl -b | grep sway`

## 📄 License

MIT
