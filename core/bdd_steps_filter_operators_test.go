package core

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// ── Filter operator step state ───────────────────────────────────────────────

type filterOperatorContext struct {
	dc            *domainContext
	validationErr error
	sqlClause     string
	sqlArgs       []any
	translateErr  error
}

func newFilterOperatorContext(dc *domainContext) *filterOperatorContext {
	return &filterOperatorContext{dc: dc}
}

// ── When steps ──────────────────────────────────────────────────────────────

func (fo *filterOperatorContext) iValidateOperatorForPropertyType(operator, propType string) {
	fo.validationErr = ValidateFilterOperator(propType, operator)
}

func (fo *filterOperatorContext) iTranslateFilterPropertyOperatorValue(property, operator, value string) {
	rule := FilterRule{Property: property, Operator: operator, Value: value}
	fo.sqlClause, fo.sqlArgs, fo.translateErr = FilterRuleToSQL(rule)
}

// ── Then steps ──────────────────────────────────────────────────────────────

func (fo *filterOperatorContext) theOperatorValidationShouldPass() error {
	if fo.validationErr != nil {
		return fmt.Errorf("expected validation to pass, got %v", fo.validationErr)
	}
	return nil
}

func (fo *filterOperatorContext) theOperatorValidationShouldFail() error {
	if fo.validationErr == nil {
		return fmt.Errorf("expected validation to fail, got nil")
	}
	return nil
}

func (fo *filterOperatorContext) theSQLClauseShouldContain(substr string) error {
	if fo.translateErr != nil {
		return fmt.Errorf("translation error: %v", fo.translateErr)
	}
	if !strings.Contains(fo.sqlClause, substr) {
		return fmt.Errorf("expected SQL clause to contain %q, got %q", substr, fo.sqlClause)
	}
	return nil
}

func (fo *filterOperatorContext) theSQLArgsShouldHaveNValues(n int) error {
	if fo.translateErr != nil {
		return fmt.Errorf("translation error: %v", fo.translateErr)
	}
	if len(fo.sqlArgs) != n {
		return fmt.Errorf("expected %d SQL args, got %d", n, len(fo.sqlArgs))
	}
	return nil
}

// ── Init ────────────────────────────────────────────────────────────────────

func initFilterOperatorSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	fo := newFilterOperatorContext(dc)

	// When
	ctx.Step(`^I validate operator "([^"]*)" for property type "([^"]*)"$`, fo.iValidateOperatorForPropertyType)
	ctx.Step(`^I translate filter property "([^"]*)" operator "([^"]*)" value "([^"]*)"$`, fo.iTranslateFilterPropertyOperatorValue)

	// Then
	ctx.Step(`^the operator validation should pass$`, fo.theOperatorValidationShouldPass)
	ctx.Step(`^the operator validation should fail$`, fo.theOperatorValidationShouldFail)
	ctx.Step(`^the SQL clause should contain "([^"]*)"$`, fo.theSQLClauseShouldContain)
	ctx.Step(`^the SQL args should have (\d+) values?$`, fo.theSQLArgsShouldHaveNValues)
}
