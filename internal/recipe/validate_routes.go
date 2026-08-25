// Route rules: what a Recipe answers, and whether anything can reach it.

package recipe

import (
	"fmt"
	"slices"
	"strings"
)

// validateRoutes checks the routes a Recipe declares, and the paths they share.
//
// Extracted from Validate, which was 991 lines of sequential rules in one
// function. Each of these appends through the same add, so the order and the
// wording of every problem are unchanged.
func (r *Recipe) validateRoutes(add func(string, ...any)) {
	if len(r.Routes) == 0 {
		add("at least one route is required")
	}

	seen := map[string]bool{}

	for i, route := range r.Routes {
		where := fmt.Sprintf("route %d (%s %s)", i+1, route.Method, route.Path)

		if route.Method == "" {
			add("%s: method is required", where)
		} else if route.Method != strings.ToUpper(route.Method) {
			add("%s: method must be upper case", where)
		}

		if !strings.HasPrefix(route.Path, "/") {
			add("%s: path must start with /", where)
		}

		// Selects is part of the identity: a GraphQL Recipe is several routes
		// on one path, told apart by what the query asks for, and two of them
		// selecting the same field would still be the duplicate this rule
		// exists to catch. The headers a route matches on are part of it for
		// the same reason -- an AWS Recipe is every operation on POST /,
		// told apart by X-Amz-Target.
		// And selects_body for the third time: Gemini is one path answering
		// two shapes depending on what the prompt asked for, which is two
		// routes rather than one repeated.
		key := route.Method + " " + route.Path + " " + route.Selects + " " + route.SelectsBody
		for _, name := range sortedKeys(route.MatchesHeader) {
			key += " " + name + "=" + route.MatchesHeader[name]
		}

		// And the query parameters, for the same reason again: two routes on
		// one path telling themselves apart by what the caller asked to expand
		// are two routes, not a duplicate.
		for _, name := range sortedKeys(route.MatchesQuery) {
			key += " ?" + name + "=" + route.MatchesQuery[name]
		}
		if seen[key] {
			add("%s: duplicate route", where)
		}
		seen[key] = true

		// A route that only ever fails reaches no operation and touches no
		// resource, so requiring either would be asking for a declaration that
		// cannot mean anything.
		if route.Error != "" {
			if _, ok := r.Errors[route.Error]; !ok {
				add("%s: unknown error %q", where, route.Error)
			}

			if route.Operation != "" || route.Resource != "" {
				add("%s: declares error %q and an operation or resource, and it can only ever do the first", where, route.Error)
			}

			continue
		}

		if route.Operation == "" {
			add("%s: operation is required", where)
		} else if !contains(validOperations, route.Operation) {
			add("%s: operation %q must be one of %s", where, route.Operation, strings.Join(validOperations, ", "))
		}

		if route.Resource == "" {
			add("%s: resource is required", where)
		} else if _, ok := r.Resources[route.Resource]; !ok {
			add("%s: unknown resource %q", where, route.Resource)
		}

		if !contains(validDeletedBody, route.DeletedBody) {
			add("%s: deleted_body %q must be one of %s", where, route.DeletedBody, strings.Join(validDeletedBody[1:], ", "))
		}

		if route.DeletedBody != "" && route.Operation != "delete" {
			add("%s: declares deleted_body on a %s, which never deletes anything", where, route.Operation)
		}

		// Offset and page numbering are only paging if the runtime reads the
		// parameter the caller sends, and there is no universal spelling: one
		// provider says per_page, another PageSize, another count. Without
		// both names the runtime reads "limit" and the style's own word, which
		// is right for some and wrong for plenty, and the wrongness is
		// invisible — the page size is ignored, one full page comes back, and
		// the caller's paging loop runs once and passes.
		//
		// 146 routes declared a style with no names, which was harmless while
		// nothing read the style and became a claim the moment something did.
		//
		// "-" is the third thing a name can be, beside a spelling and an
		// omission: the provider accepts no name for this at all. Datadog's
		// event listing fixes the page at a thousand and takes no size
		// parameter, so honouring an invented one would let code that sends
		// it work here and be ignored in production.
		if s := route.Pagination.Style; s == "offset" || s == "page" {
			if route.Pagination.LimitParam == "" || route.Pagination.CursorParam == "" {
				add("%s: declares %s paging without naming the provider's parameters, so the runtime would read `limit` and %q, which is a guess; set limit_param and cursor_param", where, s, s)
			}
		}

		if len(route.Returns) > 0 {
			if route.EmptyBody {
				add("%s: declares returns and empty_body, and there is no body for the fields to be in", where)
			}

			if resource, ok := r.Resources[route.Resource]; ok {
				// A listing whose entries carry no identifier is a listing
				// nothing can be fetched from, and a create that does not say
				// what it made leaves the caller with no way to ask again.
				// Both are easy to write by accident: returns names the
				// fields to keep, and the identifier is not one of them
				// unless it is asked for.
				//
				// Three Recipes shipped that way before this existed. Asana's
				// task listing handed back nothing but a name, Supabase's
				// create stopped telling anyone the project's ref, and every
				// case still passed, because no case asserted an identifier
				// on a trimmed route.
				//
				// A get may leave it out. Braze does: its details endpoint
				// answers without the identifier the caller just used to ask.
				// A resource that never puts its identifier on the wire says
				// so with field "-", and is not held to this.
				if route.Operation == "list" || route.Operation == "create" {
					if resource.ID.Field != "-" && !slices.Contains(route.Returns, "id") {
						add("%s: returns does not name id, so this %s answers with records nothing can address; add \"id\", or say field: \"-\" if this provider really does not send one",
							where, route.Operation)
					}
				}

				for _, field := range route.Returns {
					// Constants count. They are stamped onto the record and
					// trimmed with everything else, so a route that answers
					// with one has to name it. Braze stamps "message":
					// "success" on every object, and leaving it out of a
					// trimmed response drops the only field a Braze client
					// checks.
					_, constant := resource.Constants[field]

					if _, declared := resource.Fields[field]; !declared && !constant && field != "id" {
						// The identifier is held as "id" whatever the provider
						// calls it on the wire, and returns names the record's
						// own keys rather than the rendered ones. Writing the
						// wire name here is the obvious mistake, so say so
						// instead of reporting an unknown field.
						if field == resource.ID.Field {
							add("%s: returns %q, which is what the identifier is called on the wire; name it \"id\", because the trim runs before the rename",
								where, field)

							continue
						}

						add("%s: returns %q, which is not a field on resource %q", where, field, route.Resource)
					}
				}
			}
		}

		for _, name := range route.Scope {
			if name == "id" {
				add("%s: id cannot be a scope parameter", where)
				continue
			}

			if !strings.Contains(route.Path, "{"+name+"}") {
				add("%s: scope %q does not appear in the path", where, name)
			}

			if resource, ok := r.Resources[route.Resource]; ok {
				if _, declared := resource.Fields[name]; !declared {
					add("%s: scope %q is not a field on resource %q", where, name, route.Resource)
				}
			}
		}

		if in := route.Pagination.In; in != "" && in != "query" && in != "body" {
			add("%s: pagination in %q must be query or body", where, in)
		}

		if !contains(validPagination, route.Pagination.Style) {
			add("%s: pagination.style %q is not supported", where, route.Pagination.Style)
		}

		if route.IDFrom != "" {
			source, name, ok := strings.Cut(route.IDFrom, ":")

			switch {
			case route.IDFrom == "auth":
				// Names no location, because there is nothing to read. The
				// only thing that can be wrong with it is a path that also
				// carries an identifier, which would mean the request does
				// carry one after all.
				if strings.Contains(route.Path, "{id}") {
					add("%s: id_from auth says the request carries no identifier, and the path carries one", where)
				}
			case !ok || name == "":
				add("%s: id_from %q must look like query:channel, body:channel or auth", where, route.IDFrom)
			case source != "query" && source != "body":
				add("%s: id_from source %q must be query, body or auth", where, source)
			case strings.Contains(route.Path, "{id}"):
				add("%s: id_from and an {id} path parameter cannot both apply", where)
			}
		}

		if route.Operation != "list" && route.Operation != "create" &&
			route.IDFrom == "" && !strings.Contains(route.Path, "{id}") {
			add("%s: a %s needs an {id} in the path or an id_from", where, route.Operation)
		}

		// A hidden identifier is defensible on a resource fetched one at a
		// time by a key that lives in the path. It is not defensible on a
		// collection: a page of records with nothing to tell them apart
		// cannot be what the provider sends, because nobody could address
		// the second one.
		for i, f := range route.Filters {
			fwhere := fmt.Sprintf("%s: filters[%d]", where, i)

			if route.Operation != "list" {
				add("%s: filters only narrow a listing, and this is a %s", fwhere, route.Operation)
			}

			if f.Param == "" || f.Field == "" {
				add("%s: needs both a param and a field", fwhere)

				continue
			}

			resource, ok := r.Resources[route.Resource]
			if !ok {
				continue
			}

			if _, declared := resource.Fields[f.Field]; !declared {
				add("%s: filters on %q, which is not a field on resource %q", fwhere, f.Field, route.Resource)

				continue
			}

			if f.All != "" && f.All == f.Default {
				add("%s: the value that turns the filter off is also its default, so it never applies", fwhere)
			}

			for value, options := range f.Values {
				if len(options) == 0 {
					add("%s: expands %q into nothing, so asking for it would hide every record", fwhere, value)
				}

				if f.All != "" && value == f.All {
					add("%s: expands %q, which is also the value that turns the filter off", fwhere, value)
				}
			}

			if f.Default == "" {
				continue
			}

			// A default that names a bucket has to be checked against the
			// bucket's members rather than against the word, because no
			// record's field ever holds the word.
			covered := map[string]bool{}

			if options, grouped := f.Values[f.Default]; grouped {
				for _, option := range options {
					covered[option] = true
				}
			} else {
				covered[f.Default] = true
			}

			// A default that excludes nothing is a default nobody can see. It
			// is the same rule as a scope: the behaviour is only described
			// when a fixture holds a record the filter must hide, because
			// otherwise the listing looks identical with the filter and
			// without it, and a mutation removing it passes.
			excluded := false

			for _, resources := range r.Fixtures {
				for _, record := range resources[route.Resource] {
					if value, present := record[f.Field]; present && !covered[fmt.Sprint(value)] {
						excluded = true
					}
				}
			}

			if !excluded {
				add("%s: defaults to %q and no fixture holds a %s the default would hide, so the listing looks the same with the filter and without it",
					fwhere, f.Default, route.Resource)
			}
		}

		// A receipt has no fields, so nothing filters the request on the way
		// in and the whole payload comes back out. That is the helpful kind
		// of wrong: a client reads its own request back and believes the
		// provider confirmed it, locally, for as long as the test suite is
		// the only thing calling.
		if route.Operation == "create" && len(route.Returns) == 0 &&
			len(r.Resources[route.Resource].Fields) == 0 {
			add("%s: resource %q declares no fields, so this create would echo the request body back; set returns to say what the provider actually sends",
				where, route.Resource)
		}

		for _, name := range route.Beside {
			bwhere := fmt.Sprintf("%s: beside %q", where, name)

			if route.Operation != "list" {
				add("%s: only a listing carries other collections, and this is a %s", bwhere, route.Operation)
			}

			if name == route.Resource {
				add("%s: is the route's own resource, so it would be written twice into the same body", bwhere)

				continue
			}

			beside, ok := r.Resources[name]
			if !ok {
				add("%s: is not a resource in this Recipe", bwhere)

				continue
			}

			// Without its own name it lands wherever the route's resource
			// landed, and one collection overwrites the other in silence.
			if beside.Collection == "" {
				add("%s: needs a collection name of its own, or it would overwrite the collection beside it", bwhere)

				continue
			}

			if beside.Collection == r.Resources[route.Resource].Collection {
				add("%s: names the same collection as %q, so one would overwrite the other", bwhere, route.Resource)
			}

			// The scope is applied to every collection in the body, so a
			// resource that cannot be partitioned by it would be returned
			// whole to whoever asked for one slice of it.
			for _, field := range route.Scope {
				if _, declared := beside.Fields[field]; !declared {
					add("%s: has no %q field, so the scope this route applies would not narrow it", bwhere, field)
				}
			}
		}

		// A hidden identifier on a listing is usually wrong, and the
		// exception is a positional array: a collection where position is
		// the identity because there is nothing else. Cohere's embeddings
		// are exactly that, and it is the trap rather than an oversight --
		// the only thing tying a vector to its input is the index, so a
		// client that filters or reorders inputs anywhere in the pipeline
		// pairs the wrong vector with the wrong document.
		//
		// The test is whether anything fetches one on its own. If a route
		// does, the identifier exists and hiding it from the listing is
		// inconsistent. If none does, position is all there is and saying so
		// is the honest description.
		if route.Operation == "list" && r.Resources[route.Resource].ID.Field == "-" &&
			r.Resources[route.Resource].ID.CarriedBy == "" && fetchedByID(r, route.Resource) {
			add("%s: resource %q hides its identifier and %s fetches one by id, so the listing withholds an identifier that exists; name the field that carries it with carried_by, or stop hiding it",
				where, route.Resource, fetchRoute(r, route.Resource))
		}
	}
}
