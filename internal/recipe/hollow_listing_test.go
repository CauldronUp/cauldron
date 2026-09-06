package recipe

import (
	"testing"
)

// A listing shown only empty is answered but not described.
//
// UnshownListing asks whether any case gets a successful answer out of the
// route. That is the right first question and it is not the last one: a case
// that asserts the collection is empty proves the route exists, the status it
// answers with and the shape of the envelope around it, and says nothing at
// all about a record. Several Recipes have no fixture holding one, so the
// empty answer is the only evidence available -- worth having, worth writing
// down, and worth counting apart from the listings that show a record.
func TestAListingShownOnlyEmptyIsCounted(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{
			Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list",
		}},
		Resources: map[string]Resource{"thing": {Collection: "data"}},
		Conformance: []Case{{
			Name:    "the listing answers, and this fixture has none",
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect:  Expectation{Status: 200, Absent: []string{"data[0]"}},
		}},
	}

	if n := r.UnshownListing(); n != 0 {
		t.Fatalf("a listing with a successful case: counted %d unshown, want 0", n)
	}

	if n := r.HollowListing(); n != 1 {
		t.Fatalf("a listing shown only empty: counted %d, want 1", n)
	}

	r.Conformance = append(r.Conformance, Case{
		Name:    "the listing answers with a thing in it",
		Request: Request{Method: "GET", Path: "/v1/things"},
		Expect:  Expectation{Status: 200, Body: map[string]any{"data[0].id": "t_1"}},
	})

	if n := r.HollowListing(); n != 0 {
		t.Errorf("a listing with a record asserted: counted %d hollow, want 0", n)
	}
}

// A listing nothing answers at all is not hollow. It is unshown, which is a
// different line in the report, and counting it twice would double the debt.
func TestAnUnshownListingIsNotHollow(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{Method: "GET", Path: "/v1/things", Operation: "list"}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect:  Expectation{Status: 401},
		}},
	}

	if n := r.HollowListing(); n != 0 {
		t.Errorf("a listing with no successful case: counted %d hollow, want 0", n)
	}
}

// A case asserting something outside the collection does not describe a record
// either. The question is whether anything is asserted *in* the collection.
func TestAnEnvelopeAssertionIsNotARecord(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{
			Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list",
		}},
		Resources: map[string]Resource{"thing": {Collection: "data"}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect:  Expectation{Status: 200, Body: map[string]any{"has_more": false}},
		}},
	}

	if n := r.HollowListing(); n != 1 {
		t.Errorf("a listing asserting only its envelope: counted %d hollow, want 1", n)
	}
}

// A bare listing has no envelope, so anything asserted is in the collection.
//
// Alpaca's account and SES's account are bare and collapse_single: with one
// record the response *is* the record, and the only way to assert a field on it
// is at the top level. Requiring an index there asked for a shape those
// responses never have.
func TestABareListingHasNothingOutsideTheCollection(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{
			Method: "GET", Path: "/v2/account", Resource: "account", Operation: "list",
			List: &ListResponse{Style: "bare", CollapseSingle: true},
		}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v2/account"},
			Expect:  Expectation{Status: 200, Body: map[string]any{"id": "acct_1"}},
		}},
	}

	if n := r.HollowListing(); n != 0 {
		t.Errorf("a collapsed bare listing asserting a field: counted %d hollow, want 0", n)
	}

	r.Conformance[0].Expect.Body = nil

	if n := r.HollowListing(); n != 1 {
		t.Errorf("a bare listing asserting nothing: counted %d hollow, want 1", n)
	}
}
