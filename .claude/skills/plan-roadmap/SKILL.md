---
name: plan-roadmap
description: Use when the user asks to "plan a roadmap", "review milestone scope", "rebalance milestones", "check if milestone is too large", or wants to adjust issue distribution across versions to fit a one-week release cadence.
---

# Plan Roadmap

Review all milestones' scope, analyze issue relationships, and rebalance to fit a one-week release cadence.

## Language

All issue comments and milestone descriptions MUST be written in **English**.

## Principles

- One minor version = one week of work
- Every release must deliver visible user value — avoid pure tech-debt or pure testing milestones
- Epics and their children should stay in the same milestone (or move together)
- Blocked issues cannot ship before their blockers
- Discussion-tagged issues need design first — don't schedule for immediate delivery

## Process

### Step 1: Gather

Fetch all milestones and their issues at once:

```bash
# All milestones
gh api repos/typemd/typemd/milestones --jq '.[] | {title, number, open_issues, closed_issues}'

# Open issues for each milestone (run in parallel)
gh issue list --milestone "v<VERSION>" --state open --json number,title,labels,body --limit 50
# Repeat for every milestone
```

Run the `gh issue list` commands for all milestones in parallel to minimize latency.

### Step 2: Analyze relationships

For each issue across all milestones, query parent, sub-issues, and blocking relationships:

```bash
gh api graphql -f query='query {
  repository(owner:"typemd", name:"typemd") {
    issue(number: <N>) {
      number title state
      parent { number title state milestone { title } }
      subIssues(first: 20) { nodes { number title state milestone { title } } }
      blockedBy(first: 10) { nodes { number title state milestone { title } } }
      blocking(first: 10) { nodes { number title state milestone { title } } }
    }
  }
}'
```

Build a dependency map:

- **Epic groups** — parent + children should move together
- **Block chains** — if A blocks B, A must ship first (same or earlier milestone)
- **Cross-milestone blockers** — issues blocked by items in a later milestone (flag as problem)

### Step 3: Estimate effort

Categorize each issue by effort:

| Size | Criteria | Estimate |
|------|----------|----------|
| **Small** | Single file change, clear scope, no design needed | < 1 day |
| **Medium** | Multiple files, some design, tests needed | 1–2 days |
| **Large** | Cross-package, significant design, new patterns | 3–5 days |

**Estimation heuristics:**

- `discussion` label → needs design phase, add 1 day
- Epic with open children → sum children's estimates
- Already-closed children of epic → reduce parent estimate
- `chore` label → typically small
- New command/feature → medium unless it touches core data model
- Dependency upgrade → medium-to-large (risk of breakage)

**Weekly budget:** ~5 working days. Target 70–80% utilization (3.5–4 days of estimated work) to allow for unexpected complexity.

### Step 4: Propose rebalancing

Present a report covering all milestones:

**Overview:**

| Milestone | Issues | Est. Effort | Budget | Status |
|-----------|--------|-------------|--------|--------|
| v0.X.0 | N open | X days | 3.5–4 days | over/under/ok |

**Per-milestone detail:**

For each milestone, list:

| # | Title | Size | Labels | Dependencies |
|---|-------|------|--------|-------------|
| ... | ... | S/M/L | ... | blocks #N, blocked by #N, child of #N |

**Proposed changes:**

1. **Keep** — issues that fit their milestone's budget
2. **Move** — issues to move between milestones, with reason
3. **Close** — issues already done or obsolete
4. **Flag** — cross-milestone dependency problems (e.g. blocker in a later milestone than the issue it blocks)

**Release value check:**

After effort-based rebalancing, review each milestone's release value. Every release should have a clear theme and at least one user-visible feature. Flag milestones that contain only:

- Tech debt (dependency upgrades, refactoring)
- Testing (BDD, unit tests)
- Internal tooling with no user-facing change

When a milestone lacks user value, mix in a feature issue from an adjacent milestone. Prefer pairing related items (e.g. a dependency upgrade + a feature that benefits from it) over arbitrary mixing.

Present a value summary table:

| Milestone | Theme | User-Visible Value |
|-----------|-------|--------------------|
| v0.X.0 | ... | what users get in this release |

**Rebalancing rules:**

- Every milestone must ship at least one user-visible improvement
- Prefer keeping small, independent issues (quick wins)
- Move large epics with all their children together
- Keep blockers in earlier milestones than the issues they block
- `discussion` issues go to later milestones unless design is already done
- Don't split an epic across milestones (parent and children stay together)
- When mixing in features to add value, prefer natural pairings (e.g. upgrade + enhancement that depends on it)

### Step 5: Confirm and execute

Use AskUserQuestion to confirm the plan. Options:

- **"Looks good, execute"** — proceed
- **"I want to adjust"** — iterate

Once confirmed, execute all changes across milestones:

```bash
# Move issues between milestones
gh issue edit <number> --milestone "v<TARGET>"

# Close obsolete issues
gh issue close <number> --comment "<reason>"

# Create new milestone if needed
gh api repos/typemd/typemd/milestones -f title="v<VERSION>" -f description="<description>"
```

### Step 6: Report

Present final state of all affected milestones:

```bash
gh api repos/typemd/typemd/milestones --jq '.[] | "\(.title): \(.open_issues) open, \(.closed_issues) closed"'
```
