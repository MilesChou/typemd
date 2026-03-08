# CLAUDE.md

## Project Overview

typemd is a local-first CLI knowledge management tool. Objects (books, people, ideas) are stored as Markdown files with YAML frontmatter, connected by Relations. SQLite provides indexing.

## Architecture

- **core/** — Core library: objects, types, relations, index
- **cmd/** — CLI commands (Cobra)
- **tui/** — Terminal UI (Bubble Tea)
- **mcp/** — MCP server
- **web/** — Web UI (future)
- **app/** — Desktop app via Wails (future)
- **websites/** — Non-Go websites (site, docs, blog)
- **marketplace/** — Claude Marketplace plugins (future)

## Data Model

- Objects identified by `type/<slug>-<ulid>` (e.g. `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`)
- Type schemas: `.typemd/types/*.yaml`
- Relations defined as properties in type schemas
- Wiki-links: `[[type/name-ulid]]` syntax in markdown body, with backlink tracking
- SQLite index: `.typemd/index.db`
- Object files: `objects/<type>/<name>.md`

## Build & Test

```bash
go build ./...
go test ./...
go run ./cmd/tmd [command]
```

## Testing

This project uses two layers of testing:

- **BDD (Godog)** — Define behaviors, establish shared vocabulary, and describe what a feature does from the user's perspective. Gherkin `.feature` files live in `<package>/features/`. BDD scenarios focus on **what**, not implementation details.
- **Unit tests** — Verify precise logic: edge cases, output formats, exact values, error conditions. Traditional Go `testing` style.

When deciding where a test belongs: if it defines a behavior or names a concept, write a BDD scenario. If it validates an implementation detail (e.g. JSON format, lowercase ULID, flag edge cases), write a unit test.

### BDD scope by package

| Package | Testing approach |
|---------|-----------------|
| `core/` | BDD (`core/features/`) + unit tests |
| `tui/`  | BDD (`tui/features/`, planned) + unit tests |
| `web/`  | BDD (`web/features/`, future) |
| `cmd/`  | Minimal — CLI commands delegate to `core/`, covered by core BDD scenarios |
| `mcp/`  | Unit tests — BDD TBD |
