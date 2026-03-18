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

## Sorting

Queries support sorting by property values. When used via the TUI [view mode](/tui/tui#view-mode--table-display), sort rules are defined in the view configuration:

```yaml
sort:
  - property: rating
    direction: desc
```

Sort direction is either `asc` (ascending) or `desc` (descending). Multiple sort rules are applied in order (first rule is the primary sort).

## Filter operators

Views support type-aware filter operators beyond simple `key=value` equality. Each property type has a defined set of valid operators:

| Property type | Operators |
|---------------|-----------|
| `string`, `url` | `is`, `is_not`, `contains`, `does_not_contain`, `starts_with`, `ends_with`, `is_empty`, `is_not_empty` |
| `number` | `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `is_empty`, `is_not_empty` |
| `date`, `datetime` | `eq`, `before`, `after`, `on_or_before`, `on_or_after`, `is_empty`, `is_not_empty` |
| `select` | `is`, `is_not`, `is_empty`, `is_not_empty` |
| `multi_select`, `relation` | `contains`, `does_not_contain`, `is_empty`, `is_not_empty` |
| `checkbox` | `is`, `is_not` |

Filter rules are defined in view configurations:

```yaml
filter:
  - property: status
    operator: is
    value: reading
  - property: rating
    operator: gt
    value: "3"
```

Multiple filters are combined with AND logic.

## See also

- [Search](/basics/search) — full-text search across all content
- [Views](/advanced/file-structure#views) — view configuration format
- [tmd query](/tui/query) — CLI reference
- [Data Model](/developers/data-model#querying-architecture) — querying details
