Feature: View configuration
  ViewConfig defines how objects of a type are filtered, sorted, grouped, and displayed.

  # ── ViewConfig creation ─────────────────────────────────────

  Scenario: Create a ViewConfig with all fields
    Given a view config "by-rating" with layout "list"
    And the view has filter property "status" operator "is" value "reading"
    And the view has sort property "rating" direction "desc"
    And the view has group_by "genre"
    Then the view name should be "by-rating"
    And the view layout should be "list"
    And the view should have 1 filter rule
    And the view should have 1 sort rule
    And the view group_by should be "genre"

  Scenario: Create a minimal ViewConfig
    Given a view config "default" with layout "list"
    Then the view name should be "default"
    And the view should have 0 filter rules
    And the view should have 0 sort rules
    And the view group_by should be ""

  # ── YAML serialization ──────────────────────────────────────

  Scenario: Serialize ViewConfig to YAML
    Given a view config "by-rating" with layout "list"
    And the view has sort property "rating" direction "desc"
    When I serialize the view config to YAML
    Then the view YAML should contain "name: by-rating"
    And the view YAML should contain "layout: list"
    And the view YAML should contain "property: rating"
    And the view YAML should not contain "filter:"
    And the view YAML should not contain "group_by:"

  Scenario: Deserialize ViewConfig from YAML
    Given view YAML content:
      """
      name: reading-now
      layout: list
      filter:
        - property: status
          operator: is
          value: reading
      sort:
        - property: name
          direction: asc
      group_by: genre
      """
    When I deserialize the view YAML
    Then the deserialized view name should be "reading-now"
    And the deserialized view should have 1 filter rule
    And the deserialized view should have 1 sort rule
    And the deserialized view group_by should be "genre"
