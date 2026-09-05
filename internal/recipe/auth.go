// Authentication: how a provider decides a request is yours.

package recipe

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// Upstream records which real API version this Recipe targets. Without it, a
// Recipe silently rots as the provider moves on.
type Upstream struct {
	// API is the provider's own version label: v1, v2, 2026-06-30. Ninety-five
	// Recipes here say v1 and fifty say v2, which is the point -- it is the
	// provider's word, not ours.
	API string `yaml:"api"`
	// Docs is where a person reads about it.
	Docs string `yaml:"docs"`
	// Provider says whose API this is, which API alone cannot.
	//
	// Two Recipes both saying "v1" describe two different APIs, and until this
	// existed nothing in the format related a provider's v1 to its v2. The
	// absence was not academic: the MBTA's v2 clients reach realtime.mbta.com,
	// ClinicalTrials.gov's reach clinicaltrials.gov/ct2 and DataCite's minting
	// clients reach mds.datacite.org, and every one of those exclusions lived
	// in prose beside a detection entry where nothing could check it.
	//
	// It is optional, because most Recipes describe the only version there is
	// and inventing a name for a provider with one API adds nothing.
	Provider string `yaml:"provider"`
	// Supersedes names versions of the same API this one replaced, and what
	// identifies each on the wire.
	//
	// The host is the part that does work. A client is told to be a v2 client
	// by the host it talks to, and detection has nothing else to go on.
	Supersedes []Superseded `yaml:"supersedes"`
	// Spec is the URL of the provider's own machine-readable description.
	//
	// A Recipe is written once and the provider keeps moving. Its conformance
	// suite cannot notice, because that suite asserts what the Recipe says
	// rather than what the provider does: rename a field upstream and every
	// case still passes, green and wrong.
	Spec string `yaml:"spec"`
	// SpecHash fingerprints the parts of that description this Recipe claims,
	// as openapi.Fingerprint computes them.
	//
	// Deliberately not a checksum of the file. Providers republish these
	// constantly, and a scan that reports drift on every publish gets switched
	// off, after which the change that mattered arrives unannounced.
	SpecHash string `yaml:"spec_hash"`
	// SpecSeen is the date the fingerprint was taken, in the form every other
	// piece of dated evidence here uses.
	//
	// A fingerprint without one says a description once looked like this and
	// not when, so nobody can tell a Recipe checked yesterday from one checked
	// three years ago.
	SpecSeen string `yaml:"spec_seen"`
}

// Superseded is one earlier version of the same API, and what identifies it.
type Superseded struct {
	// Version is the provider's label for it, in the same words API uses.
	Version string `yaml:"version"`
	// Host is what a client of that version talks to, which is the only thing
	// detection can tell versions apart by.
	Host string `yaml:"host"`
	// Note is why it matters, for a person reading the Recipe.
	Note string `yaml:"note"`
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
	// Also names other carriers the same credential may arrive in.
	//
	// One host often takes one secret several ways. Kickbox accepts its key as
	// a query parameter or a bearer header. Watchmode takes three channels.
	// Clearbit takes a bearer or Basic with a blank password. Geoapify
	// documents a header as an alternative to its query parameter. Nine Recipes
	// recorded that and served one channel, which is the wrong direction to be
	// wrong in: a client using the other was refused by the fake and accepted
	// by the provider.
	//
	// An alternative inherits everything it does not name, so a Recipe usually
	// writes only the carrier. They are tried only when the primary holds
	// nothing, because getting the primary wrong is a different mistake from
	// using another channel, and the first verdict is the one about what the
	// caller actually did.
	//
	// This is for one credential in several places. Two genuinely different
	// credentials on one route -- ClickHouse accepts an undocumented JWT beside
	// its own scheme -- is a different thing and is not this.
	Also []Auth `yaml:"also"`
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
	// AbsentError names the errors entry raised when a request carries no
	// credential at all -- no header, no query parameter, nothing.
	//
	// Most providers answer the same way whether you sent nothing or sent
	// rubbish, and for those this stays empty and one message covers both.
	// A dozen Recipes here had to write a header comment apologising for
	// discarding the other half of a provider's answer, which is what a
	// missing feature looks like from the outside.
	//
	// The distinction is not cosmetic. Pipedream answers an absent credential
	// with 404 "record not found" -- describing the caller's own account as
	// missing -- and a junk one with a plain 401, so the more complete mistake
	// gets the better diagnosis. Turso reports an absent header as a malformed
	// token, "invalid number of segments" for a credential nobody sent.
	AbsentError string `yaml:"absent_error"`
	// MalformedError names the errors entry raised when a credential is
	// present but cannot be the right shape: the declared prefix is missing,
	// or a declared pattern does not match.
	//
	// Separate from AbsentError because several providers really do answer in
	// three tiers rather than two. Make: "User is not logged in." for no
	// header, "Invalid token header." for a malformed one, "Not authorized."
	// for a well-formed token it does not know. SingleStore: "Unauthorized",
	// then a hex-decoding complaint naming the offending byte, then a length
	// complaint. Collapsing those to one message loses the middle answer,
	// which is the one that tells a caller their credential never left the
	// keyboard intact.
	MalformedError string `yaml:"malformed_error"`
	// RejectedError names the errors entry raised when a credential is the
	// right shape and is not one this Recipe holds.
	//
	// This is the verdict a real integration actually hits -- a key that was
	// rotated, or copied from the wrong environment -- and the reason it needs
	// naming as much as the other two is FireHydrant, whose sentences run the
	// other way round. "This endpoint requires you to be authenticated." is
	// what an absent credential gets, which is also the generic wording, and
	// "The bearer token you provided is invalid or expired." is reserved for a
	// token that was genuinely presented. Without this field the interesting
	// half of that pair is the half that cannot be served.
	RejectedError string `yaml:"rejected_error"`
	// Unentitled names the keys that authenticate and are refused anyway,
	// and UnentitledError names the error they raise.
	//
	// A key on the wrong plan, without the scope, or asking for a region its
	// tier does not include. The provider knows who the caller is and says no
	// to this particular thing, which is the one refusal the other three
	// verdicts cannot express -- they are all about the credential, and this
	// is about the request.
	//
	// Keys named here are never accepted, so a Recipe cannot declare the same
	// key in both places and expect it to work; the validator refuses that.
	Unentitled      []string `yaml:"unentitled"`
	UnentitledError string   `yaml:"unentitled_error"`
	// AfterRouting checks the credential only once the request has matched a
	// route, so a path the provider does not have answers 404 and a method it
	// does not support answers 405, whether or not a credential was sent.
	//
	// The default is false, which checks first and refuses everything
	// unauthenticated -- the behaviour every Recipe was written against, and
	// the commoner arrangement.
	//
	// Providers genuinely disagree about this and 159 Recipes here have had to
	// say so in prose. Airbyte checks the credential first, always, and
	// answered a byte-identical 401 to every path real or invented. Fireworks
	// resolves the route and the method first, so an unrouted path 404s and a
	// wrong method 405s with nothing sent at all. Temporal routes first;
	// Mezmo's two hosts do opposite things to each other.
	//
	// This is a boolean, and the arrangements below are the ones it cannot
	// express. The list has been added to three times and the sentence above
	// it was never corrected: it said five arrangements and three recorded,
	// over a list of twelve, for long enough that individual Recipes wrote
	// down whichever count was current when they were written -- Akeyless
	// calls itself a seventh shape, Clio a sixth, Sendcloud a sixth. Those
	// numbers are all of different lists.
	//
	// So no count here. A list that grows is a list, and the number in front
	// of it is one more thing to forget:
	//
	//   Census   splits on the HTTP method. A GET reaches the gate and
	//            answers 401 in plain text; a DELETE on the same path finds
	//            no route for that verb and answers 404 in JSON.
	//   Rootly   splits on whether the route exists -- credential first for
	//            every path it has, router first for every path it does not,
	//            which is not something a caller can know in advance.
	//   NocoDB   splits on whether the route names an identifier. A route
	//            with an id resolves the id first and 404s with the id echoed
	//            back, credential or no credential; a route without one
	//            checks the credential first.
	//   Baserow  is a third thing again: URL pattern first, then the
	//            credential, then the specific id and the method. A wrong
	//            method on a real path answers 401 rather than 405.
	//   CloudTalk splits on which kind of mistake was made: it checks the
	//            credential before an unrouted path, and lets a wrong method
	//            skip the credential entirely to answer a real 405. Its own
	//            second host is routing-first throughout, so one vendor
	//            disagrees with itself.
	//   Flexport splits on whether the verb is one it declares anywhere. PUT,
	//            DELETE and OPTIONS hit a Rails CSRF page -- 422, HTML --
	//            whatever the credential, while an unrouted GET does not 404
	//            at all: it redirects to the marketing site.
	//   Sendcloud splits on whether any credential was presented at all --
	//            not on the path, not on the method, and not on whether the
	//            credential was any good. Two pieces of software answer here,
	//            its own application and an OAuth2 gateway in front of it,
	//            which is what makes that the dividing line.
	//   Timescale splits on which verdict it reached. An absent or malformed
	//            credential is refused before routing; a well-formed one that
	//            is simply wrong reaches the router and gets a real 405. So
	//            the order depends on how badly you got the credential wrong.
	//   Northflank splits within itself. Its shallow routes check the
	//            credential first; its nested build route checks the
	//            **body's shape** first and answers 400 to a request
	//            carrying no credential at all. So whether you are told
	//            your credential is missing depends on whether your JSON
	//            parsed, and the answer differs between two routes on one
	//            host -- which no Recipe-wide setting can express at all.
	//   LiveKit  splits on which credential was presented, which is Timescale's
	//            split inverted. A **bad** credential answers 401 everywhere,
	//            unrouted paths and wrong methods included; an **absent** one
	//            lets routing answer instead. So sending nothing gets a more
	//            informative reply than sending something wrong.
	//   Genius   splits within one verdict. Its absent-credential case depends
	//            on routing -- 401 on a real path, 403 on one that does not
	//            exist -- while its rejected-credential case is
	//            routing-blind and answers identically either way. Two
	//            different programs answer, which is why: its own application
	//            and the OAuth library in front of it.
	//   Canada Post is not about routing at all. It validates the tracking
	//            number's **shape** -- twelve, thirteen or sixteen
	//            alphanumeric characters, its three real formats -- before it
	//            examines the credential, and does so even when a wrong
	//            credential is presented. A domain constraint on a path
	//            parameter running ahead of authentication is a third axis
	//            this boolean does not have, and an id pattern here is
	//            checked strictly after the credential.
	//
	// Several of these have since been found again on other providers, which
	// is the useful signal: Sendcloud's shape turned up on Treasury Prime and
	// Clio, Northflank's on Melio and Google Maps, Genius's on Trakt, Papaya
	// Global and Cloudflare Stream, and Canada Post's on Cloudflare Stream and
	// Electricity Maps. A shape found twice is a mechanism worth building; a
	// shape found once is a note. Only the second kind is safe to leave here.
	//
	// The ones below would each need the credential's state to survive a
	// route match without gating it -- and Northflank's would need it to
	// survive body validation as well, per route. That is a larger change
	// than a boolean and is not obviously worth making for twelve providers.
	// What is worth doing is what those Recipes do: say so in the file, so
	// the next reader knows the emulator is approximating rather than that
	// the provider is simple.
	AfterRouting bool `yaml:"after_routing"`
}

// Verdict is what a credential check concluded, beyond accepted or not.
//
// A bool was enough while every rejection produced the same 401. It stopped
// being enough the moment a Recipe wanted to say which of a provider's two
// sentences applied, and returning the reason costs nothing: a Recipe that
// declares neither error still gets one message for all three verdicts, byte
// for byte what it got before.
type Verdict int

// The verdicts a credential check can reach.
const (
	// Accepted: the credential is one this Recipe holds.
	Accepted Verdict = iota
	// Absent: nothing was presented. Not an empty string -- the header or
	// query parameter carrying the credential is not there at all.
	Absent
	// Malformed: something was presented and it cannot be a credential,
	// because the declared prefix or pattern rules it out -- or because the
	// prefix is all there was. A bare "Bearer " was sent by somebody, so it
	// is not absent, and there is no value to have been wrong about, so it
	// is not rejected.
	Malformed
	// Rejected: the right shape, and not a credential this Recipe holds.
	Rejected
	// Unentitled: a credential this Recipe holds, for something it may not
	// have. Authentication succeeded and authorisation did not.
	//
	// This is a different question from the other three and providers keep
	// answering it separately. Apollo.io's own description declares a 403 for
	// a key that is valid and out of scope. Electricity Maps' free tier is
	// entitled to one zone and refuses the others. Fitbit publishes an
	// insufficient-scope message whose blanks its own troubleshooting table
	// never fills in. Three sightings, all documented, none of them reachable
	// with the three verdicts above -- because every one of those is decided
	// by the shape or the value of what was sent, and this one is decided by
	// what the sender is allowed to do.
	//
	// A Recipe reaches it by naming those keys in Auth.Unentitled. They
	// authenticate, so they are not rejected; they are refused anyway, so
	// they are not accepted. That makes a documented 403 something a
	// conformance case can send a request for, rather than something only a
	// fault injection can demonstrate the shape of -- and arming proves a
	// shape while saying nothing about what triggers it.
	Unentitled
)

// ErrorFor names the errors entry a verdict raises.
//
// Falling back to authentication_error for every verdict is what keeps this
// backward compatible: a Recipe declaring neither field cannot tell the three
// apart, which is exactly how every Recipe behaved before the fields existed.
func (a Auth) ErrorFor(v Verdict) string {
	switch v {
	case Absent:
		if a.AbsentError != "" {
			return a.AbsentError
		}
	case Malformed:
		if a.MalformedError != "" {
			return a.MalformedError
		}
	case Rejected:
		if a.RejectedError != "" {
			return a.RejectedError
		}
	case Unentitled:
		if a.UnentitledError != "" {
			return a.UnentitledError
		}
	case Accepted:
	}

	return "authentication_error"
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
