package recipe

import (
	"strings"
	"testing"
)

// A Recipe that names a credential must hold one.
//
// This engine reads an auth block with no keys and no pattern as "accept
// anything", so that an author can model routes first and tighten the
// credential later. Naming a scheme and a header is not modelling routes
// first, though -- it reads as a claim that the credential is checked, and
// two Recipes made exactly that claim while enforcing nothing.
//
// Attio and CometChat both declared a scheme, named the header it travels in,
// wrote paragraphs about it, and listed no key. Both served every request that
// presented no credential at all. CometChat's file said in three separate
// places that its routes were "fixture-gated on the one apikey this Recipe
// holds"; it held none. Code written against either would have passed with no
// credential and failed against the provider, which is the direction of
// wrongness this whole project exists to avoid.
//
// scheme: none already says "this surface needs no credential", so the
// route-first case keeps a way to say itself, and it says it in words rather
// than by omission.
func TestASchemeWithNothingToCheckIsRefused(t *testing.T) {
	r := checking()
	r.Auth.Keys = nil

	err := r.Validate()
	if err == nil {
		t.Fatal("a Recipe naming a scheme and holding no key validated, so it would have authenticated nothing")
	}

	if !strings.Contains(err.Error(), "auth") {
		t.Errorf("the complaint does not mention auth: %v", err)
	}
}

// A pattern is the other way to hold a credential -- AWS signs every call, so
// there is no fixed value -- and it satisfies the rule on its own.
func TestAPatternIsEnoughToBeCheckingSomething(t *testing.T) {
	r := checking()
	r.Auth.Keys = nil
	r.Auth.Pattern = "^AWS4-HMAC-SHA256 "

	if err := r.Validate(); err != nil {
		t.Errorf("a Recipe checking a shape rather than a value was refused: %v", err)
	}
}

// And a Recipe that checks nothing on purpose still says so.
func TestSchemeNoneStillNeedsNoKeys(t *testing.T) {
	r := checking()
	r.Auth.Keys = nil
	r.Auth.Scheme = "none"

	if err := r.Validate(); err != nil {
		t.Errorf("scheme: none was refused for holding no key: %v", err)
	}
}

// An override that changes the carrier says what happens to the prefix.
//
// Empty means inherit, which is right for everything else and wrong for a
// prefix: a prefix belongs to a carrier, and a route or an alternative moving
// to another one almost never wants the old one's. Doppler found this the
// expensive way -- its Basic channel inherited "Bearer " and refused every key
// legitimately sent through it, making the fake stricter than the provider.
//
// "-" clears it, which is what clear already means for a list key and an
// identifier field. Requiring one of the two turns a silent stricter-fake into
// a refusal to boot.
func TestAnOverrideChangingTheCarrierMustSayWhatHappensToThePrefix(t *testing.T) {
	r := checking()
	r.Auth.Also = []Auth{{Scheme: "header", Header: "X-Api-Key"}}

	if err := r.Validate(); err == nil {
		t.Fatal("an alternative on another carrier silently inherited a prefix that carrier does not use")
	}

	r.Auth.Also = []Auth{{Scheme: "header", Header: "X-Api-Key", Prefix: "-"}}

	if err := r.Validate(); err != nil {
		t.Errorf("clearing the prefix explicitly was still refused: %v", err)
	}

	r.Auth.Also = []Auth{{Scheme: "header", Header: "X-Api-Key", Prefix: "Token "}}

	if err := r.Validate(); err != nil {
		t.Errorf("naming the carrier's own prefix was refused: %v", err)
	}
}

// The same for a route, which resolves through the same inheritance.
func TestARouteChangingTheCarrierMustSayWhatHappensToThePrefix(t *testing.T) {
	r := checking()
	r.Routes[0].Auth = &Auth{Scheme: "header", Header: "X-Report-Key", Keys: []string{"another"}}

	if err := r.Validate(); err == nil {
		t.Fatal("a route on another carrier silently inherited a prefix that carrier does not use")
	}

	r.Routes[0].Auth.Prefix = "-"

	if err := r.Validate(); err != nil {
		t.Errorf("clearing the prefix explicitly was still refused: %v", err)
	}
}

// Staying on the same carrier inherits the prefix, which is the ordinary case:
// a second surface with its own secret and the same shape of token.
func TestAnOverrideOnTheSameCarrierInheritsThePrefixSilently(t *testing.T) {
	r := checking()
	r.Routes[0].Auth = &Auth{Keys: []string{"another"}}

	if err := r.Validate(); err != nil {
		t.Errorf("a route naming only its own key was refused: %v", err)
	}
}

func checking() *Recipe {
	return &Recipe{
		Name:       "checking",
		Capability: "identity",
		Version:    "0.1.0",
		Upstream:   Upstream{API: "v1"},
		Auth: Auth{
			Scheme: "bearer",
			Prefix: "Bearer ",
			Keys:   []string{"a-key-the-recipe-holds"},
		},
		Resources: map[string]Resource{
			"thing": {ID: ID{Style: "opaque"}, Fields: map[string]Field{"id": {Type: "string"}}},
		},
		Routes: []Route{
			{Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list"},
		},
	}
}
