package core

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// ── Resolve steps ───────────────────────────────────────────────────────────

func (dc *domainContext) iResolveTheObjectByItsFullID() {
	if dc.currentObject == nil {
		dc.lastErr = fmt.Errorf("no current object")
		return
	}
	dc.resolvedID, dc.lastErr = dc.vault.ResolveID(dc.currentObject.ID)
}

func (dc *domainContext) theResolvedIDShouldMatchTheOriginal() error {
	if dc.resolvedID != dc.currentObject.ID {
		return fmt.Errorf("resolved ID = %q, want %q", dc.resolvedID, dc.currentObject.ID)
	}
	return nil
}

func (dc *domainContext) iResolveTheObjectByPrefix(prefix string) {
	obj, err := dc.vault.ResolveObject(prefix)
	dc.lastErr = err
	if err == nil {
		dc.retrieved = obj
	}
}

func (dc *domainContext) theResolvedObjectShouldMatchTheCreatedOne() error {
	if dc.retrieved == nil {
		return fmt.Errorf("no resolved object")
	}
	if dc.currentObject == nil {
		return fmt.Errorf("no current object to compare")
	}
	if dc.retrieved.ID != dc.currentObject.ID {
		return fmt.Errorf("resolved ID = %q, want %q", dc.retrieved.ID, dc.currentObject.ID)
	}
	return nil
}

func (dc *domainContext) iResolveTheObjectByAPartialULIDPrefix() {
	if dc.currentObject == nil {
		dc.lastErr = fmt.Errorf("no current object")
		return
	}
	// Use type + display name + first 4 chars of ULID as partial prefix
	displayName := dc.currentObject.DisplayName()
	ulidPart := strings.TrimPrefix(dc.currentObject.Filename, displayName+"-")
	partial := ulidPart[:4]
	prefix := dc.currentObject.Type + "/" + displayName + "-" + partial
	obj, err := dc.vault.ResolveObject(prefix)
	dc.lastErr = err
	if err == nil {
		dc.retrieved = obj
	}
}

func (dc *domainContext) anAmbiguousMatchErrorShouldOccurWithNCandidates(expected int) error {
	if dc.lastErr == nil {
		return fmt.Errorf("expected AmbiguousMatchError, got nil")
	}
	ambErr, ok := dc.lastErr.(*AmbiguousMatchError)
	if !ok {
		return fmt.Errorf("expected *AmbiguousMatchError, got %T: %v", dc.lastErr, dc.lastErr)
	}
	if len(ambErr.Matches) != expected {
		return fmt.Errorf("candidates = %d, want %d", len(ambErr.Matches), expected)
	}
	return nil
}

func initResolveSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^I resolve the object by its full ID$`, dc.iResolveTheObjectByItsFullID)
	ctx.Step(`^the resolved ID should match the original$`, dc.theResolvedIDShouldMatchTheOriginal)
	ctx.Step(`^I resolve the object by prefix "([^"]*)"$`, dc.iResolveTheObjectByPrefix)
	ctx.Step(`^the resolved object should match the created one$`, dc.theResolvedObjectShouldMatchTheCreatedOne)
	ctx.Step(`^I resolve the object by a partial ULID prefix$`, dc.iResolveTheObjectByAPartialULIDPrefix)
	ctx.Step(`^an ambiguous match error should occur with (\d+) candidates$`, dc.anAmbiguousMatchErrorShouldOccurWithNCandidates)
}
