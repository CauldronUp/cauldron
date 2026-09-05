package recipe

import (
	"strings"
	"testing"
)

// A Recipe cannot call its own published key malformed.
//
// Shape decides whether a credential is worth comparing, and the keys are what
// it gets compared against, so a shape excluding them describes a provider
// nobody has: every request refused, including the one the Recipe's own
// fixtures and README tell a reader to send.
func TestAShapeMustMatchTheKeysTheRecipePublishes(t *testing.T) {
	r := checking()
	r.Auth.Shape = "^pdl_[a-z]+_apikey_[A-Za-z0-9]+$"

	err := r.Validate()
	if err == nil {
		t.Fatal("a shape excluding the Recipe's own key validated")
	}

	if !strings.Contains(err.Error(), "a-key-the-recipe-holds") {
		t.Errorf("the complaint does not name the key it would refuse: %v", err)
	}

	r.Auth.Keys = []string{"pdl_sdbx_apikey_cauldron000"}

	if err := r.Validate(); err != nil {
		t.Errorf("a shape matching the key was still refused: %v", err)
	}
}

// A key listed as unentitled is held to the same rule. It is a credential this
// Recipe expects to be sent and answers a particular way, so a shape calling
// it malformed makes that answer unreachable.
func TestAShapeMustMatchTheUnentitledKeysToo(t *testing.T) {
	r := checking()
	r.Auth.Shape = "^key-[a-z]+$"
	r.Auth.Keys = []string{"key-good"}
	r.Auth.Unentitled = []string{"a-different-thing-entirely"}
	r.Auth.UnentitledError = "no_plan"
	r.Errors = map[string]Error{"no_plan": {Status: 403, Message: "Your plan does not include this."}}

	if err := r.Validate(); err == nil {
		t.Fatal("a shape that makes the unentitled answer unreachable validated")
	}
}

// Shape and pattern together say two different things about the same value,
// and the pattern wins, so the shape would decide nothing.
func TestAShapeAndAPatternTogetherAreRefused(t *testing.T) {
	r := checking()
	r.Auth.Keys = nil
	r.Auth.Pattern = "^AWS4-HMAC-SHA256 "
	r.Auth.Shape = "^AWS4"

	if err := r.Validate(); err == nil {
		t.Fatal("a Recipe declaring both a shape and a pattern validated")
	}
}

// A shape with no keys behind it refuses every request, which is not what any
// of the four providers that wanted this do.
func TestAShapeWithNoKeysIsRefused(t *testing.T) {
	r := checking()
	r.Auth.Keys = nil
	r.Auth.Shape = "^test-"

	if err := r.Validate(); err == nil {
		t.Fatal("a shape with nothing to compare against validated")
	}
}

// And an unreadable one is refused rather than silently never matching.
func TestAnUnreadableShapeIsRefused(t *testing.T) {
	r := checking()
	r.Auth.Shape = "^(unclosed"

	if err := r.Validate(); err == nil {
		t.Fatal("an invalid regular expression validated as a shape")
	}
}

// Two cases cannot share a name.
//
// The name is how the verify report refers to a case and how a reader finds
// it. Three Recipes shipped one twice -- Airtable, Chargebee and CircleCI each
// carried the same paging case verbatim in two places, one of them differing
// by a trailing blank line -- and nothing caught it: the counts read one
// higher than the Recipe's real coverage, the report printed the sentence
// twice, and a reader chasing a failure had two candidates for it.
func TestTwoConformanceCasesCannotShareAName(t *testing.T) {
	r := checking()
	r.Conformance = []Case{
		{Name: "a listing comes back", Request: Request{Method: "GET", Path: "/v1/things"}},
		{Name: "a listing comes back", Request: Request{Method: "GET", Path: "/v1/things"}},
	}

	err := r.Validate()
	if err == nil {
		t.Fatal("a Recipe carrying the same case name twice validated")
	}

	if !strings.Contains(err.Error(), "a listing comes back") {
		t.Errorf("the complaint does not name the case: %v", err)
	}
}
