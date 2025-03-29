# Troubleshooting sway.flem

This document provides solutions for common issues you might encounter when using sway.flem.

## Sway Not Running

**Issue:** Error message "Sway window manager is not running"

**Solution:** Make sure you're running sway.flem within a Sway session.
sway.flem cannot operate without Sway running.

## Configuration File Not Found

**Issue:** Error loading config file

**Solution:** Check the path to your configuration file.
Use the `--config` option to specify the correct path,
or place your configuration file at the default location `./sway.flem.yaml`.

## Applications Not Starting

**Issue:** Applications defined in the configuration don't start

**Solutions:**
- Check that the command is correct and works when run manually
- Run sway.flem with the `--verbose` flag to see debug information
- Make sure the application name matches what's in your config (app_id or class)

## Windows Not Properly Configured

**Issue:** Windows start but aren't positioned or sized correctly

**Solutions:**
- Verify your size/position syntax is correct
- Check if the app's window type supports the operation (some windows resist resizing)
- Increase the delay value for slow-starting applications
- Run with `--verbose` to check for errors during configuration

## Finding the Correct App Name

**Issue:** Not sure what name to use for an application

**Solution:** Use this command to see all app_ids and classes:

```bash
swaymsg -t get_tree | grep -E 'app_id|class'
```

Launch your application, then run this command to see what app_id or class it uses.

## Workspace Not Moving to Correct Output

**Issue:** Workspace shows on the wrong monitor

**Solutions:**
- Check that the output name is correct (use `swaymsg -t get_outputs` to list outputs)
- Some applications may override workspace assignments with their own rules
- Try using the correct output name (e.g., `DP-1`, `HDMI-A-1`) rather than generic names

## Close Unmatched Not Working

**Issue:** Windows that aren't in the config aren't being closed

**Solution:** Make sure `close_unmatched: true` is set for the workspace,
and verify that the window's app_id/class doesn't match any app names in your config.

## Changes Not Taking Effect After Reload

**Issue:** Configuration changes don't take effect after reloading Sway

**Solution:** Make sure your `exec_always` command in Sway config is correct.
Try manually running sway.flem with the correct config path.
