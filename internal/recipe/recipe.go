// Package recipe defines the Recipe format — the description of how Cauldron
// emulates one external dependency.
//
// A Recipe is not a mock. A mock returns a shape; a Recipe models behaviour:
// what resources exist, how they change, what the provider emits afterwards,
// and how it fails. The format is declarative on purpose, so that the majority
// of Recipes are data a contributor can read and review rather than code.
package recipe

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Recipe is a complete emulation description for one provider.
type Recipe struct {
	Name    string `yaml:"recipe"`
	Version string `yaml:"version"`
	// Capability is what kind of thing this provider is, so a hundred Recipes
	// can be found by what they do rather than by whether you remember the
	// company's name.
	//
	// One word from a fixed list, not a free string. The value of a category
	// is that two people reaching for it independently land on the same one,
	// and a free string gives you "payments", "payment", "billing" and "money"
	// within a month. Adding a word is a deliberate change to the list.
	//
	// Deliberately not part of the Recipe's name. Renaming stripe to
	// payments.stripe would break every configuration and every command
	// anybody has already written, to buy a grouping a field gives for
	// nothing.
	Capability string              `yaml:"capability"`
	Upstream   Upstream            `yaml:"upstream"`
	Auth       Auth                `yaml:"auth"`
	Resources  map[string]Resource `yaml:"resources"`
	Routes     []Route             `yaml:"routes"`
	Webhooks   Webhooks            `yaml:"webhooks"`
	Responses  Responses           `yaml:"responses"`
	Errors     map[string]Error    `yaml:"errors"`
	Fixtures   map[string]Fixture  `yaml:"fixtures"`
	// RequiredHeaders are headers a request must carry, mapped to the error
	// name to raise when one is missing. Forgetting Notion-Version is the
	// classic Notion integration bug, and a fake that does not enforce it lets
	// code ship that fails on the first real call.
	RequiredHeaders map[string]RequiredHeader `yaml:"required_headers"`
	// Conformance is the evidence that this Recipe resembles the real provider.
	Conformance []Case `yaml:"conformance"`
}

// Fixture is a named seed dataset: resource name to a list of records.
type Fixture map[string][]map[string]any

// Load reads and validates a Recipe from a YAML file.
func Load(path string) (*Recipe, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(contents)
}

func Parse(contents []byte) (*Recipe, error) {
	var r Recipe

	decoder := yaml.NewDecoder(strings.NewReader(string(contents)))
	decoder.KnownFields(true)

	if err := decoder.Decode(&r); err != nil {
		return nil, fmt.Errorf("recipe is not valid YAML: %w", err)
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	return &r, nil
}

// ValidationError collects every problem with a Recipe, so an author sees all
// of them at once instead of fixing one and rerunning.
type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	if len(e.Problems) == 1 {
		return "invalid recipe: " + e.Problems[0]
	}

	return fmt.Sprintf("invalid recipe:\n  - %s", strings.Join(e.Problems, "\n  - "))
}
