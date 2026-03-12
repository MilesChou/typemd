package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

// ── Name property steps ─────────────────────────────────────────────────────

func (dc *domainContext) iSetTheObjectNameTo(name string) {
	dc.currentObject.Properties["name"] = name
}

func (dc *domainContext) iRemoveTheNamePropertyFromTheObject() {
	delete(dc.currentObject.Properties, "name")
}

func (dc *domainContext) getNameShouldReturn(expected string) error {
	got := dc.currentObject.GetName()
	if got != expected {
		return fmt.Errorf("GetName() = %q, want %q", got, expected)
	}
	return nil
}

func (dc *domainContext) getNameShouldReturnTheDisplayName() error {
	got := dc.currentObject.GetName()
	expected := dc.currentObject.DisplayName()
	if got != expected {
		return fmt.Errorf("GetName() = %q, want DisplayName() = %q", got, expected)
	}
	return nil
}

func (dc *domainContext) theSyncedObjectShouldHaveNameMatchingItsDisplayName() error {
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("get object: %v", err)
	}
	name, _ := obj.Properties["name"].(string)
	expected := obj.DisplayName()
	if name != expected {
		return fmt.Errorf("synced name = %q, want DisplayName() = %q", name, expected)
	}
	return nil
}

func (dc *domainContext) theSyncedObjectShouldHaveName(expected string) error {
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("get object: %v", err)
	}
	name, _ := obj.Properties["name"].(string)
	if name != expected {
		return fmt.Errorf("synced name = %q, want %q", name, expected)
	}
	return nil
}

func initNameSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^I set the object name to "([^"]*)"$`, dc.iSetTheObjectNameTo)
	ctx.Step(`^I remove the name property from the object$`, dc.iRemoveTheNamePropertyFromTheObject)
	ctx.Step(`^GetName should return "([^"]*)"$`, dc.getNameShouldReturn)
	ctx.Step(`^GetName should return the DisplayName$`, dc.getNameShouldReturnTheDisplayName)
	ctx.Step(`^the synced object should have name matching its DisplayName$`, dc.theSyncedObjectShouldHaveNameMatchingItsDisplayName)
	ctx.Step(`^the synced object should have name "([^"]*)"$`, dc.theSyncedObjectShouldHaveName)
}
