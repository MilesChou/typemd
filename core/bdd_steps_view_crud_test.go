package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
	"gopkg.in/yaml.v3"
)

// ── View CRUD step state ─────────────────────────────────────────────────────

type viewCrudContext struct {
	dc          *domainContext
	vc          *viewConfigContext // shares the view config context
	views       []ViewConfig
	loadedView  *ViewConfig
	defaultView *ViewConfig
}

func newViewCrudContext(dc *domainContext, vc *viewConfigContext) *viewCrudContext {
	return &viewCrudContext{dc: dc, vc: vc}
}

// ── Given steps ─────────────────────────────────────────────────────────────

func (vcc *viewCrudContext) aSavedViewForTypeWithSort(viewName, typeName, sortProp, sortDir string) {
	view := ViewConfig{
		Name:   viewName,
		Layout: ViewLayoutList,
		Sort:   []SortRule{{Property: sortProp, Direction: sortDir}},
	}
	data, _ := yaml.Marshal(&view)
	dir := filepath.Join(vcc.dc.vault.TypesDir(), typeName, "views")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, viewName+".yaml"), data, 0644)
}

// ── When steps ──────────────────────────────────────────────────────────────

func (vcc *viewCrudContext) iListViewsForType(typeName string) {
	views, err := vcc.dc.vault.ListViews(typeName)
	vcc.views = views
	vcc.dc.lastErr = err
}

func (vcc *viewCrudContext) iLoadViewForType(viewName, typeName string) {
	view, err := vcc.dc.vault.LoadView(typeName, viewName)
	vcc.loadedView = view
	vcc.dc.lastErr = err
}

func (vcc *viewCrudContext) iSaveViewForType(typeName string) {
	vcc.dc.lastErr = vcc.dc.vault.SaveView(typeName, vcc.vc.view)
}

func (vcc *viewCrudContext) iDeleteViewForType(viewName, typeName string) {
	vcc.dc.lastErr = vcc.dc.vault.DeleteView(typeName, viewName)
}

func (vcc *viewCrudContext) iGetTheDefaultViewForType(typeName string) {
	vcc.defaultView = vcc.dc.vault.DefaultView(typeName)
}

// ── Then steps ──────────────────────────────────────────────────────────────

func (vcc *viewCrudContext) iShouldHaveNViews(n int) error {
	if len(vcc.views) != n {
		return fmt.Errorf("expected %d views, got %d", n, len(vcc.views))
	}
	return nil
}

func (vcc *viewCrudContext) theLoadedViewNameShouldBe(expected string) error {
	if vcc.loadedView == nil {
		return fmt.Errorf("no view loaded")
	}
	if vcc.loadedView.Name != expected {
		return fmt.Errorf("expected loaded view name %q, got %q", expected, vcc.loadedView.Name)
	}
	return nil
}

func (vcc *viewCrudContext) loadingViewForTypeShouldSucceed(viewName, typeName string) error {
	_, err := vcc.dc.vault.LoadView(typeName, viewName)
	if err != nil {
		return fmt.Errorf("expected LoadView(%q, %q) to succeed, got %v", typeName, viewName, err)
	}
	return nil
}

func (vcc *viewCrudContext) loadingViewForTypeShouldFail(viewName, typeName string) error {
	_, err := vcc.dc.vault.LoadView(typeName, viewName)
	if err == nil {
		return fmt.Errorf("expected LoadView(%q, %q) to fail, got nil", typeName, viewName)
	}
	return nil
}

func (vcc *viewCrudContext) theDefaultViewNameShouldBe(expected string) error {
	if vcc.defaultView.Name != expected {
		return fmt.Errorf("expected default view name %q, got %q", expected, vcc.defaultView.Name)
	}
	return nil
}

func (vcc *viewCrudContext) theDefaultViewLayoutShouldBe(expected string) error {
	if string(vcc.defaultView.Layout) != expected {
		return fmt.Errorf("expected default view layout %q, got %q", expected, vcc.defaultView.Layout)
	}
	return nil
}

func (vcc *viewCrudContext) theDefaultViewShouldSortBy(property, direction string) error {
	if len(vcc.defaultView.Sort) == 0 {
		return fmt.Errorf("default view has no sort rules")
	}
	s := vcc.defaultView.Sort[0]
	if s.Property != property || s.Direction != direction {
		return fmt.Errorf("expected sort by %s %s, got %s %s", property, direction, s.Property, s.Direction)
	}
	return nil
}

// ── Init ────────────────────────────────────────────────────────────────────

func initViewCrudSteps(ctx *godog.ScenarioContext, dc *domainContext, vc *viewConfigContext) {
	vcc := newViewCrudContext(dc, vc)

	// Given
	ctx.Step(`^a saved view "([^"]*)" for type "([^"]*)" with sort "([^"]*)" "([^"]*)"$`, vcc.aSavedViewForTypeWithSort)

	// When
	ctx.Step(`^I list views for type "([^"]*)"$`, vcc.iListViewsForType)
	ctx.Step(`^I load view "([^"]*)" for type "([^"]*)"$`, vcc.iLoadViewForType)
	ctx.Step(`^I save view for type "([^"]*)"$`, vcc.iSaveViewForType)
	ctx.Step(`^I delete view "([^"]*)" for type "([^"]*)"$`, vcc.iDeleteViewForType)
	ctx.Step(`^I get the default view for type "([^"]*)"$`, vcc.iGetTheDefaultViewForType)

	// Then
	ctx.Step(`^I should have (\d+) views?$`, vcc.iShouldHaveNViews)
	ctx.Step(`^the loaded view name should be "([^"]*)"$`, vcc.theLoadedViewNameShouldBe)
	ctx.Step(`^loading view "([^"]*)" for type "([^"]*)" should succeed$`, vcc.loadingViewForTypeShouldSucceed)
	ctx.Step(`^loading view "([^"]*)" for type "([^"]*)" should fail$`, vcc.loadingViewForTypeShouldFail)
	ctx.Step(`^the default view name should be "([^"]*)"$`, vcc.theDefaultViewNameShouldBe)
	ctx.Step(`^the default view layout should be "([^"]*)"$`, vcc.theDefaultViewLayoutShouldBe)
	ctx.Step(`^the default view should sort by "([^"]*)" "([^"]*)"$`, vcc.theDefaultViewShouldSortBy)
}
