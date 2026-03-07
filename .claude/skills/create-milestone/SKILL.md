---
name: create-milestone
description: This skill should be used when the user asks to "create a milestone", "plan a release", "organize issues into a milestone", "what should the next version be", or discusses grouping open issues into a release target.
---

# Create GitHub Milestone

Plan and create a GitHub milestone for the TypeMD project through iterative Q&A.

Do NOT create the milestone until the user explicitly confirms. Analyze existing issues first, then propose a milestone through conversation.

## Process

### Step 1: Survey open issues

Fetch all open issues without a milestone:

```bash
gh issue list --state open --json number,title,labels,milestone
```

Also fetch existing milestones for context:

```bash
gh api repos/typemd/typemd/milestones
```

Present a summary of unassigned issues grouped by component label (core, cli, tui, mcp, web).

### Step 2: Propose milestone

Based on the unassigned issues, propose a milestone:

- **Title** — version number following semver (e.g. `v0.2.0`)
- **Description** — one sentence summarizing the release goal
- **Candidate issues** — which unassigned issues fit this milestone

Explain the reasoning behind the grouping. Then use AskUserQuestion to collect feedback on three topics at once:

1. Does the version number and goal make sense? (options: "OK", "I have a different idea")
2. Which issues should be included or excluded? (options: "All good", "I want to adjust")
3. Should a due date be set? (options: "No due date", "Set a due date")

### Step 3: Draft and confirm

Present the full milestone draft:

- **Title**
- **Description**
- **Due date** — or none
- **Issues to assign** — list with number and title

Use AskUserQuestion to confirm: "This is the milestone I'll create. Anything to adjust?" (options: "Looks good, create it", "I want to adjust")

Only proceed after the user selects "Looks good, create it".

### Step 4: Create milestone and assign issues

Create the milestone:

```bash
gh api repos/typemd/typemd/milestones -f title="<title>" -f description="<description>" -f state="open"
```

Add `--field due_on=<YYYY-MM-DDT00:00:00Z>` if a due date was set.

Then assign each confirmed issue:

```bash
gh issue edit <number> --milestone "<title>"
```

### Step 5: Confirm

Return the milestone URL and list of assigned issues to the user.
