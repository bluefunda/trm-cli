# AGENTS.md

Instructions for AI coding agents working on trm-cli.

## Project Overview

Go 1.24 CLI for the BlueFunda bluerequests change/release management platform. All requests flow through `trm-bff` via gRPC — the CLI never talks to NATS or backend services directly.

Binary name: `requests`
Module: `github.com/bluefunda/trm-cli`
Entry point: `cmd/requests/main.go`
Config location: `~/.trm/config.yaml`

## Build and Test Commands

```bash
make build          # go mod tidy + build binary to ./requests
make vet            # go vet ./...
make fmt            # gofmt -w .
make tidy           # go mod tidy
make test           # go test -v -race -count=1 ./...
make test-cover     # Coverage report
make proto          # Regenerate protobuf bindings (requires protoc)
make snapshot       # goreleaser snapshot (test release build)
```

### Validation Sequence

```bash
make fmt && make vet && make test && make build
```

All must pass with zero errors before committing.

## Project Structure

```
cmd/requests/main.go            # Entry point: invokes root command
api/proto/
  bff.proto                     # BFFService definition — keep in sync with trm-bff
  bff/
    bff.pb.go                   # Generated — do not edit by hand
    bff_grpc.pb.go              # Generated — do not edit by hand
internal/
  cmd/
    root.go                     # Root cobra command, global flags (--bff, --realm, --output)
    helpers.go                  # bffConn(), printer(), token refresh utilities
    login.go                    # OAuth2 device authorization flow
    health.go                   # gRPC health check
    version.go                  # Version display
    user.go                     # User info command
    cr.go                       # Change request + comment commands (list/get/create/update/delete/stage)
    events.go                   # Event subscribe/publish commands
    rpc.go                      # Low-level request-reply command
  grpc/
    conn.go                     # gRPC connection, TLS auto-detect, auth interceptors
  auth/
    auth.go                     # OAuth2 device authorization grant (RFC 8628)
  config/
    config.go                   # ~/.trm/config.yaml loader; token storage
  ui/
    output.go                   # Printer: table/json/quiet output modes
scripts/
  generate-proto.sh             # Runs protoc with module= path mode
```

## BFFService RPC Surface

| Command group          | RPC                      | Description                          |
|------------------------|--------------------------|--------------------------------------|
| `requests user info`   | GetUserInfo              | Current user from Keycloak           |
| `requests events sub`  | SubscribeEvents          | Stream realm-scoped events           |
| `requests events pub`  | PublishEvent             | Publish event to NATS subject        |
| `requests rpc request` | RequestReply             | NATS request-reply                   |
| `requests cr list`     | ListChangeRequests       | List change requests with filters    |
| `requests cr get`      | GetChangeRequest         | Get a single change request by ID    |
| `requests cr create`   | CreateChangeRequest      | Create a new change request          |
| `requests cr update`   | UpdateChangeRequest      | Update fields on a change request    |
| `requests cr delete`   | DeleteChangeRequest      | Archive a change request             |
| `requests cr stage`    | UpdateChangeRequestStage | Advance the workflow stage           |
| `requests cr comment list`   | ListComments       | List comments on a change request    |
| `requests cr comment add`    | AddComment         | Add a comment                        |
| `requests cr comment update` | UpdateComment      | Edit a comment                       |
| `requests cr comment delete` | DeleteComment      | Delete a comment                     |

## Adding a New Command

1. Create `internal/cmd/<name>.go` with a cobra command variable
2. Register it in `root.go` via `rootCmd.AddCommand(<name>Cmd)`
3. Use `bffConn()` from `helpers.go` to get the gRPC connection
4. All RPC calls must go through `trm-bff` — never connect to NATS or backend directly

## Proto Workflow

`api/proto/bff.proto` must stay in sync with `trm-bff/api/proto/bff.proto`. When adding RPCs:

1. Add the RPC and message types to both `bff.proto` files
2. Run `make proto` in trm-cli to regenerate client bindings
3. Implement the client call in `internal/cmd/*.go`
4. Implement the server handler in trm-bff `internal/transport/grpc/handler.go`

## Code Conventions

### Command pattern

```go
var fooCmd = &cobra.Command{
    Use:   "foo",
    Short: "Short description",
    RunE: func(cmd *cobra.Command, args []string) error {
        conn, cfg, err := bffConn()
        if err != nil {
            return err
        }
        defer conn.Close()
        // ...
    },
}
```

Register subcommands in the parent's `init()`. Register top-level commands in `root.go` → `init()`.

### Output contract
- **stdout**: Data only (tables, JSON) — used for piping
- **stderr**: Status messages via `p.OK()`, `p.Info()`, `p.Error()`, `p.Warn()`
- Support `table`, `json`, `quiet` output formats for commands that return data

### Error handling
- Use `RunE` (not `Run`) — return errors, do not `os.Exit()`
- Wrap errors: `fmt.Errorf("context: %w", err)`

## Git and Branch Conventions

Follows `bluefunda` org-level standards:

- **Conventional Commits**: `feat:`, `fix:`, `chore:`, `docs:`, `perf:`
- **Branch naming**: `<type>/<short-description>`
- **PRs**: squash-merged to `main`; title must use conventional commit format
- **CI**: `bluefunda/release-foundry` shared workflows
- **Release**: Release Please → GoReleaser → GitHub Release (linux/darwin amd64/arm64 + deb/rpm) + Homebrew tap
- **Required secrets**: `GH_PAT`, `HOMEBREW_TAP_TOKEN`

## Dependencies

- `github.com/spf13/cobra` — CLI framework
- `github.com/fatih/color` — Terminal colour output
- `google.golang.org/grpc` — gRPC client
- `google.golang.org/protobuf` — Protobuf runtime
- `gopkg.in/yaml.v3` — Config file parsing

## Do NOT

- Connect to NATS, trm-backend-go, or Keycloak directly from the CLI — all traffic goes through `trm-bff` via gRPC
- Edit `api/proto/bff/bff.pb.go` or `bff_grpc.pb.go` by hand — regenerate with `make proto`
- Commit tokens, credentials, or `~/.trm/` contents
- Modify `.github/workflows/` without explicit request
