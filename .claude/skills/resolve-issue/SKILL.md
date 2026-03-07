---
name: resolve-issue
description: This skill should be used when the user asks to "resolve an issue", "work on issue #N", "fix #N", "implement #N", "close #N", "tackle #N", "pick up #N", "start working on #N", "what should I work on next", or references a specific GitHub issue number they want to work on. Can also auto-select the best issue when no number is specified.
allowed-tools:
  - Bash(gh api repos/typemd/typemd/milestones:*)
  - Bash(gh issue list:*)
  - Bash(gh issue view:*)
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
- **"Start over"** — discard previous progress and begin from Preflight

If no progress comments exist, start from Preflight.

## Comment Format

All phases (Phase 1–3) write comments to the issue. Preflight is a lightweight step that does not require issue comments.

Use this format:

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

## Preflight

Preflight covers all lightweight preparation steps before the main phases begin. No issue comments are written during this stage.

### Issue Selection (when no issue number is specified)

If the user does not specify an issue number, automatically select the best issue to work on.

**Step 1: Find the nearest milestone**

```bash
gh api repos/typemd/typemd/milestones --jq 'sort_by(.due_on // "9999") | .[0] | {title, number, due_on}'
```

If no milestone exists, fall back to all open issues.

**Step 2: List open issues in that milestone**

```bash
gh issue list --milestone "<milestone_title>" --state open --json number,title,labels,assignees,body --limit 50
```

**Step 3: Rank issues by priority**

Evaluate each issue using these criteria (highest priority first):

1. **Blocker** — blocks other issues (look for "blocks #N" or "blocked by #N" references in issue bodies, or issues labeled `blocker` / `priority:critical`)
2. **High value** — labeled `priority:high`, or is a bug affecting core functionality
3. **Low effort, high impact** — small scope issues that unblock progress (labeled `good first issue`, `quick win`, or estimated as small)
4. **Dependencies resolved** — issues whose blockers are already closed

**Step 4: Present top 3 candidates**

Ask the user via AskUserQuestion with the top 3 recommended issues:

- **"#N: \<title\>"** — for each candidate, include a one-line reason why it's recommended (e.g., "Blocks 3 other issues", "Critical bug", "Quick win for milestone X")

The user selects one, then proceed to **Check Issue State** with that issue number.

### Check Issue State

Verify the issue is actionable:

```bash
gh issue view <number> --json state,linkedBranches
```

- If the issue is **closed**, inform the user and stop.
- If there is already an **open PR** linked to this issue, inform the user and ask whether to continue or stop.

### Understand the Issue

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

### Branch Strategy

If a branch matching `fix/issue-<N>-*` or `feat/issue-<N>-*` already exists, inform the user and ask whether to reuse it or create a new one.

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

## Phases

### Phase 1: Design

Invoke `superpowers:brainstorming` skill to explore the design space.

The brainstorming skill will:

1. Explore project context
2. Ask clarifying questions
3. Propose 2-3 approaches
4. Present design for user approval

**IMPORTANT:** When brainstorming invokes `superpowers:writing-plans`, the plan output should NOT be saved to `docs/plans/`. Instead, capture the full design and implementation plan to write into the issue comment.

**Comment content:** The complete design — architecture decisions, approach chosen, implementation plan with steps.

### Phase 2: Implement

Execute the implementation plan from Phase 1. Choose the appropriate approach:

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

### Phase 3: Verify and Ship

#### 3a. Verify

Invoke `superpowers:verification-before-completion` to confirm:

- All tests pass
- No regressions
- Implementation matches the plan

#### 3b. Update Documentation

Invoke `update-doc` skill to compare implementation changes against all documentation and fix discrepancies before committing.

#### 3c. Commit and Push

Invoke `git:commit-push` skill to commit and push changes.

#### 3d. Open PR

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

#### 3e. Final Comment

**Comment content:** Test results summary, PR URL, final status.

Mark the comment status as `✅ Completed`.

### Done

Present the PR URL to the user. The issue will be automatically closed when the PR is merged (via `Closes #N`).
