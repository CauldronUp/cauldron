// Errors: how a provider says no.

package recipe

// Error is a named failure mode that `cauldron fault` can inject.
type Error struct {
	Status int    `yaml:"status"`
	Code   string `yaml:"code"`
	// Type is the provider's error category, which is often a much smaller set
	// than the codes. Stripe has four types and dozens of codes, and client
	// libraries switch on the type. Empty falls back to the code, which is
	// wrong often enough that every Recipe should set it.
	Type    string `yaml:"type"`
	Message string `yaml:"message"`
	// Style overrides the Recipe-wide error envelope for this failure alone,
	// because a provider can answer two shapes and the npm registry does.
	//
	// Checked against registry.npmjs.org on 2026-08-22: a package that does
	// not exist answers {"error":"Not found"}, and a version that does not
	// exist on a package that does answers the bare JSON string "version not
	// found: 99.99.99". Same status, same registry, one object and one
	// string. Code reading body.error off the second finds undefined, and
	// code that reports body.error as the reason reports "undefined".
	Style string `yaml:"style"`
	// Key overrides the Recipe-wide envelope key for this failure alone,
	// because a provider can answer two failures in two different places and
	// Shopify's GraphQL API does.
	//
	// A GraphQL request that cannot be served at all -- a bad token, a
	// throttled shop, a malformed query -- comes back as {"errors": [...]} at
	// the top level. A request that was served and refused on business
	// grounds comes back as
	// {"data": {"productCreate": {"userErrors": [...]}}}, nested under the
	// mutation's own name, with no top-level errors at all. Both are HTTP
	// 200, and a client that checks only one of the two channels misses every
	// failure in the other.
	//
	// A dotted name nests, the same way the Recipe-wide key does.
	Key string `yaml:"key"`
	// MessageField overrides the Recipe-wide field carrying the sentence, for
	// this failure alone.
	//
	// Shopify needs it and the reason is worth stating. A throttled GraphQL
	// request answers 200 with {"errors": [{"message": ...}]} -- an array of
	// objects. A request with a bad token answers 401 with
	// {"errors": "[API] Invalid API key or access token"} -- the same key,
	// holding a bare string. So errors[0].message reads the sentence on one
	// and reads the character "[" on the other, because indexing a string in
	// JavaScript succeeds, and .message on that is undefined. Nothing throws
	// and nothing is logged.
	//
	// Without this a Recipe could describe one of the two and would be
	// claiming the other does not happen.
	MessageField string `yaml:"message_field"`
	// CodeField overrides the Recipe-wide field carrying the code, for this
	// failure alone, and "-" removes it.
	//
	// The same asymmetry MessageField already fixed, on the other half of the
	// pair. A Recipe with two error envelopes could say where the prose lives
	// in each and could not say that one of them has no code at all -- so the
	// failure carried its own name as a code, which is worse than a wrong
	// code because it looks like a real one.
	//
	// VTEX is the case. Its OMS endpoints answer {"error": {code, message,
	// exception}} and its newer document endpoints answer an RFC 9110 problem
	// detail -- {"type", "title", "status", "traceId"} -- which has no code
	// anywhere in it. A client switching on a code has nothing to switch on,
	// and that is the thing worth serving.
	CodeField string            `yaml:"code_field"`
	// TypeField overrides the Recipe-wide field carrying the category, for
	// this failure alone, and "-" removes it.
	//
	// The last third of the same asymmetry MessageField and CodeField have
	// already closed. The error envelope names three fields and two of them
	// could move per failure, which Clio found by having a provider that
	// needs the third: it answers one shape to a request carrying no
	// credential and another to one carrying a wrong credential, and only
	// the second has a property named "type". That Recipe reached it by
	// turning type_field off across the whole file and pointing code_field
	// at "type" on the one error -- the right bytes, from a knob named for
	// something else, with a comment left behind to explain the
	// substitution.
	TypeField string            `yaml:"type_field"`
	Headers   map[string]string `yaml:"headers"`
	// Fields are extra body properties this failure carries, merged over the
	// Recipe-wide ones. Dropbox describes each failure with its own nested
	// union, so a single set of constants would make every error claim to be
	// the same one.
	Fields map[string]any `yaml:"fields"`
	// Empty serves a status line and no body at all.
	//
	// An error with no message falls back to generic wording, which is right
	// almost always -- a Recipe that forgot to write one should not silently
	// serve nothing. But six Recipes here have had to record that their
	// provider genuinely answers zero bytes, and could not reproduce it:
	// Raygun's 404, PropelAuth's 401, Backblaze's and Nile's and Tigris's
	// auto-generated 405s, and Snipcart's credential failures, which carry no
	// Content-Type either.
	//
	// Silence is a real answer and a hostile one. A client calling .json() on
	// nothing throws, exactly as it does on prose, and the Recipes that meet
	// it should be able to say so rather than serve a helpful sentence the
	// provider never sent.
	Empty bool `yaml:"empty"`
}
