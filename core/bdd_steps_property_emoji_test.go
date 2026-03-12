package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// ── Property emoji steps ────────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithPropertyHavingEmoji(typeName, propName, emoji string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: string\n    emoji: %s\n", typeName, propName, emoji)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingUniqueEmojis(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
    emoji: 📖
  - name: rating
    type: number
    emoji: ⭐
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingDuplicateEmojis(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
    emoji: 👤
  - name: author
    type: string
    emoji: 👤
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithSomePropertiesMissingEmojis(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
  - name: author
    type: string
  - name: rating
    type: number
    emoji: ⭐
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func initPropertyEmojiSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with property "([^"]*)" having emoji "([^"]*)"$`, dc.aTypeSchemaWithPropertyHavingEmoji)
	ctx.Step(`^a type schema "([^"]*)" with properties having unique emojis$`, dc.aTypeSchemaWithPropertiesHavingUniqueEmojis)
	ctx.Step(`^a type schema "([^"]*)" with properties having duplicate emojis$`, dc.aTypeSchemaWithPropertiesHavingDuplicateEmojis)
	ctx.Step(`^a type schema "([^"]*)" with some properties missing emojis$`, dc.aTypeSchemaWithSomePropertiesMissingEmojis)
}
