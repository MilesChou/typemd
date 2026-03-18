package core

import (
	"testing"
)

func TestValidateFilterOperator_AllStringOps(t *testing.T) {
	ops := []string{"is", "is_not", "contains", "does_not_contain", "starts_with", "ends_with", "is_empty", "is_not_empty"}
	for _, op := range ops {
		if err := ValidateFilterOperator("string", op); err != nil {
			t.Errorf("string/%s should be valid: %v", op, err)
		}
	}
}

func TestValidateFilterOperator_AllNumberOps(t *testing.T) {
	ops := []string{"eq", "neq", "gt", "gte", "lt", "lte", "is_empty", "is_not_empty"}
	for _, op := range ops {
		if err := ValidateFilterOperator("number", op); err != nil {
			t.Errorf("number/%s should be valid: %v", op, err)
		}
	}
}

func TestValidateFilterOperator_AllDateOps(t *testing.T) {
	ops := []string{"eq", "before", "after", "on_or_before", "on_or_after", "is_empty", "is_not_empty"}
	for _, op := range ops {
		if err := ValidateFilterOperator("date", op); err != nil {
			t.Errorf("date/%s should be valid: %v", op, err)
		}
	}
}

func TestValidateFilterOperator_AllCheckboxOps(t *testing.T) {
	ops := []string{"is", "is_not"}
	for _, op := range ops {
		if err := ValidateFilterOperator("checkbox", op); err != nil {
			t.Errorf("checkbox/%s should be valid: %v", op, err)
		}
	}
}

func TestValidateFilterOperator_InvalidCrossType(t *testing.T) {
	cases := []struct {
		propType string
		operator string
	}{
		{"string", "gt"},
		{"number", "contains"},
		{"checkbox", "gt"},
		{"select", "contains"},
		{"date", "contains"},
	}
	for _, tc := range cases {
		if err := ValidateFilterOperator(tc.propType, tc.operator); err == nil {
			t.Errorf("%s/%s should be invalid", tc.propType, tc.operator)
		}
	}
}

func TestValidateFilterOperator_UnknownType(t *testing.T) {
	if err := ValidateFilterOperator("unknown", "is"); err == nil {
		t.Error("unknown type should fail")
	}
}

func TestFilterRuleToSQL_UnknownOperator(t *testing.T) {
	rule := FilterRule{Property: "x", Operator: "nope", Value: "y"}
	_, _, err := FilterRuleToSQL(rule)
	if err == nil {
		t.Error("unknown operator should fail")
	}
}

func TestFilterRuleToSQL_IsEmpty_NoArgs(t *testing.T) {
	rule := FilterRule{Property: "author", Operator: "is_empty"}
	_, args, err := FilterRuleToSQL(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("is_empty should have 0 args, got %d", len(args))
	}
}

func TestFilterRuleToSQL_Contains_WrapsWithPercent(t *testing.T) {
	rule := FilterRule{Property: "title", Operator: "contains", Value: "Go"}
	_, args, err := FilterRuleToSQL(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 1 || args[0] != "%Go%" {
		t.Errorf("contains should wrap value with %%, got %v", args)
	}
}
