# Configuration Examples

This document provides various example configurations for common use cases.

## Basic Examples

### Simple Horizontal Split

```yaml
workspaces:
  "1":
    layout: "h"
    apps:
      - name: "firefox"
        size: "70ppt"
      - name: "terminal"
        size: "30ppt"
```

### Simple Vertical Split

```yaml
workspaces:
  "1":
    layout: "v"
    apps:
      - name: "firefox"
        size: "70ppt"
      - name: "terminal"
        size: "30ppt"
```

### Tabbed Layout

```yaml
workspaces:
  "1":
    layout: "tabbed"
    apps:
      - name: "firefox"
        size: "33ppt"
      - name: "chromium"
        size: "33ppt"
      - name: "brave"
        size: "34ppt"
```

### Stacking Layout

```yaml
workspaces:
  "1":
    layout: "stacking"
    apps:
      - name: "firefox"
        size: "33ppt"
      - name: "terminal"
        size: "33ppt"
      - name: "code"
        size: "34ppt"
```

## Advanced Examples

### Development Environment

A development environment with a code editor, terminal, and browser:

```yaml
workspaces:
  "dev":
    layout: "splith"
    apps:
      - name: "code"
        size: "60ppt"
    container:
      split: "splitv"
      size: "40ppt"
      apps:
        - name: "terminal"
          cmd: "alacritty"
          size: "50ppt"
        - name: "firefox"
          size: "50ppt"
```

### Communication Workspace

A workspace for chat applications:

```yaml
workspaces:
  "comms":
    layout: "splith"
    apps:
      - name: "slack"
        size: "33ppt"
      - name: "discord"
        size: "33ppt"
      - name: "telegram"
        cmd: "telegram-desktop"
        size: "34ppt"
```

### Multi-Monitor Setup

Using names for workspaces to help with organization:

```yaml
workspaces:
  "1:code":
    layout: "splith"
    apps:
      - name: "code"
        size: "70ppt"
      - name: "terminal"
        size: "30ppt"

  "2:web":
    layout: "tabbed"
    apps:
      - name: "firefox"
        size: "100ppt"
      - name: "chromium"
        size: "100ppt"

  "3:comms":
    layout: "splitv"
    apps:
      - name: "slack"
        size: "50ppt"
      - name: "discord"
        size: "50ppt"
```

### Complex Nested Layout

A complex layout with multiple nested containers:

```yaml
workspaces:
  "2":
    layout: "splith"
    apps:
      - name: "app1"
        size: "15ppt"
      - name: "app2"
        size: "15ppt"
    container:
      split: "splitv"
      size: "70ppt"
      apps:
        - name: "app3"
          size: "15ppt"
        - name: "app4"
          size: "15ppt"
      container:
        split: "splith"
        size: "70ppt"
        apps:
          - name: "app5"
            size: "30ppt"
          - name: "app6"
            size: "70ppt"
```

Visual representation of this layout:

```
+-------------------------------+
|        |                      |
| app1   |         app3         |
|        |                      |
+--------+----------------------+
|        |                      |
| app2   |         app4         |
|        |                      |
+--------+----------+-----------+
|        |   app5   |   app6    |
|        |          |           |
+--------+----------+-----------+
```

### Media Workspace

A workspace for media consumption:

```yaml
workspaces:
  "media":
    layout: "splitv"
    apps:
      - name: "spotify"
        size: "20ppt"
    container:
      split: "splith"
      size: "80ppt"
      apps:
        - name: "vlc"
          size: "70ppt"
        - name: "terminal"
          cmd: "alacritty -e ncmpcpp"
          size: "30ppt"
```

## Custom Application Commands

Using custom commands to launch applications with specific parameters:

```yaml
workspaces:
  "custom":
    layout: "splith"
    apps:
      - name: "firefox"
        cmd: "firefox --private-window"
        size: "50ppt"
      - name: "terminal"
        cmd: "alacritty --working-directory ~/projects"
        size: "50ppt"
```
