# Rook

Production-safe SQL notebook desktop app built with Go + Wails v3 + SvelteKit.

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
