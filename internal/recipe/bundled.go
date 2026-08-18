package recipe

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

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
	// The suggestion belongs on the error rather than at each call site,
	// because every one of them wants it and the one that had it printed all
	// hundred and twenty-seven names instead.
	if close := Suggest(e.Name); len(close) > 0 {
		return fmt.Sprintf("no recipe named %q; did you mean %s?", e.Name, humanList(close))
	}

	return fmt.Sprintf("no recipe named %q; run 'cauldron recipe list' to see what there is", e.Name)
}

// humanList renders a short list the way a sentence would.
func humanList(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return items[0] + " or " + items[1]
	default:
		return strings.Join(items[:len(items)-1], ", ") + " or " + items[len(items)-1]
	}
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
