Feature: Type directory structure
  Type schemas support directory format (.typemd/types/<name>/schema.yaml)
  in addition to single-file format (.typemd/types/<name>.yaml).
  Auto-migration upgrades single files to directory format on load.

  # ── Load from directory format ──────────────────────────────

  Scenario: Load type schema from directory format
    Given a vault is ready
    And a type schema directory "movie" with schema content:
      """
      name: movie
      emoji: "\U0001F3AC"
      properties:
        - name: rating
          type: number
      """
    When I load type "movie"
    Then no error should occur
    And the loaded schema should have emoji "🎬"
    And the loaded schema should have 1 property

  Scenario: Directory format takes precedence over single file
    Given a vault is ready
    And a type schema file "movie" exists on disk
    And a type schema directory "movie" with schema content:
      """
      name: movie
      emoji: "\U0001F3AC"
      properties: []
      """
    When I load type "movie"
    Then no error should occur
    And the loaded schema should have emoji "🎬"

  # ── Auto-migration ──────────────────────────────────────────

  Scenario: Single file auto-migrated to directory on load
    Given a vault is ready
    And a type schema file "project" exists on disk
    When I load type "project"
    Then no error should occur
    And the type schema directory "project" should exist
    And the type schema single file "project" should not exist

  # ── ListTypes with mixed formats ────────────────────────────

  Scenario: ListTypes discovers directory format
    Given a vault is ready
    And a type schema directory "movie" with schema content:
      """
      name: movie
      properties: []
      """
    When I list all types
    Then the type list should contain "movie"

  # ── SaveType writes directory format ────────────────────────

  Scenario: SaveType creates directory format
    Given a vault is ready
    And a type schema "article" with no extra fields
    And the schema has a "title" string property
    When I save the type schema
    Then no error should occur
    And the type schema directory "article" should exist

  Scenario: SaveType removes old single file
    Given a vault is ready
    And a type schema file "draft" exists on disk
    And a type schema "draft" with no extra fields
    When I save the type schema
    Then no error should occur
    And the type schema directory "draft" should exist
    And the type schema single file "draft" should not exist

  # ── DeleteType removes directory ────────────────────────────

  Scenario: Delete type removes entire directory
    Given a vault is ready
    And a type schema directory "scratch" with schema content:
      """
      name: scratch
      properties: []
      """
    When I delete type "scratch"
    Then no error should occur
    And the type schema directory "scratch" should not exist
