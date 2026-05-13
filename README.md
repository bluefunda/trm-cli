# requests CLI

Command line interface for the [BlueFunda](https://bluefunda.com) bluerequests platform — event-driven change and release management for SAP operations.

The binary is `requests`. Use it to manage change requests, subscribe to events, manage users, and interact with the TRM platform via gRPC — all requests flow through `trm-bff`, never directly to the backend.

## Installation

### One-line installer (macOS and Linux)

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/bluefunda/trm-cli/main/install.sh)"
```

Installs to `/usr/local/bin` if writable, otherwise `~/.local/bin`. Override with `REQUESTS_INSTALL_DIR`.

### Homebrew (macOS)

```bash
brew tap bluefunda/tap
brew install --cask requests
```

### Debian / Ubuntu

```bash
curl -sL https://github.com/bluefunda/trm-cli/releases/latest/download/requests_<version>_linux_amd64.deb -o requests.deb
sudo dpkg -i requests.deb
```

### RHEL / Fedora / Rocky

```bash
sudo dnf install https://github.com/bluefunda/trm-cli/releases/latest/download/requests_<version>_linux_amd64.rpm
```

### From GitHub Releases

Download the binary for your platform from the [releases page](https://github.com/bluefunda/trm-cli/releases/latest).

| Platform      | Archive                                  |
|---------------|------------------------------------------|
| macOS (ARM64) | `requests_<version>_darwin_arm64.zip`    |
| macOS (AMD64) | `requests_<version>_darwin_amd64.zip`    |
| Linux (AMD64) | `requests_<version>_linux_amd64.tar.gz`  |
| Linux (ARM64) | `requests_<version>_linux_arm64.tar.gz`  |

### From Source

```bash
go install github.com/bluefunda/trm-cli/cmd/requests@latest
```

## Quick Start

```bash
# Authenticate (opens browser for OAuth device flow)
requests login

# Check connection
requests health

# List change requests
requests cr list

# List for a specific project, with filters
requests cr list --project <project-id> --status open --severity high

# Get a single change request
requests cr get <id>

# Create a change request
requests cr create --description "Enable new payment gateway" --project <id> --type "feature" --severity "medium"

# Update a change request
requests cr update <id> --status "in-progress" --assignee john.doe

# Advance stage
requests cr stage <id> --stage "qa"

# List comments
requests cr comment list <cr-id>

# Add a comment
requests cr comment add <cr-id> --message "Approved for QA deployment"

# Show current user
requests user info
```

## Commands

### Auth & Connectivity

| Command | Description |
|---------|-------------|
| `requests login` | Authenticate via OAuth2 device flow |
| `requests health` | Check gRPC connection to the backend |
| `requests user info` | Show current user details |
| `requests version` | Print CLI version |

### Change Requests

| Command | Description |
|---------|-------------|
| `requests cr list [flags]` | List change requests (with optional filters) |
| `requests cr get <id>` | Get a change request by ID |
| `requests cr create [flags]` | Create a new change request |
| `requests cr update <id> [flags]` | Update fields on a change request |
| `requests cr delete <id>` | Archive a change request |
| `requests cr stage <id> --stage <stage>` | Update the workflow stage |

#### `cr list` flags

| Flag | Description |
|------|-------------|
| `--project` | Filter by project ID |
| `--description` | Filter by description (partial match) |
| `--status` | Filter by status |
| `--type` | Filter by request type |
| `--severity` | Filter by severity |
| `--assignee` | Filter by assignee |
| `--archived` | Include archived change requests |

#### `cr create` / `cr update` flags

| Flag | Description |
|------|-------------|
| `--description` | Change request description |
| `--project` | Project ID |
| `--type` | Request type |
| `--severity` | Severity level |
| `--assignee` | Assignee username |
| `--status` | *(update only)* New status |

### Comments

| Command | Description |
|---------|-------------|
| `requests cr comment list <cr-id>` | List comments on a change request |
| `requests cr comment add <cr-id> --message "..."` | Add a comment |
| `requests cr comment update <comment-id> --message "..."` | Edit a comment |
| `requests cr comment delete <comment-id>` | Delete a comment |

### Events (low-level)

| Command | Description |
|---------|-------------|
| `requests events subscribe [pattern]` | Stream realm-scoped events (Ctrl-C to stop) |
| `requests events publish <subject> <data>` | Publish an event to a realm-scoped subject |
| `requests rpc request <subject> <data>` | Low-level request-reply |

Run `requests <command> --help` for full options.

## Authentication

requests CLI uses the **OAuth2 device authorization flow**:

1. `requests login` requests a device code
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
