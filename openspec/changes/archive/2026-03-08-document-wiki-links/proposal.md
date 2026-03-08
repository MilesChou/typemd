## Why

Wiki-links are a core content linking mechanism in typemd, allowing users to reference other objects inline using `[[type/name-ulid]]` syntax in markdown body, with automatic backlink tracking. This feature is fully implemented (PR #89), but lacks a formal OpenSpec specification documenting its behavioral contract. Adding a spec ensures a clear baseline for future modifications.

## What Changes

- Establish formal specification for the existing wiki-links and backlinks feature
- No code changes — documentation only

## Capabilities

### New Capabilities

- `wiki-links`: Wiki-link syntax parsing, storage, backlink tracking, broken link detection, and rendering

### Modified Capabilities

(none)

## Impact

- No code impact
- Creates `openspec/specs/wiki-links/spec.md` as a behavioral contract
