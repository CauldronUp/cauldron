package recipe

import (
	"fmt"
	"strings"
	"testing"
)

// A route that exists to model the empty answer is not missing evidence.
//
// Five routes across four Recipes are variants of a real listing with a query
// that matches nothing: Jisho's words-no-match and words-no-keyword, OpenAIRE's
// publications-no-match, iNaturalist's observation id of twelve nines, Deezer's
// search-empty-parameter. They exist because the provider answers differently
// when nothing matches -- OpenAIRE sends its total as the string "0" where a
// populated answer sends the number 1 -- and there is no record in them to
// describe, by construction.
//
// HollowListing was counting all five, and a count that can never reach zero
// teaches its reader to stop looking at it.
func TestARouteThatModelsTheEmptyAnswerIsNotHollow(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{
			Method: "GET", Path: "/v1/things-no-match", Resource: "thing",
			Operation: "list", Empty: true,
		}},
		Resources: map[string]Resource{"thing": {Collection: "data"}},
		Conformance: []Case{{
			Name:    "nothing matches",
			Request: Request{Method: "GET", Path: "/v1/things-no-match"},
			Expect:  Expectation{Status: 200, Absent: []string{"data[0]"}},
		}},
	}

	if n := r.HollowListing(); n != 0 {
		t.Errorf("a route that models the empty answer: counted %d hollow, want 0", n)
	}

	r.Routes[0].Empty = false

	if n := r.HollowListing(); n != 1 {
		t.Errorf("the same route without the declaration: counted %d hollow, want 1", n)
	}
}

// The declaration has to be true of the file, or it is a way of silencing the
// count rather than explaining it.
func TestEmptyIsRefusedOnARouteThatAnswersWithRecords(t *testing.T) {
	r := &Recipe{
		Name:    "example",
		Version: "0.1.0",
		Routes: []Route{{
			Method: "GET", Path: "/v1/things", Resource: "thing",
			Operation: "list", Empty: true,
		}},
		Resources: map[string]Resource{"thing": {ID: ID{Style: "opaque"}, Collection: "data"}},
		Conformance: []Case{{
			Name:    "the listing answers",
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect:  Expectation{Status: 200, Body: map[string]any{"data[0].id": "t_1"}},
		}},
	}

	var problems []string
	r.validateRoutes(func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	})

	if len(problems) == 0 {
		t.Fatal("a route declaring empty while a case asserts a record in it was accepted")
	}

	if !strings.Contains(strings.Join(problems, "\n"), "empty") {
		t.Errorf("the problems do not mention the declaration: %v", problems)
	}
}

// And it is a claim about a listing. Nothing else has a collection to be empty.
func TestEmptyIsRefusedOnANonListing(t *testing.T) {
	r := &Recipe{
		Name:    "example",
		Version: "0.1.0",
		Routes: []Route{{
			Method: "GET", Path: "/v1/things/{id}", Resource: "thing",
			Operation: "get", Empty: true,
		}},
		Resources: map[string]Resource{"thing": {ID: ID{Style: "opaque"}}},
	}

	var problems []string
	r.validateRoutes(func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	})

	if len(problems) == 0 {
		t.Error("empty on a non-listing was accepted")
	}
}
