# sway.flem Configuration Guide

This document provides detailed information about configuring sway.flem to manage your Sway workspaces.

## Configuration File Format

sway.flem uses YAML for its configuration file format.
The default location is `./sway.flem.yaml` in the current directory,
but you can specify a different location using the `--config` command-line option.

## Complete Configuration Structure

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
        position: "center"  # (optional)
        floating: true  # Set window to floating mode (optional)
        posts:  # Commands to run after app launches
          - "cmd 1"
          - "cmd 2"
        delay: 1  # Delay in seconds before configuring (optional)
```

## Finding App Names

The `name` field in app configurations should match the app_id (for Wayland applications) or class (for X11 applications). You can find these values using:

```bash
swaymsg -t get_tree | grep -E 'app_id|class'
```

## Finding Launch Commands

If you are using a launcher like tofi you can run `tofi-drun` in the terminal,
launch the application, and you should get the command use to launch the application
in the output of the command.

## Layout Types

sway.flem supports the following layout types:

- `splith` - Horizontal split (windows arranged side by side)
- `splitv` - Vertical split (windows stacked vertically)
- `tabbed` - Tabbed layout (only one window visible at a time, with tabs to switch)
- `stacking` - Stacked layout (similar to tabbed, but with window titles stacked)

## Position Values

The following position values are supported:

- `center` or `middle` - Center the window
- `top` - Top of the screen
- `bottom` - Bottom of the screen
- `left` - Left side of the screen
- `right` - Right side of the screen
- `pointer`, `cursor`, or `mouse` - Position at the cursor
- Custom coordinates (e.g., `100 200` for x=100, y=200)

## Size Specification

Sizes can be specified in:

- Pixels: `800 600` (width 800px, height 600px)
- Percentage points: `50ppt 100ppt` (50% width, 100% height)

## Post-Launch Commands

Post-launch commands (`posts`) are executed after the application is launched and configured. These can be useful for:

- Setting up the application environment
- Running commands within the application
- Triggering additional actions

## Workspace Output Assignment

Assigning workspaces to specific outputs allows sway.flem to organize your workspaces across multiple monitors.

Use the `output` field in a workspace configuration to specify which monitor should display the workspace.

> [!NOTE]
>
> This should be done in your SwayWM config. But if for whatever reason you need it, it's here.
