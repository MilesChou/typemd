package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// ── Pinned property steps ───────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithPropertyHavingPin(typeName, propName string, pin int) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: string\n    pin: %d\n", typeName, propName, pin)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingUniquePins(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: status
    type: string
    pin: 1
  - name: rating
    type: number
    pin: 2
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingDuplicatePins(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: status
    type: string
    pin: 1
  - name: rating
    type: number
    pin: 1
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithSomePropertiesUnpinned(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
  - name: author
    type: string
  - name: status
    type: string
    pin: 1
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func initPinnedSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with property "([^"]*)" having pin (-?\d+)$`, dc.aTypeSchemaWithPropertyHavingPin)
	ctx.Step(`^a type schema "([^"]*)" with properties having unique pins$`, dc.aTypeSchemaWithPropertiesHavingUniquePins)
	ctx.Step(`^a type schema "([^"]*)" with properties having duplicate pins$`, dc.aTypeSchemaWithPropertiesHavingDuplicatePins)
	ctx.Step(`^a type schema "([^"]*)" with some properties unpinned$`, dc.aTypeSchemaWithSomePropertiesUnpinned)
}
