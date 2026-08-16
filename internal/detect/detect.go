// Package detect inspects a project directory and works out what it needs to
// run: language runtimes, backing services, and which third-party APIs it talks
// to. Everything is inferred from files that are already in the repository —
// Cauldron never requires its own config file.
package detect

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// Kind classifies a detected requirement.
type Kind string

const (
	// KindRuntime is a language runtime, e.g. PHP 8.5 or Node 24.
	KindRuntime Kind = "runtime"
	// KindService is a backing service Cauldron boots itself, e.g. PostgreSQL.
	KindService Kind = "service"
	// KindRecipe is a third-party API Cauldron emulates, e.g. Stripe.
	KindRecipe Kind = "recipe"
)

// Requirement is a single thing the project needs.
type Requirement struct {
	Kind Kind
	// Name is the canonical identifier: "php", "postgres", "stripe".
	Name string
	// Version is the resolved or requested version, when one is knowable.
	Version string
	// Source records which file produced this requirement, so `cauldron doctor`
	// can explain itself rather than being a black box.
	Source string
	// Evidence is the specific token that triggered detection, e.g. the package
	// name "stripe/stripe-php".
	Evidence string
}

// Project is the result of inspecting a directory.
type Project struct {
	Root         string
	Framework    string
	Requirements []Requirement
	// Unmatched holds dependencies that look like third-party API clients but
	// have no Recipe yet. Surfacing these is deliberate: silently falling back
	// to the real network is how a test suite starts lying.
	Unmatched []Requirement
}

// ErrNoProject is returned when a directory contains nothing Cauldron can read.
var ErrNoProject = errors.New("no recognisable project manifest found")

// detector reads one manifest format.
type detector interface {
	// file is the manifest this detector consumes, relative to the root.
	file() string
	// detect parses the manifest contents and appends to the project.
	detect(contents []byte, p *Project) error
}

func detectors() []detector {
	return []detector{
		composerDetector{},
		npmDetector{},
		goModDetector{},
	}
}

// Detect inspects root and returns everything it can determine.
func Detect(root string) (*Project, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fs.ErrInvalid
	}

	p := &Project{Root: root}
	found := false

	for _, d := range detectors() {
		path := filepath.Join(root, d.file())

		contents, err := os.ReadFile(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}

		found = true

		if err := d.detect(contents, p); err != nil {
			return nil, err
		}
	}

	if !found {
		return nil, ErrNoProject
	}

	p.normalise()

	return p, nil
}

// add records a requirement, keeping the first version we learn for a given
// kind+name pair. Manifests are read in a deliberate order, so an earlier,
// more authoritative source wins.
func (p *Project) add(r Requirement) {
	for i, existing := range p.Requirements {
		if existing.Kind == r.Kind && existing.Name == r.Name {
			if existing.Version == "" && r.Version != "" {
				p.Requirements[i].Version = r.Version
			}
			return
		}
	}

	p.Requirements = append(p.Requirements, r)
}

func (p *Project) addUnmatched(r Requirement) {
	for _, existing := range p.Unmatched {
		if existing.Evidence == r.Evidence {
			return
		}
	}

	p.Unmatched = append(p.Unmatched, r)
}

// normalise sorts requirements into a stable order so that output is
// deterministic — the same project always prints the same plan.
func (p *Project) normalise() {
	order := map[Kind]int{KindRuntime: 0, KindService: 1, KindRecipe: 2}

	sort.SliceStable(p.Requirements, func(i, j int) bool {
		a, b := p.Requirements[i], p.Requirements[j]
		if order[a.Kind] != order[b.Kind] {
			return order[a.Kind] < order[b.Kind]
		}
		return a.Name < b.Name
	})

	sort.SliceStable(p.Unmatched, func(i, j int) bool {
		return p.Unmatched[i].Evidence < p.Unmatched[j].Evidence
	})
}

// Of returns every requirement of a given kind.
func (p *Project) Of(kind Kind) []Requirement {
	var out []Requirement

	for _, r := range p.Requirements {
		if r.Kind == kind {
			out = append(out, r)
		}
	}

	return out
}

// Has reports whether a requirement of the given kind and name was detected.
func (p *Project) Has(kind Kind, name string) bool {
	for _, r := range p.Requirements {
		if r.Kind == kind && r.Name == name {
			return true
		}
	}

	return false
}
