// Package project reads and writes the file a project uses to say which
// Recipes it wants.
//
// Detection covers the common case and cannot cover all of it. A project that
// talks to a provider over raw HTTP has no dependency to find, and so does one
// using any of the two dozen Recipes no package maps to yet. Without somewhere
// to write it down, the answer was to retype the names on every command, which
// is the kind of friction that ends with somebody not using the tool.
package project

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileName is the file a project keeps its answer in.
const FileName = "cauldron.yaml"

// Config is what that file holds.
type Config struct {
	// Recipes is the whole list, not an addition to what detection found.
	//
	// Once the file exists it is the answer, because a list that only added
	// to detection could never take anything away, and the surprise being
	// fixed is usually something detection got wrong rather than something
	// it missed. Writing the file for the first time copies what detection
	// found into it, so nothing is lost by starting one.
	Recipes []string `yaml:"recipes"`
}

// Path is where the file lives for a project directory.
func Path(dir string) string {
	return filepath.Join(dir, FileName)
}

// Load reads the config. A missing file is not an error: it means the project
// has not said anything, and detection answers instead.
func Load(dir string) (*Config, bool, error) {
	contents, err := os.ReadFile(Path(dir))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}

		return nil, false, err
	}

	var config Config

	decoder := yaml.NewDecoder(strings.NewReader(string(contents)))
	decoder.KnownFields(true)

	if err := decoder.Decode(&config); err != nil {
		return nil, false, fmt.Errorf("%s is not valid: %w", FileName, err)
	}

	return &config, true, nil
}

// Save writes the config, sorted, so a file edited by hand and a file written
// by the tool look the same and a diff shows what changed rather than where
// the tool happened to append.
func Save(dir string, config *Config) error {
	sort.Strings(config.Recipes)

	var out strings.Builder

	out.WriteString("# Which providers this project talks to.\n")
	out.WriteString("#\n")
	out.WriteString("# Cauldron reads this instead of guessing from the manifests once the\n")
	out.WriteString("# file exists, so a provider it cannot detect belongs here, and one it\n")
	out.WriteString("# detects wrongly can be removed.\n")
	out.WriteString("#\n")
	out.WriteString("# cauldron add <name>     adds one\n")
	out.WriteString("# cauldron recipe list    shows what there is\n")
	out.WriteString("\nrecipes:\n")

	for _, name := range config.Recipes {
		out.WriteString("  - " + name + "\n")
	}

	return os.WriteFile(Path(dir), []byte(out.String()), 0o644)
}

// Has reports whether a recipe is already listed.
func (c *Config) Has(name string) bool {
	for _, held := range c.Recipes {
		if held == name {
			return true
		}
	}

	return false
}
