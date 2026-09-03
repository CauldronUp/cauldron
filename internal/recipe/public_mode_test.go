package recipe

import (
	"testing"

	"gopkg.in/yaml.v3"
)

// The bare boolean keeps meaning what it meant, and the named modes read.
//
// A typo must not quietly become "not public": that would turn a declared
// exemption into a gated route, which fails as a 401 on a path the provider
// serves to anyone -- the failure this field exists to prevent. So anything
// that is not a boolean and not a mode this understands is an error rather
// than a shrug.
func TestPublicReadsABooleanAndTheNamedModes(t *testing.T) {
	for _, c := range []struct {
		yaml   string
		always bool
		absent bool
		refuse bool
	}{
		{"true", true, false, false},
		{"false", false, false, false},
		{"always", true, false, false},
		{"when-absent", false, true, false},
		{"when_absent", false, false, true},
		{"yes-please", false, false, true},
	} {
		var mode PublicMode

		err := yaml.Unmarshal([]byte(c.yaml), &mode)

		if c.refuse {
			if err == nil {
				t.Errorf("public: %s was accepted and should not be", c.yaml)
			}

			continue
		}

		if err != nil {
			t.Errorf("public: %s was refused: %v", c.yaml, err)

			continue
		}

		if mode.Always != c.always || mode.WhenAbsent != c.absent {
			t.Errorf("public: %s read as %+v", c.yaml, mode)
		}
	}
}

// Exempts is the whole difference between the two modes, so it is asserted
// directly rather than only through a served request.
func TestWhichVerdictsEachModeExcuses(t *testing.T) {
	always := PublicMode{Always: true}
	absent := PublicMode{WhenAbsent: true}
	none := PublicMode{}

	for _, v := range []Verdict{Accepted, Absent, Malformed, Rejected, Unentitled} {
		if !always.Exempts(v) {
			t.Errorf("the unconditional mode did not excuse verdict %d", v)
		}

		if none.Exempts(v) {
			t.Errorf("an undeclared mode excused verdict %d", v)
		}
	}

	if !absent.Exempts(Absent) {
		t.Error("when-absent did not excuse an absent credential")
	}

	for _, v := range []Verdict{Malformed, Rejected, Unentitled} {
		if absent.Exempts(v) {
			t.Errorf("when-absent excused verdict %d, which was presented", v)
		}
	}
}
