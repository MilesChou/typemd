---
title: Shared Properties
description: Reusable property definitions referenced by multiple type schemas.
sidebar:
  order: 3.7
---

Shared Properties let you define a property once and reuse it across multiple type schemas. Instead of duplicating the same property definition in every type, you define it in a central file and reference it with the `use` keyword.

## Why Shared Properties?

Imagine you have a `due_date` property that appears in `project`, `task`, and `milestone` types. Without Shared Properties, you'd copy the same definition three times — and any change (like adding an emoji) would need to be updated in all three places.

With Shared Properties, you define `due_date` once and reference it everywhere:

```yaml
# .typemd/properties.yaml
properties:
  - name: due_date
    type: date
    emoji: "\U0001F4C5"
  - name: priority
    type: select
    options:
      - value: high
      - value: medium
      - value: low
```

## File format

Shared Properties are defined in `.typemd/properties.yaml` at the root of your vault's `.typemd/` directory. The file uses a top-level `properties` array:

```yaml
properties:
  - name: due_date
    type: date
    emoji: "\U0001F4C5"
  - name: priority
    type: select
    options:
      - value: high
      - value: medium
      - value: low
  - name: assignee
    type: relation
    target: person
    multiple: true
```

Each entry supports all the same fields as a regular type schema property: `name`, `type`, `emoji`, `pin`, `options`, `target`, `default`, `multiple`, `bidirectional`, and `inverse`.

The file is optional. If `.typemd/properties.yaml` does not exist or contains no `properties` array, TypeMD treats it as an empty set of shared properties.

## Referencing shared properties

In a type schema, use the `use` keyword to reference a shared property by name:

```yaml
# .typemd/types/project.yaml
name: project
emoji: "\U0001F4CB"
properties:
  - name: title
    type: string
  - use: due_date
  - use: priority
  - name: budget
    type: number
```

When the type is loaded, each `use` entry is resolved to the full property definition from `.typemd/properties.yaml`. The resolved type schema behaves exactly as if the properties were defined inline.

You can mix `use` entries with regular property definitions. The order in the `properties` array is preserved.

## Overriding fields

When referencing a shared property via `use`, you can override two fields:

| Field | Description |
|-------|-------------|
| `pin` | Override the pin position for this type |
| `emoji` | Override the display emoji for this type |

All other fields (`type`, `options`, `target`, `default`, `multiple`, `bidirectional`, `inverse`, `template`) cannot be overridden — they come from the shared definition.

### Override example

```yaml
# .typemd/types/task.yaml
name: task
properties:
  - use: due_date
    pin: 1
    emoji: "\U0001F5D3\uFE0F"
  - use: priority
```

In this example, the `task` type uses the shared `due_date` property but pins it to position 1 and changes its emoji. The `priority` property is used as-is from the shared definition.

## Validation rules

Shared properties are validated using the same rules as type schema properties, plus additional constraints:

- **Unique names** — No two shared properties can have the same name
- **No reserved names** — The name `name` cannot be used (it is a reserved [system property](/concepts/data-model#system-properties))
- **Valid types** — Property type must be a valid [predefined property type](/concepts/predefined-property)
- **Options required** — `select` and `multi_select` properties must include an `options` array
- **Target required** — `relation` properties must include a `target` field

When a type schema uses `use`:

- **`use` and `name` are mutually exclusive** — A property entry cannot have both `use` and `name`
- **Reference must exist** — The referenced shared property must exist in `.typemd/properties.yaml`
- **No name conflicts** — A local property cannot have the same name as a referenced shared property
- **No duplicate references** — The same shared property cannot be referenced twice in one type schema
- **Only `pin` and `emoji` overrides** — Setting any other field on a `use` entry is an error

Run `tmd type validate` to check all shared properties and type schema references.

## Shared Properties vs. inline properties

| | Shared Property | Inline Property |
|---|----------------|-----------------|
| Defined in | `.typemd/properties.yaml` | `.typemd/types/<type>.yaml` |
| Reusable | Yes — referenced via `use` | No — scoped to one type |
| Customizable per type | `pin` and `emoji` only | Fully customizable |
| Use case | Properties shared across multiple types | Properties unique to one type |

**Rule of thumb**: if a property appears in two or more types with the same definition, make it a shared property.
