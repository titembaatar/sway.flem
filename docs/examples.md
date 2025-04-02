# Configuration Examples

This document provides various example configurations for common use cases.

## Basic Examples

### Simple Horizontal Split

```yaml
workspaces:
  "1":
    layout: "h"
    containers:
      - app: "app1"
        size: "70ppt"
      - app: "app2"
        size: "30ppt"
```

### Simple Vertical Split

```yaml
workspaces:
  "1":
    layout: "v"
    containers:
      - app: "app1"
        size: "70ppt"
      - app: "app2"
        size: "30ppt"
```

### Tabbed Layout

```yaml
workspaces:
  "1":
    layout: "tabbed"
    containers:
      - app: "app1"
        size: "33ppt"
      - app: "app2"
        size: "33ppt"
      - app: "app3"
        size: "34ppt"
```

### Stacking Layout

```yaml
workspaces:
  "1":
    layout: "stacking"
    containers:
      - app: "app1"
        size: "33ppt"
      - app: "app2"
        size: "33ppt"
      - app: "app3"
        size: "34ppt"
```

## Advanced Examples

### Development Environment

A development environment with a code editor, terminal, and browser:

```yaml
workspaces:
  "dev":
    layout: "splith"
    containers:
      - app: "editor"
        size: "60ppt"
      - split: "splitv"
        size: "40ppt"
        containers:
          - app: "terminal"
            cmd: "custom-terminal"
            size: "50ppt"
          - app: "browser"
            size: "50ppt"
```

### Communication Workspace

A workspace for chat applications:

```yaml
workspaces:
  "comms":
    layout: "splith"
    containers:
      - app: "chat1"
        size: "33ppt"
      - app: "chat2"
        size: "33ppt"
      - app: "chat3"
        cmd: "custom-chat-app"
        size: "34ppt"
```

### Multi-Monitor Setup

Using names for workspaces to help with organization, and focusing specific workspaces on different monitors:

```yaml
# Focus workspaces on different monitors at the end
focus:
  - "2:web"    # Focus this workspace on the second monitor
  - "1:code"   # Focus this workspace on the first monitor

workspaces:
  "1:code":
    layout: "splith"
    containers:
      - app: "editor"
        size: "70ppt"
      - app: "terminal"
        size: "30ppt"

  "2:web":
    layout: "tabbed"
    containers:
      - app: "browser1"
        size: "100ppt"
      - app: "browser2"
        size: "100ppt"

  "3:comms":
    layout: "splitv"
    containers:
      - app: "chat1"
        size: "50ppt"
      - app: "chat2"
        size: "50ppt"
```

This example assumes you've already configured your Sway to assign workspaces to specific outputs, like:

```
# In your Sway config (~/.config/sway/config)
workspace "1:code" output DP-1
workspace "2:web" output DP-2
workspace "3:comms" output DP-1
```

### Complex Nested Layout

A complex layout with multiple nested containers:

```yaml
workspaces:
  "2":
    layout: "splith"
    containers:
      - app: "app1"
        size: "15ppt"
      - app: "app2"
        size: "15ppt"
      - split: "splitv"
        size: "70ppt"
        containers:
          - app: "app3"
            size: "15ppt"
          - app: "app4"
            size: "15ppt"
          - split: "splith"
            size: "70ppt"
            containers:
              - app: "app5"
                size: "30ppt"
              - app: "app6"
                size: "70ppt"
```

Visual representation of this layout:

```
┌────┬────┬────────────────────────┐
│    │    │          app3          │
│    │    ├────────────────────────┤
│    │    │          app4          │
│    │    ├────────┬───────────────┤
│app1│app2│        │               │
│    │    │        │               │
│    │    │  app5  │     app6      │
│    │    │        │               │
│    │    │        │               │
└────┴────┴────────┴───────────────┘
```

### Media Workspace

A workspace for media consumption:

```yaml
workspaces:
  "media":
    layout: "splitv"
    containers:
      - app: "music-player"
        size: "20ppt"
      - split: "splith"
        size: "80ppt"
        containers:
          - app: "video-player"
            size: "70ppt"
          - app: "media-control"
            cmd: "terminal -e media-controller"
            size: "30ppt"
```

## Custom Application Commands and Post Actions

Using custom commands to launch applications with specific parameters and post-launch actions:

```yaml
workspaces:
  "custom":
    layout: "splith"
    containers:
      - app: "browser"
        cmd: "browser --private-window"
        size: "50ppt"
        post:
          - "browser --new-tab resource1"
          - "browser --new-tab resource2"
      - app: "terminal"
        cmd: "terminal --working-directory ~/projects"
        size: "50ppt"
        post:
          - "terminal -e 'echo Welcome to your workspace'"
```
