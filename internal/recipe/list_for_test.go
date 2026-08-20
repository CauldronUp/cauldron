package recipe

import "testing"

// A provider's listings do not always share a shape. Clerk's users and
// sessions answer with bare arrays and its organisations with
// {data, total_count}; Algolia's browse carries a cursor its search has not.
// A Recipe-wide envelope has to be wrong about one of them, and the wrongness
// is the expensive kind: code written against the emulator reads
// response.data.map(...) and receives an array from the provider, where .data
// is undefined.

func TestARouteWithoutAnOverrideInheritsEverything(t *testing.T) {
	r := Recipe{Responses: Responses{List: ListResponse{Style: "bare", CountField: "total"}}}

	got := r.ListFor(Route{})
	if got.Style != "bare" || got.CountField != "total" {
		t.Errorf("a route with no list block changed the envelope: %+v", got)
	}
}

func TestARouteOverridesOnlyWhatItNames(t *testing.T) {
	r := Recipe{Responses: Responses{List: ListResponse{Style: "bare", CountField: "total", LimitField: "per_page"}}}

	got := r.ListFor(Route{List: &ListResponse{Style: "wrapped", Key: "data"}})

	if got.Style != "wrapped" || got.Key != "data" {
		t.Errorf("the override did not apply: %+v", got)
	}

	// The fields it did not name are still the Recipe's.
	if got.CountField != "total" || got.LimitField != "per_page" {
		t.Errorf("the override cleared fields it never mentioned: %+v", got)
	}
}

func TestADashClearsAnInheritedField(t *testing.T) {
	// A route saying the provider sends nothing there, which is different
	// from saying nothing at all.
	r := Recipe{Responses: Responses{List: ListResponse{Style: "wrapped", CountField: "total_count"}}}

	got := r.ListFor(Route{List: &ListResponse{CountField: "-"}})
	if got.CountField != "" {
		t.Errorf(`"-" did not clear the field: %q`, got.CountField)
	}

	if got.Style != "wrapped" {
		t.Errorf("clearing one field disturbed another: %+v", got)
	}
}

func TestABooleanCanOnlyBeTurnedOn(t *testing.T) {
	// An unset boolean and a false one are the same value in YAML, so a route
	// can add a header the Recipe does not send and cannot take one away.
	// Guessing which was meant is how a Recipe ends up asserting something
	// nobody wrote.
	r := Recipe{Responses: Responses{List: ListResponse{}}}

	if got := r.ListFor(Route{List: &ListResponse{LinkHeader: true}}); !got.LinkHeader {
		t.Error("a route could not turn the Link header on")
	}

	on := Recipe{Responses: Responses{List: ListResponse{LinkHeader: true}}}
	if got := on.ListFor(Route{List: &ListResponse{}}); !got.LinkHeader {
		t.Error("an empty list block turned the Link header off")
	}
}

func TestTheRecipeIsNotMutatedByAnOverride(t *testing.T) {
	// ListFor returns a value, and a route that overrides must not change the
	// envelope every other route inherits.
	r := Recipe{Responses: Responses{List: ListResponse{Style: "bare"}}}

	_ = r.ListFor(Route{List: &ListResponse{Style: "wrapped", Key: "data"}})

	if r.Responses.List.Style != "bare" || r.Responses.List.Key != "" {
		t.Errorf("one route's override leaked into the Recipe: %+v", r.Responses.List)
	}
}
