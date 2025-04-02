# Command Line Options

sway.flem provides several command line options to control its behavior.

## Basic Usage

```
flem sway -config <config-file>
```

## Available Options

| Option | Description |
|--------|-------------|
| `-config <file>` | Path to the configuration file (required) |
| `-version` | Show version information and exit |
| `-verbose` | Enable verbose logging |
| `-debug` | Enable debug mode with extra logging |
| `-dry-run` | Validate configuration without making changes |

## Detailed Explanation

### `-config <file>` (Required)

Specifies the path to the YAML configuration file that defines your workspace layout. This is the only required option.

Example:
```bash
flem sway -config ~/.config/sway/config.yml
```

### `-version`

Displays the current version of sway.flem and exits.

Example:
```bash
flem sway -version
# Output: flem sway v0.1.0
```

### `-verbose`

Increases the logging level to include informational messages. This is useful for understanding what sway.flem is doing during setup.

Example:
```bash
flem sway -config ~/.config/sway/config.yml -verbose
```

### `-debug`

Enables debug mode, which provides detailed logging about all operations. This is useful for troubleshooting issues.

Example:
```bash
flem sway -config ~/.config/sway/config.yml -debug
```

### `-dry-run`

Performs a validation of the configuration file without actually making any changes to Sway. This is useful for testing whether your configuration is valid before applying it.

Example:
```bash
flem sway -config ~/.config/sway/config.yml -dry-run
```

## Log Levels

sway.flem has the following log levels, from least to most verbose:

1. **ERROR** - Only error messages (default)
2. **WARN** - Error and warning messages
3. **INFO** - Error, warning, and informational messages (enabled with `-verbose`)
4. **DEBUG** - All messages including detailed debug information (enabled with `-debug`)

## Environment Integration

### Sway Config Integration

You can add sway.flem to your Sway config to automatically set up your workspaces when Sway starts:

```
# ~/.config/sway/config
exec flem sway -config ~/.config/sway/config.yml
```

### Running on Demand

You can also create a keybinding in your Sway config to run sway.flem on demand:

```
# ~/.config/sway/config
bindsym $mod+Shift+r exec flem sway -config ~/.config/sway/config.yml
```

This would allow you to reset your workspace layout by pressing the defined key combination.

### Using with Multiple Monitors

For multi-monitor setups, you can use the `focus` setting in your configuration to
ensure the right workspaces are focused on each monitor when sway.flem completes:

```yaml
# In your config.yml file
focus:
  - "6"  # Focus workspace 6 on one monitor
  - "1"  # Focus workspace 1 on another monitor

workspaces:
  # Your workspace configurations...
```

This works well if you've already assigned workspaces to specific monitors in your Sway config:

```
# In your Sway config (~/.config/sway/config)
workspace 1 output DP-1
workspace 2 output DP-1
workspace 6 output DP-2
workspace 7 output DP-2
```
