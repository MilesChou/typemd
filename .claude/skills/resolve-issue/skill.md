---
name: resolve-issue
description: This skill should be used when the user asks to "resolve an issue", "work on issue #N", "fix #N", "implement #N", "close #N", "tackle #N", "pick up #N", "start working on #N", or references a specific GitHub issue number they want to work on.
---

# Resolve Issue

Orchestrate the full lifecycle of resolving a GitHub issue — from reading the issue to opening a PR — with user confirmation at each phase.

All progress is tracked via **issue comments**, enabling resume from any interruption point.

## Language

All issue comments MUST be written in **English**.

## Resume Detection

Before starting, check if work has already begun:

```bash
gh issue view <number> --json comments --jq '.comments[].body'
```

Look for comments matching the pattern `## 🔄 Phase N:`. Find the latest one and check its status:

- `✅ Completed` → resume from the **next** phase
- `⏸️ Paused` or `❌ Blocked` → resume from **that** phase

If prior progress is detected, present a summary and ask the user via AskUserQuestion:

- **"Continue from Phase N"** — resume from the detected point
- **"Start over"** — discard previous progress and begin from Phase 1

If no progress comments exist, start from Phase 1.

## Comment Format

Every phase writes a comment to the issue upon completion. Use this format:

```markdown
## 🔄 Phase N: <Phase Name>

**Status:** ✅ Completed

<phase-specific content>

---
_Updated by Claude Code at YYYY-MM-DD HH:MM (UTC)_
```

If pausing mid-phase, write the comment with `⏸️ Paused` and include what was completed so far.

```bash
gh issue comment <number> --body "<comment>"
```

## Preflight Checks

Before entering Phase 1, verify the issue is actionable:

```bash
gh issue view <number> --json state,linkedBranches
```

- If the issue is **closed**, inform the user and stop.
- If there is already an **open PR** linked to this issue, inform the user and ask whether to continue or stop.
- If a branch matching `fix/issue-<N>-*` or `feat/issue-<N>-*` already exists, inform the user and ask whether to reuse it or create a new one.

## Phases

### Phase 1: Understand the Issue

Read the issue and confirm understanding with the user.

```bash
gh issue view <number> --json title,body,issueType,labels,milestone,assignees
```

Present a summary:

- **Title**
- **Type** (Bug / Feature / Task / Epic)
- **Labels**
- **Milestone**
- **Key requirements** extracted from the body

Ask the user via AskUserQuestion:

- **"Looks correct"** — proceed
- **"I have additional context"** — let the user add info before proceeding

**Comment content:** Requirements summary and any additional context from the user.

### Phase 2: Branch Strategy

Ask the user how to set up the working environment via AskUserQuestion:

- **"Worktree (isolated)"** — invoke `superpowers:using-git-worktrees` skill
- **"Branch in current repo"** — create a branch directly

Branch naming convention:

- Bug → `fix/issue-<N>-<slug>`
- Feature / Task / Epic → `feat/issue-<N>-<slug>`

Where `<slug>` is a short kebab-case summary derived from the issue title (max 5 words).

```bash
git checkout -b <branch-name>
```

**Comment content:** Branch name, strategy chosen (worktree or branch).

### Phase 3: Design

Invoke `superpowers:brainstorming` skill to explore the design space.

The brainstorming skill will:

1. Explore project context
2. Ask clarifying questions
3. Propose 2-3 approaches
4. Present design for user approval

**IMPORTANT:** When brainstorming invokes `superpowers:writing-plans`, the plan output should NOT be saved to `docs/plans/`. Instead, capture the full design and implementation plan to write into the issue comment.

**Comment content:** The complete design — architecture decisions, approach chosen, implementation plan with steps.

### Phase 4: Implement

Execute the implementation plan from Phase 3. Choose the appropriate approach:

- **TDD** (`superpowers:test-driven-development`) — when the task has clear inputs/outputs or is fixing a bug with a reproducible case
- **Subagent-driven** (`superpowers:subagent-driven-development`) — when the plan has 3+ sequential steps that each produce testable output
- **Parallel agents** (`superpowers:dispatching-parallel-agents`) — when the plan has 2+ independent tasks with no shared state (e.g., separate packages, separate files)

If unsure, default to sequential implementation without invoking a sub-skill.

At key decision points, check with the user before proceeding.

When pausing mid-implementation, write a comment listing:

- Files created or modified so far
- Steps completed from the plan
- Remaining steps

**Comment content:** Summary of what was implemented, files changed, any deviations from the plan.

### Phase 5: Verify and Ship

#### 5a. Verify

Invoke `superpowers:verification-before-completion` to confirm:

- All tests pass
- No regressions
- Implementation matches the plan

#### 5b. Update Documentation

Invoke `update-doc` skill to compare implementation changes against all documentation and fix discrepancies before committing.

#### 5c. Commit and Push

Invoke `git:commit-push` skill to commit and push changes.

#### 5d. Open PR

Create a pull request that references the issue:

```bash
gh pr create --title "<title>" --body "$(cat <<'EOF'
## Summary

<1-3 bullet points describing what was done>

## Issue

Closes #<N>

## Test Plan

<how to verify the changes>

---
🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

#### 5e. Final Comment

**Comment content:** Test results summary, PR URL, final status.

Mark the comment status as `✅ Completed`.

### Done

Present the PR URL to the user. The issue will be automatically closed when the PR is merged (via `Closes #N`).
