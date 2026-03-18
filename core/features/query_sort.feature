Feature: Query with sort
  QueryService supports sorting results by property values.

  Scenario: Query with sort by name ascending
    Given a vault is ready
    And a "book" object named "zebra-book" exists
    And a "book" object named "alpha-book" exists
    When I query with filter "type=book" sorted by "name" "asc"
    Then the sorted results should have 2 objects
    And the first sorted result name should come before the second alphabetically

  Scenario: Query with sort by created_at descending
    Given a vault is ready
    And a "book" object named "old-book" exists
    And a "book" object named "new-book" exists
    When I query with filter "type=book" sorted by "created_at" "desc"
    Then the sorted results should have 2 objects

  Scenario: Query without sort returns results
    Given a vault is ready
    And a "book" object named "test-book" exists
    When I query with filter "type=book" and no sort
    Then the sorted results should have 1 object
