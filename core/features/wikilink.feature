Feature: Links and backlinks
  Objects can reference each other using [[type/name-ulid]] wiki-link syntax in their body.
  The system tracks these links and their backlinks in the database.

  Scenario: Links are parsed and stored on sync
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And a "person" object named "robert-martin" exists
    And "clean-code" body contains a wiki-link to "robert-martin"
    When I sync the index
    Then "clean-code" should have 1 wiki-link
    And the wiki-link target should be "robert-martin"
    And "robert-martin" should have 1 backlink from "clean-code"

  Scenario: Broken link has empty resolved ID
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And "clean-code" body contains a wiki-link to "person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj"
    When I sync the index
    Then "clean-code" should have 1 wiki-link
    And the wiki-link should have an empty resolved ID

  Scenario: Links are updated on re-sync
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And a "person" object named "alice" exists
    And a "person" object named "bob" exists
    And "clean-code" body contains a wiki-link to "alice"
    And I sync the index
    When I change "clean-code" wiki-link to "bob"
    And I sync the index
    Then "clean-code" wiki-link should point to "bob"
    And "alice" should have 0 backlinks

  Scenario: Links support display text
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And a "person" object named "robert-martin" exists
    And "clean-code" body contains a wiki-link to "robert-martin" with display text "Uncle Bob"
    When I sync the index
    Then "clean-code" should have 1 wiki-link
    And the wiki-link display text should be "Uncle Bob"

  Scenario: Backlinks from multiple sources
    Given a vault is ready with note schemas
    And a "note" object named "target" exists
    And a "note" object named "alpha" exists
    And a "note" object named "beta" exists
    And "alpha" body contains a wiki-link to "target"
    And "beta" body contains a wiki-link to "target"
    When I sync the index
    Then "target" should have 2 backlinks

  Scenario: Duplicate link targets are deduplicated
    Given a vault is ready with note schemas
    And a "note" object named "source" exists
    And a "note" object named "target" exists
    And "source" body contains duplicate wiki-links to "target"
    When I sync the index
    Then "source" should have 1 wiki-link
    And the wiki-link target should be "target"

  Scenario: Links are cleaned up when source object is deleted
    Given a vault is ready with note schemas
    And a "note" object named "writer" exists
    And a "note" object named "reader" exists
    And "reader" body contains a wiki-link to "writer"
    And I sync the index
    When I delete the object "reader" from disk
    And I sync the index
    Then "writer" should have 0 backlinks

  Scenario: Backlinks appear in display properties
    Given a vault is ready with note schemas
    And a "note" object named "origin" exists
    And a "note" object named "destination" exists
    And "origin" body contains a wiki-link to "destination"
    And I sync the index
    Then "destination" should have a "backlinks" display property from "origin"

  Scenario: Links without display text render with stripped ULID
    Given a vault is ready with note schemas
    And a "note" object named "my-note" exists
    And a "note" object named "doc" exists
    And "doc" body contains a wiki-link to "my-note"
    When I render the body of "doc"
    Then the rendered body should contain "note/my-note"
    And the rendered body should not contain "[["
