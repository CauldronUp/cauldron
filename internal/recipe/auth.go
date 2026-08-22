// Authentication: how a provider decides a request is yours.

package recipe

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// Upstream records which real API version this Recipe targets. Without it, a
// Recipe silently rots as the provider moves on.
type Upstream struct {
	API  string `yaml:"api"`
	Docs string `yaml:"docs"`
}

// RequiredHeader is one header a request must carry.
//
// It reads from YAML either as a bare error name, meaning every request needs
// the header, or as a mapping with a methods list, meaning only those methods
// do. The second form exists because Greenhouse only wants On-Behalf-Of on a
// write: reads work without it, so an integration passes every test it has and
// then gets a 403 the first time it tries to change something.
type RequiredHeader struct {
	// Error names the error to raise when the header is missing.
	Error string `yaml:"error"`
	// Methods limits the requirement to those HTTP methods. Empty means all.
	Methods []string `yaml:"methods"`
}

// UnmarshalYAML accepts either a bare error name or the full mapping.
func (h *RequiredHeader) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		return value.Decode(&h.Error)
	}

	type plain RequiredHeader

	return value.Decode((*plain)(h))
}

// Applies reports whether the header is required for this HTTP method.
func (h RequiredHeader) Applies(method string) bool {
	if len(h.Methods) == 0 {
		return true
	}

	for _, m := range h.Methods {
		if strings.EqualFold(m, method) {
			return true
		}
	}

	return false
}

// Auth describes how the provider authenticates callers.
type Auth struct {
	// Scheme is one of: bearer, basic, header, query, none.
	//
	// A query credential travels in the URL, which is worth reproducing
	// precisely because it is a bad idea: URLs end up in access logs, browser
	// history and error reports. Trello and a good deal of older software do
	// it anyway, and an emulator that quietly accepted a header instead would
	// hide the exposure.
	Scheme string `yaml:"scheme"`
	// Header is the header carrying the credential, when scheme is header.
	Header string `yaml:"header"`
	// Param is the query parameter carrying the credential, when the scheme
	// is query.
	Param string `yaml:"param"`
	// Prefix is stripped from the credential before comparison, e.g. "Bearer ".
	Prefix string `yaml:"prefix"`
	// Credential says which half of a basic credential carries the secret:
	// "username" (the default, which is what Twilio does with the account SID)
	// or "password" (Mailgun, whose username is the constant "api"). Checking
	// the wrong half means a bad key is never rejected at all.
	Credential string `yaml:"credential"`
	// Keys are the credentials the emulator accepts. Test keys only — a Recipe
	// must never carry a real secret.
	Keys []string `yaml:"keys"`
	// Pattern accepts any credential matching this regular expression, for
	// schemes where the value is computed per request and cannot be compared
	// against a fixed list.
	//
	// AWS signs every request with SigV4, so the Authorization header is
	// different each time and there is no key to hold. Verifying the signature
	// would mean implementing the algorithm, which is not what this project is
	// for. Checking the shape catches the failure that actually happens —
	// credentials not configured, or the header missing entirely — and the
	// Recipe header has to say plainly that a wrongly signed request is
	// accepted. Silence about that would be worse than the gap.
	Pattern string `yaml:"pattern"`
}

// ValidAuthSchemes returns the credential schemes a Recipe may declare.
//
// Exported so the runtime's test suite can assert that every scheme the
// validator accepts is one the handler actually checks. The two are separate
// pieces of code that have to agree and nothing else makes them: adding a
// scheme here without adding a case there would silently authorise every
// request against every Recipe using it.
func ValidAuthSchemes() []string {
	out := make([]string, len(validSchemes))
	copy(out, validSchemes)

	return out
}
