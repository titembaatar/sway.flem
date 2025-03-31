# Nested Layout Configurations

This document explains how to use the nested container feature in sway.flem to create complex window layouts.

## Container Structure

Containers allow you to define nested groups of applications with different layouts. This gives you much more control over how your workspace is organized.

### Basic Container Syntax

Here's the basic structure of a container configuration:

```yaml
container:
  layout: "split_type"  # splith or splitv only for containers
  size: "70ppt"        # Size of the container relative to parent's layout
  apps:
    - name: "app1"
      size: "50ppt"    # Size relative to container's layout
    - name: "app2"
      size: "50ppt"    # Size relative to container's layout
  container:  # Optional nested container
    layout: "another_split"
    apps:
      - name: "app3"
      - name: "app4"
```

> [!NOTE]
>
> For containers, only `splith` and `splitv` are valid values for the layout field,
> as these correspond to the Sway split commands.
> The full layout options (`tabbed` and `stacking`) are available yet.

### How Containers Work

When sway.flem processes a container:

1. It launches the first app in the container
2. It sets the container layout (splith, splitv, etc.)
3. It launches the remaining apps in the container
4. If a nested container is present, it repeats the process recursively

This approach leverages Sway's automatic focusing behavior to create the desired layout structure.

## Example Configuration

Here's an example of a workspace with nested containers that creates a 3-column layout
with the third column containing multiple rows and a split in the bottom row:

```yaml
workspaces:
  2:
    layout: "splith"
    apps:
      - name: "kitty1"
      - name: "kitty2"
    container:
      layout: "splitv"
      apps:
        - name: "kitty3"
        - name: "kitty4"
      container:
        layout: "splith"
        apps:
          - name: "kitty5"
          - name: "kitty6"
```

This configuration will result in the following layout:

```
┌───────────┬───────────┬───────────┐
│           │           │           │
│           │           │     3     │
│           │           │           │
│           │           ├───────────┤
│           │           │           │
│     1     │     2     │     4     │
│           │           │           │
│           │           ├─────┬─────┤
│           │           │     │     │
│           │           │  5  │  6  │
│           │           │     │     │
└───────────┴───────────┴─────┴─────┘
```

## Tips for Working with Nested Layouts

1. **Order matters**: Apps are launched in the order they appear in the configuration.
2. **Layout/Split selection**:
   - For workspaces, you can use: `splith`, `splitv`, `tabbed`, or `stacking`
   - For containers, you can only use: `splith` or `splitv` (for now)
3. **Testing**: Start with a simple nested structure and build up to more complex layouts.
4. **Sizing**: See [Container Sizing section](#container-sizing)

## Container Sizing

How size works with nested containers:

1. **Container Size**: The `size` property on a container applies to the first app within that container and is relative to the parent's layout:
   - If parent has `splith` layout, the container's size is applied as width
   - If parent has `splitv` layout, the container's size is applied as height

2. **App Size**: The `size` property on an app is relative to its container's layout:
   - In a `splith` container, app size affects width
   - In a `splitv` container, app size affects height

### Example:

```yaml
workspaces:
  2:
    layout: "splith"  # Horizontal layout at workspace level
    apps:
      - name: "app1"
        size: "20ppt"  # 20% of workspace width
    container:
      layout: "splitv"  # Vertical layout for container
      size: "80ppt"     # 80% of workspace width (because parent is splith)
      apps:
        - name: "app2"
          size: "70ppt"  # 70% of container height (because container is splitv)
        - name: "app3"
          size: "30ppt"  # 30% of container height
```

This creates a workspace where:
- App1 takes 20% of the width
- The container takes 80% of the width
- Within the container, app2 takes 70% of the height and app3 takes 30%

