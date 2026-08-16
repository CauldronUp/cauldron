package runtime

import (
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
}

// route is a compiled recipe.Route ready to match against a request.
type route struct {
	method   string
	segments []segment
	spec     recipe.Route
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
			spec:     spec,
		})
	}

	return out
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

				out = append(out, segment{
					param:  part[open+1 : close],
					prefix: part[:open],
					suffix: part[close+1:],
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
	parts := splitPath(path)

	var (
		best      route
		bestVars  map[string]string
		bestScore = -1
		found     bool
	)

	for _, candidate := range r.routes {
		if !strings.EqualFold(candidate.method, method) {
			continue
		}

		vars, score, ok := candidate.matches(parts)
		if !ok {
			continue
		}

		if score > bestScore {
			best, bestVars, bestScore, found = candidate, vars, score, true
		}
	}

	return best, bestVars, found
}

// allowedMethods reports which methods a path supports, so a mismatched method
// can return 405 with a useful Allow header rather than a bare 404.
func (r *router) allowedMethods(path string) []string {
	parts := splitPath(path)

	var out []string

	seen := map[string]bool{}

	for _, candidate := range r.routes {
		if _, _, ok := candidate.matches(parts); ok && !seen[candidate.method] {
			out = append(out, candidate.method)
			seen[candidate.method] = true
		}
	}

	return out
}

func (rt route) matches(parts []string) (map[string]string, int, bool) {
	if len(parts) != len(rt.segments) {
		return nil, 0, false
	}

	vars := map[string]string{}
	score := 0

	for i, seg := range rt.segments {
		if seg.param != "" {
			value := parts[i]

			if !strings.HasPrefix(value, seg.prefix) || !strings.HasSuffix(value, seg.suffix) {
				return nil, 0, false
			}

			value = value[len(seg.prefix) : len(value)-len(seg.suffix)]
			if value == "" {
				return nil, 0, false
			}

			vars[seg.param] = value

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
