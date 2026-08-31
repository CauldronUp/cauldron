// Adding a provider's own description to a Recipe, without letting it
// overwrite one.

package recipe

import "testing"

// The bargain this makes is worth stating before the tests that hold it up.
//
// A description carries breadth: Box declares 186 paths where its Recipe models
// 7, OpenAI declares 182 where its Recipe models 5. Hand-writing the rest will
// never happen, and a client that only needs POST /v1/orders to exist and
// answer something plausibly shaped is well served by a schema.
//
// A description does not carry truth. Kraken's says error is an array of
// string; it cannot say that the array is empty on success and that an empty
// array is true. Deezer's declares 200: Artist; it cannot say that every
// failure is also a 200. Asana's omits completed and due_on from the task
// listing entirely -- they are opt_fields, real and undeclared -- so a mock
// built from that description serves tasks without them, and the Recipe that
// pins them is right where the description is wrong.
//
// So: the description may add routes the Recipe does not have, and may never
// change one it does. Anything derived is marked, counted and reported, because
// a user must never have to guess which kind of answer they are looking at.
func base() *Recipe {
	return &Recipe{
		Name:       "provider",
		Version:    "0.1.0",
		Capability: "docs",
		Upstream:   Upstream{API: "v1"},
		Auth:       Auth{Scheme: "none"},
		Resources: map[string]Resource{
			"order": {
				ID:         ID{Style: "opaque"},
				Collection: "orders",
				Fields:     map[string]Field{"total": {Type: "integer"}},
			},
		},
		Routes: []Route{
			{Method: "GET", Path: "/v1/orders", Resource: "order", Operation: "list"},
		},
		Errors:   map[string]Error{"not_found": {Status: 404, Message: "no"}},
		Fixtures: map[string]Fixture{"empty": {}},
	}
}

func drafted() *Recipe {
	return &Recipe{
		Name:       "provider",
		Version:    "0.1.0",
		Capability: "docs",
		Upstream:   Upstream{API: "v1"},
		Auth:       Auth{Scheme: "none"},
		Resources: map[string]Resource{
			// The same name as the base's, and not the same shape.
			"order": {ID: ID{Style: "opaque"}, Fields: map[string]Field{"amount": {Type: "string"}}},
			// One the base has never heard of, and -- as a drafted resource does --
			// naming no collection.
			"refund": {ID: ID{Style: "opaque"}, Fields: map[string]Field{"reason": {Type: "string"}}},
		},
		Routes: []Route{
			// Already modelled, and modelled better.
			{Method: "GET", Path: "/v1/orders", Resource: "order", Operation: "list"},
			// Not modelled, and its resource collides with one that is.
			{Method: "POST", Path: "/v1/orders", Resource: "order", Operation: "create"},
			// Not modelled, and free to add.
			{Method: "GET", Path: "/v1/refunds", Resource: "refund", Operation: "list"},
		},
		Errors:   map[string]Error{"not_found": {Status: 400, Message: "different"}},
		Fixtures: map[string]Fixture{"empty": {}},
	}
}

func TestTheRecipeAlwaysWinsOnARouteItAlreadyModels(t *testing.T) {
	merged, report := Augment(base(), drafted())

	var found int

	for _, route := range merged.Routes {
		if route.Method == "GET" && route.Path == "/v1/orders" {
			found++

			if route.Derived {
				t.Error("the description replaced a route the Recipe already had")
			}
		}
	}

	if found != 1 {
		t.Errorf("the route appears %d times, want 1", found)
	}

	if report.Kept != 1 {
		t.Errorf("report kept %d, want 1", report.Kept)
	}
}

func TestARouteOnlyTheDescriptionHasIsAddedAndMarked(t *testing.T) {
	merged, report := Augment(base(), drafted())

	var refunds *Route

	for i, route := range merged.Routes {
		if route.Path == "/v1/refunds" {
			refunds = &merged.Routes[i]
		}
	}

	if refunds == nil {
		t.Fatalf("the description's route was not added: %+v", merged.Routes)
	}

	if !refunds.Derived {
		t.Error("a route from the description is not marked as derived")
	}

	if _, ok := merged.Resources["refund"]; !ok {
		t.Error("the resource the added route needs was not copied")
	}

	if report.Added != 1 {
		t.Errorf("report added %d, want 1", report.Added)
	}
}

// A resource name that means one thing in the Recipe and another in the
// description cannot be reconciled, and guessing which is meant would put a
// field on the wire that no provider sends. The route is dropped and said so.
func TestARouteWhoseResourceCollidesIsSkippedAndReported(t *testing.T) {
	merged, report := Augment(base(), drafted())

	for _, route := range merged.Routes {
		if route.Method == "POST" && route.Path == "/v1/orders" {
			t.Error("a route was added on a resource whose name the Recipe already uses differently")
		}
	}

	if report.Skipped != 1 {
		t.Errorf("report skipped %d, want 1", report.Skipped)
	}

	if len(report.Reasons) != 1 {
		t.Fatalf("report gave %d reasons, want 1", len(report.Reasons))
	}

	// The Recipe's own resource is untouched.
	if _, ok := merged.Resources["order"].Fields["amount"]; ok {
		t.Error("the description's resource overwrote the Recipe's")
	}

	if _, ok := merged.Resources["order"].Fields["total"]; !ok {
		t.Error("the Recipe's own resource lost a field")
	}
}

func TestTheDescriptionNeverChangesAnErrorTheRecipeDeclares(t *testing.T) {
	merged, _ := Augment(base(), drafted())

	if got := merged.Errors["not_found"]; got.Status != 404 || got.Message != "no" {
		t.Errorf("the description overwrote an error: %+v", got)
	}
}

// Nothing to add is not a failure, and must not quietly alter the Recipe.
func TestAugmentingWithNothingReturnsTheRecipeUnchanged(t *testing.T) {
	for _, extra := range []*Recipe{nil, {Name: "provider"}} {
		merged, report := Augment(base(), extra)

		if len(merged.Routes) != 1 {
			t.Errorf("routes became %d, want 1", len(merged.Routes))
		}

		if report.Added != 0 || report.Skipped != 0 {
			t.Errorf("report says added %d skipped %d, want 0 and 0", report.Added, report.Skipped)
		}
	}
}

// A derived Recipe has to say where it came from, or the report is the only
// place the provenance lives and it disappears the moment anybody saves the
// merged Recipe to a file.
func TestTheMergedRecipeRecordsWhereTheAdditionsCameFrom(t *testing.T) {
	extra := drafted()
	extra.Upstream.Spec = "https://example.test/openapi.json"

	merged, _ := Augment(base(), extra)

	if merged.DerivedFrom != "https://example.test/openapi.json" {
		t.Errorf("the merged Recipe does not name its description: %q", merged.DerivedFrom)
	}
}

// The guarantee the whole feature rests on: reaching for a description can add
// to a Recipe and can never break one.
//
// Adyen is what found it. Its Recipe wraps listings without a key, so every
// resource must name its collection; the drafted resources name none, and the
// merged Recipe would not validate. A Recipe that will not validate will not
// serve, so a user who asked for more endpoints would have got none at all.
func TestAugmentNeverReturnsARecipeThatWillNotValidate(t *testing.T) {
	r := base()
	// Wrapped with no key: every resource now owes a collection name.
	r.Responses.List.Style = "wrapped"

	extra := drafted()

	merged, report := Augment(r, extra)

	if err := merged.Validate(); err != nil {
		t.Fatalf("augmenting produced a Recipe that will not validate: %v", err)
	}

	// And it did so by making the addition work, not by refusing to add.
	if report.Added == 0 {
		t.Errorf("nothing was added, so the guarantee was met by giving up: %+v", report)
	}
}

// The collection name is read from the provider's own path rather than guessed
// from the resource name, because the path is what the provider actually calls
// the collection.
func TestACopiedResourceTakesItsCollectionNameFromThePath(t *testing.T) {
	r := base()
	r.Responses.List.Style = "wrapped"

	merged, _ := Augment(r, drafted())

	if got := merged.Resources["refund"].Collection; got != "refunds" {
		t.Errorf("collection is %q, want %q, taken from /v1/refunds", got, "refunds")
	}
}

// Asana wraps every listing in "data" -- all three of its written resources say
// so -- and its Recipe declares no list key, so the wrapper is each resource's
// collection name. A derived resource taking its collection from the path
// therefore served {"workspaces": []} where Asana sends {"data": [...]}, which
// is a wrong envelope rather than a missing one, and worse for a client that
// parses.
//
// Where the written Recipe is unanimous about its collection name, that is
// observed evidence about the provider and a derived resource should follow it.
// Where it is not unanimous, there is nothing to follow and the path stands.
func TestADerivedResourceFollowsAUnanimousCollectionName(t *testing.T) {
	r := base()
	r.Responses.List.Style = "wrapped"
	// Two resources agreeing, as Asana has three.
	r.Resources["order"] = Resource{
		ID:         ID{Style: "opaque"},
		Collection: "data",
		Fields:     map[string]Field{"total": {Type: "integer"}},
	}
	r.Resources["invoice"] = Resource{
		ID:         ID{Style: "opaque"},
		Collection: "data",
		Fields:     map[string]Field{"total": {Type: "integer"}},
	}
	r.Routes = append(r.Routes, Route{Method: "GET", Path: "/v1/invoices", Resource: "invoice", Operation: "list"})

	merged, _ := Augment(r, drafted())

	if got := merged.Resources["refund"].Collection; got != "data" {
		t.Errorf("collection is %q, want %q: the Recipe is unanimous and the derived resource ignored it", got, "data")
	}
}

// And when the Recipe disagrees with itself there is nothing to follow, so the
// path is the only evidence left.
func TestADerivedResourceFallsBackToThePathWhenTheRecipeDisagrees(t *testing.T) {
	r := base()
	r.Responses.List.Style = "wrapped"
	r.Resources["order"] = Resource{
		ID: ID{Style: "opaque"}, Collection: "orders",
		Fields: map[string]Field{"total": {Type: "integer"}},
	}
	r.Resources["invoice"] = Resource{
		ID: ID{Style: "opaque"}, Collection: "results",
		Fields: map[string]Field{"total": {Type: "integer"}},
	}
	r.Routes = append(r.Routes, Route{Method: "GET", Path: "/v1/invoices", Resource: "invoice", Operation: "list"})

	merged, _ := Augment(r, drafted())

	if got := merged.Resources["refund"].Collection; got != "refunds" {
		t.Errorf("collection is %q, want %q, taken from the path", got, "refunds")
	}
}
