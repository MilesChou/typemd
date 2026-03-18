Feature: View CRUD
  Views can be listed, loaded, saved, deleted, and defaulted via the Vault facade.

  # ── ListViews ───────────────────────────────────────────────

  Scenario: List views for type with saved views
    Given a vault is ready
    And a type schema directory "book" with schema content:
      """
      name: book
      properties: []
      """
    And a saved view "default" for type "book" with sort "name" "asc"
    And a saved view "by-rating" for type "book" with sort "rating" "desc"
    When I list views for type "book"
    Then I should have 2 views

  Scenario: List views for type with no views
    Given a vault is ready
    When I list views for type "note"
    Then I should have 0 views

  # ── LoadView ────────────────────────────────────────────────

  Scenario: Load an existing view
    Given a vault is ready
    And a type schema directory "book" with schema content:
      """
      name: book
      properties: []
      """
    And a saved view "by-rating" for type "book" with sort "rating" "desc"
    When I load view "by-rating" for type "book"
    Then no error should occur
    And the loaded view name should be "by-rating"

  Scenario: Load a non-existent view returns error
    Given a vault is ready
    When I load view "missing" for type "book"
    Then an error should occur

  # ── SaveView ────────────────────────────────────────────────

  Scenario: Save a new view
    Given a vault is ready
    And a type schema directory "book" with schema content:
      """
      name: book
      properties: []
      """
    And a view config "reading-now" with layout "list"
    And the view has sort property "name" direction "asc"
    When I save view for type "book"
    Then no error should occur
    And loading view "reading-now" for type "book" should succeed

  # ── DeleteView ──────────────────────────────────────────────

  Scenario: Delete an existing view
    Given a vault is ready
    And a type schema directory "book" with schema content:
      """
      name: book
      properties: []
      """
    And a saved view "temp" for type "book" with sort "name" "asc"
    When I delete view "temp" for type "book"
    Then no error should occur
    And loading view "temp" for type "book" should fail

  Scenario: Delete a non-existent view returns error
    Given a vault is ready
    When I delete view "missing" for type "book"
    Then an error should occur

  # ── DefaultView ─────────────────────────────────────────────

  Scenario: Default view when no saved default exists
    Given a vault is ready
    When I get the default view for type "book"
    Then the default view name should be "default"
    And the default view layout should be "list"
    And the default view should sort by "name" "asc"

  Scenario: Default view uses saved default when available
    Given a vault is ready
    And a type schema directory "book" with schema content:
      """
      name: book
      properties: []
      """
    And a saved view "default" for type "book" with sort "rating" "desc"
    When I get the default view for type "book"
    Then the default view should sort by "rating" "desc"
