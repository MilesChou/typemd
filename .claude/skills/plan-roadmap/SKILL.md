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
- Parent issues are epic trackers (organizational only) — place them in the same milestone as their first child
- Children of an epic may span multiple milestones; dependency order is between siblings, not parent → child
- Blocked issues cannot ship before their blockers (check `blockedBy` between sibling issues)
- Discussion-tagged issues need design first — don't schedule for immediate delivery
- Themes drive issue selection — not just effort balancing
- When a milestone is under budget, prefer non-business tasks (testing, infra, docs) over adding more features

## Process

### Step 1: Gather

Fetch all milestones, their issues, and unassigned issues. Use GraphQL for milestone issue lists — it returns authoritative milestone assignments unlike `gh issue list --milestone` which can be unreliable:

```bash
# All milestones
gh api repos/typemd/typemd/milestones --jq '.[] | {title, number, open_issues, closed_issues}'

# Open issues for each milestone via GraphQL (run all in parallel)
gh api graphql -f query='query {
  repository(owner:"typemd", name:"typemd") {
    milestone(number: <N>) {
      title
      issues(first: 30, states: OPEN) {
        nodes { number title labels(first: 5) { nodes { name } }
          parent { number title state milestone { title } }
        }
      }
    }
  }
}'

# Open issues without a milestone
gh issue list --search "no:milestone is:open" --json number,title,labels --limit 100
```

### Step 2: Analyze relationships

For all milestone issues, query parent, sub-issues, and blocking relationships in parallel:

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
- **Cross-milestone blockers** — flag as problems

### Step 3: Define themes for the next 5 releases

Look at the next 5 upcoming milestones. For each one, propose a **release theme** — a short, user-facing headline that:

- Describes the user-visible value in a concrete, exciting way (e.g. "Type System is Here", "Built-in Types", "TUI Editing")
- Follows naturally from the previous release
- Respects dependency order

**Headline style:** Think product update cards like Capacities — feature-forward, verb-driven, not abstract. "Navigation & Discovery" beats "Usability Improvements". Present all 5 themes together and discuss with the user before proceeding. Adjust based on feedback.

**If a milestone has 0 open issues**, flag it as ready to release — don't force issues into it.

### Step 4: Populate each themed milestone

Work through each milestone **one at a time** with the user:

1. **Start from existing issues** already assigned to that milestone
2. **Scan all unassigned issues** for relevance to the theme — propose good candidates
3. **Move out** issues that don't fit the theme (to a later milestone or backlog)
4. **Check effort budget** — target 70–80% of weekly capacity (3.5–4 days)

**When under budget:** Look for non-business tasks first — testing, infra, documentation — before adding more features. These are good fillers that don't dilute the milestone's theme.

**Effort sizing:**

| Size | Criteria | Estimate |
|------|----------|----------|
| **Small** | Single file change, clear scope | < 1 day |
| **Medium** | Multiple files, some design, tests needed | 1–2 days |
| **Large** | Cross-package, significant design, new patterns | 3–5 days |

Heuristics: `discussion` label → add 1 day; new command → medium; epic tracker → sum children.

**Issue fitness criteria** (for adding from backlog):
- Complements or extends the milestone's theme
- Small or medium effort
- Not blocked by unresolved discussions
- Dependencies resolved in same or earlier milestone
- Epic tracker (parent): place in same milestone as its first child
- Sibling dependencies: if child A blocks child B, A must be in same or earlier milestone than B

### Step 5: Confirm and execute

Confirm each milestone's issue list with the user before executing moves. Work one milestone at a time — confirm theme and issues together, then execute immediately before moving to the next.

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
