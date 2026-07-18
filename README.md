# Rook

Production-safe SQL notebook desktop app built with Go + Wails v3 + SvelteKit.

## Linux Prerequisites (Wails Runtime)
Install the required native packages before running the app.

- Ubuntu/Debian:
  `sudo apt update && sudo apt install -y build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev libsoup-3.0-dev`
- Fedora:
  `sudo dnf install -y gcc gcc-c++ make pkgconf-pkg-config gtk4-devel webkitgtk6.0-devel libsoup3-devel`
- Arch:
  `sudo pacman -S --needed base-devel pkgconf gtk4 webkitgtk-6.0 libsoup3`

Verify with:
`pkg-config --modversion gtk4 webkitgtk-6.0 libsoup-3.0`

## Run Instructions
1. Install frontend dependencies:
   `npm --prefix frontend install`
2. Run the desktop app:
   `go run ./cmd/rook`

## Build Instructions
1. Build frontend assets:
   `npm --prefix frontend run build`
2. Build the Go desktop binary:
   `go build -o bin/rook ./cmd/rook`

## MVP Scope
- Notebook SQL cells (split by two newlines)
- Expanded SQL safety scenarios across destructive, performance, keyword, and injection-risk patterns
- Configurable safety rules in settings (thresholds, strict PROD mode, read-only mode, FORCE policy)
- Optional safety bypass using SQL comment token `--FORCE` (policy-controlled)
- Transaction-oriented execution service with per-cell cancel + rollback-on-cancel semantics
- Auto commit/rollback timer controls (default 5s)
- Light/dark mode and local settings persistence
- Basic notebook file manager UI scaffold
- `.snb` notebook read/write format

## Implementation vs Future Scope
| Area | Implemented | Future Scope |
|---|---|---|
| SQL editor | Notebook SQL cells with multi-cell execution support | Richer editing UX (linting presets, snippets, schema-aware autocomplete) |
| Safety system | Rule-based validation, strict PROD mode, read-only mode, optional `--FORCE` bypass | Policy packs per environment/team and audit reporting dashboard |
| Execution | Transaction-aware execution with per-cell cancel + rollback-on-cancel | Parallel cell execution orchestration and query queue prioritization |
| Session controls | Auto commit/rollback timers with local settings persistence | Team-shared defaults and profile-based execution templates |
| Notebook files | Local `.snb` read/write with basic file manager scaffold | Versioning, metadata tags, workspace folders, and cloud sync |
