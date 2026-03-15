---
title: tmd object create
description: Create a new object from a type schema.
sidebar:
  order: 2.5
---

Creates a new object file (Markdown + YAML frontmatter) based on the type schema.

```bash
tmd object create <type> [name]
tmd object create book clean-code
tmd object create book clean-code -t review
tmd object create person robert-martin
tmd object create journal
```

If the type has a name template, the `name` argument is optional — the name will be auto-generated from the template. When no name is provided and no name template is defined, an error is returned.

The command generates a `.md` file under `objects/<type>/` with all schema-defined properties set to their default values (or `null` if no default is specified). The object is also added to the SQLite index.

Each object is assigned a unique ULID suffix, so multiple objects of the same type can share the same name (e.g. two `book/clean-code-<ulid>` files with different ULIDs). If the type does not exist, an error is returned.

## Flags

| Flag | Description |
|------|-------------|
| `-t`, `--template` | Template name to use (from `templates/<type>/`) |

## Templates

If the type has templates defined in `templates/<type>/`, they are resolved automatically:

- **No templates** — proceeds without a template
- **Single template** — auto-applied without prompting
- **Multiple templates** — interactive prompt for selection (or use `-t` to skip the prompt)

Template frontmatter properties override schema defaults. Template body becomes the initial object body. Auto-managed system properties (`created_at`, `updated_at`) in templates are ignored.

