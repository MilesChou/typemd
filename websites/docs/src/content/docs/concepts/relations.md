---
title: Relations
description: How Objects connect to each other through typed links.
sidebar:
  order: 4
---

A Relation is a structured link between two Objects. It's how you express that a book has an author, a project has members, or an idea connects to a meeting.

## What is a Relation?

A Relation is a property of type `relation` defined in a Type schema. Unlike a plain hyperlink or wiki-link, a Relation is:

- **Named** — it has a property name like `author` or `books`, not just a URL
- **Typed** — it knows the target Type (e.g. `author` points to a `person`)
- **Optionally bidirectional** — linking a book to an author can automatically link the author back to the book

## Unidirectional vs. bidirectional

A **unidirectional** Relation is a one-way link. Setting `author` on a book points to a person, but the person Object doesn't know about it.

A **bidirectional** Relation keeps both sides in sync. When you link a book's `author` to a person, the person's `books` property is automatically updated — and vice versa.

```yaml
# book.yaml
- name: author
  type: relation
  target: person
  bidirectional: true
  inverse: books

# person.yaml
- name: books
  type: relation
  target: book
  multiple: true
  bidirectional: true
  inverse: author
```

## Single vs. multiple

- **Single-value** — A book has one `author`. Re-linking overwrites the previous value.
- **Multiple-value** — A person can have many `books`. Linking appends to the list.

This is controlled by the `multiple` field in the schema.

## Relations vs. wiki-links

TypeMD supports both Relations and wiki-links (`[[type/slug]]`). They serve different purposes:

| | Relation | Wiki-link |
|---|----------|-----------|
| Defined in | Type schema (frontmatter) | Markdown body |
| Structured | Yes — named, typed, queryable | No — freeform inline reference |
| Bidirectional | Configurable per schema | Backlinks tracked automatically |
| Use case | Formal connections (author, project members) | Informal references (see also, mentioned in) |

Use Relations for connections that are part of your data model. Use wiki-links for casual references in your notes.
