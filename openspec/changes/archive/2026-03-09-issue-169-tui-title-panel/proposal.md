## Why

The object title is currently embedded inside the body panel, rendered by `renderBodyHeader()` as a plain text line with a dash separator. This makes it visually indistinct from body content and provides no clear, persistent context header. A dedicated title panel above the body and properties panels provides immediate visual identification of the current object and establishes the structural foundation for pinned properties (#168).

## What Changes

- Add a new title panel that spans above both the body and properties panels (right side only — left list panel is unaffected)
- Title panel displays: type emoji + type name + object display name (e.g., "📖 book · Clean Code")
- Title panel has its own rounded border, visually distinct from body and properties panels
- Remove `renderBodyHeader()` from body panel rendering — title moves to the dedicated panel
- Adjust vertical height calculations: body and properties panels shrink by the title panel's height (1 line content + 2 border = 3 lines)
- When no object is selected, the title panel is hidden

### Layout

Before:

```
╭──────────╮╭──────────────────╮╭──────────╮
│ ▼ book   ││ book/clean-code  ││Properties│
│   Clean..││ ──────────────── ││ ──────── │
│   GoLang ││ Body content...  ││ author:  │
│ ▼ person ││ ...              ││   Bob    │
╰──────────╯╰──────────────────╯╰──────────╯
```

After:

```
╭──────────╮╭─────────────────────────────╮
│ ▼ book   ││ 📖 book · Clean Code        │
│   Clean..│╰─────────────────────────────╯
│   GoLang │╭─────────────────╮╭──────────╮
│ ▼ person ││ Body content... ││Properties│
│   Alice  ││ ...             ││ author:  │
│          ││                 ││   Bob    │
╰──────────╯╰─────────────────╯╰──────────╯
```

## Capabilities

### New Capabilities
- `tui-title-panel`: Dedicated title panel in TUI detail view showing type emoji, type name, and object display name, positioned above body and properties panels on the right side

### Modified Capabilities

## Impact

- `tui/detail.go`: Add `renderTitlePanel()`, remove title from `renderBody()` and `renderBodyHeader()`
- `tui/app.go`: Adjust `View()` layout composition — title panel sits above body+props row; update `contentH` calculation to subtract title panel height
- `tui/theme.go`: No changes expected (reuses existing border/focus styles)
