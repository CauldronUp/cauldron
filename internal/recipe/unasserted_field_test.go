package recipe

import (
	"testing"
)

// A field name nothing asserts could be renamed to anything.
//
// The format already refuses that for the two envelopes: responses.list.* and
// responses.error.* have to be asserted somewhere or the Recipe does not load,
// on the grounds that a name nothing checks is a name nobody would notice
// changing. The records inside those envelopes -- the fields a client actually
// reads -- had no such rule, and they are the larger half by an order of
// magnitude.
func TestAFieldNoCaseAssertsIsCounted(t *testing.T) {
	r := &Recipe{
		Resources: map[string]Resource{"thing": {
			Fields: map[string]Field{
				"name":  {Type: "string"},
				"shape": {Type: "string"},
			},
		}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect:  Expectation{Status: 200, Body: map[string]any{"data[0].name": "a thing"}},
		}},
	}

	if n := r.UnassertedField(); n != 1 {
		t.Fatalf("one of two fields asserted: counted %d unasserted, want 1", n)
	}

	r.Conformance[0].Expect.Body["data[0].shape"] = "round"

	if n := r.UnassertedField(); n != 0 {
		t.Errorf("both fields asserted: counted %d unasserted, want 0", n)
	}
}

// A field the wire never carries cannot be asserted, and is not debt.
//
// in: "-" is how a Recipe says the record holds the field and the response does
// not send it -- a partition that lives in the path, mostly. Counting those
// would ask for evidence of something the Recipe has said does not happen.
func TestAFieldThatNeverGoesOnTheWireIsNotCounted(t *testing.T) {
	r := &Recipe{
		Resources: map[string]Resource{"thing": {
			Fields: map[string]Field{"app_name": {Type: "string", In: "-"}},
		}},
	}

	if n := r.UnassertedField(); n != 0 {
		t.Errorf("a field declared off the wire: counted %d unasserted, want 0", n)
	}
}

// And the name to look for is the one on the wire.
//
// A field stored as title and sent as rendered is asserted by a case naming
// rendered; looking for title would report debt against a name no response
// carries, and the Recipe would be told to assert something that cannot appear.
func TestTheAssertedNameIsTheWireName(t *testing.T) {
	r := &Recipe{
		Resources: map[string]Resource{"post": {
			Fields: map[string]Field{"title_rendered": {Type: "string", As: "rendered"}},
		}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/posts"},
			Expect:  Expectation{Status: 200, Body: map[string]any{"title.rendered": "Hello"}},
		}},
	}

	if n := r.UnassertedField(); n != 0 {
		t.Errorf("a renamed field asserted under its wire name: counted %d unasserted, want 0", n)
	}
}

// A name is a name *of something*, and the case has to be about that thing.
//
// The first version of this counter asked whether any case in the file asserted
// the name anywhere, which is the right question for an envelope -- there is one
// envelope per Recipe -- and the wrong one for a record. Nutritionix declares
// name on four resources and a case about one of them credited all four; across
// the corpus, 476 names were credited by a case about something else entirely.
//
// The name is credited by a case whose route serves the resource. A resource
// with no route of its own is the exception: it appears only nested inside
// another's response, so an assertion anywhere is an assertion about it.
func TestAFieldIsAssertedByACaseAboutItsOwnResource(t *testing.T) {
	r := &Recipe{
		Routes: []Route{
			{Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list"},
			{Method: "GET", Path: "/v1/others", Resource: "other", Operation: "list"},
		},
		Resources: map[string]Resource{
			"thing": {Fields: map[string]Field{"colour": {Type: "string"}}},
			"other": {Fields: map[string]Field{"colour": {Type: "string"}}},
		},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/others"},
			Expect:  Expectation{Status: 200, Body: map[string]any{"data[0].colour": "red"}},
		}},
	}

	if n := r.UnassertedField(); n != 1 {
		t.Fatalf("a colour asserted on the other listing only: counted %d unasserted, want 1", n)
	}

	r.Conformance = append(r.Conformance, Case{
		Request: Request{Method: "GET", Path: "/v1/things"},
		Expect:  Expectation{Status: 200, Body: map[string]any{"data[0].colour": "blue"}},
	})

	if n := r.UnassertedField(); n != 0 {
		t.Errorf("both listings asserting their own colour: counted %d unasserted, want 0", n)
	}
}

// A resource served by no route of its own lives inside another's response, so
// an assertion anywhere in the file is an assertion about it.
func TestANestedOnlyResourceIsAssertedAnywhere(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list"}},
		Resources: map[string]Resource{
			"thing":  {Fields: map[string]Field{"id": {Type: "string"}}},
			"client": {Fields: map[string]Field{"trading_name": {Type: "string"}}},
		},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect: Expectation{Status: 200, Body: map[string]any{
				"data[0].id":                  "t_1",
				"data[0].client.trading_name": "Acme",
			}},
		}},
	}

	if n := r.UnassertedField(); n != 0 {
		t.Errorf("a nested-only resource asserted through its parent: counted %d unasserted, want 0", n)
	}
}

// And the case has to succeed, for the same reason a paging parameter's does:
// a refusal never reached the response, so it says nothing about its shape.
func TestAFieldAssertedOnlyByAFailingCaseIsUnasserted(t *testing.T) {
	r := &Recipe{
		Routes:    []Route{{Method: "GET", Path: "/v1/things", Resource: "thing", Operation: "list"}},
		Resources: map[string]Resource{"thing": {Fields: map[string]Field{"colour": {Type: "string"}}}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/things"},
			Expect:  Expectation{Status: 404, Body: map[string]any{"data[0].colour": "red"}},
		}},
	}

	if n := r.UnassertedField(); n != 1 {
		t.Errorf("a field asserted only by a 404: counted %d unasserted, want 1", n)
	}
}

// A nested field is asserted at the path it actually occupies.
//
// Front's message carries author_id, declared in: author and as: id, so the
// response has author.id -- and the message's own id is a different field at
// the top level. Looking for the leaf name alone credited author_id with every
// assertion of any id anywhere on that route, which is the collision the as:
// field exists to make possible in the first place.
func TestANestedFieldNeedsItsWholePath(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{Method: "GET", Path: "/v1/messages", Resource: "message", Operation: "list"}},
		Resources: map[string]Resource{"message": {
			Fields: map[string]Field{
				"id":        {Type: "string"},
				"author_id": {Type: "string", In: "author", As: "id"},
			},
		}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/messages"},
			Expect:  Expectation{Status: 200, Body: map[string]any{"data[0].id": "msg_1"}},
		}},
	}

	if n := r.UnassertedField(); n != 1 {
		t.Fatalf("only the top-level id asserted: counted %d unasserted, want 1", n)
	}

	r.Conformance[0].Expect.Body["data[0].author.id"] = "usr_1"

	if n := r.UnassertedField(); n != 0 {
		t.Errorf("both asserted at their own paths: counted %d unasserted, want 0", n)
	}
}

// A deeper nest is the same rule with more of it.
func TestADeeperNestNeedsItsWholePath(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{Method: "GET", Path: "/v1/transfers", Resource: "transfer", Operation: "list"}},
		Resources: map[string]Resource{"transfer": {
			Fields: map[string]Field{
				"destinationCardID": {Type: "string", In: "destination.card", As: "cardID"},
			},
		}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/v1/transfers"},
			Expect:  Expectation{Status: 200, Body: map[string]any{"data[0].card.cardID": "crd_1"}},
		}},
	}

	if n := r.UnassertedField(); n != 1 {
		t.Fatalf("asserted at a path missing the outer object: counted %d unasserted, want 1", n)
	}

	r.Conformance[0].Expect.Body = map[string]any{"data[0].destination.card.cardID": "crd_1"}

	if n := r.UnassertedField(); n != 0 {
		t.Errorf("asserted at its whole path: counted %d unasserted, want 0", n)
	}
}

// A nest may name an element of a list, and the index is not part of the name.
//
// Ory Kratos declares a login flow's csrf_group as in: ui.nodes[0].attributes,
// as: group, so the response has ui.nodes[0].attributes.group. Assertion paths
// have their [n] stripped before comparison and the declared chain did not, so
// "nodes[0]" was compared against "nodes" and twenty-three names on one
// resource read as unasserted while a case asserted every one of them.
func TestAnIndexInANestIsNotPartOfTheName(t *testing.T) {
	r := &Recipe{
		Routes: []Route{{Method: "GET", Path: "/self-service/login/flows", Resource: "login_flow", Operation: "get"}},
		Resources: map[string]Resource{"login_flow": {
			Fields: map[string]Field{
				"csrf_group": {Type: "string", In: "ui.nodes[0].attributes", As: "group"},
			},
		}},
		Conformance: []Case{{
			Request: Request{Method: "GET", Path: "/self-service/login/flows"},
			Expect: Expectation{Status: 200, Body: map[string]any{
				"ui.nodes[0].attributes.group": "default",
			}},
		}},
	}

	if n := r.UnassertedField(); n != 0 {
		t.Errorf("a field nested through a list element: counted %d unasserted, want 0", n)
	}

	r.Conformance[0].Expect.Body = map[string]any{"ui.nodes[0].group": "default"}

	if n := r.UnassertedField(); n != 1 {
		t.Errorf("asserted at a path missing the attributes object: counted %d unasserted, want 1", n)
	}
}
