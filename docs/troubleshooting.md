# Troubleshooting Guide

## Common Issues

### Configuration Errors

#### Invalid Configuration File
**Symptoms**:
- Error when running `flem sway`
- Configuration not loading

**Solutions**:
1. Verify YAML syntax
2. Use `-dry-run` flag to validate configuration
3. Check indentation and structure

**Example of Validation**:
```bash
flem sway -config config.yml -dry-run
```

#### Incorrect Layout Types
**Symptoms**:
- Unexpected workspace layout
- Layout not applied correctly

**Supported Layouts**:
- `splith` (horizontal)
- `splitv` (vertical)
- `tabbed`
- `stacking`

### Application Launch Problems

#### Application Not Starting
**Symptoms**:
- Specified application fails to launch
- No error in sway.flem output

**Troubleshooting Steps**:
1. Verify application is installed
2. Check application command in configuration
3. Use full path to application
4. Run with `-debug` for detailed logs

**Example Configuration**:
```yaml
containers:
  - app: "firefox"
    cmd: "/usr/bin/firefox"  # Use full path if needed
```

### Workspace and Layout Issues

#### Unexpected Container Sizes
**Symptoms**:
- Containers not sized as expected
- Uneven workspace distribution

**Troubleshooting**:
- Verify size specifications
- Use percentage points (`ppt`)
- Ensure total percentages don't exceed 100%

**Correct Example**:
```yaml
containers:
  - app: "app1"
    size: "50ppt"  # Exactly 50%
  - app: "app2"
    size: "50ppt"  # Exactly 50%
```

> [!NOTE]
>
> The resizing is handle by Sway, over 2 "containers" resizing one "container" will resize the
> others. I did not find any workaround. Try changing sizes to get close to what you want.

### Logging and Debugging

#### Enabling Verbose Logging
```bash
flem sway -config config.yml -verbose
```

#### Enabling Debug Logging
```bash
flem sway -config config.yml -debug
```

## Advanced Troubleshooting

### Checking Sway Compatibility
1. Verify Sway version
2. Ensure `swaymsg` is available
3. Check Wayland compatibility

### Dependency Verification
- Go version 1.24.1+
- Sway window manager
- Required system utilities

## Reporting Issues

1. Gather log output
2. Prepare configuration file (remove sensitive information)
3. Open an issue on GitHub
   - Include version information
   - Provide detailed steps to reproduce

## Performance Considerations

- Minimize number of applications in a workspace
- Use reasonable delay between application launches
- Avoid overly complex nested container structures
