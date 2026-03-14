---
name: create-blog
description: Use when the user wants to write a blog post — technical articles, devlogs, opinion pieces, or any non-release blog content. Triggers on "write a blog post", "write an article", "create blog post". Do NOT use for release notes — those are handled by the `create-release` skill.
---

# Create Blog

Write a blog post for the TypeMD blog. Handles topic exploration, source verification, zh-tw drafting, user review, English sync, and file creation.

## Input

The user provides a topic or direction. If unclear, ask for clarification before proceeding.

## Process

1. **Research** — Understand the topic by reading source code
2. **Draft** — Write the zh-tw version and present for review
3. **Sync** — Create the English version via `sync-blog` skill
4. **Finish** — Stage files (do NOT commit — leave that to the user)

---

## 1. Research

Before writing any technical content, read the actual source code to ensure accuracy.

```bash
# Find relevant source files
# Read implementation, interfaces, key types
# Check BDD feature files for behaviors and examples
# Check example vault for correct YAML/config syntax
```

**Rules:**
- Never guess code structure, API signatures, or YAML syntax from memory — always verify against real code
- Read BDD `.feature` files for concrete scenarios and examples
- Read `examples/book-vault/` for correct schema formats
- If the topic involves architecture, read the key interfaces and their implementations

---

## 2. Draft

### File naming

```
websites/blog/src/content/posts/zh-tw/<slug>.md
```

The slug should be kebab-case, descriptive, and SEO-friendly. Use dashes, not dots (e.g., `cqrs-architecture-of-typemd-core.md`, not `cqrs.architecture.md`).

Discuss the slug with the user before writing if there's any ambiguity.

### Frontmatter

```yaml
---
title: "<concise, descriptive title>"
description: "<1-2 sentence teaser for SEO and social sharing>"
date: <today's date>
tags: [<relevant tags>]
---
```

Common tags (zh-tw → en):
- 架構 → architecture
- 開發日誌 → devlog
- 理念 → philosophy
- 教學 → tutorial
- 發布 → release

### Writing guidelines

- Write in **Traditional Chinese (zh-tw)**
- Follow the [Capacities](https://capacities.io/whats-new) writing style: conversational, problem-solution framing, user-benefit focused, short punchy section titles
- Use **we**, not I
- Technical terms stay as-is (CQRS, DDD, Go, SQLite, frontmatter, etc.)
- Code examples must be **verified against actual source code** — never fabricate
- Keep it engaging — tell a story, not just list facts
- Open with a hook or question that draws the reader in
- End with a takeaway or forward-looking thought

### Present draft for review

**Always present the full draft to the user in the chat before writing the file.** Only write the file after the user confirms or provides feedback.

---

## 3. Sync

After the user approves the zh-tw post, use the `sync-blog` skill to create the English version.

---

## 4. Finish

Stage the new files but do NOT commit or push. The user decides when and how to commit (possibly bundled with other changes).

```bash
git add websites/blog/src/content/posts/zh-tw/<slug>.md \
       websites/blog/src/content/posts/en/<slug>.md
```

Present:
- The file paths created
- Remind the user the files are staged but not committed
