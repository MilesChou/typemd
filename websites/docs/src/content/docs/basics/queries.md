---
title: Queries
description: Structured filtering of Objects by property values.
sidebar:
  order: 5
---

Queries let you filter Objects by their property values using structured conditions. Unlike [search](/basics/search), which matches free-text across all content, queries target specific properties with exact values.

## How it works

Queries use `key=value` format to match Objects. When multiple conditions are provided, they are combined with AND logic — only Objects matching all conditions are returned.

The `type` key is a special filter that matches the Object's type name. All other keys match against frontmatter property values.

## CLI usage

```bash
# Filter by type
tmd query "type=book"

# Filter by type and property
tmd query "type=book status=reading"

# Output as JSON
tmd query "type=book" --json
```

| Option | Description |
|--------|-------------|
| `--json` | Output results in JSON format |

## Examples

| Query | Matches |
|-------|---------|
| `type=book` | All book Objects |
| `type=book status=reading` | Books with status "reading" |
| `type=person role=author` | People with role "author" |

## See also

- [Search](/basics/search) — full-text search across all content
- [tmd query](/tui/query) — CLI reference
- [Data Model](/developers/data-model#querying-architecture) — querying details
