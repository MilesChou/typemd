---
title: Wiki-links
description: Inline references between Objects using [[target]] syntax.
sidebar:
  order: 5
---

A Wiki-link is an inline reference to another Object, written directly in your Markdown body. It's the simplest way to connect ideas as you write.

## What is a Wiki-link?

When you mention another Object in your notes, you can wrap its ID in double brackets:

```markdown
---
title: Go in Action
---

# Notes

Great introduction to Go concurrency patterns.
See also [[note/go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc]].
```

You can also add display text:

```markdown
See [[note/go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc|Go Concurrency Patterns]] for details.
```

## Backlinks

Wiki-links are tracked automatically. If Object A contains `[[B]]`, then B knows that A references it — this is called a **backlink**.

Backlinks appear as a built-in property when viewing an Object:

```
backlinks: ⟵ A
```

You don't need to maintain backlinks manually. TypeMD computes them from the wiki-link syntax in your Markdown body.

## Wiki-links vs. Relations

TypeMD has two ways to connect Objects. They serve different purposes:

| | Wiki-link | Relation |
|---|-----------|----------|
| Defined in | Markdown body | Type schema (frontmatter) |
| Structured | No — freeform inline reference | Yes — named, typed, queryable |
| Bidirectional | Via backlinks (read-only, automatic) | Configurable per schema |
| Schema required | No | Yes (`type: relation`) |
| Use case | Informal references (see also, mentioned in) | Formal connections (author, project members) |

**Rule of thumb**: use Relations for connections that are part of your data model. Use wiki-links for casual references in your notes.
