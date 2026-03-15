Feature: Object relations
  Objects can be linked through typed relations defined in type schemas.
  Relations support bidirectional linking with inverse properties.

  Background:
    Given a vault is ready with relation schemas

  Scenario: Link objects with bidirectional relation
    Given a "book" object named "golang-in-action" exists
    And a "person" object named "alan-donovan" exists
    When I link "golang-in-action" to "alan-donovan" via "author"
    Then the "author" property of "golang-in-action" should reference "alan-donovan"
    And the "books" property of "alan-donovan" should contain "golang-in-action"

  Scenario: Link with type mismatch fails
    Given a "book" object named "book-a" exists
    And a "book" object named "book-b" exists
    When I link "book-a" to "book-b" via "author"
    Then an error should occur

  Scenario: Overwrite single-value relation
    Given a "book" object named "test" exists
    And a "person" object named "alan" exists
    And a "person" object named "brian" exists
    When I link "test" to "alan" via "author"
    And I link "test" to "brian" via "author"
    Then the "author" property of "test" should reference "brian"

  Scenario: Unlink both directions
    Given a "book" object named "test" exists
    And a "person" object named "alan" exists
    And I link "test" to "alan" via "author"
    When I unlink "test" from "alan" via "author" with both flag
    Then the "author" property of "test" should be empty
    And the "books" property of "alan" should be empty

  Scenario: List relations for a linked object
    Given a "book" object named "test" exists
    And a "person" object named "alan" exists
    When I link "test" to "alan" via "author"
    Then listing relations for "test" should return 2 entries

  Scenario: List relations for an unlinked object
    Given a "book" object named "test" exists
    Then listing relations for "test" should return 0 entries

  Scenario: Link object to tag via system property
    Given a "tag" object named "go" exists
    And a "book" object named "golang-book" exists
    When I link "golang-book" to "go" via "tags"
    Then the "tags" property of "golang-book" should contain "go"

  Scenario: Unlink tag from object via system property
    Given a "tag" object named "go" exists
    And a "book" object named "golang-book" exists
    And I link "golang-book" to "go" via "tags"
    When I unlink "golang-book" from "go" via "tags" without both flag
    Then the "tags" property of "golang-book" should be empty

  Scenario: Relation property without target is invalid
    Given a type schema "article" with a relation property missing target
    When I validate all schemas
    Then schema "article" should have errors

  Scenario: Link to non-existent object fails
    Given a "book" object named "lonely" exists
    When I link "lonely" to a non-existent object via "author"
    Then an error should occur

  Scenario: Link with unknown relation name fails
    Given a "book" object named "test-book" exists
    And a "person" object named "test-person" exists
    When I link "test-book" to "test-person" via "nonexistent"
    Then an error should occur

  Scenario: Append to multiple-value relation
    Given a "book" object named "book-x" exists
    And a "book" object named "book-y" exists
    And a "person" object named "writer" exists
    When I link "book-x" to "writer" via "author"
    And I link "book-y" to "writer" via "author"
    Then the "books" property of "writer" should have 2 items

  Scenario: Duplicate link is rejected
    Given a "person" object named "dup-person" exists
    And a "book" object named "dup-book" exists
    And I link "dup-person" to "dup-book" via "books"
    When I link "dup-person" to "dup-book" via "books"
    Then an error should occur

  Scenario: Inverse property must exist in target schema
    Given a type schema "article" with bidirectional relation to missing inverse
    And an "article" object named "test-article" exists
    And a "person" object named "test-target" exists
    When I link "test-article" to "test-target" via "reviewer"
    Then an error should occur

  Scenario: Unlink without both flag leaves inverse intact
    Given a "book" object named "intact-book" exists
    And a "person" object named "intact-person" exists
    And I link "intact-book" to "intact-person" via "author"
    When I unlink "intact-book" from "intact-person" via "author" without both flag
    Then the "author" property of "intact-book" should be empty
    And the "books" property of "intact-person" should contain "intact-book"

  Scenario: Unlink one from multiple-value relation
    Given a "book" object named "multi-a" exists
    And a "book" object named "multi-b" exists
    And a "person" object named "multi-person" exists
    And I link "multi-a" to "multi-person" via "author"
    And I link "multi-b" to "multi-person" via "author"
    When I unlink "multi-a" from "multi-person" via "author" with both flag
    Then the "books" property of "multi-person" should contain "multi-b"
    And the "books" property of "multi-person" should have 1 items

  Scenario: Reverse relation appears in display properties
    Given a "book" object named "display-book" exists
    And a "person" object named "display-person" exists
    And I link "display-book" to "display-person" via "author"
    When I build display properties for "display-person"
    Then the display properties should contain a reverse relation with indicator "←"
