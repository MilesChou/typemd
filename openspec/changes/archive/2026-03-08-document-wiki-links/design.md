## Context

Wiki-links allow users to reference other objects inline in markdown body using `[[type/name-ulid]]` syntax. The feature was implemented in PR #89 and is fully operational. This change documents the existing architecture — no code changes are involved.

Current implementation spans:
- `core/wikilink.go` — Parser, DB queries, rendering, sync logic
- `core/vault.go` — `wikilinks` table schema
- `core/sync.go` — Wiki-link extraction during index sync
- `core/display.go` — Backlinks as built-in display property
- `core/validate.go` — Broken link detection

## Goals / Non-Goals

**Goals:**
- Document the behavioral contract for wiki-links as a formal OpenSpec spec
- Capture all existing behavior: parsing, storage, backlinks, rendering, validation

**Non-Goals:**
- Code changes or refactoring
- New features (e.g., type inference for link targets, auto-complete)

## Decisions

### Full object ID for link targets
Wiki-link targets use full object IDs including ULID suffix (e.g., `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6]]`). This was a design change from the original proposal which considered DisplayID format. Full IDs simplify target resolution to a direct lookup.

### Separate `wikilinks` table
Wiki-links are stored in a dedicated `wikilinks` table, separate from schema-defined `relations`. This keeps the two linking mechanisms independent — relations are structured (defined in type schemas), while wiki-links are freeform (written in markdown body).

### In-memory target resolution during sync
`syncWikiLinks` accepts a `knownIDs` map to resolve targets without N+1 DB queries. Unresolved targets get an empty `to_id`, marking them as broken links.

### Backlinks as built-in display property
Backlinks appear as a system-level `backlinks` property in `DisplayProperty`, appended after schema properties and reverse relations. This is not a user-defined property — it's computed from the `wikilinks` table.

## Risks / Trade-offs

- [Full ID syntax is verbose] → Users must use exact IDs; no fuzzy matching or type inference. This trades usability for implementation simplicity.
- [No cascading delete from DB foreign keys] → Wiki-link cleanup relies on application-level sync logic rather than DB cascades.
