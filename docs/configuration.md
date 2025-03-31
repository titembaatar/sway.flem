# Configuration Reference

sway.flem uses a YAML configuration file to define your workspace layout. This document provides a comprehensive reference for all configuration options.

## Top-Level Structure

```yaml
workspaces:
  "1":  # Workspace name
    layout: "splith"  # Workspace layout
    apps:  # Top-level applications
      - name: "app1"
        size: "30ppt"
    container:  # Optional nested container
      split: "splitv"
      size: "70ppt"
      apps:
        # Container applications
  "2":  # Another workspace
    # ...
```

## Configuration Fields

### Workspaces

The `workspaces` section is the top-level container for all workspace configurations. Each key in this object represents a workspace name.

```yaml
workspaces:
  "1": { ... }    # Workspace 1
  "2": { ... }    # Workspace 2
  "code": { ... } # Workspace named "code"
```

### Workspace Configuration

Each workspace has the following properties:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `layout` | string | Yes | The layout for the workspace |
| `apps` | array | No | List of applications to launch directly in the workspace |
| `container` | object | No | A nested container within the workspace |

### Layout Types

The `layout` and `split` properties support the following values:

| Value | Alternatives | Description |
|-------|--------------|-------------|
| `splith` | `horizontal`, `h` | Horizontal split, windows side by side |
| `splitv` | `vertical`, `v` | Vertical split, windows stacked top to bottom |
| `tabbed` | `tab`, `t` | Tabbed layout, windows are tabs |
| `stacking` | `stack`, `s` | Stacking layout, windows are stacked |

### Apps Configuration

The `apps` array contains objects with the following properties:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | string | Yes | Name of the application (used for identification) |
| `size` | string | Yes | Size of the window (format: "NNppt" or "NNpx") |
| `cmd` | string | No | Custom command to launch the application |

Example:
```yaml
apps:
  - name: "firefox"
    size: "70ppt"
  - name: "terminal"
    cmd: "alacritty"
    size: "30ppt"
```

### Container Configuration

The `container` object has the following properties:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `split` | string | Yes | The layout for the container |
| `size` | string | Yes | Size of the container |
| `apps` | array | No | List of applications to launch in the container |
| `container` | object | No | A nested container within this container |

Containers can be nested to create complex layouts.

## Sizes

Sizes can be specified in:

- Percentage points: `"50ppt"` - relative to the parent container
- Pixels: `"800px"` - absolute size in pixels

## Complete Example

Here's a comprehensive example with multiple workspaces and nested containers:

```yaml
workspaces:
  "1":
    layout: "h"
    apps:
      - name: "firefox"
        size: "70ppt"
      - name: "terminal"
        size: "30ppt"

  "2":
    layout: "splith"
    apps:
      - name: "slack"
        size: "15ppt"
      - name: "spotify"
        size: "15ppt"
    container:
      split: "splitv"
      size: "70ppt"
      apps:
        - name: "code"
          size: "70ppt"
        - name: "terminal"
          cmd: "alacritty"
          size: "30ppt"

  "3":
    layout: "tabbed"
    apps:
      - name: "firefox"
        size: "33ppt"
      - name: "chromium"
        size: "33ppt"
      - name: "brave"
        size: "34ppt"
```
