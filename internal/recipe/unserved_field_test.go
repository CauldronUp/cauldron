package recipe

import (
	"testing"
)

// A field the emulator can never send is worse than one nothing asserts.
//
// UnassertedField counts names with no evidence behind them. This counts names
// the fake will not produce at all: the Recipe says a customer carries
// default_currency, a fixture holds a customer, and the response omits the key
// entirely. Code written against Cauldron reads undefined and finds out what
// the provider sends the first time it talks to one.
//
// The rule is the sandbox's own, read off the shaping it does. A declared field
// reaches the wire when a fixture record sets it, or it has a default, or it is
// null_when_unset, or its type is one the sandbox stamps from the clock and it
// has not said stamped: false. Nothing else appears.
//
// Only resources some fixture actually holds are counted. A resource nothing
// seeds has no record for any of its fields to be missing from, and the
// listings counters already say so in plainer words.
func TestAFieldNoFixtureSetsIsCounted(t *testing.T) {
	r := &Recipe{
		Resources: map[string]Resource{"thing": {
			Fields: map[string]Field{
				"name":   {Type: "string"},
				"colour": {Type: "string"},
			},
		}},
		Fixtures: map[string]Fixture{
			"one": {"thing": []map[string]any{{"name": "a thing"}}},
		},
	}

	if n := r.UnservedField(); n != 1 {
		t.Fatalf("one field set and one not: counted %d unserved, want 1", n)
	}

	r.Fixtures["one"]["thing"][0]["colour"] = "red"

	if n := r.UnservedField(); n != 0 {
		t.Errorf("both fields set: counted %d unserved, want 0", n)
	}
}

// The three ways a field reaches the wire without a fixture setting it.
func TestADefaultANullAndAStampAllReachTheWire(t *testing.T) {
	for what, field := range map[string]Field{
		"a default":       {Type: "string", Default: "unknown"},
		"a null":          {Type: "string", NullWhenUnset: true},
		"a stamp":         {Type: "datetime"},
		"a ms stamp":      {Type: "timestamp_ms"},
		"a stamp as text": {Type: "timestamp_ms_string"},
	} {
		r := &Recipe{
			Resources: map[string]Resource{"thing": {Fields: map[string]Field{"when": field}}},
			Fixtures:  map[string]Fixture{"one": {"thing": []map[string]any{{}}}},
		}

		if n := r.UnservedField(); n != 0 {
			t.Errorf("%s: counted %d unserved, want 0", what, n)
		}
	}
}

// A stamped type that says stamped: false is a field whose absence means
// something, and the sandbox leaves it absent on purpose.
func TestAStampTurnedOffDoesNotReachTheWire(t *testing.T) {
	off := false

	r := &Recipe{
		Resources: map[string]Resource{"thing": {
			Fields: map[string]Field{"deleted_at": {Type: "datetime", Stamped: &off}},
		}},
		Fixtures: map[string]Fixture{"one": {"thing": []map[string]any{{}}}},
	}

	if n := r.UnservedField(); n != 1 {
		t.Errorf("a stamp turned off: counted %d unserved, want 1", n)
	}
}

// A field the wire never carries is not missing from it.
func TestAFieldOffTheWireIsNotUnserved(t *testing.T) {
	r := &Recipe{
		Resources: map[string]Resource{"thing": {
			Fields: map[string]Field{"app_name": {Type: "string", In: "-"}},
		}},
		Fixtures: map[string]Fixture{"one": {"thing": []map[string]any{{}}}},
	}

	if n := r.UnservedField(); n != 0 {
		t.Errorf("a field declared off the wire: counted %d unserved, want 0", n)
	}
}

// And a resource nothing seeds is somebody else's count.
func TestAResourceNoFixtureHoldsIsNotCountedHere(t *testing.T) {
	r := &Recipe{
		Resources: map[string]Resource{"thing": {
			Fields: map[string]Field{"name": {Type: "string"}},
		}},
		Fixtures: map[string]Fixture{"empty": {}},
	}

	if n := r.UnservedField(); n != 0 {
		t.Errorf("a resource no fixture holds: counted %d unserved, want 0", n)
	}
}
