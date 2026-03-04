---
name: create-issue
description: This skill should be used when the user asks to "create an issue", "open an issue", "file a bug", "request a feature", "add a task", or discusses a problem or idea that should be tracked as a GitHub issue.
---

# Create GitHub Issue

Create a GitHub issue for the TypeMD project through iterative Q&A.

Do NOT create the issue until the user explicitly confirms. Ask questions one at a time to refine the idea before writing anything.

## Process

### Step 1: Understand the idea

Ask the user what they want to achieve. One question at a time:

1. **What's the problem or idea?** — Understand the motivation first.
2. **What's the expected behavior?** — Clarify what "done" looks like.
3. **Any constraints or context?** — Related issues, technical considerations, scope.

Use multiple choice (via AskUserQuestion) when possible. Keep it conversational — skip questions that are already answered from context.

### Step 2: Determine labels

Each issue gets **two labels** — one type and one component.

**Type label** (pick one):

| Label | When to use |
|-------|-------------|
| `enhancement` | New feature or improvement |
| `bug` | Something isn't working |
| `documentation` | Docs changes |

**Component label** (pick one):

| Label | Scope |
|-------|-------|
| `core` | Core library — objects, types, relations, index (`core/`) |
| `cli` | CLI commands (`cmd/`) |
| `tui` | Terminal UI (`tui/`) |
| `mcp` | MCP server (`mcp/`) |
| `web` | Web UI (`web/`) |

Suggest labels based on context. Ask for confirmation if ambiguous.

### Step 3: Assign milestone

Fetch open milestones via `gh api repos/typemd/typemd/milestones` and present them as options using AskUserQuestion. Always include a "None" option for issues that don't belong to any milestone.

### Step 4: Draft and confirm

Present the full issue draft to the user before creating:

- **Title** — concise, plain language, no prefix
- **Labels** — type + component
- **Milestone** — selected milestone or none
- **Body** — using the template below

Ask: "This is the issue I'll create. Anything to adjust?"

Only proceed after the user confirms.

### Step 5: Create issue

```bash
gh issue create --title "<title>" --label "<type>,<component>" --milestone "<milestone>" --body "<body>"
```

Omit `--milestone` if the user selected "None".

Body template:

```markdown
## Summary

<1-3 sentences describing the issue>

## Current Behavior

<what happens now, if applicable>

## Expected Behavior

<what should happen>
```

### Step 6: Confirm

Return the issue URL to the user.
