# TRM CLI

Command line interface for the [BlueFunda](https://bluefunda.com) TRM (Transport/Release Management) platform — event-driven change and release management for SAP operations.

The binary is `trm`. Use it to subscribe to events, publish messages, manage users, and interact with the BFF gRPC API from the terminal.

## Installation

### One-line installer (macOS and Linux)

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/bluefunda/trm-cli/main/install.sh)"
```

Installs to `/usr/local/bin` if writable, otherwise `~/.local/bin`. Override with `TRM_INSTALL_DIR`.

### Homebrew (macOS)

```bash
brew tap bluefunda/tap
brew install --cask trm
```

### Debian / Ubuntu

```bash
curl -sL https://github.com/bluefunda/trm-cli/releases/latest/download/trm_<version>_linux_amd64.deb -o trm.deb
sudo dpkg -i trm.deb
```

### RHEL / Fedora / Rocky

```bash
sudo dnf install https://github.com/bluefunda/trm-cli/releases/latest/download/trm_<version>_linux_amd64.rpm
```

### From GitHub Releases

Download the binary for your platform from the [releases page](https://github.com/bluefunda/trm-cli/releases/latest).

| Platform      | Archive                            |
|---------------|------------------------------------|
| macOS (ARM64) | `trm_<version>_darwin_arm64.zip`   |
| macOS (AMD64) | `trm_<version>_darwin_amd64.zip`   |
| Linux (AMD64) | `trm_<version>_linux_amd64.tar.gz` |
| Linux (ARM64) | `trm_<version>_linux_arm64.tar.gz` |

### From Source

```bash
go install github.com/bluefunda/trm-cli/cmd/trm@latest
```

## Quick Start

```bash
# Authenticate (opens browser for OAuth device flow)
trm login

# Check connection
trm health

# Subscribe to all events
trm events subscribe

# Subscribe to a specific pattern
trm events subscribe "orders.>"

# Publish an event
trm events publish "orders.created" '{"id":"123"}'

# Show current user
trm user info
```

## Commands

| Command | Description |
|---------|-------------|
| `trm login` | Authenticate via OAuth2 device flow |
| `trm health` | Check gRPC connection to the backend |
| `trm events subscribe [pattern]` | Stream realm-scoped events (Ctrl-C to stop) |
| `trm events publish <subject> <data>` | Publish an event to a realm-scoped subject |
| `trm rpc request <subject> <data>` | Low-level request-reply |
| `trm user info` | Show current user details |
| `trm version` | Print CLI version |

Run `trm <command> --help` for full options.

## Authentication

TRM CLI uses the **OAuth2 device authorization flow**:

1. `trm login` requests a device code
2. Your browser opens the verification URL
3. The CLI polls for authorization completion
4. Tokens are stored locally in `~/.trm/`

Tokens are refreshed automatically — you only need to log in once.

## Configuration

```yaml
# ~/.trm/config.yaml
bff_url: grpc.bluefunda.com:443
```

## License

See [LICENSE](./LICENSE).
