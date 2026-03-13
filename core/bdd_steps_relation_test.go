package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// ── Relation steps ──────────────────────────────────────────────────────────

func (dc *domainContext) aVaultIsReadyWithRelationSchemas() {
	dc.aVaultIsReady()

	bookSchema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
`)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "book.yaml"), bookSchema, 0644)

	personSchema := []byte(`name: person
properties:
  - name: name
    type: string
  - name: books
    type: relation
    target: book
    multiple: true
    bidirectional: true
    inverse: author
`)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "person.yaml"), personSchema, 0644)
}

func (dc *domainContext) iLinkToVia(sourceName, targetName, relation string) {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		dc.lastErr = fmt.Errorf("object %q or %q not found", sourceName, targetName)
		return
	}
	dc.lastErr = dc.vault.LinkObjects(source.ID, relation, target.ID)
}

func (dc *domainContext) iLinkTheFirstBookToTheSecondBookVia(relation string) {
	if dc.prevObject == nil || dc.currentObject == nil {
		dc.lastErr = fmt.Errorf("need at least 2 objects")
		return
	}
	dc.lastErr = dc.vault.LinkObjects(dc.prevObject.ID, relation, dc.currentObject.ID)
}

func (dc *domainContext) thePropertyOfShouldReference(prop, ownerName, targetName string) error {
	owner := dc.objects[ownerName]
	target := dc.objects[targetName]
	if owner == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", ownerName, targetName)
	}
	obj, err := dc.vault.GetObject(owner.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if obj.Properties[prop] != target.ID {
		return fmt.Errorf("%s.%s = %v, want %q", ownerName, prop, obj.Properties[prop], target.ID)
	}
	return nil
}

func (dc *domainContext) thePropertyOfShouldContain(prop, ownerName, targetName string) error {
	owner := dc.objects[ownerName]
	target := dc.objects[targetName]
	if owner == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", ownerName, targetName)
	}
	obj, err := dc.vault.GetObject(owner.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	items, ok := obj.Properties[prop].([]any)
	if !ok {
		return fmt.Errorf("%s.%s type = %T, want []any", ownerName, prop, obj.Properties[prop])
	}
	for _, item := range items {
		if item == target.ID {
			return nil
		}
	}
	return fmt.Errorf("%s.%s does not contain %q", ownerName, prop, target.ID)
}

func (dc *domainContext) thePropertyOfShouldBeEmpty(prop, ownerName string) error {
	owner := dc.objects[ownerName]
	if owner == nil {
		return fmt.Errorf("object %q not found", ownerName)
	}
	obj, err := dc.vault.GetObject(owner.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val := obj.Properties[prop]
	if val != nil {
		return fmt.Errorf("%s.%s = %v, want nil", ownerName, prop, val)
	}
	return nil
}

func (dc *domainContext) unlinkObjects(sourceName, targetName, relation string, both bool) {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		dc.lastErr = fmt.Errorf("object %q or %q not found", sourceName, targetName)
		return
	}
	dc.lastErr = dc.vault.UnlinkObjects(source.ID, relation, target.ID, both)
}

func (dc *domainContext) iUnlinkFromViaWithBothFlag(sourceName, targetName, relation string) {
	dc.unlinkObjects(sourceName, targetName, relation, true)
}

func (dc *domainContext) iUnlinkFromViaWithoutBothFlag(sourceName, targetName, relation string) {
	dc.unlinkObjects(sourceName, targetName, relation, false)
}

func (dc *domainContext) listingRelationsForShouldReturnNEntries(name string, expected int) error {
	obj := dc.objects[name]
	if obj == nil {
		return fmt.Errorf("object %q not found", name)
	}
	rels, err := dc.vault.ListRelations(obj.ID)
	if err != nil {
		return fmt.Errorf("ListRelations error: %v", err)
	}
	if len(rels) != expected {
		return fmt.Errorf("relations count = %d, want %d", len(rels), expected)
	}
	return nil
}

func initRelationSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a vault is ready with relation schemas$`, dc.aVaultIsReadyWithRelationSchemas)
	ctx.Step(`^I link "([^"]*)" to "([^"]*)" via "([^"]*)"$`, dc.iLinkToVia)
	ctx.Step(`^I link the first book to the second book via "([^"]*)"$`, dc.iLinkTheFirstBookToTheSecondBookVia)
	ctx.Step(`^the "([^"]*)" property of "([^"]*)" should reference "([^"]*)"$`, dc.thePropertyOfShouldReference)
	ctx.Step(`^the "([^"]*)" property of "([^"]*)" should contain "([^"]*)"$`, dc.thePropertyOfShouldContain)
	ctx.Step(`^the "([^"]*)" property of "([^"]*)" should be empty$`, dc.thePropertyOfShouldBeEmpty)
	ctx.Step(`^I unlink "([^"]*)" from "([^"]*)" via "([^"]*)" with both flag$`, dc.iUnlinkFromViaWithBothFlag)
	ctx.Step(`^I unlink "([^"]*)" from "([^"]*)" via "([^"]*)" without both flag$`, dc.iUnlinkFromViaWithoutBothFlag)
	ctx.Step(`^listing relations for "([^"]*)" should return (\d+) entries$`, dc.listingRelationsForShouldReturnNEntries)
}
