Feature: Filter operators
  Filter operators are type-aware and translate to SQL conditions.

  # ── Operator validation ─────────────────────────────────────

  Scenario: Valid string operator accepted
    When I validate operator "contains" for property type "string"
    Then the operator validation should pass

  Scenario: Invalid string operator rejected
    When I validate operator "gt" for property type "string"
    Then the operator validation should fail

  Scenario: Valid number operator accepted
    When I validate operator "gt" for property type "number"
    Then the operator validation should pass

  Scenario: Valid date operator accepted
    When I validate operator "after" for property type "date"
    Then the operator validation should pass

  Scenario: Valid select operator accepted
    When I validate operator "is" for property type "select"
    Then the operator validation should pass

  Scenario: Invalid select operator rejected
    When I validate operator "contains" for property type "select"
    Then the operator validation should fail

  Scenario: Valid checkbox operator accepted
    When I validate operator "is" for property type "checkbox"
    Then the operator validation should pass

  Scenario: Invalid checkbox operator rejected
    When I validate operator "gt" for property type "checkbox"
    Then the operator validation should fail

  Scenario: Valid relation operator accepted
    When I validate operator "contains" for property type "relation"
    Then the operator validation should pass

  # ── SQL translation ─────────────────────────────────────────

  Scenario: "is" operator generates equality SQL
    When I translate filter property "status" operator "is" value "reading"
    Then the SQL clause should contain "= ?"
    And the SQL args should have 1 value

  Scenario: "contains" operator generates LIKE SQL
    When I translate filter property "author" operator "contains" value "Tolk"
    Then the SQL clause should contain "LIKE ?"
    And the SQL args should have 1 value

  Scenario: "gt" operator generates CAST comparison SQL
    When I translate filter property "rating" operator "gt" value "4"
    Then the SQL clause should contain "CAST("
    And the SQL clause should contain "> ?"

  Scenario: "is_empty" operator generates null check SQL
    When I translate filter property "author" operator "is_empty" value ""
    Then the SQL clause should contain "IS NULL"
    And the SQL args should have 0 values

  Scenario: "after" operator generates date comparison SQL
    When I translate filter property "published" operator "after" value "2025-01-01"
    Then the SQL clause should contain "> ?"
    And the SQL args should have 1 value
