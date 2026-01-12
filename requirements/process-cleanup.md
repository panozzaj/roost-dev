# Process Cleanup Improvements

## Problem

Some apps (Next.js, webpack, etc.) spawn child processes that don't get killed properly when roost-dev stops them. This leaves orphaned processes that:
- Hold onto ports
- Consume resources
- Require manual cleanup or wrapper scripts

## Current Implementation

roost-dev already implements process group killing:

```go
// Process started in its own process group
cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

// Kill() sends SIGTERM/SIGKILL to entire group
syscall.Kill(-pgid, syscall.SIGTERM)
time.Sleep(100 * time.Millisecond)
syscall.Kill(-pgid, syscall.SIGKILL)
```

## Why It Sometimes Fails

1. **100ms grace period too short** - Complex apps need time to shut down gracefully before SIGKILL
2. **Escaped process groups** - Node.js `spawn({ detached: true })` creates new process groups that escape the `-pgid` kill
3. **Double-forked daemons** - Some processes daemonize themselves, escaping the original process group

## Proposed Solutions

### Option A: Increase grace period (simplest)

Change the 100ms wait to 2-3 seconds.

**Tradeoffs:**
- (+) Simple one-line change
- (+) Helps most apps
- (-) Slower restarts for quick apps
- (-) Doesn't help if children escaped to new process group

### Option B: Configurable stop_timeout

Add config option per-app:

```yaml
stop_timeout: 5s
```

**Tradeoffs:**
- (+) Flexible for different apps
- (+) Power users can tune it
- (-) More config surface
- (-) Default still needs to be reasonable

### Option C: Recursive process tree kill

Use `pgrep -P $PID` to find all descendants and kill them recursively, regardless of process group.

```bash
# Find all descendants of a PID
pgrep -P $PID
```

**Tradeoffs:**
- (+) Most robust - catches escaped children
- (+) Works even with detached Node processes
- (-) External dependency on pgrep (in macOS base system)
- (-) More complex implementation

### Option D: Document wrapper script pattern

No code changes. Document that problematic apps can use wrapper scripts:

```bash
#!/bin/bash
trap 'kill -TERM -$$' SIGTERM SIGINT SIGHUP EXIT
exec npm run dev
```

**Tradeoffs:**
- (+) No code changes
- (+) Puts burden on edge cases
- (-) Poor UX for common tools like Next.js

## Recommendation

1. **Short term:** Option A - increase grace period to 2 seconds
2. **Optional:** Add Option B for per-app tuning via config
3. **Long term:** Option C if simpler fixes don't help enough

## Affected Tools

Known tools that may have this issue:
- Next.js (spawns webpack workers)
- Webpack dev server
- Node.js apps with child_process
- Python with multiprocessing
- Any tool using detached child processes
