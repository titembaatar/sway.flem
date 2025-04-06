# sway.flem Roadmap

## Current Version: v0.1.0

## Short-Term Improvements

### Application Launch Intelligence
- [x] Implement Sway tree inspection to detect existing applications
- [x] Add logic to skip launching already running applications
- [x] Add configuration option to control post-command behavior (`rerun-post`)
- [ ] Provide configuration options to force re-launch or skip

### Refresh and Recovery Features
- [ ] Add a "refresh" command to reconfigure existing workspace
- [ ] Implement smart diff between current and desired state
- [ ] Provide options for incremental or full workspace reset

## Near-Term Improvements

### Environment Management
- [ ] Create a mechanism to save current workspace state
- [ ] Develop restore functionality for saved workspace configurations
- [ ] Support partial or full workspace restoration

### Configuration Enhancements
- [ ] Add more flexible application matching (process name, window class)
- [ ] Support for conditional launches
- [ ] Improve error handling and reporting for configuration issues

### Performance and Reliability
- [ ] Optimize application launch sequence
- [ ] Add more robust error recovery mechanisms
- [ ] Improve logging for debugging workspace setup

## 📝 Contribution

As a solo project, contributions and feedback are always welcome!
