# Command Line Interface

## Basic Usage

```bash
flem sway -config <config-file>
```

## Available Options

| Option | Description | Type | Default |
|--------|-------------|------|---------|
| `-config` | Path to the configuration file | String | Required |
| `-version` | Display version information | Flag | - |
| `-verbose` | Enable verbose logging | Flag | Disabled |
| `-debug` | Enable debug mode with detailed logging | Flag | Disabled |
| `-dry-run` | Validate configuration without making changes | Flag | Disabled |

## Detailed Option Reference

### `-config`
- **Required**: Yes
- **Usage**: Specifies the path to the YAML configuration file
- **Example**:
  ```bash
  flem sway -config ~/.config/sway/workspace.yml
  ```

### `-version`
- **Usage**: Displays the current version of sway.flem
- **Example**:
  ```bash
  flem sway -version
  # Output: flem v0.1.0
  ```

### `-verbose`
- **Usage**: Increases logging verbosity
- **Logging Level**: Includes informational messages
- **Recommended**: For understanding setup process details

### `-debug`
- **Usage**: Enables comprehensive debug logging
- **Logging Level**: Includes all debug and informational messages
- **Recommended**: For troubleshooting configuration or launch issues

### `-dry-run`
- **Usage**: Validates configuration without applying changes
- **Helpful For**:
  - Checking configuration syntax
  - Verifying workspace layout
  - Catching potential errors before execution

## Usage Examples

### Basic Configuration
```bash
flem sway -config ~/workspace.yml
```

### Verbose Logging
```bash
flem sway -config ~/workspace.yml -verbose
```

### Debug Mode
```bash
flem sway -config ~/workspace.yml -debug
```

### Dry Run Validation
```bash
flem sway -config ~/workspace.yml -dry-run
```

## Integration with Sway Config

You can integrate sway.flem into your Sway configuration:

```
# In ~/.config/sway/config
exec flem sway -config ~/.config/sway/workspace.yml
```

Or create a keybinding:

```
# Reset workspace layout
bindsym $mod+Shift+w exec flem sway -config ~/.config/sway/workspace.yml
```

## Logging Levels

1. **ERROR**: Critical issues only (default)
2. **WARN**: Warnings and errors
3. **INFO**: Informational messages (with `-verbose`)
4. **DEBUG**: Detailed diagnostic information (with `-debug`)
