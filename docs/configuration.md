# Configuration Reference

sway.flem uses a YAML configuration file to define your workspace layout.
This document provides a comprehensive reference for all configuration options.

## Top-Level Structure

```yaml
focus:   # Optional list of workspaces to focus at the end
  - "6"  # First workspace to focus (e.g., on second monitor)
  - "1"  # Second workspace to focus (e.g., on first monitor)

workspaces:
  1:                   # Workspace name
    layout: splith     # Workspace layout
    containers:        # List of containers in this workspace
      - app: app1
        size: 30
      - app: app2
        size: 70
  2:                   # Another workspace
    layout: splitv
    containers:
      - app: app3
        split: splitv  # Only needed for nested containers
        size: 70
        containers:    # Nested containers
          - app: app4
            size: 50
          - app: app5
            size: 50
```

## Configuration Fields

### Focus

The `focus` field allows you to specify which workspaces should be focused at the end of
the setup process. This is particularly useful for multi-monitor setups where you want
specific workspaces to be visible on each monitor when setup completes.

```yaml
focus:
  - 6  # Focus workspace 6 first (e.g., on second monitor)
  - 1  # Then focus workspace 1 (e.g., on first monitor)
```

The workspaces will be focused in the order specified.
The last workspace in the list will have keyboard focus.

### Workspaces

The `workspaces` section is the top-level container for all workspace configurations.
Each key in this object represents a workspace name.

```yaml
workspaces:
  1: { ... }      # Workspace 1
  2: { ... }      # Workspace 2
  "code": { ... } # Workspace named "code"
```

### Workspace Configuration

Each workspace has the following properties:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `layout` | string | Yes | The layout for the workspace |
| `containers` | array | Yes | List of containers in the workspace |

### Layout Types

The `layout` and `split` properties support the following values:

| Value | Alternatives | Description |
|-------|--------------|-------------|
| `splith` | `horizontal`, `h` | Horizontal split, windows side by side |
| `splitv` | `vertical`, `v` | Vertical split, windows stacked top to bottom |
| `tabbed` | `tab`, `t` | Tabbed layout, windows are tabs |
| `stacking` | `stack`, `s` | Stacking layout, windows are stacked |

### Container Configuration

Each container can have the following properties:

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `app` | string | Yes (for app containers) | Name or identifier of the application |
| `cmd` | string | No | Custom command to launch the application |
| `size` | string | No | Size of the container (format: "NN" (ppt) or "NNppt" or "NNpx") |
| `delay` | integer | No | Delay in seconds to wait after launching the app |
| `post` | array | No | List of commands to execute after the app launches |
| `split` | string | Yes (for nested containers) | The layout for nested containers |
| `containers` | array | Yes (for nested containers) | List of child containers |

A container must either be an app container (with the `app` property) or
a nested container (with the `containers` property), but not both.

Example app container:
```yaml
- app: "app1"
  cmd: "custom-app1-command"
  size: "70ppt"
  post:
    - "command1"
    - "command2"
```

Example nested container:
```yaml
- split: "splitv"
  size: "70ppt"
  containers:
    - app: "app2"
      size: "50ppt"
    - app: "app3"
      size: "50ppt"
```

## Sizes

Sizes can be specified in:

- Percentage points: `"50"` `"50ppt"` - relative to the parent container
- Pixels: `"800px"` - absolute size in pixels

Size is optional.
If not specified, the application or container will use the default size allocated by Sway.

## Post-Launch Commands

The `post` property allows you to execute additional commands after an application has been launched.

This is useful for:
- Opening multiple tabs or windows
- Sending keys or commands to the launched application
- Triggering follow-up actions

Example:
```yaml
- app: "app1"
  size: "60ppt"
  post:
    - "app1 --action1"
    - "app1 --action2"
```

## Complete Example

Here's a comprehensive example with multiple workspaces, nested containers,
post-launch commands, and focusing specific workspaces:

```yaml
# Focus on these workspaces at the end (useful for multi-monitor setups)
focus:
  - "6"  # Focus on workspace 6 first (e.g., on second monitor)
  - "1"  # Then focus on workspace 1 (e.g., on first monitor)

workspaces:
  1:
    layout: "h"
    containers:
      - app: "app1"
        size: "70ppt"
        post:
          - "app1 --action1"
          - "app1 --action2"
      - app: "app2"
        size: "30ppt"
        post:
          - "app2-helper-command"

  2:
    layout: "splith"
    containers:
      - app: "app3"
        size: "15ppt"
      - app: "app4"
        size: "15ppt"
      - split: "splitv"
        size: "70ppt"
        containers:
          - app: "app5"
            size: "70ppt"
          - app: "app6"
            cmd: "custom-app6-command"
            size: "30ppt"

  6:
    layout: "splitv"
    containers:
      - app: "app7"
        cmd: "app7 --option"
        post:
          - "app7 --action1"
          - "app7 --action2"
      - app: "app8"
        cmd: "custom-app8-command"
```
