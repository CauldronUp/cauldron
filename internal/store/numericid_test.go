package store

import "testing"

// A fixture that pins a numeric identifier must keep it.
//
// The doc comment on Create has promised this since it was written -- "fixtures
// pin their IDs so that seeded data is quotable in docs and tests" -- and for
// numeric ids it was not true. The identifier was read with a string type
// assertion, so a YAML integer yielded the empty string and the record was
// handed a freshly minted id instead, silently.
//
// Ten shipped Recipes pin numeric ids. Paystack's bank list pins 1, 7, 9, 18,
// 21, 31 and 76, and served 1 through 7: the right number of records, the wrong
// identifier on six of them, and nothing anywhere to notice it. Its conformance
// suite passed throughout, because no case asserted an id -- which is the whole
// reason this test is in the store rather than in a Recipe.
func TestAFixtureKeepsTheNumericIdentifierItPinned(t *testing.T) {
	s := New(1)
	s.DeclareStyle("bank", "numeric", "", 0)

	for _, id := range []any{1, 7, 76} {
		if _, err := s.Create("bank", Record{"id": id, "name": "Example"}); err != nil {
			t.Fatalf("create %v: %v", id, err)
		}
	}

	for _, want := range []string{"1", "7", "76"} {
		if _, err := s.Get("bank", want); err != nil {
			t.Errorf("record pinned as %s was not stored under it: %v", want, err)
		}
	}
}

// The int64 and float64 forms a YAML decoder can produce have to work too, or
// the fix depends on which library parsed the file.
func TestANumericIdentifierIsKeptWhateverItsGoType(t *testing.T) {
	for name, id := range map[string]any{
		"int":     44,
		"int64":   int64(44),
		"float64": float64(44),
	} {
		s := New(1)
		s.DeclareStyle("thing", "numeric", "", 0)

		if _, err := s.Create("thing", Record{"id": id}); err != nil {
			t.Fatalf("%s: create: %v", name, err)
		}

		if _, err := s.Get("thing", "44"); err != nil {
			t.Errorf("%s: pinned id was not kept: %v", name, err)
		}
	}
}

// A float that is not a whole number is not an identifier anybody pinned on
// purpose, and rendering 1.5 as "1.5" or "2" are both guesses. Minting instead
// is the honest answer.
func TestAFractionalIdentifierIsNotTreatedAsPinned(t *testing.T) {
	s := New(1)
	s.DeclareStyle("thing", "numeric", "", 0)

	stored, err := s.Create("thing", Record{"id": 1.5})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if got, _ := stored["id"].(string); got == "1.5" {
		t.Error("a fractional id was treated as a pinned identifier")
	}
}

// Still mints when nothing was pinned, which is the behaviour every Recipe
// without fixture ids depends on.
func TestARecordWithNoIdentifierStillGetsOne(t *testing.T) {
	s := New(1)
	s.DeclareStyle("thing", "numeric", "", 0)

	stored, err := s.Create("thing", Record{"name": "Example"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if id, _ := stored["id"].(string); id == "" {
		t.Error("a record with no pinned id was not given one")
	}
}
