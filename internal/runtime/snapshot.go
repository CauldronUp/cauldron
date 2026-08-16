package runtime

import (
	"fmt"
	"time"

	"github.com/CauldronUp/cauldron/internal/store"
)

// SnapshotVersion is bumped whenever the format changes incompatibly, so an
// old snapshot fails with a sentence rather than restoring wrong data.
const SnapshotVersion = 1

// Snapshot is the complete state of one sandbox.
//
// It exists so that "here is the exact state where it reproduces" can be a
// file attached to a bug report rather than a paragraph of instructions that
// nobody can follow.
type Snapshot struct {
	Version int    `json:"version"`
	Recipe  string `json:"recipe"`
	// RecipeVersion guards against restoring into a Recipe that has since
	// changed shape underneath the data.
	RecipeVersion string       `json:"recipe_version"`
	Fixture       string       `json:"fixture,omitempty"`
	Clock         time.Time    `json:"clock"`
	State         store.Export `json:"state"`
}

// Snapshot captures the sandbox.
func (s *Sandbox) Snapshot() Snapshot {
	return Snapshot{
		Version:       SnapshotVersion,
		Recipe:        s.recipe.Name,
		RecipeVersion: s.recipe.Version,
		Fixture:       s.Fixture(),
		Clock:         s.clock.Now(),
		State:         s.store.Export(),
	}
}

// ErrSnapshotMismatch means a snapshot cannot be restored into this sandbox.
type ErrSnapshotMismatch struct {
	Want, Got string
	Reason    string
}

func (e *ErrSnapshotMismatch) Error() string {
	return fmt.Sprintf("this snapshot %s (%s, not %s)", e.Reason, e.Got, e.Want)
}

// Restore replaces the sandbox's state with a snapshot.
func (s *Sandbox) Restore(snapshot Snapshot) error {
	if snapshot.Version != SnapshotVersion {
		return &ErrSnapshotMismatch{
			Reason: "was taken by a different version of Cauldron",
			Want:   fmt.Sprintf("format %d", SnapshotVersion),
			Got:    fmt.Sprintf("format %d", snapshot.Version),
		}
	}

	if snapshot.Recipe != s.recipe.Name {
		return &ErrSnapshotMismatch{
			Reason: "belongs to another recipe",
			Want:   s.recipe.Name,
			Got:    snapshot.Recipe,
		}
	}

	// A differing Recipe version is a warning, not a refusal: Recipes evolve,
	// and refusing would make old bug reports unopenable. Restoring into a
	// Recipe that has removed a resource still fails, in Import.
	if err := s.store.Import(snapshot.State); err != nil {
		return err
	}

	if !snapshot.Clock.IsZero() {
		s.clock.Set(snapshot.Clock)
	}

	s.faults.Clear()
	s.log.reset()
	s.webhooks.reset()

	s.mu.Lock()
	s.fixture = snapshot.Fixture
	s.mu.Unlock()

	return nil
}

// StaleRecipe reports whether a snapshot was taken against a different version
// of this Recipe, so the caller can warn without refusing.
func (s *Sandbox) StaleRecipe(snapshot Snapshot) bool {
	return snapshot.RecipeVersion != "" && snapshot.RecipeVersion != s.recipe.Version
}
