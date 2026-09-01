package runtime

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// segment is one part of a route path. A parameter segment captures its value,
// and may sit between literal text: Shopify and Twilio both write
// /orders/{id}.json, so a parameter is not always a whole segment.
type segment struct {
	literal string
	param   string
	prefix  string
	suffix  string
	// greedy says this parameter swallows the rest of the path, slashes and
	// all. Declared {name...}, which is Go's own ServeMux spelling, and
	// permitted only as the last segment for the same reason ServeMux permits
	// it only there: anywhere else there is nothing to say where it stops.
	//
	// Identifiers with slashes in them are commoner than the shape of a URL
	// path suggests. A DOI is "10.1145/3510003" and Semantic Scholar takes it
	// raw. A Guardian content id is "world/2023/jan/01/some-headline" -- five
	// segments that are one identifier. A Hugging Face repo is "org/name".
	// Percent-encoding the slash does not help: Semantic Scholar answers a
	// %2F with an error naming the encoding, so the client cannot escape its
	// way out of the problem either.
	greedy bool
}

// route is a compiled recipe.Route ready to match against a request.
type route struct {
	method   string
	segments []segment
	// slash records whether the declared path ends in one, because for some
	// providers it is part of the path rather than decoration.
	//
	// Sixty-one routes across five Recipes -- Sentry, PostHog, GoCardless
	// Bank Account Data, Uploadcare and now Saleor -- declare a trailing
	// slash, and every one of them is a Django or Django-shaped API where the
	// URL pattern is anchored and the slash is required. The router trimmed
	// it from both sides, so all sixty-one accepted the path without it and
	// the claim in the Recipe was decoration too.
	//
	// The failure that hides is the quiet kind. Django answers a slash-less
	// path with a 301 to the slash, and a client that follows a redirect
	// after a POST commonly drops the body -- so the request arrives, empty,
	// at the right URL, and the server refuses it for a reason that has
	// nothing to do with the redirect. Locally it worked.
	slash bool
	spec  recipe.Route
}

// router matches requests to routes.
type router struct {
	routes []route
}

// newRouter compiles a Recipe's routes. Compilation happens once at boot so
// matching a request never allocates a parser.
func newRouter(r *recipe.Recipe) *router {
	out := &router{}

	for _, spec := range r.Routes {
		out.routes = append(out.routes, route{
			method:   strings.ToUpper(spec.Method),
			segments: compilePath(spec.Path),
			slash:    len(spec.Path) > 1 && strings.HasSuffix(spec.Path, "/"),
			spec:     spec,
		})
	}

	return out
}

// mentions reports whether a GraphQL query names a field, as a whole word.
//
// A substring match is not good enough and the difference is not theoretical:
// selects "me" matched `query { viewer { name email } }`, because "name"
// contains it. Short root fields are common -- me, user, node, team -- and an
// accidental match sends a request to the wrong fixture, which is a bug this
// project exists to catch rather than commit.
//
// A GraphQL name is a letter or underscore followed by letters, digits or
// underscores, so a match is a real one when neither neighbour could be part
// of the same name.
func mentions(query, field string) bool {
	if query == "" || field == "" {
		return false
	}

	for from := 0; ; {
		at := strings.Index(query[from:], field)
		if at < 0 {
			return false
		}

		at += from
		end := at + len(field)

		if boundaryBefore(query, at) && boundaryAfter(query, end) {
			return true
		}

		from = at + 1
	}
}

// boundaryBefore reports whether the character before a match ends a word.
//
// A percent-escape is a delimiter wearing a disguise. In a form-encoded body
// the punctuation around a value is escaped, so a marker sitting immediately
// after a quote or a brace is preceded by the last hex digit of %22 or %7B --
// which is a digit, which is a name character, so the whole-word check refused
// a match that is plainly a whole word.
//
// Expensify is what found it: its API takes a JSON document inside a form
// field, so every structural character in that document arrives percent-
// encoded, and three cases had to be rewritten to route by a JSON body they
// do not really send.
func boundaryBefore(s string, at int) bool {
	if at == 0 {
		return true
	}

	if !isNameRune(s[at-1]) {
		return true
	}

	// Two hex digits and a percent sign: the escape is the boundary.
	return at >= 3 && s[at-3] == '%' && isHexDigit(s[at-2]) && isHexDigit(s[at-1])
}

// boundaryAfter is the plain check, and stays plain on purpose.
//
// The escape problem is one-sided. A percent-escape that FOLLOWS a match
// starts with '%', which is not a name character, so the ordinary test
// already calls it a boundary. Only the escape before a match hides its
// delimiter behind a hex digit.
func boundaryAfter(s string, end int) bool {
	return end == len(s) || !isNameRune(s[end])
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isNameRune(c byte) bool {
	return c == '_' ||
		(c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
}

func compilePath(path string) []segment {
	var out []segment

	for _, part := range strings.Split(strings.Trim(path, "/"), "/") {
		if part == "" {
			continue
		}

		if open := strings.Index(part, "{"); open >= 0 {
			if close := strings.Index(part[open:], "}"); close >= 0 {
				close += open

				name := part[open+1 : close]

				greedy := strings.HasSuffix(name, "...")
				if greedy {
					name = strings.TrimSuffix(name, "...")
				}

				out = append(out, segment{
					param:  name,
					prefix: part[:open],
					suffix: part[close+1:],
					greedy: greedy,
				})

				continue
			}
		}

		out = append(out, segment{literal: part})
	}

	return out
}

// match finds the route for a method and path, returning any captured
// parameters.
//
// Literal segments beat parameter segments, so /v1/customers/search would win
// over /v1/customers/{id} regardless of declaration order. Providers routinely
// have both, and matching on declaration order alone is a subtle trap.
func (r *router) match(method, path string) (route, map[string]string, bool) {
	return r.matchSelecting(method, path, "", "", nil, nil)
}

// matchSelecting is match with the request's GraphQL query in hand, so routes
// that share a path can be told apart by what the body asks for.
//
// A route declaring selects matches only when the query mentions that word,
// and beats an equally-scoring route that declares nothing, so a Recipe can
// have several selecting routes and one fallback for everything else.
func (r *router) matchSelecting(method, path, query, body string, params url.Values, headers http.Header) (route, map[string]string, bool) {
	parts := splitPath(path)

	// Whether the caller ended the path in a slash, which for some providers
	// is part of the path. See route.slash.
	slash := len(path) > 1 && strings.HasSuffix(path, "/")

	var (
		best     route
		bestVars map[string]string
		// No starting value, because a score can be negative: a greedy
		// segment scores below zero so that a route spelling a request out
		// always beats one that swallows everything. Seeded at -1, the first
		// greedy route to match tied with the sentinel and lost to it, and
		// the request fell through to a 405 naming a method it had already
		// used. found is the guard instead.
		bestScore int
		found     bool
	)

	for _, candidate := range r.routes {
		if !strings.EqualFold(candidate.method, method) {
			continue
		}

		if candidate.slash != slash {
			continue
		}

		vars, score, ok := candidate.matches(parts)
		if !ok {
			continue
		}

		if candidate.spec.Selects != "" {
			if !mentions(query, candidate.spec.Selects) {
				continue
			}

			// A route that asked for this query beats one that takes
			// anything, however the paths scored.
			score += 1000
		}

		// The same again for the providers whose answer depends on what the
		// caller wrote rather than on where they sent it.
		if candidate.spec.SelectsBody != "" {
			if !mentions(body, candidate.spec.SelectsBody) {
				continue
			}

			score += 1000
		}

		// The same bargain for the APIs that name the operation in a header
		// rather than in the path or the body.
		if len(candidate.spec.MatchesHeader) > 0 {
			if !carries(headers, candidate.spec.MatchesHeader) {
				continue
			}

			score += 1000
		}

		// And the same again for the APIs where asking for more changes the
		// shape rather than the endpoint. Clover's ?expand=lineItems and
		// Asana's opt_fields both work this way: one path, two shapes, and the
		// compact one is what a client gets for forgetting.
		if len(candidate.spec.MatchesQuery) > 0 {
			if !asksFor(params, candidate.spec.MatchesQuery) {
				continue
			}

			// Per parameter, so a route matching two of them beats a route
			// matching one. This used to be a flat bonus, which made two
			// matching routes tie and handed the request to whichever was
			// declared first -- silently, and with nothing in either route
			// saying it depended on the order.
			//
			// Datamuse is what found it. /words?rel_rhy=blue answers a
			// syllable count and no tags; adding md=dpsr answers both plus
			// definitions. Modelled as two routes, the second request matched
			// them both and got the first one's answer, so the Recipe declared
			// five fields and served three.
			score += 1000 * len(candidate.spec.MatchesQuery)
		}

		if !found || score > bestScore {
			best, bestVars, bestScore, found = candidate, vars, score, true
		}
	}

	return best, bestVars, found
}

// allowedMethods reports which methods a path supports, so a mismatched method
// can return 405 with a useful Allow header rather than a bare 404.
func (r *router) allowedMethods(path string) []string {
	parts := splitPath(path)
	slash := len(path) > 1 && strings.HasSuffix(path, "/")

	var out []string

	seen := map[string]bool{}

	for _, candidate := range r.routes {
		// A route that answers only certain queries says nothing about which
		// methods a path supports. Counting it turned an unmodelled GraphQL
		// query into 405 with Allow: POST, which tells a client to change
		// the method it already got right. Not being modelled is a 404.
		if candidate.spec.Selects != "" || candidate.spec.SelectsBody != "" ||
			len(candidate.spec.MatchesHeader) > 0 || len(candidate.spec.MatchesQuery) > 0 {
			continue
		}

		if candidate.slash != slash {
			continue
		}

		if _, _, ok := candidate.matches(parts); ok && !seen[candidate.method] {
			out = append(out, candidate.method)
			seen[candidate.method] = true
		}
	}

	return out
}

func (rt route) matches(parts []string) (map[string]string, int, bool) {
	// A greedy last segment relaxes the length check in one direction only:
	// it can stand for several path segments but never for none, so a route
	// ending {id...} still needs something after the literals to swallow.
	if last := len(rt.segments) - 1; last >= 0 && rt.segments[last].greedy {
		if len(parts) < len(rt.segments) {
			return nil, 0, false
		}
	} else if len(parts) != len(rt.segments) {
		return nil, 0, false
	}

	vars := map[string]string{}
	score := 0

	for i, seg := range rt.segments {
		if seg.param != "" {
			value := parts[i]

			if seg.greedy {
				value = strings.Join(parts[i:], "/")
			}

			if !strings.HasPrefix(value, seg.prefix) || !strings.HasSuffix(value, seg.suffix) {
				return nil, 0, false
			}

			// Both can hold at once on a string too short to contain them
			// both, when they overlap: /v1/aba{id}aba matched against
			// /v1/ababa starts with "aba" and ends with "aba" in five
			// characters. Slicing that panicked in the request path, and
			// net/http recovers per connection, so the client saw an EOF with
			// no response and the only trace was a stack on a stderr nobody
			// was reading. No shipped Recipe declares a two-sided parameter,
			// which is why this went unnoticed rather than why it was safe.
			if len(value) < len(seg.prefix)+len(seg.suffix) {
				return nil, 0, false
			}

			value = value[len(seg.prefix) : len(value)-len(seg.suffix)]
			if value == "" {
				return nil, 0, false
			}

			vars[seg.param] = value

			// A greedy parameter is the least specific thing a path can
			// declare, so it must never outscore a route that spells the same
			// request out. /paper/search and /paper/{id...} both match
			// /paper/search; the literal has to win, and would anyway on its
			// own segment score, but only until somebody writes a greedy
			// route with more literals in front of it. Scoring it negative
			// keeps that from being a race the specific route can lose.
			if seg.greedy {
				score--

				continue
			}

			// Literal text around a parameter is still evidence of a more
			// specific route, so /orders/{id}.json beats /orders/{id}.
			if seg.prefix != "" || seg.suffix != "" {
				score++
			}

			continue
		}

		if seg.literal != parts[i] {
			return nil, 0, false
		}

		score++
	}

	return vars, score, true
}

func splitPath(path string) []string {
	var out []string

	for _, part := range strings.Split(strings.Trim(path, "/"), "/") {
		if part != "" {
			out = append(out, part)
		}
	}

	return out
}

// carries reports whether a request's headers hold every value a route asked
// for.
//
// The comparison is exact rather than a substring: X-Amz-Target names one
// operation, and secretsmanager.ListSecrets is not secretsmanager.ListSecretsX
// however alike they read. Matching loosely here would route a request to a
// neighbouring operation and answer it confidently.
func carries(headers http.Header, want map[string]string) bool {
	if headers == nil {
		return false
	}

	for name, value := range want {
		if headers.Get(name) != value {
			return false
		}
	}

	return true
}

// asksFor reports whether the request's query carries every parameter a route
// declares, with the declared value among that parameter's comma-separated
// members.
//
// Membership rather than equality, because that is what the providers send:
// ?expand=lineItems,payments asks for two things in one parameter, and a
// route selecting either should answer it. A parameter with a single value is
// the one-member case of the same rule, so nothing needs to say which it is.
func asksFor(params url.Values, want map[string]string) bool {
	for name, value := range want {
		if !among(params[name], value) {
			return false
		}
	}

	return true
}

// among reports whether want appears as a comma-separated member of any of
// the supplied values.
func among(values []string, want string) bool {
	for _, value := range values {
		for _, member := range strings.Split(value, ",") {
			if strings.TrimSpace(member) == want {
				return true
			}
		}
	}

	return false
}
