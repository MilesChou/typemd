## 1. Title Panel Rendering

- [ ] 1.1 Add `renderTitlePanel(obj, typeName, emoji, width)` function in `tui/detail.go` that returns the title content string: `emoji type · DisplayName` (or `type · DisplayName` when no emoji)
- [ ] 1.2 Add unit tests for `renderTitlePanel` covering: with emoji, without emoji, nil object (empty string)

## 2. Layout Integration

- [ ] 2.1 Store type name and emoji on the model (or derive from current `typeGroup`) so `View()` can pass them to `renderTitlePanel`
- [ ] 2.2 Update `View()` in `tui/app.go` to compose the right-side layout vertically: title panel above body+props row using `lipgloss.JoinVertical`
- [ ] 2.3 Calculate title panel width dynamically: body width + props width (+ borders) when props visible, body width only when props hidden
- [ ] 2.4 Reduce `contentH` by 3 (title panel height) when an object is selected, so body and props panels shrink accordingly

## 3. Remove Old Title from Body

- [ ] 3.1 Remove `renderBodyHeader()` call from `renderBody()` and from edit mode rendering in `View()`
- [ ] 3.2 Update or remove `renderBodyHeader()` function and any tests that depend on it
- [ ] 3.3 Remove `bodyEditHeaderLines` constant adjustment since title is no longer in body panel

## 4. Edge Cases and Polish

- [ ] 4.1 Hide title panel when no object is selected (keep existing placeholder in body panel)
- [ ] 4.2 Verify title panel width updates correctly when properties panel is toggled (`p` key) or resized (`[`/`]` keys)
- [ ] 4.3 Verify edit mode still works correctly without the old body header
