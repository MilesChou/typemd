# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [v0.2.0] - 2026-03-11

### Added

- Property Type System — define 9 property types (`string`, `text`, `number`, `bool`, `date`, `datetime`, `url`, `enum`, `relation`) in type schemas (#8)
- Shared Properties — define reusable property definitions in `.typemd/properties.yaml` and reference them via `use` in type schemas (#188)
- Emoji on Types — add an `emoji` field to type schemas for visual identification in the TUI (#145)
- Emoji on Properties — add an `emoji` field to property schemas for compact display (#144)
- TUI Title Panel — dedicated header showing the type emoji and object name when viewing an object (#169)
- TUI Pinned Properties — mark properties as `pinned: true` in schema for prominent display in the TUI detail view (#168)
- TUI Session Persistence — cursor position, selected object, and panel state are restored across TUI restarts (#82)
- `--readonly` flag — launch the TUI in read-only mode to disable all editing (#107)
- `--reindex` flag — global flag to force rebuild the SQLite index on startup, replacing the `tmd reindex` subcommand (#159)
- Prefix Matching — resolve objects by a short prefix of their ULID suffix instead of the full ID (#72)
- Homebrew Installation — install via `brew install typemd/tap/tmd` (#140)

### Changed

- `name` Property — now a required system property automatically populated from the object slug; type schemas cannot define a property named `name` (#187)
- TUI Object List — type emoji shown in group headers alongside type name (#163)
- Undefined Properties — properties not declared in the type schema are silently filtered during sync (#174)

### Fixed

- Relation Display — ULID suffix stripped from relation property display values

[v0.2.0]: https://github.com/typemd/typemd/releases/tag/v0.2.0

## [v0.1.0] - 2026-03-08

### Added

- Objects & Types — define typed schemas in YAML, create objects as Markdown files with `tmd object create` (#18)
- ULID filenames — unique suffix for conflict-free object naming (#48)
- Relations — bidirectional links via `tmd relation link` / `tmd relation unlink`, single-value overwrite and multi-value append
- Wiki-links & Backlinks — `[[target]]` syntax in markdown body with automatic backlink tracking (#10)
- Querying — `tmd query` for type/property filtering, `tmd search` for full-text search, both with `--json` output
- Validation — `tmd type validate` checks schema integrity, property types, orphaned relations, and broken wiki-links (#20)
- Migration — `tmd migrate` updates existing objects when schemas evolve (#22)
- Auto-reindex — SQLite index is automatically rebuilt when empty or missing (#41)
- Orphan cleanup — stale relations detected and removed during reindex (#21)
- CLI reorganization — commands grouped by resource type: `tmd object`, `tmd type`, `tmd relation` (#141)
- TUI — three-panel layout (#47), in-place body editing (#85), edit mode with visual indicator (#84), auto-save on exit (#86), help popup (#104)
- TUI display — ULID stripped from display names (#75), reduced indentation (#57), grouped object list (#43)
- MCP Server — `tmd mcp` exposes vault to AI assistants
- `.gitignore` on init — `tmd init` creates `.typemd/.gitignore` to exclude `index.db` (#1)
- `tmd` binary — `go install` produces `tmd` binary (#61)
- Documentation site with English and zh-TW support (#50, #54)
- BDD testing framework with Godog and Gherkin feature files (#111, #112)
- GitHub Actions release workflow for multi-platform binaries (#39)
- Codebase refactoring — unified naming conventions, extracted helpers, improved error handling (#56)
- Vault structure refactoring — remove `objects/` directory layer (#117)

[v0.1.0]: https://github.com/typemd/typemd/releases/tag/v0.1.0
