---
name: create-issue
description: This skill should be used when the user asks to "create an issue", "open an issue", "file a bug", "request a feature", "add a task", or discusses a problem or idea that should be tracked as a GitHub issue.
---

# Create GitHub Issue

Create a GitHub issue for the TypeMD project through iterative Q&A.

Do NOT create the issue until the user explicitly confirms. Ask questions one at a time to refine the idea before writing anything.

## Language

All issue content (title, body) MUST be written in **English**.

## Process

### Step 1: Understand the idea

Ask the user what they want to achieve. One question at a time:

1. **What's the problem or idea?** — Understand the motivation first.
2. **What's the expected behavior?** — Clarify what "done" looks like.
3. **Any constraints or context?** — Related issues, technical considerations, scope.

Use multiple choice (via AskUserQuestion) when possible. Keep it conversational — skip questions that are already answered from context.

### Step 2: Check for duplicates

Before proceeding, search existing issues for potential duplicates:

```bash
gh issue list --state all --json number,title,state,labels --limit 100
```

Compare the new idea against existing issues by title and topic. If a similar issue exists, present it to the user via AskUserQuestion with options:

- **"It's a duplicate"** — stop and link to the existing issue
- **"Related but different"** — continue creating, and reference the related issue in the body
- **"Not related"** — continue creating as normal

Skip this step only if the user has already referenced a specific issue number in their request.

### Step 3: Determine issue type

Each issue gets exactly **one issue type**. Issue types replace the old type labels (`bug`, `enhancement`, `epic`, `chore`).

| Issue Type | ID | When to use |
|---|---|---|
| Task | `IT_kwDOD9OO7M4B4kKu` | CI, refactoring, dependencies, project configuration, general tasks |
| Bug | `IT_kwDOD9OO7M4B4kKv` | Something isn't working |
| Feature | `IT_kwDOD9OO7M4B4kKw` | New feature or improvement with clear scope |
| Epic | `IT_kwDOD9OO7M4B4rT8` | High-level feature plan, will be broken into sub-issues |

Suggest the type based on context. Ask for confirmation if ambiguous.

### Step 4: Determine labels

After setting the issue type, assign **one or more component labels**. Optionally add extra labels if applicable.

**Component label** (pick one or more):

| Label | Scope |
|---|---|
| `core` | Core library — objects, types, relations, index (`core/`) |
| `cli` | CLI commands (`cmd/`) |
| `tui` | Terminal UI (`tui/`) |
| `mcp` | MCP server (`mcp/`) |
| `web` | Web UI (`web/`) |
| `app` | Desktop app via Wails (`app/`) |

**Optional extra labels**:

| Label | When to use |
|---|---|
| `discussion` | Needs discussion before implementation |
| `documentation` | Docs changes |

Suggest labels based on context. Ask for confirmation if ambiguous.

### Step 5: Assign milestone

Open milestones:

!`gh api repos/typemd/typemd/milestones --jq '.[] | "\(.number) \(.title)"'`

Present the milestones above as options using AskUserQuestion. Always include a "None" option for issues that don't belong to any milestone.

### Step 6: Relationships (optional)

Ask whether this issue has relationships to other issues. Use AskUserQuestion with options:

- **"No relationships"** — skip
- **"Sub-issue of an existing issue"** — will set a parent issue
- **"Blocked by another issue"** — will add a blocking relationship

If the user selects a relationship, ask them to specify the issue number. Look up the issue's node ID:

```bash
gh issue view <number> --json id --jq '.id'
```

Multiple relationships can be set. After the user is done, proceed to the next step.

### Step 7: Draft and confirm

Present the full issue draft to the user, then use AskUserQuestion to confirm:

- **Title** — concise, plain language, no prefix
- **Type** — issue type name
- **Labels** — component + optional extra labels
- **Milestone** — selected milestone or none
- **Relationships** — parent issue or blocking issues, if any
- **Body** — using the template below

Use AskUserQuestion with options like "Create" and "Needs changes" to get user confirmation. Only proceed after the user confirms.

### Step 8: Create issue

Use GraphQL to create the issue with the issue type set:

```bash
# Get repo and milestone IDs
REPO_ID=$(gh api repos/typemd/typemd --jq '.node_id')
MILESTONE_ID=$(gh api repos/typemd/typemd/milestones/<number> --jq '.node_id')

# Get label IDs
LABEL_IDS=$(gh api repos/typemd/typemd/labels/<name> --jq '.node_id')

gh api graphql -f query='
mutation($repoId: ID!, $title: String!, $body: String!, $typeId: ID!, $milestoneId: ID, $labelIds: [ID!], $parentId: ID) {
  createIssue(input: {
    repositoryId: $repoId
    title: $title
    body: $body
    issueTypeId: $typeId
    milestoneId: $milestoneId
    labelIds: $labelIds
    parentIssueId: $parentId
  }) {
    issue { number url }
  }
}
' \
  -f repoId="$REPO_ID" \
  -f title="<title>" \
  -f body="<body>" \
  -f typeId="<issue_type_id>" \
  -f milestoneId="$MILESTONE_ID" \
  -f labelIds='["<label_id_1>","<label_id_2>"]' \
  -f parentId="<parent_issue_id>"
```

Omit `milestoneId`, `labelIds`, or `parentId` parameters if not applicable.

**After creation**, if there are blocking relationships, add them:

```bash
# Get the newly created issue's node ID
ISSUE_ID=$(gh issue view <number> --json id --jq '.id')

gh api graphql -f query='
mutation($issueId: ID!, $blockingId: ID!) {
  addBlockedBy(input: {
    issueId: $issueId
    blockingIssueId: $blockingId
  }) {
    blockedIssue { number }
  }
}
' -f issueId="$ISSUE_ID" -f blockingId="<blocking_issue_id>"
```

Body template:

```markdown
## Summary

<1-3 sentences describing the issue>

## Current Behavior

<what happens now, if applicable>

## Expected Behavior

<what should happen>
```

### Step 9: Confirm

Return the issue URL to the user.
