---
title: Concepts Overview
description: The design philosophy behind TypeMD.
sidebar:
  order: 1
---

TypeMD is a local-first knowledge management tool built around one core idea: **think in Objects, not files**.

## Objects, not files

Most note-taking tools organize content as files in folders. TypeMD takes a different approach — your knowledge is made of **Objects**. A book, a person, an idea, a meeting — each is an Object with a defined structure and typed connections to other Objects.

The file system is just a storage layer. You interact with Objects, not filenames.

## Structure through Types

Every Object belongs to a **Type** that defines its structure. A `book` has an `author` and `status`; a `person` has a `role` and `company`. Types are defined as simple YAML schema files — you design the structure that fits your thinking.

## Connections through Relations and Links

Objects connect to each other in two ways:

- **Relations** — structured, typed links defined in the schema. A book's `author` relation points to a `person` Object. Relations can be bidirectional: updating one side automatically updates the other.
- **Links** — inline references using `[[type/slug-ulid]]` syntax in the Markdown body. Lighter than Relations, useful for ad-hoc connections. TypeMD tracks backlinks automatically.

## Local-first, Markdown-native

Everything is stored as plain Markdown files with YAML frontmatter. No cloud dependency, no lock-in. You can read, edit, and version-control your knowledge base with any tool.

## Learn more

- [Objects](/concepts/objects) — what Objects are and how they work
- [Types](/concepts/types) — defining structure with schemas
- [Relations](/concepts/relations) — structured connections between Objects
- [Links](/concepts/links) — inline references and backlinks
- [Glossary](/concepts/glossary) — definitions of all key terms
