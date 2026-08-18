package runtime

import "testing"

// Shopify and Twilio both write /orders/{id}.json. Until the router understood
// a parameter with literal text around it, every single-object route on both
// providers matched nothing and answered 404. The conformance suite found it;
// this test keeps it found.
func TestParameterWithALiteralSuffix(t *testing.T) {
	segments := compilePath("/admin/api/{version}/orders/{id}.json")

	rt := route{method: "GET", segments: segments}

	vars, _, ok := rt.matches(splitPath("/admin/api/2026-01/orders/450789469.json"))
	if !ok {
		t.Fatal("a parameter followed by .json should match")
	}

	if vars["id"] != "450789469" {
		t.Errorf("id = %q, want the value without the extension", vars["id"])
	}

	if vars["version"] != "2026-01" {
		t.Errorf("version = %q", vars["version"])
	}
}

func TestASuffixIsRequiredWhenDeclared(t *testing.T) {
	rt := route{method: "GET", segments: compilePath("/orders/{id}.json")}

	if _, _, ok := rt.matches(splitPath("/orders/450789469")); ok {
		t.Error("a path missing the declared .json should not match")
	}

	if _, _, ok := rt.matches(splitPath("/orders/450789469.xml")); ok {
		t.Error("a different extension should not match")
	}
}

func TestAnEmptyParameterDoesNotMatch(t *testing.T) {
	rt := route{method: "GET", segments: compilePath("/orders/{id}.json")}

	// Otherwise "/orders/.json" would resolve to an object with a blank id.
	if _, _, ok := rt.matches(splitPath("/orders/.json")); ok {
		t.Error("a parameter with no value should not match")
	}
}

func TestALiteralSuffixOutranksABareParameter(t *testing.T) {
	r := &router{routes: []route{
		{method: "GET", segments: compilePath("/orders/{id}")},
		{method: "GET", segments: compilePath("/orders/{id}.json")},
	}}

	_, vars, ok := r.match("GET", "/orders/12.json")
	if !ok {
		t.Fatal("no route matched")
	}

	if vars["id"] != "12" {
		t.Errorf("id = %q, want the more specific route to win", vars["id"])
	}
}

// A parameter with literal text on both sides, matched against a path short
// enough for the two to overlap.
//
// /v1/aba{id}aba against /v1/ababa: five characters that both start with "aba"
// and end with "aba", so both checks passed and the slice that followed ran
// off the end. It panicked in the request path, where net/http recovers per
// connection and the caller sees an EOF with no response at all.
//
// No shipped Recipe declares a two-sided parameter, so this was latent rather
// than live. Adding Recipes is what this project mostly does, though, and the
// validator does not refuse the shape.
func TestAParameterWrappedInTextDoesNotPanicOnAShortPath(t *testing.T) {
	rt := route{method: "GET", segments: compilePath("/v1/aba{id}aba")}

	for _, path := range []string{"/v1/ababa", "/v1/abab", "/v1/aba", "/v1/abaaba"} {
		func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Errorf("%s panicked: %v", path, recovered)
				}
			}()

			if _, _, ok := rt.matches(splitPath(path)); ok {
				t.Errorf("%s matched, and there is no identifier left in it", path)
			}
		}()
	}

	// A path long enough to carry both really does match, so the guard rejects
	// only the impossible ones.
	vars, _, ok := rt.matches(splitPath("/v1/abaXYZaba"))
	if !ok {
		t.Fatal("/v1/abaXYZaba should match")
	}

	if vars["id"] != "XYZ" {
		t.Errorf("id = %q, want the middle", vars["id"])
	}
}
