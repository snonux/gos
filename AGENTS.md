# AGENTS.md

## Project Overview

Gos (Go Social Media) is a Go-based command-line tool for scheduling and managing social media posts to Mastodon and LinkedIn. It serves as a replacement for Buffer.com, allowing users to queue and schedule posts from the terminal.

## Project Structure

```
gos/
├── cmd/
│   ├── gos/       # Main gos binary
│   └── gosc/      # Gos Composer binary (quick post composition)
├── internal/
│   ├── colour/    # Terminal color utilities
│   ├── config/    # Configuration management
│   ├── entry/     # Post entry handling
│   ├── oi/        # Output/input utilities
│   ├── platforms/ # Social media platform implementations
│   │   ├── linkedin/
│   │   ├── mastodon/
│   │   ├── noop/
│   │   └── platform.go
│   ├── prompt/    # User prompts
│   ├── queue/     # Message queue management
│   ├── schedule/  # Posting schedule logic
│   ├── summary/   # Gemini gemtext summary generation
│   ├── table/     # Table output formatting
│   ├── tags/      # Share tag parsing
│   ├── timestamp/ # Timestamp utilities
│   ├── main.go
│   ├── run.go
│   └── version.go
├── docs/          # Documentation and images
├── examples/      # Example files
└── gosdir/        # Example gos directory structure
```

## Build System

This project uses [Mage](https://magefile.org/) for build automation.

### Commands

| Command | Description |
|---------|-------------|
| `mage` or `mage build` | Build `gos` and `gosc` binaries |
| `mage install` | Build and install binaries to `~/go/bin` |
| `mage clean` | Remove built binaries |
| `mage test` | Run all tests |
| `mage fuzz` | Run fuzz tests (10s) |
| `mage lint` | Run golangci-lint |
| `mage vet` | Run `go vet` |
| `mage dev` | Run tests, vet, lint, then build with race detector |
| `mage devInstall` | Install `gopls` and `golangci-lint` |

### Before Committing

Run `mage dev` to ensure tests pass, vet and lint checks succeed, and the build completes with race detection.

## Testing

- Tests use the standard Go testing framework
- Test files follow the `*_test.go` naming convention
- Run tests: `mage test` or `go test -v ./...`
- Run fuzz tests: `mage fuzz` or `go test ./internal/entry/ -fuzz=FuzzExtractURLs -fuzztime=10s`

## Code Conventions

- **Go version**: 1.23+
- **Module path**: `codeberg.org/snonux/gos`
- **Package layout**: Internal packages under `internal/`, commands under `cmd/`
- **Error handling**: Standard Go error handling patterns
- **Dependencies**: Minimal external dependencies (fatih/color, golang.org/x packages)

### Go Coding Practices

Follow the practices defined in `/home/paul/git/conf/snippets/go/go-projects/go-projects.md`:

## Linting

Uses `golangci-lint`. Install with:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Or run `mage devInstall`.

## Key Dependencies

- `github.com/fatih/color` - Terminal colors
- `golang.org/x/oauth2` - OAuth2 for LinkedIn
- `golang.org/x/net` - HTML parsing for LinkedIn previews
- `github.com/magefile/mage` - Build automation
