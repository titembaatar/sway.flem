# Command Line Options

sway.flem provides several command line options to control its behavior.

## Basic Usage

```
sway.flem -config <config-file>
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
sway.flem -config ~/.config/sway.flem.yaml
```

### `-version`

Displays the current version of sway.flem and exits.

Example:
```bash
sway.flem -version
# Output: sway.flem v0.1.0
```

### `-verbose`

Increases the logging level to include informational messages. This is useful for understanding what sway.flem is doing during setup.

Example:
```bash
sway.flem -config ~/.config/sway.flem.yaml -verbose
```

### `-debug`

Enables debug mode, which provides detailed logging about all operations. This is useful for troubleshooting issues.

Example:
```bash
sway.flem -config ~/.config/sway.flem.yaml -debug
```

### `-dry-run`

Performs a validation of the configuration file without actually making any changes to Sway. This is useful for testing whether your configuration is valid before applying it.

Example:
```bash
sway.flem -config ~/.config/sway.flem.yaml -dry-run
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
exec sway.flem -config ~/.config/sway.flem.yaml
```

### Running on Demand

You can also create a keybinding in your Sway config to run sway.flem on demand:

```
# ~/.config/sway/config
bindsym $mod+Shift+r exec sway.flem -config ~/.config/sway.flem.yaml
```

This would allow you to reset your workspace layout by pressing the defined key combination.
