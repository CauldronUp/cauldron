// Authentication: deciding whether a request carries a credential this
// Recipe accepts.

package runtime

import (
	"crypto/subtle"
	"net/http"
	"regexp"
	"strings"
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

// authorised checks the credential according to the Recipe's auth scheme.
func (s *Sandbox) authorised(r *http.Request) bool {
	auth := s.recipe.Auth

	// A Recipe that declares no keys accepts anything, so an author can model
	// routes first and tighten auth later.
	if auth.Scheme == "" || auth.Scheme == "none" || (len(auth.Keys) == 0 && auth.Pattern == "") {
		return true
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
	case "query":
		// The credential travels in the URL. Reproducing that exactly is the
		// point: a header-based fake would hide the fact that the secret ends
		// up in access logs and browser history.
		presented = r.URL.Query().Get(auth.Param)
	case "basic":
		user, password, ok := r.BasicAuth()
		if !ok {
			return false
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
		return false
	}

	if auth.Prefix != "" {
		if !strings.HasPrefix(presented, auth.Prefix) {
			return false
		}

		presented = strings.TrimPrefix(presented, auth.Prefix)
	}

	presented = strings.TrimSpace(presented)

	// A pattern, for a credential computed per request. AWS signs every call,
	// so there is no fixed value to compare and the shape is what can be
	// checked. That catches the failure that actually happens — credentials
	// not configured at all — and the Recipe says plainly that it does not
	// verify the signature.
	if auth.Pattern != "" {
		matched, err := regexp.MatchString(auth.Pattern, presented)

		return err == nil && matched
	}

	for _, key := range auth.Keys {
		// Constant time, which matters less here than almost anywhere: the
		// keys are published fixtures. It is here because this file is read as
		// a description of how a provider behaves, and a comparison that
		// leaks its answer byte by byte is not the pattern to hand somebody
		// who is about to go and write the real thing.
		if subtle.ConstantTimeCompare([]byte(presented), []byte(key)) == 1 {
			return true
		}
	}

	return false
}
