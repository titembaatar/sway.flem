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
  rerun-post: true|false        # Optional (default: false)
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

## Application Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `app` | string | Required | Name of the application |
| `cmd` | string | Same as `app` | Custom launch command |
| `size` | string | None | Size specification |
| `delay` | integer | 0 | Delay in seconds before proceeding |
| `rerun-post` | boolean | false | Whether to re-run post commands when focusing existing applications |
| `post` | array | None | List of commands to run after launching |

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
        post:
          - "echo 'Editor loaded'"
        rerun-post: true  # Will run post commands even when focusing an existing instance

      - app: "terminal"
        size: 40
        post:
          - "echo 'Terminal loaded'"
        # rerun-post defaults to false, so post commands won't run when focusing

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
> - Post commands run only on initial launch by default; set `rerun-post: true` to run them every time
