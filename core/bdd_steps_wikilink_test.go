package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// ── Wiki-link steps ─────────────────────────────────────────────────────────

func (dc *domainContext) aVaultIsReadyWithNoteSchemas() {
	dc.aVaultIsReady()
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "book.yaml"),
		[]byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "person.yaml"),
		[]byte("name: person\nproperties:\n  - name: name\n    type: string\n"), 0644)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)
}

func (dc *domainContext) bodyContainsAWikiLinkTo(sourceName, targetName string) {
	source := dc.objects[sourceName]
	if source == nil {
		panic(fmt.Sprintf("source object %q not found", sourceName))
	}
	// If target is a known object slug, use its ID; otherwise treat as raw ID
	targetID := targetName
	if target, ok := dc.objects[targetName]; ok {
		targetID = target.ID
	}
	body := fmt.Sprintf("---\ntitle: %s\n---\n\nSee [[%s]].\n", sourceName, targetID)
	os.WriteFile(dc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (dc *domainContext) bodyContainsAWikiLinkToWithDisplayText(sourceName, targetName, displayText string) {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		panic(fmt.Sprintf("object %q or %q not found", sourceName, targetName))
	}
	body := fmt.Sprintf("---\ntitle: %s\n---\n\nBy [[%s|%s]].\n", sourceName, target.ID, displayText)
	os.WriteFile(dc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (dc *domainContext) iSyncTheIndex() {
	_, err := dc.vault.SyncIndex()
	dc.lastErr = err
}

func (dc *domainContext) shouldHaveNWikiLinks(name string, expected int) error {
	obj := dc.objects[name]
	if obj == nil {
		return fmt.Errorf("object %q not found", name)
	}
	links, err := dc.vault.ListWikiLinks(obj.ID)
	if err != nil {
		return fmt.Errorf("ListWikiLinks error: %v", err)
	}
	dc.wikiLinks = links
	if len(links) != expected {
		return fmt.Errorf("wiki-links = %d, want %d", len(links), expected)
	}
	return nil
}

func (dc *domainContext) theWikiLinkTargetShouldBe(targetName string) error {
	target := dc.objects[targetName]
	if target == nil {
		return fmt.Errorf("object %q not found", targetName)
	}
	if len(dc.wikiLinks) == 0 {
		return fmt.Errorf("no wiki-links to check")
	}
	if dc.wikiLinks[0].ToID != target.ID {
		return fmt.Errorf("wiki-link ToID = %q, want %q", dc.wikiLinks[0].ToID, target.ID)
	}
	return nil
}

func (dc *domainContext) shouldHaveNBacklinksFrom(targetName string, expected int, sourceName string) error {
	target := dc.objects[targetName]
	if target == nil {
		return fmt.Errorf("object %q not found", targetName)
	}
	backlinks, err := dc.vault.ListBacklinks(target.ID)
	if err != nil {
		return fmt.Errorf("ListBacklinks error: %v", err)
	}
	if len(backlinks) != expected {
		return fmt.Errorf("backlinks = %d, want %d", len(backlinks), expected)
	}
	source := dc.objects[sourceName]
	if source == nil {
		return fmt.Errorf("source object %q not found", sourceName)
	}
	found := false
	for _, bl := range backlinks {
		if bl.FromID == source.ID {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("no backlink from %q found", sourceName)
	}
	return nil
}

func (dc *domainContext) theWikiLinkShouldHaveAnEmptyResolvedID() error {
	if len(dc.wikiLinks) == 0 {
		return fmt.Errorf("no wiki-links to check")
	}
	if dc.wikiLinks[0].ToID != "" {
		return fmt.Errorf("ToID = %q, want empty", dc.wikiLinks[0].ToID)
	}
	return nil
}

func (dc *domainContext) iChangeWikiLinkTo(sourceName, targetName string) {
	dc.bodyContainsAWikiLinkTo(sourceName, targetName)
}

func (dc *domainContext) wikiLinkShouldPointTo(sourceName, targetName string) error {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", sourceName, targetName)
	}
	links, err := dc.vault.ListWikiLinks(source.ID)
	if err != nil {
		return fmt.Errorf("ListWikiLinks error: %v", err)
	}
	if len(links) != 1 {
		return fmt.Errorf("wiki-links = %d, want 1", len(links))
	}
	if links[0].ToID != target.ID {
		return fmt.Errorf("wiki-link ToID = %q, want %q", links[0].ToID, target.ID)
	}
	return nil
}

func (dc *domainContext) shouldHaveNBacklinks(targetName string, expected int) error {
	target := dc.objects[targetName]
	if target == nil {
		return fmt.Errorf("object %q not found", targetName)
	}
	backlinks, err := dc.vault.ListBacklinks(target.ID)
	if err != nil {
		return fmt.Errorf("ListBacklinks error: %v", err)
	}
	if len(backlinks) != expected {
		return fmt.Errorf("backlinks = %d, want %d", len(backlinks), expected)
	}
	return nil
}

func (dc *domainContext) theWikiLinkDisplayTextShouldBe(expected string) error {
	if len(dc.wikiLinks) == 0 {
		return fmt.Errorf("no wiki-links to check")
	}
	if dc.wikiLinks[0].DisplayText != expected {
		return fmt.Errorf("DisplayText = %q, want %q", dc.wikiLinks[0].DisplayText, expected)
	}
	return nil
}

func initWikiLinkSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a vault is ready with note schemas$`, dc.aVaultIsReadyWithNoteSchemas)
	ctx.Step(`^"([^"]*)" body contains a wiki-link to "([^"]*)"$`, dc.bodyContainsAWikiLinkTo)
	ctx.Step(`^"([^"]*)" body contains a wiki-link to "([^"]*)" with display text "([^"]*)"$`, dc.bodyContainsAWikiLinkToWithDisplayText)
	ctx.Step(`^I sync the index$`, dc.iSyncTheIndex)
	ctx.Step(`^"([^"]*)" should have (\d+) wiki-links?$`, dc.shouldHaveNWikiLinks)
	ctx.Step(`^the wiki-link target should be "([^"]*)"$`, dc.theWikiLinkTargetShouldBe)
	ctx.Step(`^"([^"]*)" should have (\d+) backlinks? from "([^"]*)"$`, dc.shouldHaveNBacklinksFrom)
	ctx.Step(`^the wiki-link should have an empty resolved ID$`, dc.theWikiLinkShouldHaveAnEmptyResolvedID)
	ctx.Step(`^I change "([^"]*)" wiki-link to "([^"]*)"$`, dc.iChangeWikiLinkTo)
	ctx.Step(`^"([^"]*)" wiki-link should point to "([^"]*)"$`, dc.wikiLinkShouldPointTo)
	ctx.Step(`^"([^"]*)" should have (\d+) backlinks$`, dc.shouldHaveNBacklinks)
	ctx.Step(`^the wiki-link display text should be "([^"]*)"$`, dc.theWikiLinkDisplayTextShouldBe)
}
