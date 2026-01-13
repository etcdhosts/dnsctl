# dnsctl

> Command line tool for managing DNS host records in etcd.

## Installation

### From Source

```sh
git clone https://github.com/etcdhosts/dnsctl.git
cd dnsctl
task build

# Binary will be in dist/
./dist/dnsctl --help
```

### Install to GOPATH

```sh
task install
dnsctl --help
```

## Configuration

Create `~/.dnsctl.yaml`:

```yaml
endpoints:
  - https://172.16.1.21:2379
  - https://172.16.1.22:2379
  - https://172.16.1.23:2379
key: /etcdhosts
ca: /etc/etcd/ssl/etcd-ca.pem
cert: /etc/etcd/ssl/etcd.pem
cert_key: /etc/etcd/ssl/etcd-key.pem
```

Or generate an example config:

```sh
dnsctl example > ~/.dnsctl.yaml
```

### Configuration Options

| Option | Description |
|--------|-------------|
| `endpoints` | etcd endpoints list |
| `key` | etcd key or prefix (default: `/etcdhosts`) |
| `dial_timeout` | Connection timeout (default: `5s`) |
| `req_timeout` | Request timeout (default: `5s`) |
| `ca` | CA certificate file |
| `cert` | Client certificate file |
| `cert_key` | Client key file |
| `username` | etcd username |
| `password` | etcd password |

## Usage

### List Records

```sh
# Default hosts format
dnsctl list

# JSON output
dnsctl list -o json

# YAML output
dnsctl list -o yaml

# Read from specific revision
dnsctl list -r 12345
```

### Edit Records

Use system editor to edit DNS records:

```sh
# Open in default editor ($EDITOR or vi)
dnsctl edit

# Use specific editor
EDITOR=nano dnsctl edit
EDITOR=vim dnsctl edit
```

Features:
- Auto-deduplication (removes duplicate records)
- Validates hosts format before saving
- Preserves extended attributes (weight, TTL, health check)

### View History

```sh
# List all versions (single-key mode)
dnsctl history

# Limit to recent versions
dnsctl history -n 10

# Per-host mode: list all domains
dnsctl history

# Per-host mode: view domain history
dnsctl history example.com
```

Output example:
```
REVISION        VERSION     RECORDS   MODIFIED
--------------  ----------  --------  --------------------
12350           3           5         2024-01-12 10:30:00 (latest)
12340           2           4         2024-01-12 09:15:00
12330           1           3         -
```

### Compare Versions

```sh
# Compare two revisions
dnsctl diff 12340 12350

# Compare old revision with current (0 = current)
dnsctl diff 12340 0
```

Output with colors:
- Red (`-`): removed lines
- Green (`+`): added lines

### Purge Hostname

```sh
# Delete all records for hostname
dnsctl purge example.com
```

### Other Commands

```sh
# Show version
dnsctl --version

# Show help
dnsctl --help

# Use custom config file
dnsctl -c /path/to/config.yaml list

# Generate example config
dnsctl example > ~/.dnsctl.yaml
```

## Examples

### Setup Load Balanced Service

```sh
# Use editor to configure multiple backends with weights
dnsctl edit

# In editor, add:
# 192.168.1.1 api.example.com # +etcdhosts weight=3
# 192.168.1.2 api.example.com # +etcdhosts weight=2
# 192.168.1.3 api.example.com # +etcdhosts weight=1

# Verify
dnsctl list
```

### Migrate Backend

```sh
# Use editor to modify records
dnsctl edit

# In editor:
# - Add new backend: 192.168.1.10 api.example.com # +etcdhosts weight=1
# - Delete old backend line: 192.168.1.1 api.example.com
```

### Export and Backup

```sh
# Export to JSON
dnsctl list -o json > backup.json

# Export to YAML
dnsctl list -o yaml > backup.yaml
```

## Building

```sh
# Build for current platform
task build

# Run tests
task test

# Run integration tests (requires Docker)
task test-integration

# Build all release binaries
task release

# Create GitHub release
task gh-release
```

## Related Projects

- [etcdhosts](https://github.com/etcdhosts/etcdhosts) - CoreDNS plugin
- [client-go](https://github.com/etcdhosts/client-go) - Go client library

## License

Apache 2.0
