## Context

The TUI detail view currently has a three-panel horizontal layout: left (object list), body (title + markdown), and properties. The title is rendered inside the body panel by `renderBodyHeader()` as `type/display-name` with a dash separator line. Type emoji is already available via `typeGroup.Emoji` (loaded from type schema).

The goal is to extract the title into its own dedicated panel above the body and properties panels, on the right side only.

## Goals / Non-Goals

**Goals:**
- Add a title panel spanning above both body and properties panels
- Display format: `emoji type · DisplayName` (e.g., "📖 book · Clean Code")
- Remove title rendering from the body panel
- Maintain correct vertical height calculations

**Non-Goals:**
- Pinned properties display (deferred to #168)
- Title panel as a focusable panel (no scroll, no edit — purely informational)
- Changes to the left list panel

## Decisions

### 1. Layout composition approach

**Decision:** Compose the right side vertically using `lipgloss.JoinVertical()` for title + body/props row, then join horizontally with left panel.

**Rationale:** The current layout uses `lipgloss.JoinHorizontal()` for all panels at the same level. The title panel spans body + props width, so we need a vertical join on the right side first.

```
panels = JoinHorizontal(
    leftPanel,
    JoinVertical(
        titlePanel,
        JoinHorizontal(bodyPanel, propsPanel),
    ),
)
```

**Alternative considered:** Render title as part of body panel content (current approach). Rejected because the title must span both body and properties panels.

### 2. Title panel height

**Decision:** Fixed at 3 lines total (1 content line + 2 border lines). The content area of body and properties panels is reduced by 3 lines.

**Rationale:** The title is always a single line. No scrolling or dynamic sizing needed.

### 3. Title panel width

**Decision:** Title panel width = body panel width + properties panel width + properties border. When properties are hidden, title width = body panel width only.

**Rationale:** Must exactly match the combined width of the panels below it for visual alignment.

### 4. Title content format

**Decision:** ` emoji type · DisplayName` — type emoji and type name first, then separator dot, then object display name.

**Rationale:** Groups the type context (emoji + name) together, then the specific object. Consistent with the list panel where groups show `emoji type_name`.

### 5. No-selection state

**Decision:** When no object is selected, hide the title panel entirely and keep the current "Select an object to view details" message in the body panel.

**Rationale:** The title panel has no content to show without a selected object.

## Risks / Trade-offs

- **[Vertical space]** The title panel consumes 3 lines of vertical space, reducing body/props content area. → Acceptable trade-off; the title was already taking 2 lines inside the body panel, net cost is only 1 additional line for the border.
- **[Width alignment]** Title panel width must stay in sync when properties panel is toggled or resized. → Calculate title width dynamically from current body + props widths in `View()`.
