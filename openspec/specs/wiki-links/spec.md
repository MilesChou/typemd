## ADDED Requirements

### Requirement: Wiki-link syntax is parsed from markdown body

The system SHALL parse `[[target]]` and `[[target|display text]]` syntax from object markdown body content. Target SHALL be a full object ID in `type/name-ulid` format. Duplicate targets within the same body SHALL be deduplicated, keeping the first occurrence.

#### Scenario: Simple wiki-link is parsed

- **WHEN** an object body contains `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6]]`
- **THEN** the parser extracts one wiki-link with target `person/bob-01kk3gqm8zrrbjjwkx90f727y6` and empty display text

#### Scenario: Wiki-link with display text is parsed

- **WHEN** an object body contains `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6|Uncle Bob]]`
- **THEN** the parser extracts one wiki-link with target `person/bob-01kk3gqm8zrrbjjwkx90f727y6` and display text `Uncle Bob`

#### Scenario: Duplicate targets are deduplicated

- **WHEN** an object body contains the same target `[[book/clean-code-01abc]]` twice
- **THEN** the parser returns only one wiki-link for that target

### Requirement: Wiki-links are stored in the database on sync

The system SHALL extract wiki-links from each object body during index sync and store them in the `wikilinks` table. Each sync SHALL replace all existing wiki-links for that object (delete + insert).

#### Scenario: Wiki-links are created on first sync

- **WHEN** an object body contains a wiki-link to an existing object and the index is synced
- **THEN** a wikilink record is stored with the source object as `from_id` and the resolved target as `to_id`

#### Scenario: Wiki-links are updated on re-sync

- **WHEN** an object's wiki-link target changes and the index is re-synced
- **THEN** old wikilink records for that object are removed and new ones are inserted

#### Scenario: Wiki-links to deleted objects are cleaned up

- **WHEN** an object that is the source of wiki-links is deleted and the index is synced
- **THEN** all wikilink records with that object as `from_id` are removed

### Requirement: Broken wiki-links have empty resolved ID

When a wiki-link target does not match any existing object ID, the system SHALL store the wikilink with an empty `to_id` field, preserving the original target text.

#### Scenario: Link to non-existent object

- **WHEN** an object body contains `[[person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj]]` and no such object exists
- **THEN** the wikilink record has an empty `to_id` and target `person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj`

### Requirement: Backlinks are queryable

The system SHALL provide a way to query all objects that link to a given object (backlinks) via the `wikilinks` table.

#### Scenario: Single backlink

- **WHEN** object A contains a wiki-link to object B
- **THEN** querying backlinks for object B returns object A

#### Scenario: Multiple backlinks

- **WHEN** objects A and C both contain wiki-links to object B
- **THEN** querying backlinks for object B returns both A and C

### Requirement: Backlinks are displayed as a built-in property

The system SHALL display backlinks as a system-level `backlinks` property in object detail views. This property SHALL appear after schema-defined properties and reverse relations.

#### Scenario: Object with backlinks shows them in display properties

- **WHEN** object B has backlinks from objects A and C
- **THEN** object B's display properties include a `backlinks` entry listing A and C

#### Scenario: Object without backlinks omits the property

- **WHEN** object B has no backlinks
- **THEN** object B's display properties do not include a `backlinks` entry

### Requirement: Wiki-links are rendered with display text

When rendering markdown body for display, the system SHALL replace wiki-link syntax with human-readable text. `[[target|text]]` SHALL render as the display text. `[[target]]` SHALL render as the DisplayID (target with ULID suffix stripped).

#### Scenario: Render link with display text

- **WHEN** body contains `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6|Uncle Bob]]`
- **THEN** it renders as `Uncle Bob`

#### Scenario: Render link without display text

- **WHEN** body contains `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6]]`
- **THEN** it renders as `person/bob` (ULID stripped)

### Requirement: Broken wiki-links are detected by validation

The `tmd validate` command SHALL report wiki-links whose targets do not resolve to existing objects.

#### Scenario: Broken link is reported

- **WHEN** an object contains a wiki-link `[[person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj]]` that does not resolve
- **THEN** validation reports `<object-id>: broken wiki-link [[person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj]]`

#### Scenario: Valid links pass validation

- **WHEN** all wiki-links in the vault resolve to existing objects
- **THEN** wiki-link validation reports no errors
