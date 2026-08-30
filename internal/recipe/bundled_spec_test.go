// The Recipes that name a provider's own description, and what they owe.

package recipe

import "testing"

// A shipped Recipe that names a description and has not fingerprinted it is
// the worst of both states: the scan reports it every day as unrecorded, which
// is a line nobody can act on, and a reader of the Recipe sees a spec URL and
// reasonably assumes something is checking it.
//
// The validator cannot enforce this, because an unrecorded fingerprint is
// exactly what a Recipe looks like between adding the URL and running the
// scan. It is a rule about what may ship, not about what may be written.
func TestEveryShippedRecipeThatNamesADescriptionHasFingerprintedIt(t *testing.T) {
	for _, name := range Bundled() {
		r, err := Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		if r.Upstream.Spec == "" {
			continue
		}

		if r.Upstream.SpecHash == "" {
			t.Errorf("%s names a description and has not fingerprinted it: run 'cauldron drift --record %s'", name, name)
		}

		if r.Upstream.SpecSeen == "" {
			t.Errorf("%s fingerprinted a description without recording when", name)
		}
	}
}

// A superseded version is knowledge about which clients must not be offered
// this Recipe, and it is only worth writing down if it names the host that
// tells them apart. The validator already refuses one without a host; this
// checks the shipped Recipes actually carry the knowledge rather than an
// empty list.
func TestASupersededVersionThatShipsNamesItsHost(t *testing.T) {
	found := 0

	for _, name := range Bundled() {
		r, err := Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		for _, was := range r.Upstream.Supersedes {
			found++

			if was.Host == "" {
				t.Errorf("%s supersedes %q and names no host", name, was.Version)
			}

			if was.Note == "" {
				t.Errorf("%s supersedes %q and does not say why it matters", name, was.Version)
			}
		}
	}

	if found == 0 {
		t.Error("no shipped Recipe records a version it replaced, and at least three of them know one")
	}
}
