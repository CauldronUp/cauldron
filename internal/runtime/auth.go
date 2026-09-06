// Authentication: deciding whether a request carries a credential this
// Recipe accepts.

package runtime

import (
	"bytes"
	"crypto/subtle"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// missingHeader reports the first required header a request does not carry.
//
// A header may be required only for some methods, because that is how several
// providers behave: Greenhouse wants On-Behalf-Of on a write and ignores it on
// a read, so an integration that only ever reads in its tests meets the
// requirement for the first time in production.
func (s *Sandbox) missingHeader(r *http.Request) (header, errorName string, ok bool) {
	for name, required := range s.recipe.RequiredHeaders {
		if r.Header.Get(name) != "" || !required.Applies(r.Method) {
			continue
		}

		raise := required.Error
		if raise == "" {
			raise = "parameter_missing"
		}

		return name, raise, false
	}

	return "", "", true
}

// credential checks the credential according to the Recipe's auth scheme, and
// says which way it failed.
//
// The three failing verdicts exist so a Recipe can serve a provider's own
// distinct sentences. Most providers have only one, and for those every
// verdict resolves to the same error; see Auth.ErrorFor.
//
// What counts as absent is the carrier, not the value: the header or query
// parameter is not on the request at all. "Bearer " with nothing after it was
// sent by somebody, so it is malformed rather than absent, which is the line
// the providers themselves draw -- Make answers a bare missing header with
// "User is not logged in." and a broken one with "Invalid token header."
// credential judges a request against an Auth and any alternatives it names.
//
// A provider often accepts one secret through more than one carrier. Kickbox
// takes its key as a query parameter or a bearer header, Watchmode takes three
// channels, Clearbit takes a bearer or Basic with a blank password. Nine
// Recipes described that in prose and served one channel, so a client written
// against the other was refused by a fake and accepted by the provider -- the
// direction of wrongness that lets code ship broken.
//
// The alternatives are tried only when the primary carrier holds nothing.
// Somebody who put a credential in the primary and got it wrong has made a
// different mistake from somebody who used another channel, and the first
// verdict is the one about what they did.
//
// The verdict returned is the most informative one seen. If the primary is
// absent and an alternative was present and wrong, the caller presented
// something somewhere, and saying "absent" would describe a request nobody
// made.
func (s *Sandbox) credential(r *http.Request, auth recipe.Auth) (recipe.Verdict, string) {
	verdict, presented := s.judge(r, auth)
	if verdict == recipe.Accepted || len(auth.Also) == 0 {
		return verdict, presented
	}

	for _, alternative := range auth.Also {
		// An alternative inherits everything it does not name, so a Recipe
		// says only what differs -- usually just the carrier.
		merged := inherit(auth, alternative)
		merged.Also = nil

		alternate, sent := s.judge(r, merged)
		if alternate == recipe.Accepted {
			return alternate, sent
		}

		// The better-informed verdict wins. A carrier that could not read what
		// was sent knows less about it than one that could: Clearbit takes a
		// bearer or Basic in the same Authorization header, so a bearer value
		// is unreadable to the Basic check and merely wrong to the bearer one.
		// Reporting "malformed" there describes the carrier's confusion rather
		// than the caller's mistake.
		if informativeness(alternate) > informativeness(verdict) {
			verdict, presented = alternate, sent
		}
	}

	return verdict, presented
}

// inherit resolves an Auth written as a difference against the one it sits
// under, so a Recipe says only what changes.
//
// Two places write a partial Auth. An alternative carrier under auth.also
// usually names nothing but its own header. A route's own credential usually
// names nothing but its own keys -- Customer.io's app surface and its track
// surface take the same shape of bearer token and different secrets.
//
// The route case is why this exists as a function. That path replaced the
// Recipe's Auth wholesale, while the comment beside it said a route inherits
// what it does not name. Wholesale replacement is the more dangerous of the
// two readings: a route naming only after_routing, or only its own refusal
// sentence, would have replaced the Recipe's credential with one holding no
// scheme -- and an Auth with no scheme accepts every request. A route
// adjusting one detail of how its credential is checked would have switched
// the check off entirely.
func inherit(base, override recipe.Auth) recipe.Auth {
	merged := base

	if override.Scheme != "" {
		merged.Scheme = override.Scheme
	}

	if override.Header != "" {
		merged.Header = override.Header
	}

	if override.Param != "" {
		merged.Param = override.Param
	}

	// "-" clears it. A carrier that takes the bare secret under a primary that
	// takes a prefixed one had no way to say so -- Doppler's Basic channel
	// inherited "Bearer " and refused every key legitimately sent through it,
	// which made the fake stricter than the provider. That is the failure
	// auth.also exists to prevent, reintroduced one level down. "-" is what
	// clear already means for a list key and an identifier field.
	switch override.Prefix {
	case "":
	case "-":
		merged.Prefix = ""
	default:
		merged.Prefix = override.Prefix
	}

	if override.Credential != "" {
		merged.Credential = override.Credential
	}

	if len(override.Keys) > 0 {
		merged.Keys = override.Keys
	}

	if override.Pattern != "" {
		merged.Pattern = override.Pattern
	}

	if override.Shape != "" {
		merged.Shape = override.Shape
	}

	if override.ShapeError != "" {
		merged.ShapeError = override.ShapeError
	}

	if len(override.Unentitled) > 0 {
		merged.Unentitled = override.Unentitled
	}

	// The refusal sentences, which a second surface often does have its own
	// of: Unleash's frontend routes answer a different wording from its admin
	// ones for the same missing token.
	if override.AbsentError != "" {
		merged.AbsentError = override.AbsentError
	}

	if override.MalformedError != "" {
		merged.MalformedError = override.MalformedError
	}

	if override.RejectedError != "" {
		merged.RejectedError = override.RejectedError
	}

	if override.UnentitledError != "" {
		merged.UnentitledError = override.UnentitledError
	}

	// Ordering is a boolean, so an override can only turn it on -- there is no
	// way to write "and this one resolves the credential first" under a Recipe
	// that resolves routing first. Nothing has needed the other direction.
	if override.AfterRouting {
		merged.AfterRouting = true
	}

	// The alternatives are inherited like everything else, and replaced only
	// when the override names its own. A route that names its own key does not
	// stop the Recipe's secret arriving the ways the Recipe says it arrives --
	// and replacing the list outright made a route override quietly strip
	// every alternative, so the fake went stricter than its provider on
	// exactly the routes somebody had described more carefully.
	//
	// The alternative-carrier loop clears this itself after calling here,
	// because alternatives do not nest.
	if len(override.Also) > 0 {
		merged.Also = override.Also
	}

	return merged
}

func (s *Sandbox) judge(r *http.Request, auth recipe.Auth) (recipe.Verdict, string) {

	// A Recipe that declares no keys accepts anything, so an author can model
	// routes first and tighten auth later. Validation refuses that arrangement
	// for a Recipe that also names a scheme -- naming one reads as a claim the
	// credential is checked, and four Recipes made that claim while enforcing
	// nothing -- so what reaches here is a Recipe still being written.
	if auth.Scheme == "" || auth.Scheme == "none" || (len(auth.Keys) == 0 && auth.Pattern == "" && auth.Shape == "") {
		return recipe.Accepted, ""
	}

	var presented string

	switch auth.Scheme {
	case "bearer", "header":
		// One branch, because the two schemes differ only in convention: a
		// bearer token usually carries a prefix and a header credential
		// usually does not, and Prefix is applied below either way.
		//
		// bearer used to read Authorization and ignore auth.header entirely,
		// which made the field inert on 127 of the 130 Recipes that declare
		// it. All 127 say "Authorization", so nothing on the wire changes --
		// but a Recipe describing a bearer token in some other header was
		// silently served on Authorization instead, and a mutation renaming
		// the header could not be caught by any case.
		header := auth.Header
		if header == "" {
			header = "Authorization"
		}

		presented = r.Header.Get(header)

		if presented == "" {
			return recipe.Absent, presented
		}
	case "body":
		// The credential travels in the request body, beside the filters.
		// Canny does this: every read is a POST and the key is a field called
		// apiKey in the same object as boardID and limit.
		//
		// Reproducing it exactly matters for the same reason the query scheme
		// exists. A body credential cannot be set once as a default header, so
		// every call site carries it and the one that forgets is a well-formed
		// request that is refused; and the secret lands in anything that logs
		// request bodies, which is the logging switched on during an incident
		// by somebody not thinking about credentials. Serving it from a header
		// would hide both.
		presented = bodyValue(r, auth.Param)

		if presented == "" {
			return recipe.Absent, presented
		}
	case "query":
		// The credential travels in the URL. Reproducing that exactly is the
		// point: a header-based fake would hide the fact that the secret ends
		// up in access logs and browser history.
		presented = r.URL.Query().Get(auth.Param)

		if presented == "" {
			return recipe.Absent, presented
		}
	case "basic":
		user, password, ok := r.BasicAuth()
		if !ok {
			// No Authorization header at all is absent. One that is there and
			// will not parse -- some other scheme, or base64 that is not --
			// was sent by somebody who tried, and is malformed.
			if r.Header.Get("Authorization") == "" {
				return recipe.Absent, presented
			}

			return recipe.Malformed, presented
		}

		// Providers disagree about which half carries the secret. Twilio puts
		// the account SID in the username; Mailgun's username is the constant
		// "api" and the key is the password. Checking the wrong half means a
		// bad key is never rejected, so the Recipe says which.
		presented = user

		if auth.Credential == "password" {
			presented = password
		}
	default:
		// Validation should have refused an unknown scheme long before a
		// request arrived, and today it does: the empty and none cases are
		// handled above, and the validator only allows the four cased here. So
		// this is unreachable, and it used to accept every request anyway.
		//
		// Nothing couples that switch to the list of valid schemes, though.
		// Adding a fifth scheme to the list without adding a case here would
		// silently authorise every request against every Recipe using it, from
		// a one-line change, with no test that would fail. Failing closed is
		// the only safe direction for a branch whose whole job is to be
		// unreachable.
		return recipe.Rejected, presented
	}

	if auth.Prefix != "" {
		if !strings.HasPrefix(presented, auth.Prefix) {
			return recipe.Malformed, presented
		}

		presented = strings.TrimPrefix(presented, auth.Prefix)
	}

	presented = strings.TrimSpace(presented)

	// A scheme with nothing after it is malformed, not merely wrong. The
	// comment on this function has said so since the verdicts were written --
	// somebody sent "Bearer ", so it is not absent -- and the code did not:
	// the prefix matched, the value became empty, and an empty string fell
	// through to the key comparison and came out rejected.
	//
	// Bright Data is what found it. It answers "Auth method is not supported"
	// to a bare "Bearer ", which is its malformed sentence, where this
	// returned the verdict for a credential it had actually looked at.
	//
	// Providers do disagree about this one -- Fireworks answers "The API key
	// you provided is invalid." to the same request -- and that costs nothing
	// here, because a Recipe that does not name a malformed error falls
	// through to the same message it always did.
	if presented == "" {
		return recipe.Malformed, presented
	}

	// A shape, for the providers that read what a credential looks like before
	// they read what it says. Airtable, Modern Treasury, Opsgenie and Paddle
	// all answer one sentence to a credential of the wrong shape and another
	// to a well-formed one nobody issued, and this is the gate that lets a
	// Recipe say both: failing it is malformed, passing it goes on to the
	// comparison below rather than authenticating on the spot.
	//
	// An unparseable expression is a validation error, so it cannot arrive
	// here. If one somehow did, refusing is the safe direction: a shape that
	// cannot be evaluated has not been satisfied.
	if auth.Shape != "" {
		matched, err := regexp.MatchString(auth.Shape, presented)
		if err != nil || !matched {
			return recipe.Misshapen, presented
		}
	}

	// A pattern, for a credential computed per request. AWS signs every call,
	// so there is no fixed value to compare and the shape is what can be
	// checked. That catches the failure that actually happens — credentials
	// not configured at all — and the Recipe says plainly that it does not
	// verify the signature.
	if auth.Pattern != "" {
		matched, err := regexp.MatchString(auth.Pattern, presented)
		if err == nil && matched {
			return recipe.Accepted, presented
		}

		// A pattern describes the shape, so failing it is a shape failure.
		// Recipes using a pattern because the credential is computed per
		// request -- AWS signs every call -- have no fixed value to be wrong
		// about, so malformed is the only honest verdict available to them.
		return recipe.Malformed, presented
	}

	// Checked before the accepted keys, so a Recipe that names one key in both
	// places gets the refusal rather than the acceptance. The validator
	// refuses that Recipe outright, and the order here decides what happens
	// to one already in memory -- a fake that quietly authenticates a key its
	// own file says is not entitled is the worse of the two failures.
	for _, key := range auth.Unentitled {
		if subtle.ConstantTimeCompare([]byte(presented), []byte(key)) == 1 {
			return recipe.Unentitled, presented
		}
	}

	for _, key := range auth.Keys {
		// Constant time, which matters less here than almost anywhere: the
		// keys are published fixtures. It is here because this file is read as
		// a description of how a provider behaves, and a comparison that
		// leaks its answer byte by byte is not the pattern to hand somebody
		// who is about to go and write the real thing.
		if subtle.ConstantTimeCompare([]byte(presented), []byte(key)) == 1 {
			return recipe.Accepted, presented
		}
	}

	return recipe.Rejected, presented
}

// informativeness ranks what a verdict tells the caller about their own
// request, for choosing between the answers several carriers give.
//
// Absent says only that this carrier was not the one used. Malformed says
// something arrived that this carrier could not read. Misshapen says it was
// read and its shape ruled it out. Rejected and unentitled say it was read,
// its shape allowed it, and it was refused on its value -- which is the most
// any of them knows.
func informativeness(v recipe.Verdict) int {
	switch v {
	case recipe.Absent:
		return 0
	case recipe.Malformed:
		return 1
	case recipe.Misshapen:
		return 2
	case recipe.Rejected, recipe.Unentitled:
		return 3
	case recipe.Accepted:
		// Never compared -- an accepted credential returns before this.
		return 4
	}

	return 0
}

// bodyValue reads a named field from a JSON or form-encoded request body,
// leaving the body readable for whoever handles the request afterwards.
//
// Both encodings, because the APIs that put a key in the body are split
// between them: Canny documents form parameters, and plenty of others in the
// same shape send JSON.
func bodyValue(r *http.Request, name string) string {
	if name == "" {
		return ""
	}

	if value, ok := jsonBody(r)[name]; ok {
		if text, ok := value.(string); ok {
			return text
		}

		return fmt.Sprint(value)
	}

	if r.Body == nil {
		return ""
	}

	raw, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		return ""
	}

	r.Body = io.NopCloser(bytes.NewReader(raw))

	values, err := url.ParseQuery(string(raw))
	if err != nil {
		return ""
	}

	return values.Get(name)
}
