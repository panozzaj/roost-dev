# roost-dev

After making changes that affect runtime behavior, rebuild and restart:

```bash
GOPATH=/Users/anthony/go go build -o ~/go/bin/roost-dev ./cmd/roost-dev && launchctl kickstart -k gui/$(id -u)/com.roost-dev
```
