// Authentication: deciding whether a request carries a credential this
// Recipe accepts.

package runtime

import (
	"crypto/subtle"
	"net/http"
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
func (s *Sandbox) credential(r *http.Request, auth recipe.Auth) (recipe.Verdict, string) {

	// A Recipe that declares no keys accepts anything, so an author can model
	// routes first and tighten auth later.
	if auth.Scheme == "" || auth.Scheme == "none" || (len(auth.Keys) == 0 && auth.Pattern == "") {
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
