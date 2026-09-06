package recipe

import (
	"testing"
)

// A listing nothing answers successfully is a description with no evidence
// under it.
//
// The Recipe can declare a collection, a page size, an envelope, a cursor field
// and a count field, and every case touching that route can be checking a
// failure -- no credential, a wrong method, an unknown path. The file reads as
// a description of a working endpoint and nothing has ever seen one.
//
// This turned up once already, sideways: the rule that a response field name
// needs a case asserting it found "two of them ... with no successful list case
// at all". Two was what that angle could see. Counted head-on it is 168.
func TestAListingWithNoSuccessfulCaseIsCounted(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{
			Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list",
		}},
		Conformance: []Case{{
			Name:    "a request without a key is refused",
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect:  Expectation{Status: 401},
		}},
	}

	if n := r.UnshownListing(); n != 1 {
		t.Fatalf("a listing whose only case is a 401: counted %d, want 1", n)
	}

	r.Conformance = append(r.Conformance, Case{
		Name:    "the listing answers",
		Request: Request{Method: "GET", Path: "/v1/things"},
		Expect:  Expectation{Status: 200},
	})

	if n := r.UnshownListing(); n != 0 {
		t.Errorf("a listing with a 200 case: counted %d, want 0", n)
	}
}

// A case with no status expects 200, which is the format's own default, so it
// counts as showing the listing.
func TestAnUnstatedStatusCountsAsTwoHundred(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{Method: "GET", Path: "/v1/things", Operation: "list"}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/things"},
		}},
	}

	if n := r.UnshownListing(); n != 0 {
		t.Errorf("a case with no declared status counted as unshown: %d", n)
	}
}

// An arming case installs a failure on purpose. A 200 beside one describes the
// arm rather than the listing, so it does not count as evidence the listing
// works.
func TestAnArmingCaseIsNotEvidenceTheListingWorks(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{Method: "GET", Path: "/v1/things", Operation: "list"}},
		Conformance: []Case{{
			Arm:     "rate_limited",
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect:  Expectation{Status: 200},
		}},
	}

	if n := r.UnshownListing(); n != 1 {
		t.Errorf("an armed case counted as showing the listing: %d", n)
	}
}

// A path placeholder matches whatever a case put in its position, or every
// scoped listing in the collection would count as unshown while having a case
// that answers it.
func TestAScopedListingIsShownByAConcreteCase(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{
			Method: "GET", Path: "/v1/apps/{app_id}/things", Operation: "list",
		}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/apps/app_42/things"},
			Expect:  Expectation{Status: 200},
		}},
	}

	if n := r.UnshownListing(); n != 0 {
		t.Errorf("a scoped listing with a concrete case counted as unshown: %d", n)
	}
}

// Only listings. A get, a create and a delete are somebody else's count.
func TestOnlyListingsCountAsUnshown(t *testing.T) {
	r := &Recipe{
		Routes: []Route{
			{Method: "GET", Path: "/v1/things/{id}", Operation: "get"},
			{Method: "POST", Path: "/v1/things", Operation: "create"},
		},
	}

	if n := r.UnshownListing(); n != 0 {
		t.Errorf("a get and a create counted as unshown listings: %d", n)
	}
}
