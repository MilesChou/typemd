---
name: sync-blog
description: Use when blog posts have been added or updated in zh-tw, when user asks to sync or translate blog content to English, or before releases to ensure both language versions are in sync.
---

# Sync Blog

Sync blog posts from Chinese (zh-TW, source of truth) to English.

## Paths

- Source: `websites/blog/src/content/posts/zh-tw/`
- Target: `websites/blog/src/content/posts/en/`

## Process

1. **Diff** — Compare file lists between `zh-tw/` and `en/`. Identify:
   - New posts in zh-tw with no English counterpart
   - Posts where zh-tw has been modified more recently than en

2. **Translate** — For each post needing sync:
   - Read the zh-tw version
   - If English version exists, read it too and identify what changed
   - Translate content to natural English (not literal translation)
   - Preserve: filename, frontmatter structure, code blocks, markdown formatting, links
   - Translate: `title`, `description`, `tags`, and body text

3. **Review** — Present a summary of changes to the user before writing files

## Translation Guidelines

- Write natural English, not word-for-word translation
- Keep the author's voice and tone — informal, first-person
- Technical terms stay as-is (TypeMD, CLI, ULID, frontmatter, etc.)
- Proper nouns stay as-is
- Tags should be translated to English equivalents (e.g. 開發日誌 → devlog, 理念 → philosophy)
- Frontmatter `date` must stay the same
