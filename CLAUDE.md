# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gos is a command-line social media scheduling tool written in Go that replaces Buffer.com. It allows users to schedule and manage posts across multiple social media platforms (Mastodon, LinkedIn, and a "Noop" pseudo-platform for tracking).

### Key Architecture Components

- **Entry System**: Text files in `~/.gosdir/` represent posts, with filename tags controlling platform targeting and scheduling behavior
- **Platform Abstraction**: `internal/platforms/platform.go` defines the common interface, with platform-specific implementations in subdirectories
- **Queue Management**: Posts move through lifecycle stages: `.txt` → `.queued` → `.posted` in `~/.gosdir/db/platforms/PLATFORM/`
- **Tag System**: Both filename tags (e.g., `post.share:mastodon.prio.txt`) and inline tags control post behavior
- **OAuth2 Flow**: LinkedIn uses OAuth2 authentication stored in `~/.config/gos/gos.json`

### Core Workflow

1. Users create `.txt` files in `gosDir` (default `~/.gosdir/`)
2. `queue.Run()` processes files and moves them to platform-specific queues
3. `schedule.Run()` selects posts based on targets, priorities, and timing rules
4. Platform implementations handle actual posting
5. Posted files are marked with `.posted` extension and timestamp

## Development Commands

### Build and Test
```bash
# Build both binaries
go-task build
# or manually:
go build -o gos cmd/gos/main.go
go build -o gosc cmd/gosc/main.go

# Run tests
go-task test
# or manually:
go test -v ./...

# Development build with race detection
go-task dev
```

### Code Quality
```bash
# Lint code
go-task lint
# or manually:
golangci-lint run

# Vet code
go-task vet
# or manually:
go vet ./...

# Run fuzzing (for specific packages)
go-task fuzz
```

### Installation
```bash
# Install to ~/go/bin/
go-task install
```

## Code Structure Notes

- **Main entry points**: `cmd/gos/main.go` (main app) and `cmd/gosc/main.go` (composer)
- **Configuration**: `internal/config/` handles CLI args and JSON config file management
- **Platform plugins**: Each platform in `internal/platforms/` implements the common `Post()` interface
- **File processing**: `internal/entry/` handles parsing text files and extracting tags
- **Scheduling logic**: `internal/schedule/` manages timing, targets, and post selection
- **Tag parsing**: `internal/tags/` handles both filename and inline tag extraction

## Platform Integration

When adding new platforms:
1. Create new directory under `internal/platforms/`
2. Implement the `Post(ctx, args, sizeLimit, entry)` interface
3. Add platform alias to `platforms.go` aliases map
4. Handle authentication/configuration in platform-specific code

## Configuration Management

- Config file: `~/.config/gos/gos.json` contains API keys and OAuth tokens
- Database: `~/.gosdir/db/platforms/` contains queued and posted files
- Cache: `~/.gosdir/cache/` (configurable via `--cacheDir`)