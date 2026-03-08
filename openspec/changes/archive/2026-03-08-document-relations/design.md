## Context

Relations provide structured linking between objects. Unlike wiki-links (freeform inline references), relations are defined in type schemas as `type: relation` properties with explicit target types, cardinality, and optional bidirectional behavior. Data is persisted in both object frontmatter (YAML) and the SQLite `relations` table.

Implementation spans:
- `core/relation.go` — Link/unlink logic, relation queries, helper functions
- `core/type_schema.go` — `Property` struct with relation fields (target, multiple, bidirectional, inverse)
- `core/display.go` — Reverse relation display in object detail views
- `core/vault.go` — `relations` table schema
- `core/sync.go` — Relation sync during indexing

## Goals / Non-Goals

**Goals:**
- Document the behavioral contract for object relations as a formal OpenSpec spec
- Capture all existing behavior: schema definition, linking, unlinking, bidirectional, validation, display

**Non-Goals:**
- Code changes or refactoring
- New relation features (e.g., cascading delete, relation metadata)

## Decisions

### Dual storage: frontmatter + SQLite
Relations are stored in both the object's YAML frontmatter (as property values) and the `relations` DB table. Frontmatter is the source of truth for file portability; the DB provides fast querying and reverse lookups.

### Single-value overwrites, multiple-value appends
A single-value relation (e.g., `author`) overwrites on re-link. A multiple-value relation (e.g., `books`) appends and rejects duplicates with `errDuplicateRelation`.

### Bidirectional via explicit inverse property
Bidirectional relations require both sides to declare their inverse. When linking A→B via a bidirectional property, the system automatically writes the inverse B→A. This keeps both files consistent.

### Type target validation at link time
`LinkObjects` validates that the target object's type matches the relation's `target` field before persisting. This prevents invalid cross-type links.

### Reverse relations in display
`BuildDisplayProperties` appends reverse relations (where the object is the `to_id`) after schema-defined properties. These are shown with `←` arrows in the UI.

## Risks / Trade-offs

- [Dual storage can drift] → Sync/reindex rebuilds the DB from frontmatter to reconcile.
- [Unlink without --both leaves inverse side stale] → By design: unlink is one-directional unless `both` flag is set. This gives users fine-grained control.
