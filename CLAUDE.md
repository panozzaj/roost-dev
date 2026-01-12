# roost-dev

## Managing roost-dev during development

There are TWO versions of roost-dev:
1. **Installed version**: `/Users/anthony/go/bin/roost-dev` (managed by launchd)
2. **Local build**: `/Users/anthony/Documents/dev/roost-dev/roost-dev`

### Option 1: Use air for hot-reloading (recommended)

[air](https://github.com/air-verse/air) watches for file changes and auto-rebuilds. Config is in `.air.toml`.

```bash
# Install air (one-time)
go install github.com/air-verse/air@latest

# Stop the launchd-managed version and any running instances
launchctl bootout gui/$(id -u)/com.roost-dev 2>/dev/null || true
pkill -9 roost-dev 2>/dev/null || true

# Run with air (auto-rebuilds on file changes)
air
```

### Option 2: Use restart script

```bash
# Rebuild and restart (uses launchd)
./scripts/restart.sh
```

### Debug request handling:

```bash
# View server request logs
curl -s "http://roost-dev.test/api/server-logs" | jq -r '.[]'
```

## Code patterns

- **Use non-blocking operations in HTTP handlers.** For process management, prefer `StartAsync()` over `Start()` in API handlers so responses return immediately. The dashboard polls for status updates.
- **Avoid holding mutexes while waiting.** Release locks before any operation that could block (network calls, waiting for ports, etc.).
- **Always background server processes.** When starting roost-dev from bash, use `run_in_background: true` or append `&` to avoid blocking the conversation. Use `tee` to capture output: `/path/to/roost-dev 2>&1 | tee ./tmp/roost-dev.log &`

## UI patterns

- **Icon buttons must have hover tooltips.** Any button that uses an icon (instead of or in addition to text) must have a `title` attribute providing a descriptive tooltip explaining what the button does.
