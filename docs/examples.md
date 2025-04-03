# Configuration Examples

## Basic Layouts

### Horizontal Split
```yaml
workspaces:
  1:
    layout: h
    containers:
      - app: "app1"
        size: 70
      - app: "app2"
        size: 30
```

### Vertical Split
```yaml
workspaces:
  1:
    layout: v
    containers:
      - app: "top-app"
        size: 70
      - app: "bottom-app"
        size: 30
```

## Development Environment

### Code Workspace with Multiple Containers
```yaml
workspaces:
  "dev":
    layout: h
    containers:
      - app: "code-editor"
        size: 60
      - split: v
        size: 40
        containers:
          - app: "terminal"
            size: 50
          - app: "browser"
            size: 50
```

## Multimedia Workspace

### Media and Control Layout
```yaml
workspaces:
  "media":
    layout: v
    containers:
      - app: "music-player"
        size: 20
      - split: h
        size: 80
        containers:
          - app: "video-player"
            size: 70
          - app: "media-control"
            size: 30
```

## Multi-Monitor Setup

### Workspace Focus and Complex Layout
```yaml
focus:
  - "2:web"    # Focus on the web workspace first
  - "1:code"   # Then focus on the code workspace

workspaces:
  "1:code":
    layout: h
    containers:
      - app: "editor"
        size: 70
      - app: "terminal"
        size: 30

  "2:web":
    layout: "tabbed"
    containers:
      - app: "browser1"
      - app: "browser2"

  "3:comms":
    layout: v
    containers:
      - app: "chat1"
        size: 50
      - app: "chat2"
        size: 50
```

## Advanced Nested Containers

### Complex Workspace Layout
```yaml
workspaces:
  2:
    layout: h
    containers:
      - app: "app1"
        size: 15
      - app: "app2"
        size: 15
      - split: v
        size: 70
        containers:
          - app: "app3"
            size: 15
          - app: "app4"
            size: 15
          - split: h
            size: 70
            containers:
              - app: "app5"
                size: 30
              - app: "app6"
                size: 70
```

## Custom Application Launch

### With Custom Commands and Post-Launch Actions
```yaml
workspaces:
  "custom":
    layout: h
    containers:
      - app: "browser"
        cmd: "firefox --private-window"
        size: 50
        post:
          - "firefox --new-tab resource1"
          - "firefox --new-tab resource2"
      - app: "terminal"
        cmd: "kitty --working-directory ~/projects"
        size: 50
        post:
          - "echo 'Welcome to your workspace'"
```

> [!NOTE]
>
> - Experiment with different layouts and sizes
> - The `focus` section helps manage multi-monitor setups
> - Custom commands and post-launch actions provide flexibility
