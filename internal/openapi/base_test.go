package openapi

import "testing"

func TestBasePathIsReadFromTheServer(t *testing.T) {
	for _, c := range []struct {
		name string
		url  string
		want string
	}{
		{"a plain server", "https://api.box.com/2.0", "/2.0"},
		{"a trailing slash", "https://api.example.com/v1/", "/v1"},
		{"no path at all", "https://api.github.com", ""},
		{"root", "https://api.example.com/", ""},
		{"several segments", "https://api.example.com/api/v2", "/api/v2"},

		// Chargebee's, which is a template in the authority and a literal
		// path. url.Parse reads {protocol} as an opaque scheme and cannot
		// split it, so this is the case the fallback exists for.
		{"templated authority", "{protocol}://{site}.{environment}:{port}/api/v2", "/api/v2"},

		// A prefix that varies per request is not a prefix.
		{"templated path", "https://api.example.com/{version}", ""},
		{"empty", "", ""},
	} {
		if got := basePathOf(c.url); got != c.want {
			t.Errorf("%s: basePathOf(%q) = %q, want %q", c.name, c.url, got, c.want)
		}
	}
}

func TestTheFirstUsableServerWins(t *testing.T) {
	doc := &Document{Servers: []Server{
		{URL: "https://api.example.com"},
		{URL: "https://api.example.com/v3"},
	}}

	// The first declares no prefix, which is not the same as declaring none
	// usefully: a description whose first server is the bare host and whose
	// second carries the version is describing the same API twice.
	//
	// Both are the same host, which is what makes borrowing the version safe.
	// A second server on a different host is a different API and is not
	// borrowed from -- see base_first_server_test.go, and Printify, which is
	// why that distinction exists.
	if got := BasePath(doc); got != "/v3" {
		t.Errorf("BasePath = %q, want /v3", got)
	}
}

func TestNoServersMeansNoBase(t *testing.T) {
	if got := BasePath(&Document{}); got != "" {
		t.Errorf("BasePath = %q, want empty", got)
	}
}
