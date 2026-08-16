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
