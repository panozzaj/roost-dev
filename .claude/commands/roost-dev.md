# Set up roost-dev environment for this project

Set up this project to run well in roost-dev (local development environment manager).

## Getting Started

First, run these commands to understand roost-dev:

- `roost-dev --help` - See available CLI commands
- `roost-dev docs` - Read full documentation

Config files are in: `~/.config/roost-dev/`

## Preferences

- For projects with both frontend and backend, use a multi-service setup
- The web/frontend should be the default service with `default: true` (accessed at `<project>.test`)
- The API/backend should be a separate service named `api` (accessed at `api-<project>.test`)
- Use `depends_on: [api]` on the web service to ensure the backend starts first
- For Vite projects:
    - Add the roost-dev host to `server.allowedHosts` in vite.config.ts
    - Set the API proxy target to `http://api-<project>.test`
