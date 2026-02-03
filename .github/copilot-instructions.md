# Copilot instructions for Kbase

## Big picture
- This repo currently contains only a PRD; no Go source files yet. Treat [PRD.md](PRD.md) as the authoritative spec for architecture, UX, and data flow.
- Kbase is a local-first TUI command knowledge base: load YAML on startup, keep in-memory list, filter via fuzzy search, and render list/detail views using Bubble Tea + Lip Gloss.
- Data flow: read ~/.config/kbase/commands.yaml (or user-specified path) → parse YAML → filter by platform (runtime.GOOS) → fuzzy-search across cmd/desc/tags → render list/detail → copy cmd to clipboard on Ctrl+c → allow edit via $EDITOR and reload.

## Data model (from PRD)
- YAML entries: cmd (required), desc (required), tags (optional), platform (optional), examples (optional).
- Platform filtering: if platform list is non-empty and does not include current OS, hide the entry.

## UX/keybinding conventions
- Ctrl+c copies the selected cmd to clipboard (override default SIGINT exit).
- Esc/Ctrl+q should quit.
- e opens $EDITOR on commands.yaml; suspend TUI, resume on editor exit, reload YAML.
- List view: left cmd, right desc (dimmed); truncate long desc in list, show full in detail pane.

## Error handling expectations
- YAML parse errors must not crash; show a red ErrorView “Parsing Failed” and allow e to fix.
- Clipboard failures in headless/SSH environments should show “Clipboard not supported” and optionally print cmd to stdout on exit.

## Dependencies/integration points
- TUI: charmbracelet/bubbletea, styling: charmbracelet/lipgloss.
- YAML: gopkg.in/yaml.v3.
- Clipboard: github.com/atotto/clipboard or golang.design/x/clipboard.
- Editor integration via $EDITOR env var.

## Repo/workflow notes
- Go module: [go.mod](go.mod) declares module github.com/BillabongJumbuck/Kbase and Go 1.25.5.
- No build/test commands are documented yet; add them here once tooling or CI is introduced.
