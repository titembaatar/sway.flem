# Troubleshooting

This guide helps you solve common issues that may arise when using sway.flem.

## Common Issues

### Configuration File Not Found

**Error Message**:
```
Failed to load configuration: config file not found: open /path/to/config.yaml: no such file or directory
```

**Solution**:
- Verify that the file path is correct
- Check file permissions
- Use an absolute path to the config file

### Invalid Configuration Format

**Error Message**:
```
Failed to decode config: yaml: line X: did not find expected key
```

**Solution**:
- Check your YAML syntax at the specified line
- Ensure proper indentation is used
- Use a YAML validator tool to check your configuration

### Unknown Layout Type

**Error Message**:
```
Configuration error: workspace '1': invalid layout type
```

**Solution**:
- Use one of the supported layout types: `splith`, `splitv`, `tabbed`, `stacking` (or their aliases)
- Check for typos in the layout name

### Missing Required Fields

**Error Message**:
```
Configuration error: workspace '1', container at index 0: app is empty
```

**Solution**:
- Ensure all required fields (app) are set for app containers
- Ensure all required fields (split, containers) are set for nested containers

### Applications Not Launching

**Issue**: Applications are defined in the config but don't launch

**Solutions**:
- Ensure the application is installed and can be launched from the command line
- Use the `-debug` flag to see detailed logs about application launching
- If using a custom command, try simplifying it to just the application name
- Check if the application needs time to start and adjust your workflow accordingly

### Windows Not Resizing Correctly

**Issue**: Windows don't get resized to the specified sizes

**Solutions**:
- Make sure sizes are specified correctly (e.g., `50ppt` for percentage)
- Some applications may override Sway's window management
- Try using the `-debug` flag to see detailed resize operations

## Debugging Tips

### Enable Debug Logging

Use the `-debug` flag to get detailed logs:

```bash
flem sway -config config.yaml -debug
```

This will print extensive information about:
- Configuration parsing
- Application launching
- Window management commands
- Error handling

### Use Dry Run Mode

The `-dry-run` flag validates your configuration without making any changes:

```bash
flem sway -config config.yaml -dry-run
```

This is useful for checking if your configuration is valid before applying it.

### Check Sway IPC

Sometimes issues arise from problems with the Sway IPC interface. You can test basic Sway communication with:

```bash
swaymsg -t get_workspaces
```

If this fails, it indicates a problem with Sway itself rather than sway.flem.

### Known Limitations

- **Application Focus**: Some applications may grab focus when launched, which can interfere with the layout creation process
- **Window Properties**: Some applications may not respect sizing commands or may have minimum size constraints
- **Resize**: Because of Sway behavior, when there is too much containers in a split, sizes are not
  apply correctly. Nothing I can do about it, try different sizes until you are satisfied.
- **Timing Issues**: If applications take a long time to start, the layout process might not be fully completed

## Getting Help

If you're experiencing issues not covered in this guide:

1. Check the [GitHub Issues](https://github.com/titembaatar/sway.flem/issues) to see if others have encountered the same problem
2. Create a new issue with:
   - Your sway.flem version (`flem sway -version`)
   - Your Sway version (`swaymsg -v`)
   - Your configuration file
   - The exact error message or behavior you're seeing
   - Debug logs can help (`flem sway -debug ...`)
