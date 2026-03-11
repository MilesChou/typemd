---
name: release-note
description: Use when preparing a release, writing release notes, creating a GitHub draft release, or updating CHANGELOG.md. Triggers on "write release notes", "prepare release", "create release", "update changelog".
---

# Release Note

Write release notes, update CHANGELOG, write a blog post, and create a GitHub draft release for a given version.

## Input

The user must provide a version number (e.g. `v0.2.0`). If not provided, ask before proceeding.

## Process

1. **Gather** — Collect closed issues and relevant commits for the version
2. **Write** — Draft release notes and CHANGELOG entry
3. **Blog** — Write a blog post in zh-tw (English will be synced via `sync-blog`)
4. **Finish** — Commit, push, and create GitHub draft release

## Language

- Release notes and CHANGELOG: **English**
- Blog post: **Traditional Chinese (zh-tw)**

---

## 1. Gather

```bash
# Get the previous release tag (for commit range)
PREV_TAG=$(git tag --sort=-version:refname | head -2 | tail -1)

# Closed issues in the milestone
gh issue list --milestone "v<VERSION>" --state closed --json number,title,labels --limit 50

# Commits since last tag — exclude chore, skill, and docs-only commits
git log ${PREV_TAG}..HEAD --oneline --no-merges
```

**Commit filtering rules** — exclude commits whose type prefix is:
- `chore` — maintenance, deps, tooling
- `docs` — documentation-only changes
- `skill` or `chore(skill)` — internal skill updates

Include: `feat`, `fix`, `refactor` (if user-visible), `perf`

If the milestone has open issues, warn the user before proceeding.

---

## 2. Write

### Release notes format

```markdown
<1-2 paragraphs of prose — what's notable, what problem this solves, what users get>

### Installation

```bash
go install github.com/typemd/typemd/cmd/tmd@v<VERSION>
```

Pre-built binaries for macOS, Linux, and Windows are available below.

### Added

- Feature Name — one-sentence description (#issue)

### Changed

- ...

### Fixed

- ...

### Documentation

- Docs: https://docs.typemd.io
- Blog: https://blog.typemd.io

### Thanks

<invite community participation, link to CONTRIBUTING.md>

**Full Changelog:** https://github.com/typemd/typemd/compare/v<PREV>...v<VERSION>
```

### Writing guidelines

- **Highlight lead** — Open with prose, not a list. What should users be excited about?
- **Group by theme, not package** — Users care about capabilities (Objects, TUI, CLI), not internal packages (core, cmd)
- **Bullet format** — `Feature Name — one-sentence description (#issue)`
- **Issue numbers required** — Every entry must reference the corresponding issue
- **Pre-release** — For `v0.x`, mark as pre-release

### CHANGELOG

Update both `CHANGELOG.md` (English) and `CHANGELOG.zh-TW.md` (Chinese) at the project root.

```markdown
## [v<VERSION>] - YYYY-MM-DD

### Added

- Feature Name — description (#issue)

### Changed / Fixed / Removed

- ...

[v<VERSION>]: https://github.com/typemd/typemd/releases/tag/v<VERSION>
```

---

## 3. Blog

Write a blog post at:

```
websites/blog/src/content/posts/zh-tw/release-<VERSION-WITHOUT-V>.md
```

e.g. `release-0.2.0.md`

**Frontmatter:**

```yaml
---
title: "TypeMD <VERSION> — <one-line theme headline>"
description: "<1-2 sentence teaser>"
date: <today's date>
tags: [發布, 開發日誌]
---
```

**Content guidelines:**

- Write in Traditional Chinese (zh-tw)
- Open with 1-2 sentences of context — why this release matters
- Explain the key features with concrete examples (code snippets, commands)
- Keep it engaging, not just a feature list — tell a story around the release theme
- End with a forward-looking sentence about what's next

After writing zh-tw, use the `sync-blog` skill to create the English version.

---

## 4. Finish

Once release notes, CHANGELOG, and blog post are ready:

### Commit and push

```bash
git add CHANGELOG.md CHANGELOG.zh-TW.md \
  websites/blog/src/content/posts/zh-tw/release-<VERSION>.md \
  websites/blog/src/content/posts/en/release-<VERSION>.md

git commit -m "chore(release): prepare v<VERSION> release notes and blog post"
git push origin main
```

### Create GitHub draft release

```bash
gh release create v<VERSION> \
  --draft \
  --prerelease \
  --title "v<VERSION>" \
  --notes "<release notes content>"
```

- Always create as **draft** — the user decides when to publish
- Mark as **prerelease** for `v0.x` versions
- If a draft already exists, update it:

```bash
gh api repos/typemd/typemd/releases/<ID> -X PATCH -f body="<notes>"
```

### Done

Present:
- The draft release URL
- The blog post file path
- Remind the user to review before publishing
