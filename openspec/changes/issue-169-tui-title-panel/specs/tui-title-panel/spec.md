## ADDED Requirements

### Requirement: Title panel displays object identity
The TUI detail view SHALL display a dedicated title panel above the body and properties panels showing the type emoji, type name, and object display name.

#### Scenario: Title panel with emoji
- **WHEN** an object of type "book" with emoji "📖" and display name "Clean Code" is selected
- **THEN** the title panel SHALL display "📖 book · Clean Code"

#### Scenario: Title panel without emoji
- **WHEN** an object of type "note" with no emoji and display name "My Note" is selected
- **THEN** the title panel SHALL display "note · My Note"

### Requirement: Title panel spans body and properties width
The title panel SHALL span the full width of both the body panel and properties panel combined.

#### Scenario: Properties visible
- **WHEN** the properties panel is visible
- **THEN** the title panel width SHALL equal the combined width of body and properties panels (including borders)

#### Scenario: Properties hidden
- **WHEN** the properties panel is hidden
- **THEN** the title panel width SHALL equal the body panel width only

### Requirement: Title panel hidden when no object selected
The title panel SHALL NOT be displayed when no object is selected.

#### Scenario: No selection
- **WHEN** no object is selected in the list
- **THEN** the title panel SHALL be hidden and the body panel SHALL display the default placeholder message

### Requirement: Body panel no longer contains title header
The body panel SHALL NOT render the object title or separator line. The body panel SHALL display only the markdown body content.

#### Scenario: Body content without title
- **WHEN** an object is selected and the body panel is rendered
- **THEN** the body panel SHALL start directly with the markdown body content, without a title line or separator

### Requirement: Title panel height is fixed
The title panel SHALL occupy exactly 3 lines of vertical space (1 content line + 2 border lines).

#### Scenario: Vertical space allocation
- **WHEN** the TUI detail view is rendered with an object selected
- **THEN** the body and properties panels SHALL have their content height reduced by 3 lines compared to the no-title-panel state
