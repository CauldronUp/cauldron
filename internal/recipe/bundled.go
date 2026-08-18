package recipe

import (
	"fmt"
	"io/fs"
	"path"
	"sort"

	"github.com/CauldronUp/cauldron/recipes"
)

// bundled holds the first-party Recipes, compiled into the binary by the
// recipes package. Cauldron ships as a single file with no runtime assets to
// locate, so `cauldron recipe info stripe` works on a fresh machine with no
// network.
var bundled = recipes.FS

const bundledRoot = "."

// Bundled returns the names of every Recipe compiled into this binary.
func Bundled() []string {
	entries, err := bundled.ReadDir(bundledRoot)
	if err != nil {
		return nil
	}

	names := make([]string, 0, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if _, err := fs.Stat(bundled, path.Join(bundledRoot, entry.Name(), "recipe.yaml")); err != nil {
			continue
		}

		names = append(names, entry.Name())
	}

	sort.Strings(names)

	return names
}

// ErrNotBundled is returned when a Recipe name is not compiled in.
type ErrNotBundled struct {
	Name string
}

func (e *ErrNotBundled) Error() string {
	return fmt.Sprintf("no recipe named %q", e.Name)
}

// Open loads a bundled Recipe by name.
func Open(name string) (*Recipe, error) {
	contents, err := bundled.ReadFile(path.Join(bundledRoot, name, "recipe.yaml"))
	if err != nil {
		return nil, &ErrNotBundled{Name: name}
	}

	return Parse(contents)
}

// Summary is the condensed view `cauldron recipe list` prints.
type Summary struct {
	Name       string
	Capability string
	Version    string
	API        string
	Resources  int
	Routes     int
	Events     int
}

// Summarise loads every bundled Recipe and reduces it to a listing row.
func Summarise() ([]Summary, error) {
	names := Bundled()
	out := make([]Summary, 0, len(names))

	for _, name := range names {
		r, err := Open(name)
		if err != nil {
			return nil, err
		}

		out = append(out, Summary{
			Name:       r.Name,
			Capability: r.Capability,
			Version:    r.Version,
			API:        r.Upstream.API,
			Resources:  len(r.Resources),
			Routes:     len(r.Routes),
			Events:     len(r.Webhooks.Events),
		})
	}

	return out, nil
}
