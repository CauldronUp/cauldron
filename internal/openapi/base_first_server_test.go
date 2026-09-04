package openapi

import "testing"

// The first server decides, and an empty answer from it is an answer.
//
// BasePath used to take the first server that yielded a path at all, skipping
// any that yielded none. That is wrong the moment a description lists more than
// one server and the primary has no prefix: OpenAPI says the first server is
// the default, so a later one's prefix is not what a client of the first sends.
//
// Printify found it. Its description gained a second server carrying
// /api/pfy/public, on a different company's domain, listed after its own
// api.printify.com which has no path. Every declared path suddenly read as
// prefixed, both of that Recipe's routes were reported absent, and its
// fingerprint moved -- for a change to a server list that has nothing to do
// with the orders API the Recipe describes.
//
// Which is the exact failure this package exists to prevent: drift reported
// against something the provider did not change to anything the Recipe claims.
// A scan that cries wolf gets switched off, and then the change that mattered
// arrives unannounced.
func TestTheFirstServerDecidesEvenWhenItHasNoPath(t *testing.T) {
	doc := &Document{Servers: []Server{
		{URL: "https://api.printify.com"},
		{URL: "https://enterprise.printful.com/api/pfy/public"},
	}}

	if got := BasePath(doc); got != "" {
		t.Errorf("BasePath = %q, want empty -- the primary server has no prefix", got)
	}
}

// A templated authority still yields its path, which is Chargebee's shape.
func TestATemplatedAuthorityStillYieldsItsPath(t *testing.T) {
	doc := &Document{Servers: []Server{
		{URL: "{protocol}://{site}.{environment}:{port}/api/v2"},
	}}

	if got := BasePath(doc); got != "/api/v2" {
		t.Errorf("BasePath = %q, want /api/v2", got)
	}
}

// A placeholder inside the path yields nothing, because a prefix that varies
// per request is not a prefix -- and now that the first server is the only one
// consulted, that nothing is final rather than a reason to look further down
// the list.
func TestATemplatedPathYieldsNothingAndDoesNotFallThrough(t *testing.T) {
	doc := &Document{Servers: []Server{
		{URL: "https://{tenant}.example.com/{version}"},
		{URL: "https://api.example.com/v1"},
	}}

	if got := BasePath(doc); got != "" {
		t.Errorf("BasePath = %q, want empty -- the second server is not the default", got)
	}
}
