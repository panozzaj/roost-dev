# roost-dev

A local development proxy for all your projects. Like [puma-dev](https://github.com/puma/puma-dev), but language-agnostic and with multi-service support.

## Features

- **Language-agnostic**: Works with any web server (Node, Ruby, Python, Elixir, Go, Rust, etc.)
- **Multi-service projects**: Run frontend and backend on separate subdomains
- **Dynamic port allocation**: No more port conflicts between projects
- **On-demand startup**: Services start when you access them
- **Process management**: Start, stop, restart from web UI or CLI
- **Static file serving**: Serve HTML files directly
- **Web dashboard**: Monitor all your services at `roost-dev.localhost`

## Installation

### From source

```bash
go install github.com/panozzaj/roost-dev/cmd/roost-dev@latest
```

### Build locally

```bash
git clone https://github.com/panozzaj/roost-dev
cd roost-dev
go build -o roost-dev ./cmd/roost-dev
```

## Quick Start

```bash
# Start on port 8080 (no sudo required)
roost-dev -http-port 8080

# Or on port 80 (requires sudo)
sudo roost-dev
```

Then open http://roost-dev.localhost:8080 to see the dashboard.

## Configuration

Place config files in `~/.config/roost-dev/`. The filename becomes the app name.

### Level 1: Port File (simplest)

Proxy to an already-running server:

```bash
echo "3000" > ~/.config/roost-dev/myapp
# Access at http://myapp.localhost
```

### Level 2: Command File (auto-start)

Let roost-dev start and manage the process:

```bash
echo "npm run dev" > ~/.config/roost-dev/myapp
# roost-dev starts it with PORT=<dynamic> and proxies to it
```

### Level 3: Static Files

Serve HTML directly:

```bash
# Point to a directory
ln -s ~/projects/my-static-site ~/.config/roost-dev/mysite

# Or point to a file
echo "~/projects/game/index.html" > ~/.config/roost-dev/game
```

### Level 4: YAML Config (multi-service)

For projects with frontend + backend:

```yaml
# ~/.config/roost-dev/myproject.yml
name: myproject
root: ~/projects/myproject

services:
  backend:
    dir: backend
    cmd: mix phx.server

  frontend:
    cmd: npm start
    env:
      REACT_APP_API_URL: http://backend-myproject.localhost
```

Access at:
- http://frontend-myproject.localhost
- http://backend-myproject.localhost

Both services can share cookies on `.myproject.localhost` - no more CORS headaches!

## How It Works

1. roost-dev listens on port 80 (or your chosen port)
2. When you visit `myapp.localhost`, it:
   - Looks up `~/.config/roost-dev/myapp`
   - If it's a command, starts it with `PORT=<dynamic>`
   - Proxies the request, preserving the `Host` header
3. Your app sees the full hostname (e.g., `admin.myapp.localhost`) and can extract subdomains

## Subdomain Support

Any subdomain is passed through to your app:

```
admin.myapp.localhost → proxied to myapp with Host: admin.myapp.localhost
api.myapp.localhost   → proxied to myapp with Host: api.myapp.localhost
```

Your app can read the `Host` header to determine which subdomain was requested.

For multi-service apps, the pattern is `{service}-{app}.localhost`:

```
frontend-myproject.localhost → frontend service
backend-myproject.localhost  → backend service
```

## CLI Options

```
-dir <path>       Configuration directory (default: ~/.config/roost-dev)
-http-port <n>    HTTP port to listen on (default: 80)
-https-port <n>   HTTPS port to listen on (default: 443)
-tld <domain>     Top-level domain to use (default: localhost)
```

## Using .test Instead of .localhost

The `.localhost` TLD works in browsers without configuration. For CLI tools (curl, etc.), use `.test` with dnsmasq:

```bash
# Install dnsmasq
brew install dnsmasq

# Configure .test to resolve to localhost
echo 'address=/.test/127.0.0.1' >> $(brew --prefix)/etc/dnsmasq.conf

# Set up resolver
sudo mkdir -p /etc/resolver
echo "nameserver 127.0.0.1" | sudo tee /etc/resolver/test

# Start dnsmasq
sudo brew services start dnsmasq

# Run roost-dev with .test TLD
roost-dev -tld test
```

## Framework Support

Most frameworks respect the `PORT` environment variable:

| Framework | Works? | Notes |
|-----------|--------|-------|
| Next.js | ✅ | Reads PORT automatically |
| Vite | ✅ | Use `npm run dev -- --port $PORT` |
| Phoenix | ✅ | Reads PORT in runtime.exs |
| Rails | ✅ | Use `-p $PORT` |
| FastAPI | ✅ | Use `--port $PORT` |
| Express | ✅ | `process.env.PORT` |

## Web Dashboard

Visit http://roost-dev.localhost to:

- See all configured apps and their status
- View which services are running
- Start/stop/restart processes
- View logs

## Comparison with Alternatives

| Feature | roost-dev | puma-dev | hotel |
|---------|-----------|----------|-------|
| Language-agnostic | ✅ | ❌ (Ruby only) | ✅ |
| Multi-service YAML | ✅ | ❌ | ❌ |
| Service env vars | ✅ | ❌ | ❌ |
| Static file serving | ✅ | ❌ | ⚠️ |
| Web UI | ✅ | ❌ | ✅ |
| Active development | ✅ | ✅ | ❌ |

## License

MIT
