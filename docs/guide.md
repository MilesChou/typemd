# TypeMD User Guide

## Core Concepts

### Object

An Object is the basic unit of TypeMD. Each object is a Markdown file with YAML frontmatter (properties) and body content.

Object IDs follow the format `type/filename`, e.g. `book/golang-in-action`.

```markdown
---
title: Go in Action
status: reading
rating: 4.5
---

# Notes

A great book about Go...
```

### Type

Every object belongs to a type. Types define property names, data types, and validation rules via schema files.

TypeMD includes three built-in types:

| Type | Properties |
|------|-----------|
| `book` | title (string), status (enum: to-read/reading/done), rating (number) |
| `person` | name (string), role (string) |
| `note` | title (string), tags (string) |

Custom type schemas go in `.typemd/types/`.

### Relation

Relations are defined as `relation`-type properties in type schemas. They support:

- **Unidirectional / Bidirectional** — bidirectional relations auto-sync both sides
- **Single / Multiple values** — multiple values stored as YAML arrays

## CLI Reference

### Global Options

| Option | Description |
|--------|-------------|
| `--vault <path>` | Vault directory path (default: current directory) |

### `tmd` (no subcommand)

Launches the TUI interactive interface.

### `tmd init`

Initializes a new vault. Creates `.typemd/` directory structure and SQLite database.

### `tmd show <object-id>`

Displays an object's full information: properties (including relations) and body.

```bash
tmd show book/golang-in-action
```

Example output:

```
book/golang-in-action

Properties
──────────
  title: Go in Action
  status: reading
  rating: 4.5
  author: → person/alan-donovan

Body
────
  # Notes
  A great book about Go...
```

Properties are displayed in schema-defined order. Relation properties use `→` for forward and `←` for inverse links.

### `tmd query <filter> [--json]`

Filter objects by properties. Conditions use `key=value` format, separated by spaces (AND logic).

```bash
tmd query "type=book"
tmd query "type=book status=reading"
tmd query "type=book" --json
```

### `tmd search <keyword> [--json]`

Full-text search across filenames, properties, and body content. Powered by SQLite FTS5.

```bash
tmd search "concurrency"
tmd search "golang" --json
```

### `tmd link <from-id> <relation> <to-id>`

Creates a relation between two objects. If the schema defines `bidirectional: true`, the inverse property is automatically updated.

```bash
tmd link book/golang-in-action author person/alan-donovan
```

### `tmd unlink <from-id> <relation> <to-id> [--both]`

Removes a relation. Use `--both` to remove the inverse side as well.

```bash
tmd unlink book/golang-in-action author person/alan-donovan --both
```

### `tmd reindex`

Scans the `objects/` directory, syncs all files to the database, and rebuilds the full-text search index. Use after manually editing files.

### `tmd mcp`

Starts an MCP (Model Context Protocol) server over stdio. AI clients (e.g. Claude Code) can query your vault through this protocol.

```bash
tmd mcp
tmd mcp --vault /path/to/vault
```

Available tools:

| Tool | Description |
|------|-------------|
| `search` | Full-text search objects, returns ID, type, and filename |
| `get_object` | Get full object detail by ID, including properties and body |

## TUI Details

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up (navigate list / scroll detail) |
| `↓` / `j` | Move down (navigate list / scroll detail) |
| `Enter` / `Space` | Select object / Toggle group |
| `Tab` | Switch focus between panels |
| `/` | Enter search mode |
| `Esc` | Exit search / Clear results |
| `q` / `Ctrl+C` | Quit |

### Auto-refresh

The TUI watches the `objects/` directory via fsnotify. When files are created, modified, or deleted, it automatically syncs the database and refreshes the view (200ms debounce), preserving the current selection when possible.

## Type Schema

### Basic Format

Create YAML files in `.typemd/types/`:

```yaml
# .typemd/types/book.yaml
name: book
properties:
  - name: title
    type: string
  - name: status
    type: enum
    values: [to-read, reading, done]
    default: to-read
  - name: rating
    type: number
```

### Property Types

| Type | Description | Example |
|------|-------------|---------|
| `string` | Text | `"Go in Action"` |
| `number` | Integer or float | `42`, `3.14` |
| `enum` | Enumerated value, requires `values` | `"reading"` |
| `relation` | Link to another object | `"person/alan"` |

### Relation Properties

```yaml
# .typemd/types/book.yaml
name: book
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person          # target type
    bidirectional: true     # auto-sync both sides
    inverse: books          # inverse property name

# .typemd/types/person.yaml
name: person
properties:
  - name: name
    type: string
  - name: books
    type: relation
    target: book
    multiple: true          # allows multiple values
    bidirectional: true
    inverse: author
```

| Field | Description |
|-------|-------------|
| `target` | Target object's type name |
| `multiple` | Whether the property holds multiple values (array) |
| `bidirectional` | Auto-sync the inverse side when linking |
| `inverse` | Property name on the target type's schema |

### Validation

TypeMD uses lenient validation:

- Only validates properties defined in the schema
- Extra properties (not in schema) are allowed
- Missing properties do not cause errors
- `enum` values must be in the `values` list
- `number` must be numeric
- `relation` targets are checked for correct type
