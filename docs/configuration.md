# Configuration Reference

## Overview

The configuration file for sway.flem is a YAML document that defines workspace layouts,
applications, and their properties.

## Configuration Structure

### Top-Level Fields

```yaml
focus:   # Optional: Workspaces to focus at the end
  - 6  # First workspace to focus
  - 1  # Second workspace to focus

workspaces:
  <workspace-name>:
    layout: <layout-type>
    containers:
      - app: <application-name>
        size: <size-specification>
```

## Workspace Configuration

### `focus`
- **Optional**: Yes
- **Type**: Array of strings
- **Description**: Specifies which workspaces to focus after setup

### Workspace Properties

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `layout` | string | Yes | Defines the workspace layout |
| `containers` | array | Yes | List of applications or nested containers |

## Layout Types

| Layout | Aliases | Description |
|--------|---------|-------------|
| `splith` | `h`, `horizontal` | Horizontal split  |
| `splitv` | `v`, `vertical` | Vertical split  |
| `tabbed` | `t`, `tab` | Tabbed layout |
| `stacking` | `s`, `stack` | Stacking layout |

## Container Configuration

### Application Container

```yaml
- app: <application-name>
  cmd: <custom-launch-command>  # Optional
  size: <size-specification>    # Optional
  delay: <launch-delay>         # Optional
  post:                         # Optional
    - <post-launch-command>
```

### Nested Container

```yaml
- split: <layout-type>
  size: <size-specification>    # Optional
  containers:
    - app: <application-name>
      size: <size-specification>
```

## Size Specification

Sizes can be specified in two formats:
- Percentage points: `50`, `50ppt`
- Pixels: `800px`

## Full Example

```yaml
focus:
  - 2
  - 1

workspaces:
  1:
    layout: h
    containers:
      - app: "editor"
        size: 60
      - app: "terminal"
        size: 40

  2:
    layout: v
    containers:
      - app: "browser"
        size: 70
      - split: h
        size: 30
        containers:
          - app: "chat1"
            size: 50
          - app: "chat2"
            size: 50
```

> [!NOTE]
>
> - Workspace names can be numbers or strings
> - Nested containers more complex layout configurations
> - Size specifications are optional
