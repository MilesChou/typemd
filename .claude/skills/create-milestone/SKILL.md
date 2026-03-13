---
name: create-milestone
description: This skill should be used when the user asks to "create a milestone", "plan a release", "organize issues into a release", "what should the next version be", or discusses grouping open issues into a release target.
---

# Create Release

Plan and create a Release issue for the TypeMD project through iterative Q&A.

A **Release issue** is a GitHub issue with the `Release` type that serves as the version planning container. Child issues are linked via GitHub's sub-issue relationship.

Do NOT create the Release issue until the user explicitly confirms. Analyze existing issues first, then propose a release through conversation.

## Language

All issue content (title, body) MUST be written in **English**.

## Issue Type IDs

| Issue Type | ID |
|---|---|
| Release | `IT_kwDOD9OO7M4B5cCA` |

## Process

### Step 1: Survey open issues

Fetch all open issues that are not sub-issues of any Release issue:

```bash
gh issue list --state open --json number,title,labels --limit 100
```

Also fetch existing Release issues for context:

```bash
gh api graphql -f query='query {
  repository(owner:"typemd", name:"typemd") {
    issues(first: 20, states: [OPEN, CLOSED], filterBy: {issueType: "Release"}, orderBy: {field: CREATED_AT, direction: DESC}) {
      nodes { number title state subIssues(first: 30) { nodes { number } } }
    }
  }
}'
```

Filter out issues that are already sub-issues of a Release. Present a summary of unassigned issues grouped by component label (core, cli, tui, mcp, web).

### Step 2: Propose release

Based on the unassigned issues, propose a release:

- **Version** — following semver (e.g. `v0.4.0`)
- **Theme** — a short, user-facing headline (e.g. "Move In", "Your Objects, Your Rules")
- **Candidate issues** — which unassigned issues fit this release

Headline style: think product update cards like Capacities — feature-forward, verb-driven, not abstract.

Explain the reasoning behind the grouping. Then use AskUserQuestion to collect feedback on two topics at once:

1. Does the version number and theme make sense? (options: "OK", "I have a different idea")
2. Which issues should be included or excluded? (options: "All good", "I want to adjust")

### Step 3: Draft and confirm

Present the full release draft:

- **Title** — `v<VERSION> — <Theme>`
- **Issues to include** — list with number and title

Use AskUserQuestion to confirm: "This is the release I'll create. Anything to adjust?" (options: "Looks good, create it", "I want to adjust")

Only proceed after the user selects "Looks good, create it".

### Step 4: Create Release issue and link sub-issues

Create the Release issue using GraphQL:

```bash
REPO_ID=$(gh api repos/typemd/typemd --jq '.node_id')

cat > /tmp/create_release.json << 'EOF'
{
  "query": "mutation($repoId: ID!, $title: String!, $body: String!, $typeId: ID!) { createIssue(input: { repositoryId: $repoId, title: $title, body: $body, issueTypeId: $typeId }) { issue { number url id } } }",
  "variables": {
    "repoId": "<REPO_ID>",
    "title": "v<VERSION> — <Theme>",
    "body": "## Theme\n\n<theme description>\n\n## Scope\n\n- #N <title>\n- #N <title>",
    "typeId": "IT_kwDOD9OO7M4B5cCA"
  }
}
EOF

gh api graphql --input /tmp/create_release.json
rm -f /tmp/create_release.json
```

Then link each issue as a sub-issue. Note: issues that already have a parent cannot be added (GitHub limits to one parent). For those issues, reference them in the body only.

```bash
PARENT_ID=$(gh issue view <release_number> --json id --jq '.id')

# For each child issue (skip if it already has a parent)
CHILD_ID=$(gh issue view <number> --json id --jq '.id')
gh api graphql -f query="mutation { addSubIssue(input: { issueId: \"$PARENT_ID\", subIssueId: \"$CHILD_ID\" }) { subIssue { number } } }"
```

### Step 5: Confirm

Return the Release issue URL and list of linked sub-issues to the user.
