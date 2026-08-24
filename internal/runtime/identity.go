// Identifiers: reading the one a request names and the one a record answers to.

package runtime

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// identifier resolves the record id for a request. Most providers put it in
// the path; RPC-shaped ones like Slack put it in the query string or the body,
// which the Recipe declares with id_from.
func (s *Sandbox) identifier(matched route, r *http.Request, vars map[string]string) (string, bool) {
	return s.resolve(matched, vars, s.rawIdentifier(matched, r, vars))
}

// resolve turns the value a route addresses a record by into the record's own
// identifier, for the routes that address it by something else.
//
// A value that matches nothing is returned unchanged, so the ordinary
// not-found path reports it. That is the right answer for the failure this
// mostly models: a receipt handle from an earlier receive is stale, and SQS
// answers a delete with it by refusing.
func (s *Sandbox) resolve(matched route, vars map[string]string, value string) (string, bool) {
	field := matched.spec.LookupBy
	if field == "" || value == "" {
		return value, true
	}

	// Within the route's scope, not across the whole collection.
	//
	// A natural key is only as unique as the thing that owns it, which is
	// what a scope is for. AfterShip keys a tracking by carrier and number
	// together because carriers mint numbers in their own namespaces and they
	// collide: the same number really does exist under usps and under fedex,
	// on two different parcels. Resolving the number globally picked
	// whichever record the store returned first and then let the scope filter
	// reject it, so one carrier answered and the other 404'd -- and which one
	// worked depended on fixture order rather than on the request.
	where := map[string]any{field: value}

	for name, scopeValue := range s.scopeVars(matched, vars) {
		where[name] = scopeValue
	}

	page, err := s.store.ListWhere(matched.spec.Resource, where, "", 0)
	if err == nil && len(page.Records) > 0 {
		return fmt.Sprint(page.Records[0]["id"]), true
	}

	// Nothing carries that value in the field this route addresses records
	// by, so there is nothing here -- and falling back to an identifier
	// lookup would be worse than useless. Supabase found that: a project has
	// a deprecated id beside the ref every path takes, and a route declared
	// to look one up by ref answered 200 for a project addressed by the
	// identifier it is not supposed to accept. The route said which field it
	// reads; a value that is not in that field is not a match, however
	// familiar it looks.
	//
	// The value comes back regardless, because the caller asked about it and
	// the failure should name it.
	return value, false
}

func (s *Sandbox) rawIdentifier(matched route, r *http.Request, vars map[string]string) string {
	if matched.spec.IDFrom == "" {
		return vars["id"]
	}

	source, name, _ := strings.Cut(matched.spec.IDFrom, ":")

	switch source {
	case "auth":
		// The caller names nothing and the provider answers about whoever
		// asked. GitHub's /user, Stripe's /v1/account, Slack's auth.test and
		// Backblaze's b2_authorize_account all work this way, and none of
		// them could be described at all before this: the format could say
		// "the id is in the path" or "the id is in the body", and the whole
		// point of these routes is that the request does not carry one.
		//
		// A sandbox holds one identity per Recipe, so this is the only record
		// in the collection. The validator refuses a fixture with more than
		// one, because there would be nothing to choose between them with.
		page, err := s.store.List(matched.spec.Resource, "", 2)
		if err != nil || len(page.Records) == 0 {
			return ""
		}

		return fmt.Sprint(page.Records[0]["id"])
	case "query":
		return r.URL.Query().Get(name)
	case "body":
		body, err := decodeBody(r)
		if err != nil {
			return ""
		}

		// A dotted name nests, because a provider that puts the identifier in
		// the body does not always put it at the top of one. DynamoDB's
		// GetItem takes {"Key": {"id": {"S": "..."}}}: the identifier is
		// three levels down, in attribute-value form, and reading only the
		// top level finds nothing at all.
		if value, ok := nestedValue(body, name); ok {
			return fmt.Sprint(value)
		}
	}

	return ""
}

// identifierName is the property the identifier goes out under, which is "id"
// unless the Recipe says otherwise.
func identifierName(spec recipe.Resource) string {
	if spec.ID.Field == "" {
		return "id"
	}

	return spec.ID.Field
}

// numeric renders an identifier as a JSON number.
//
// A value that will not parse is left as it is rather than replaced or
// dropped. That only happens when a fixture holds an identifier the declared
// type cannot describe, and the useful outcome is a response that plainly
// disagrees with the declaration rather than a silent zero.
func numeric(value any) any {
	text, ok := value.(string)
	if !ok {
		return value
	}

	if n, err := strconv.ParseInt(text, 10, 64); err == nil {
		return n
	}

	return value
}
