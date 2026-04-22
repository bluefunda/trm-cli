# AGENTS.md

Instructions for AI coding agents working on trm-cli.

## Project Overview

Go 1.24 CLI for the TRM (bluerequests) change/release management platform. Uses Cobra for commands, gRPC for backend communication with the trm-bff service, and OAuth2 device flow (Keycloak) for authentication.

Binary name: `trm`
Module: `github.com/bluefunda/trm-cli`
Config location: `~/.trm/config.yaml`

## Build and Test Commands

```bash
# MUST-RUN before submitting any change
make build          # Build binary (runs go mod tidy first)
make test           # go test -v -race -count=1 ./...
make vet            # go vet ./...
make fmt            # gofmt -w .

# Other targets
make test-cover     # Coverage report
make proto          # Regenerate protobuf code from api/proto/bff.proto
make snapshot       # goreleaser snapshot (test release build)
```

### Validation Sequence

Run these in order before committing:

```bash
make fmt
make vet
make test
make build
```

All four must pass with zero errors.

## Project Structure

```
cmd/trm/main.go              # Entry point (delegates to internal/cmd.Execute)
api/proto/
  bff.proto                  # Source-of-truth service definition (DO NOT hand-edit generated files)
  bff/                       # Generated Go code (bff.pb.go, bff_grpc.pb.go)
internal/
  cmd/                       # Cobra command tree
    root.go                  # Root command, global flags, loadConfig(), outputFormat()
    helpers.go               # bffConn(), printer(), reAuthenticate(), saveAuthTokens()
    login.go                 # OAuth device flow login
    health.go                # gRPC health check
    version.go               # Version display
    user.go                  # User commands (info)
    events.go                # Events commands (subscribe, publish)
    rpc.go                   # RPC commands (request)
  grpc/
    conn.go                  # gRPC connection, TLS auto-detect, auth interceptors
  auth/
    auth.go                  # OAuth2 device authorization grant (RFC 8628)
  config/
    config.go                # YAML config load/save, defaults, token validation
  ui/
    output.go                # Printer: table/json/quiet output modes
scripts/
  generate-proto.sh          # Protobuf code generation script
```

## TRM BFF Service

trm-bff exposes four RPCs — the CLI covers all of them:

| Command                        | RPC              | Description                              |
|--------------------------------|------------------|------------------------------------------|
| `trm user info`                | GetUserInfo      | Current user from Keycloak               |
| `trm events subscribe [pat]`   | SubscribeEvents  | Stream realm-scoped events (blocking)    |
| `trm events publish <s> <d>`   | PublishEvent     | Publish event to NATS subject            |
| `trm rpc request <s> <d>`      | RequestReply     | NATS request-reply                       |

Subjects are automatically prefixed with the realm by trm-bff.

## Safe Modification Boundaries

### Safe to modify
- `internal/cmd/*.go` — Add/modify CLI commands
- `internal/ui/*.go` — Change output formatting
- `internal/config/config.go` — Add config fields
- `internal/auth/auth.go` — Modify auth flow
- `internal/grpc/conn.go` — Modify connection/interceptor logic

### Modify with caution
- `api/proto/bff.proto` — Source-of-truth; run `make proto` after changes (requires protoc)
- `Makefile` — Build system
- `.goreleaser.yml` — Release pipeline

### Do NOT modify
- `api/proto/bff/*.pb.go` — Generated files. Run `make proto` instead.
- `.github/workflows/*.yml` — CI/CD (uses shared workflows from `bluefunda/release-foundry`)
- `cmd/trm/main.go` — Entry point; should remain a minimal delegation to `internal/cmd.Execute()`

## Code Conventions

### Command pattern
All commands follow: `trm <resource> <operation> [flags]`

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

Register sub-commands in the parent command's `init()`. Register top-level commands in `root.go` → `init()`.

### Output contract
- **stdout**: Data only (tables, JSON). Used for piping.
- **stderr**: Status messages via `p.OK()`, `p.Info()`, `p.Error()`, `p.Warn()`.
- Always support `table`, `json`, `quiet` output formats for commands that return data.

### Error handling
- Use `RunE` (not `Run`) — return errors, do not `os.Exit()`.
- Wrap errors: `fmt.Errorf("context: %w", err)`.

### gRPC calls
- Get a connection via `bffConn()` — handles auth and token refresh.
- Always `defer conn.Close()`.
- Use `trmgrpc.ContextWithTimeout()` for unary RPCs.
- For streaming RPCs use `context.WithCancel` and handle SIGINT/SIGTERM.

## Git and Branch Conventions

Follows `bluefunda` org-level standards:

- **Conventional Commits**: `feat:`, `fix:`, `perf:`, `security:`, `infra:`, `chore:`, `docs:`, `test:`
- **Branch naming**: `feat/<desc>` or `fix/<desc>`; `main` is protected — PR + CI required
- **CI**: `bluefunda/release-foundry/.github/workflows/go-ci.yml@main`
- **Release**: Release Please → GoReleaser → GitHub Release → Homebrew tap
- **Required secrets**: `GH_PAT`, `HOMEBREW_TAP_TOKEN`

## Dependencies

- `github.com/spf13/cobra` — CLI framework
- `google.golang.org/grpc` — gRPC client
- `google.golang.org/protobuf` — Protobuf runtime
- `github.com/fatih/color` — Terminal colors
- `gopkg.in/yaml.v3` — Config file parsing

Keep dependencies minimal. Run `make tidy` after any change to `go.mod`.

## Known Issues / Follow-up

- `trm-bff/api/proto/bff.proto` has `go_package = "github.com/bluefunda/cai-bff/..."` (copy-paste bug). Fix in trm-bff before next `make proto` run. The generated files already committed are correct.
