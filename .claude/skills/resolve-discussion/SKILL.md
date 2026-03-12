---
name: resolve-discussion
description: Use when resolving a GitHub issue labeled `discussion` — facilitates decision-making on open questions, documents conclusions as an issue comment, creates follow-up issues if needed, and closes the discussion. Triggered by resolve-issue when it detects a discussion-labeled issue, or directly when user asks to "resolve discussion #N", "close discussion", "wrap up discussion".
allowed-tools:
  - Bash(gh issue view:*)
  - Bash(gh issue comment:*)
  - Bash(gh issue close:*)
---

# Resolve Discussion

Facilitate decision-making for discussion issues, document conclusions, and close.

Discussion issues don't produce code — they produce **decisions**. The goal is to reach conclusions on open questions, document them, and create actionable follow-up issues if needed.

## Prerequisites

Always read the issue first to ensure context is available, regardless of how this skill was invoked:

```bash
gh issue view <number> --json title,body,labels,milestone,assignees
```

## Phase 1: Brainstorm

Use the `superpowers:brainstorming` skill to explore the discussion topic with the user:

- Ground the brainstorm in the issue's open questions and context
- Research relevant codebase areas, prior art, and constraints
- Check status of any sub-issues referenced in the issue body
- Explore trade-offs and surface hidden considerations

The brainstorming skill will interactively work through ideas with the user until a clear direction emerges.

## Phase 2: Facilitate Decisions

Once brainstorming converges, formalize decisions on each open question:

1. **List open questions** extracted from the issue body
2. **For each question**, present:
   - Summary of what was explored during brainstorming
   - Options with trade-offs (use AskUserQuestion)
   - Recommendation if evidence supports one
3. **Capture decisions** — record the user's choice for each question

## Phase 3: Document and Close

Once all questions are resolved:

### 1. Post summary comment

```bash
gh issue comment <number> --body "$(cat <<'EOF'
## Discussion Summary

### Decisions

- <decision 1>
- <decision 2>

### Follow-up

- <follow-up issue description, if any>

---
Resolved via discussion.
EOF
)"
```

### 2. Create follow-up issues

If the discussion concluded with a large feature or multiple work streams, use the `break-down-epic` skill to decompose it into actionable sub-issues.

For simpler follow-ups (individual tasks or bugs), use the `create-issue` skill for each.

### 3. Close the issue

```bash
gh issue close <number>
```

Present the summary to the user. Done.
