// Validation: what a Recipe has to say before it is allowed to ship.

package recipe

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// indexedSegment matches a path segment naming an array position, e.g. to[0].
var indexedSegment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*(\[\d+\])+$`)

// plainSegment matches an ordinary object key.
//
// Hyphens are allowed because JSON keys have them: npm sends dist-tags and
// plenty of providers send header-shaped names. This rule refused the first
// hyphenated key it ever met, which is what a pattern written against the
// Recipes that happened to exist will do. A dot is still refused, because the
// runtime splits on it and a dot inside a segment would silently mean two
// levels rather than one.
var plainSegment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*$`)

var (
	namePattern    = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	datePattern    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	validSchemes    = []string{"bearer", "basic", "header", "query", "none"}
	validOperations = []string{"create", "get", "list", "update", "delete"}
	validPagination = []string{"", "cursor", "offset", "page"}
	// The categories a Recipe may declare. Kept short on purpose: a taxonomy
	// with sixty words in it is a taxonomy nobody can hold in their head, and
	// the point of this field is that somebody looking for a payments
	// emulator finds all of them at once.
	validCapabilities = []string{
		"payments", "banking", "accounting", "tax", "payroll",
		"email", "sms", "chat", "push", "voice",
		"auth", "identity", "crm", "support", "marketing", "brokerage",
		"commerce", "shipping", "storage", "database", "queue",
		"search", "cdn", "hosting", "observability", "analytics",
		"flags", "ci", "vcs", "issues", "docs",
		"calendar", "files", "media", "ai", "signing",
		"scheduling", "hr", "forms", "cms", "infrastructure",
	}
	validDeletedBody = []string{"", "receipt", "record", "flagged", "id", "empty"}
	validCursorURL   = []string{"", "absolute", "path"}
	validSigning     = []string{"", "none", "hmac-sha256"}
	validListStyles  = []string{"", "envelope", "bare", "wrapped", "map"}
	validErrStyles   = []string{"", "nested", "flat", "list", "string_list", "text", "string"}
	validCodeTypes   = []string{"", "string", "number"}
	validIDTypes     = []string{"", "string", "number"}
	validIDStyles    = []string{"", "prefixed", "numeric", "timestamp", "opaque", "uuid", "hex", "digits"}
	validFieldTypes  = []string{"", "string", "integer", "number", "boolean", "timestamp", "timestamp_ms", "timestamp_ms_string", "datetime", "msdate", "list"}
	// The types Cauldron fills in from the sandbox clock, and therefore the
	// only ones a stamped declaration can affect.
	timeFieldTypes = []string{"timestamp", "timestamp_ms", "timestamp_ms_string", "datetime", "msdate"}

	knownMethods = map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true,
		"DELETE": true, "HEAD": true, "OPTIONS": true,
	}
)

// Parse decodes and validates a Recipe from YAML bytes.
// safeDigits is how many figures a double holds exactly. 2^53 has sixteen, so
// fifteen is safe with one to spare, and an identifier longer than that has to
// stay a string.
const safeDigits = 15

// Validate checks the Recipe is internally consistent.
func (r *Recipe) Validate() error {
	var problems []string

	add := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	}

	if r.Name == "" {
		add("recipe name is required")
	} else if !namePattern.MatchString(r.Name) {
		add("recipe name %q must be lowercase words separated by hyphens", r.Name)
	}

	if r.Version == "" {
		add("version is required")
	} else if !versionPattern.MatchString(r.Version) {
		add("version %q must look like 1.2.3", r.Version)
	}

	if r.Upstream.API == "" {
		add("upstream.api is required so the Recipe records which API version it targets")
	}

	if r.Capability == "" {
		add("capability is required: a hundred Recipes are only findable by what they do")
	} else if !contains(validCapabilities, r.Capability) {
		add("capability %q must be one of %s", r.Capability, strings.Join(validCapabilities, ", "))
	}

	if r.Auth.Scheme != "" && !contains(validSchemes, r.Auth.Scheme) {
		add("auth.scheme %q must be one of %s", r.Auth.Scheme, strings.Join(validSchemes, ", "))
	}

	if r.Auth.Credential != "" && r.Auth.Credential != "username" && r.Auth.Credential != "password" {
		add("auth.credential %q must be username or password", r.Auth.Credential)
	}

	if r.Auth.Credential != "" && r.Auth.Scheme != "basic" {
		add("auth.credential only applies to the basic scheme")
	}

	if r.Auth.Scheme == "query" && r.Auth.Param == "" {
		add("auth.param is required when the scheme is query")
	}

	if r.Auth.Param != "" && r.Auth.Scheme != "query" {
		add("auth.param only applies to the query scheme")
	}

	if len(r.Resources) == 0 {
		add("at least one resource is required")
	}

	for _, name := range sortedKeys(r.Resources) {
		resource := r.Resources[name]

		if !contains(validIDStyles, resource.ID.Style) {
			add("resource %q has id.style %q, which must be one of %s", name, resource.ID.Style, strings.Join(validIDStyles[1:], ", "))
		}

		if resource.ID.CarriedBy != "" {
			if _, declared := resource.Fields[resource.ID.CarriedBy]; !declared {
				add("resource %q says its identifier is carried by %q, which is not a field on it", name, resource.ID.CarriedBy)
			}

			// Carried by something means not sent under its own name. A
			// resource that sends both is not describing this shape, and the
			// declaration would read as though it were.
			if resource.ID.Field != "-" {
				add("resource %q says its identifier is carried by %q and also sends it as %q; carried_by is for the providers that send it in one place only",
					name, resource.ID.CarriedBy, identifierFieldName(resource.ID.Field))
			}
		}

		if !contains(validIDTypes, resource.ID.Type) {
			add("resource %q has id.type %q, which must be one of %s", name, resource.ID.Type, strings.Join(validIDTypes[1:], ", "))
		}

		// A number on the wire has to be a number in the store too, or the
		// response disagrees with the declaration and the disagreement shows
		// up as a string where a number was promised. The styles that mint
		// something with letters in it cannot be numbers, and saying so here
		// is cheaper than finding out from a response.
		if resource.ID.Type == "number" && !contains([]string{"numeric", "digits"}, resource.ID.Style) {
			add("resource %q says its identifier is a number and mints it with the %s style, which does not produce one",
				name, styleName(resource.ID.Style))
		}

		// digits exists for identifiers that are numeric and must not be
		// parsed as numbers, because they exceed what a JavaScript number can
		// hold. Declaring one of those a number is asking for the rounding bug
		// the style was added to prevent.
		//
		// The length is what decides that, not the style. A double holds every
		// integer up to 2^53, which is sixteen figures, so fifteen is safe
		// with room to spare. Datadog's monitor ids are eight and really are
		// numbers -- its own description gives an integer for the id and for
		// the deleted_monitor_id a delete answers with -- and refusing that
		// forced the Recipe to send strings where Datadog sends numbers,
		// which is the disagreement this rule exists to catch, pointed the
		// other way.
		if resource.ID.Type == "number" && resource.ID.Style == "digits" && resource.ID.Length > safeDigits {
			add("resource %q mints a %d-digit numeric string and declares it a number, which is the rounding bug the digits style exists to avoid: a double holds %d figures",
				name, resource.ID.Length, safeDigits)
		}

		// A prefix is required for the prefixed style precisely because
		// applications prefix-match identifiers. Providers whose ids carry no
		// prefix say so by declaring the opaque style, rather than leaving it
		// ambiguous whether the omission was deliberate.
		if (resource.ID.Style == "" || resource.ID.Style == "prefixed") && resource.ID.Prefix == "" {
			add("resource %q must declare an id.prefix, or use id.style opaque if the provider's identifiers have none", name)
		}

		// A resource with no fields is almost always an unfinished one, and
		// the exception is narrow enough to name: a receipt. Knock's trigger
		// answers with a workflow_run_id and nothing else, and describing
		// that as a resource with one invented field would put a property on
		// the wire the provider never sends. A receipt is created and never
		// read back, so if anything reads this resource it is not one.
		if len(resource.Fields) == 0 && len(resource.Constants) == 0 && !onlyCreated(r, name) {
			add("resource %q has no fields", name)
		}

		for _, field := range sortedKeys(resource.Fields) {
			spec := resource.Fields[field]

			// A top-level field is named by its key, so an "as" that repeats
			// the key says nothing. One that differs is the only way to put a
			// name on the wire the record cannot use for itself, and there is
			// a real case: GitHub gives an issue a number and an id, the
			// number is what a path addresses it by, and the record therefore
			// has to key on the number -- which leaves "id" unavailable for
			// the field that carries GitHub's id.
			//
			// The rule was refusing every "as" here to catch the redundant
			// kind. It catches the redundant kind now.
			if spec.As != "" && spec.In == "" && spec.As == field {
				add("resource %q field %q sets as to its own name, which is what a top-level field is called anyway", name, field)
			}

			// A field that never reaches the wire cannot be renamed on the
			// way there, and cannot be null there either.
			if spec.In == "-" {
				if spec.As != "" {
					add("resource %q field %q is not sent and also declares as, so it names a wire field that does not exist", name, field)
				}

				if spec.NullWhenUnset {
					add("resource %q field %q is not sent and also declares null_when_unset, which is a claim about a key nobody receives", name, field)
				}
			}

			if spec.NullWhenUnset && spec.Default != nil {
				add("resource %q field %q is null when unset and also has a default, so it is never unset", name, field)
			}

			if spec.NullWhenUnset && spec.Required {
				add("resource %q field %q is null when unset and also required, so it is never unset", name, field)
			}

			// Every segment of an in: path has to be something the runtime can
			// walk. A segment it cannot parse became a literal key: a field
			// nested under "to[0]" was emitted under a property actually
			// spelled "to[0]", which is a shape no provider sends, produced in
			// silence, and invisible to a conformance suite with no case
			// naming it. That has now happened three times in different
			// disguises, so it is a rule rather than a habit.
			for _, segment := range splitNonEmpty(spec.In) {
				// "-" is the whole value rather than a path: the field is not
				// sent, so there is no segment to be well formed.
				if spec.In == "-" {
					continue
				}

				if plainSegment.MatchString(segment) || indexedSegment.MatchString(segment) {
					continue
				}

				add("resource %q field %q nests under %q, and the segment %q is not a name or a name with an array index, so it would be sent as a property literally spelled that",
					name, field, spec.In, segment)
			}

			// A field named parent_thing nested under parent emits
			// parent.parent_thing, which no provider sends. It happened
			// thirty-one times across six Recipes before anything checked,
			// and every conformance case about them passed because they
			// asserted the shape the emulator produced. A provider that
			// really does repeat the parent says so with an explicit as.
			// A field named the same as its parent emits parent.parent, which
			// is usually an accident: the author meant the field to be the
			// parent object itself and wrote in: rather than leaving it off.
			// Four Recipes were written with that shape in one sitting.
			//
			// It is not always an accident. ClickUp really does send a status
			// object with a status inside it, so an explicit as: says the
			// provider means it, exactly as the sibling rule below allows.
			if spec.In != "" && spec.As == "" && strings.EqualFold(field, spec.In) {
				add("resource %q field %q nests under %q and would be sent as %s.%s; drop the in: if the field is the parent object, or set as: %s if the provider really sends that",
					name, field, spec.In, spec.In, field, field)
			}

			if spec.In != "" && spec.As == "" && strings.HasPrefix(strings.ToLower(field), strings.ToLower(spec.In)+"_") {
				add("resource %q field %q nests under %q and would be sent as %s.%s, which repeats the parent; set as: %s, or set as explicitly if the provider really sends that",
					name, field, spec.In, spec.In, field, field[len(spec.In)+1:])
			}

			if spec.Type == "list" && spec.Default != nil {
				if _, ok := spec.Default.([]any); !ok {
					add("resource %q field %q is a list, but its default is not a sequence", name, field)
				}
			}

			if !contains(validFieldTypes, spec.Type) {
				add("resource %q field %q has type %q, which must be one of %s",
					name, field, spec.Type, strings.Join(validFieldTypes[1:], ", "))
			}

			// Only time fields are ever filled in automatically, so declaring
			// stamped on anything else does nothing and reads as though it
			// does something.
			if spec.Stamped != nil && !contains(timeFieldTypes, spec.Type) {
				add("resource %q field %q declares stamped, which only applies to %s",
					name, field, strings.Join(timeFieldTypes, ", "))
			}
		}
	}

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
		key := route.Method + " " + route.Path + " " + route.Selects
		for _, name := range sortedKeys(route.MatchesHeader) {
			key += " " + name + "=" + route.MatchesHeader[name]
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

	for i, c := range r.Conformance {
		if c.Arm != "" {
			if _, ok := r.Errors[c.Arm]; !ok {
				add("conformance[%d] %q: arms %q, which is not in the errors table", i, c.Name, c.Arm)
			}

			// Arming something and then expecting a status it does not
			// produce means the fault did nothing, and the case passes while
			// proving the opposite of what it says. That is the same class of
			// mistake as a case asserting its own request back.
			//
			// This used to compare against 400, on the reasoning that an
			// armed thing is a failure and a failure is a 4xx or a 5xx. Not
			// every armed thing is a failure: Snowflake's SQL API answers the
			// same endpoint with a 202 and a statement handle when the query
			// is slow, which is an alternate path rather than an error, and a
			// case arming it has to expect the 202 it installs. Comparing
			// against the armed entry's own status catches the real mistake
			// -- arm a 429, expect a 200 -- and permits the honest one.
			if armed, ok := r.Errors[c.Arm]; ok && c.Expect.Status > 0 && c.Expect.Status != armed.Status {
				add("conformance[%d] %q: arms %q, which answers %d, and expects %d, so what it installed changed nothing",
					i, c.Name, c.Arm, armed.Status, c.Expect.Status)
			}
		}

		if !c.Expect.NoBody {
			continue
		}

		if len(c.Expect.Body) > 0 || len(c.Expect.Matches) > 0 || len(c.Expect.Absent) > 0 {
			add("conformance[%d] %q: expect.no_body claims the response is empty, so there is nothing for body, matches or absent to look at", i, c.Name)
		}
	}

	// A field name the Recipe chooses is only a claim if a case asserts it
	// where the value exists. Asserting its absence on a last page holds
	// whatever the field happens to be called, so a Recipe could declare
	// cursor_field: next_page_token, be renamed to anything, and no case
	// would notice.
	//
	// Twenty-one of these were shipped before anything checked, across twenty
	// Recipes. Two of them turned out to have no successful list case at all:
	// every case touching the collection was checking a failure, so nothing
	// asserted what a working listing looks like.
	for _, declared := range []struct{ what, name string }{
		{"cursor_field", r.Responses.List.CursorField},
		{"count_field", r.Responses.List.CountField},
		{"has_more_field", r.Responses.List.HasMoreField},
		{"complete_field", r.Responses.List.CompleteField},
		{"final_field", r.Responses.List.FinalField},
	} {
		if declared.name == "" || assertsName(r.Conformance, declared.name) {
			continue
		}

		add("responses.list.%s is %q and no conformance case asserts that name where the value exists, so renaming it would break nothing",
			declared.what, declared.name)
	}

	if r.Responses.List.HasMoreField != "" && r.Responses.List.CompleteField != "" {
		add("responses.list declares both has_more_field and complete_field, which are the same flag with opposite senses")
	}

	// The whole point of a final field is that it is not the cursor. Naming
	// them the same thing would emit one key that changes meaning depending on
	// which page you are looking at, which is worse than either field alone.
	if f := r.Responses.List.FinalField; f != "" && f == r.Responses.List.CursorField {
		add("responses.list declares cursor_field and final_field as the same name %q, so one key would mean a page token on one page and a sync token on the next", f)
	}

	// A required header is an enforcement claim, and the only thing that
	// proves it is a case that leaves the header off and is refused.
	//
	// Enforcing Notion-Version is one of the reasons this project exists:
	// forgetting it is the classic Notion integration bug, and a fake that
	// waves it through lets code ship that fails on its first real call. A
	// Recipe could declare the header, never enforce it, and every case would
	// still be green -- they all send it.
	//
	// This is the same shape as the rule above about a field name no case
	// asserts. A name nothing exercises is a name that could be anything.
	for header := range r.RequiredHeaders {
		if omitsHeader(r.Conformance, header) {
			continue
		}

		add("required_headers declares %s and no conformance case omits it and is refused, so nothing here shows it is enforced rather than merely listed", header)
	}

	// A route naming an error to raise has to name one that exists, or the
	// only sign of the typo is a 404 in the shape it was overriding.
	for _, route := range r.Routes {
		if route.NotFound == "" {
			continue
		}

		if _, ok := r.Errors[route.NotFound]; !ok {
			add("%s %s declares not_found: %q and no error by that name is defined, so it would answer in the shape it was written to replace",
				route.Method, route.Path, route.NotFound)
		}
	}

	// A filter is a claim that a listing narrows itself when asked, and the
	// parameter's name is half of it. Every case listing without the parameter
	// exercises the default; only a case sending it pins what it is called.
	for _, route := range r.Routes {
		for _, f := range route.Filters {
			if f.Param == "" || sendsParam(r.Conformance, f.Param) {
				continue
			}

			add("%s %s declares a filter on %q and no conformance case sends it, so the parameter could be named anything and every case would still pass", route.Method, route.Path, f.Param)
		}
	}

	if r.Responses.Error.Key == "-" && r.Responses.Error.Style != "list" {
		add("responses.error.key is \"-\", which removes the envelope, but that only means anything for the list style")
	}

	if r.Responses.Error.Key == "-" && len(r.Responses.Error.Fields) > 0 {
		add("responses.error.key is \"-\", so there is no envelope for responses.error.fields to sit beside")
	}

	if r.Auth.Pattern != "" {
		if _, err := regexp.Compile(r.Auth.Pattern); err != nil {
			add("auth.pattern is not a valid regular expression: %v", err)
		}

		if len(r.Auth.Keys) > 0 {
			add("auth declares both keys and a pattern, and only one can decide whether a credential is accepted")
		}
	}

	if !contains(validCodeTypes, r.Responses.Error.CodeType) {
		add("responses.error.code_type %q must be one of %s", r.Responses.Error.CodeType, strings.Join(validCodeTypes[1:], ", "))
	}

	if !contains(validErrStyles, r.Responses.Error.Style) {
		add("responses.error.style %q must be one of %s", r.Responses.Error.Style, strings.Join(validErrStyles[1:], ", "))
	}

	if r.Responses.List.EntryStyle != "" && r.Responses.List.EntryStyle != "wrapped" {
		add("responses.list.entry_style %q must be wrapped", r.Responses.List.EntryStyle)
	}

	if !contains(validCursorURL, r.Responses.List.CursorURL) {
		add("responses.list.cursor_url %q must be one of %s", r.Responses.List.CursorURL, strings.Join(validCursorURL[1:], ", "))
	}

	if !contains(validListStyles, r.Responses.List.Style) {
		add("responses.list.style %q must be one of %s", r.Responses.List.Style, strings.Join(validListStyles[1:], ", "))
	}

	if r.Responses.List.Style == "wrapped" && r.Responses.List.Key != "" {
		// The key is a fallback for resources that do not name their own
		// collection. When every one of them does, it is unreachable, and an
		// unreachable declaration reads as a description of the provider
		// while describing nothing. Five Recipes shipped with one before this
		// rule existed, each found by mutating it and watching nothing fail.
		covered := len(r.Resources) > 0

		for _, name := range sortedKeys(r.Resources) {
			if r.Resources[name].Collection == "" {
				covered = false
			}
		}

		if covered {
			add("responses.list.key is %q and every resource names its own collection, so nothing reads it",
				r.Responses.List.Key)
		}
	}

	if r.Responses.List.Style == "wrapped" && r.Responses.List.Key == "" {
		for _, name := range sortedKeys(r.Resources) {
			if r.Resources[name].Collection == "" {
				add("resource %q needs a collection name, or responses.list.key must be set, because the list style is wrapped", name)
			}
		}
	}

	if !contains(validSigning, r.Webhooks.Signing.Scheme) {
		add("webhooks.signing.scheme %q is not supported", r.Webhooks.Signing.Scheme)
	}

	if r.Webhooks.Signing.Scheme == "hmac-sha256" && r.Webhooks.Signing.Header == "" {
		add("webhooks.signing.header is required when signing is enabled")
	}

	for _, event := range r.Webhooks.Events {
		if strings.TrimSpace(event) == "" {
			add("webhooks.events contains an empty event name")
		}
	}

	for _, name := range sortedKeys(r.Errors) {
		e := r.Errors[name]

		if e.Status < 100 || e.Status > 599 {
			add("error %q has status %d, which is not a valid HTTP status", name, e.Status)
		}
	}

	for _, header := range sortedKeys(r.RequiredHeaders) {
		required := r.RequiredHeaders[header]

		if name := required.Error; name != "" {
			if _, ok := r.Errors[name]; !ok {
				add("required_headers[%s] names error %q, which is not declared", header, name)
			}
		}

		for _, method := range required.Methods {
			if !knownMethods[strings.ToUpper(method)] {
				add("required_headers[%s] limits to method %q, which is not an HTTP method", header, method)
			}
		}
	}

	for _, fixtureName := range sortedKeys(r.Fixtures) {
		for _, resourceName := range sortedKeys(r.Fixtures[fixtureName]) {
			resource, ok := r.Resources[resourceName]
			if !ok {
				add("fixture %q seeds unknown resource %q", fixtureName, resourceName)
				continue
			}

			// A route that reads its identifier from the credentials answers
			// with the one record there is. Two of them and the runtime picks
			// the first, which is the seeding order, which is not a decision
			// anybody made.
			if len(r.Fixtures[fixtureName][resourceName]) > 1 {
				for _, route := range r.Routes {
					if route.Resource == resourceName && route.IDFrom == "auth" {
						add("fixture %q seeds %d records of %q, and %s %s answers about whoever asked, so there is nothing to choose between them with",
							fixtureName, len(r.Fixtures[fixtureName][resourceName]), resourceName, route.Method, route.Path)

						break
					}
				}
			}

			// A list field declared as one and then seeded with a scalar would
			// serve the scalar and pass every case written about it, so the
			// declaration has to mean something.
			//
			// An explicit null is allowed, because several providers send one
			// and it is a distinct state from both an absent field and an
			// empty array. AssemblyAI sends utterances: null on a transcript
			// that succeeded without speaker labels, so a rule that refused
			// null would block a true claim.
			for i, record := range r.Fixtures[fixtureName][resourceName] {
				// The identifier, which is not among the declared fields and
				// so is checked on its own.
				idField := resource.ID.Field
				if idField == "" || idField == "-" {
					idField = "id"
				}

				// An identifier the provider never echoes cannot disagree
				// with anything a client can see, so its shape is free.
				// Cohere keys embeddings e1 and e2 to find them again and
				// emits neither.
				if seeded, present := record[idField]; present && resource.ID.Field != "-" {
					if why, ok := seededID(resource.ID, seeded); !ok {
						add("fixture %q record %d of %q seeds an id that %s, so seeded records and created ones would not have the same shape",
							fixtureName, i, resourceName, why)
					}
				}

				for _, field := range sortedKeys(record) {
					// A key the resource does not declare is dropped on the
					// way in, so a fixture can set it, a reader can believe
					// it, and the emulator never sends it. That silence is
					// the problem: it fooled a mutation into passing, which
					// means it could fool a conformance case the same way.
					//
					// id is always allowed, and so is whatever the Recipe
					// renamed it to, because the identifier is not declared
					// among the fields.
					if _, declared := resource.Fields[field]; !declared {
						if _, constant := resource.Constants[field]; !constant &&
							field != "id" && field != resource.ID.Field {
							add("fixture %q record %d of %q sets %q, which is not a field on that resource, so it would be dropped in silence",
								fixtureName, i, resourceName, field)
						}
					}

					if resource.Fields[field].Type != "list" || record[field] == nil {
						if seeded, ok := seededAs(resource.Fields[field].Type, record[field]); !ok {
							add("fixture %q record %d of %q declares %q as %s and seeds %s, and nothing coerces it, so %s is what goes on the wire",
								fixtureName, i, resourceName, field, resource.Fields[field].Type, seeded, seeded)
						}

						continue
					}

					if _, isList := record[field].([]any); !isList {
						add("fixture %q record %d of %q sets list field %q to something that is neither a sequence nor null",
							fixtureName, i, resourceName, field)
					}
				}
			}
		}
	}

	for i, c := range r.Conformance {
		where := fmt.Sprintf("conformance case %d", i+1)
		if c.Name != "" {
			where = fmt.Sprintf("conformance case %q", c.Name)
		} else {
			add("%s: name is required", where)
		}

		if c.Source == "" {
			add("%s: source is required, so a reader can check the claim against the provider", where)
		}

		if c.Verified != "" && !datePattern.MatchString(c.Verified) {
			add("%s: verified %q must be a date like 2026-08-15", where, c.Verified)
		}

		if c.Request.Method == "" {
			add("%s: request.method is required", where)
		} else if c.Request.Method != strings.ToUpper(c.Request.Method) {
			add("%s: request.method must be upper case", where)
		}

		if !strings.HasPrefix(c.Request.Path, "/") {
			add("%s: request.path must start with /", where)
		}

		if len(c.Request.Form) > 0 && len(c.Request.JSON) > 0 {
			add("%s: a request sends form or json, not both", where)
		}

		if c.Fixture != "" {
			if _, ok := r.Fixtures[c.Fixture]; !ok {
				add("%s: unknown fixture %q", where, c.Fixture)
			}
		}

		if c.Expect.Status == 0 {
			add("%s: expect.status is required", where)
		}

		for field, pattern := range c.Expect.Matches {
			if _, err := regexp.Compile(pattern); err != nil {
				add("%s: expect.matches[%s] is not a valid regular expression: %v", where, field, err)
			}
		}

		for header, pattern := range c.Expect.HeaderMatches {
			if _, err := regexp.Compile(pattern); err != nil {
				add("%s: expect.header_matches[%s] is not a valid regular expression: %v", where, header, err)
			}
		}

		if c.Expect.BodyMatches != "" {
			if _, err := regexp.Compile(c.Expect.BodyMatches); err != nil {
				add("%s: expect.body_matches is not a valid regular expression: %v", where, err)
			}
		}

		// An absence is a claim too: "another project's issues are not visible"
		// is exactly the kind of thing worth asserting, and it has no positive
		// half. Leaving it out of this check rejected real cases.
		// A case whose every body assertion repeats a value it sent proves
		// nothing about the emulator: it asserts its own request came back.
		// Four of these were written and only caught by deliberately breaking
		// the Recipe and noticing the case still passed, so the check is here
		// rather than in somebody's memory.
		//
		// One non-echoed claim is enough to clear it, because then something
		// the emulator decided is under test. matches and absent count, since
		// neither can be satisfied by echoing.
		if w := c.Expect.Webhook; w != nil {
			if w.None && (w.Event != "" || len(w.Body) > 0 || len(w.Matches) > 0 || len(w.Absent) > 0) {
				add("%s: expect.webhook.none claims nothing was emitted, so there is nothing for event, body, matches or absent to look at", where)
			}

			if !w.None && w.Event == "" && len(w.Body) == 0 && len(w.Matches) == 0 && len(w.Absent) == 0 {
				add("%s: expect.webhook asserts nothing, which is not the same as claiming none was emitted", where)
			}

			if w.Event != "" && !contains(r.Webhooks.Events, w.Event) {
				add("%s: expects webhook %q, which the Recipe does not declare", where, w.Event)
			}
		}

		if echoesOnly(c) {
			add("%s: every assertion repeats a value the request sent, so the case passes whatever the emulator does; assert something it decides, or add a matches or absent claim", where)
		}

		// A webhook claim counts. It is a claim about what the request caused,
		// which is evidence of exactly the kind this rule exists to demand.
		if c.Expect.Status < 400 && !c.Expect.NoBody && c.Expect.BodyMatches == "" && c.Expect.Webhook == nil &&
			len(c.Expect.Body) == 0 && len(c.Expect.Matches) == 0 &&
			len(c.Expect.Headers) == 0 && len(c.Expect.HeaderMatches) == 0 && len(c.Expect.Absent) == 0 &&
			len(c.Expect.AbsentHeaders) == 0 {
			add("%s: a case that asserts nothing about the response is not evidence of anything", where)
		}
	}

	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}

	return nil
}

// fetchedByID reports whether any route fetches one record of a resource on
// its own, which is what makes an identifier something a client can use.
func fetchedByID(r *Recipe, resource string) bool {
	for _, route := range r.Routes {
		if route.Resource != resource {
			continue
		}

		if route.Operation == "get" || route.Operation == "update" || route.Operation == "delete" {
			return true
		}
	}

	return false
}

// fetchRoute names the first route that fetches one, for the message.
func fetchRoute(r *Recipe, resource string) string {
	for _, route := range r.Routes {
		if route.Resource != resource {
			continue
		}

		if route.Operation == "get" || route.Operation == "update" || route.Operation == "delete" {
			return route.Method + " " + route.Path
		}
	}

	return "another route"
}

// onlyCreated reports whether every route touching a resource creates it.
//
// Such a resource is a receipt: the provider hands back an identifier and
// nothing else, and there is nowhere to go and read it. Knock's trigger is
// one. A resource that is read anywhere is not, and a resource with no fields
// that is read anywhere is describing a response nobody can use.
func onlyCreated(r *Recipe, resource string) bool {
	touched := false

	for _, route := range r.Routes {
		if route.Resource != resource {
			continue
		}

		touched = true

		if route.Operation != "create" {
			return false
		}
	}

	return touched
}

// identifierFieldName renders the property an identifier goes out under, for a
// message, naming the default rather than reporting an empty string.
func identifierFieldName(field string) string {
	if field == "" {
		return "id"
	}

	return field
}

// styleName renders an id style for a message, naming the default rather than
// reporting an empty string nobody wrote.
func styleName(style string) string {
	if style == "" {
		return "prefixed"
	}

	return style
}

// sortedKeys keeps iteration deterministic, so validation problems always
// appear in the same order.
// seededAs reports whether a fixture value matches the type its field
// declares, and names what it is when it does not.
//
// A field's type is documentation. Nothing in the runtime reads it: the value
// a fixture seeds is the value that goes on the wire, with whatever type YAML
// gave it. So a Recipe can declare a field a string, seed it with something
// that is not one, and serve the wrong type for as long as no case looks.
//
// The reason to check it is YAML 1.1 rather than carelessness. An unquoted
// off, no, on or yes is a boolean, and those are all real values of real
// string fields: a DigitalOcean droplet's status is one of new, active, off
// and archive, and the fixture that seeded off to say a machine was powered
// down served false to everything that read it. The comment beside it
// explained that a cost report counting running machines would miss that
// droplet, and no case asserted on it, so the claim and the wire disagreed
// from the day it was written.
//
// Only the unambiguous directions are checked. A number field seeded with a
// whole number is fine, and a timestamp is a number whichever width it lands
// in.
func seededAs(declared string, value any) (string, bool) {
	if value == nil {
		return "", true
	}

	name := "something else"

	switch value.(type) {
	case bool:
		name = "a boolean"
	case string:
		name = "a string"
	case int, int64, uint64:
		name = "a whole number"
	case float64, float32:
		name = "a number"
	case []any:
		name = "a sequence"
	case map[string]any:
		name = "a mapping"
	}

	switch declared {
	case "string", "datetime":
		_, ok := value.(string)

		return name, ok
	case "boolean":
		_, ok := value.(bool)

		return name, ok
	case "integer", "number", "timestamp", "timestamp_ms":
		switch value.(type) {
		case int, int64, uint64, float64, float32:
			return name, true
		}

		return name, false
	}

	return name, true
}

// seededID reports whether a fixture's explicit identifier has the shape the
// same resource would mint. It returns what is wrong with it, and true when
// nothing is.
//
// A fixture that seeds an id in a shape the generator would never produce puts
// two shapes in one collection: the seeded records look one way and anything
// created during a run looks another. Client code that parses or
// prefix-matches an id -- which is common enough that ID exists as a concept
// here at all -- then works against the fixtures and fails against the
// records it made itself.
//
// It was a surviving mutation that found this. Monday declares ten-digit ids
// and seeds ten-digit ids, and changing the declaration to six broke nothing:
// every case reads a seeded record, so the declared shape was never on the
// wire and could have drifted from the provider without a single case
// noticing.
func seededID(id ID, value any) (string, bool) {
	text, ok := value.(string)
	if !ok {
		// A numeric id is checked as a number elsewhere; nothing to say here.
		return "", true
	}

	if text == "" {
		return "is empty", false
	}

	switch id.Style {
	case "uuid":
		if !looksUUID(text) {
			return "is not a UUID", false
		}

	case "hex":
		body, ok := strings.CutPrefix(text, id.Prefix)
		if !ok {
			return "does not start with " + id.Prefix, false
		}

		if !all(body, isHexDigit) {
			return "is not hexadecimal", false
		}

		if id.Length > 0 && len(body) != id.Length {
			return fmt.Sprintf("is %d characters, not the declared %d", len(body), id.Length), false
		}

	case "digits":
		body, ok := strings.CutPrefix(text, id.Prefix)
		if !ok {
			return "does not start with " + id.Prefix, false
		}

		if !all(body, isDecimalDigit) {
			return "is not all digits", false
		}

		// A snowflake is a number written as text, and a number does not
		// begin with a zero.
		if body[0] == '0' {
			return "begins with a zero, which a number written as text does not", false
		}

		if id.Length > 0 && len(body) != id.Length {
			return fmt.Sprintf("is %d digits, not the declared %d", len(body), id.Length), false
		}

	case "numeric":
		if !all(text, isDecimalDigit) {
			return "is not a number", false
		}

	case "", "prefixed":
		if id.Prefix != "" && !strings.HasPrefix(text, id.Prefix) {
			return "does not start with " + id.Prefix, false
		}
	}

	return "", true
}

func looksUUID(text string) bool {
	runs := strings.Split(text, "-")
	if len(runs) != 5 {
		return false
	}

	for i, want := range []int{8, 4, 4, 4, 12} {
		if len(runs[i]) != want || !all(runs[i], isHexDigit) {
			return false
		}
	}

	return true
}

// IsPath reports whether a name is a well-formed dotted path rather than a
// literal key that happens to contain a dot.
//
// The distinction is not academic. Dropbox names a field ".tag" -- the leading
// dot is part of the name, not a separator -- so treating every dotted name as
// a path turns it into an object under an empty key. A path is at least two
// segments and every one of them is a name.
func IsPath(name string) bool {
	segments := strings.Split(name, ".")
	if len(segments) < 2 {
		return false
	}

	for _, segment := range segments {
		if !plainSegment.MatchString(segment) && !indexedSegment.MatchString(segment) {
			return false
		}
	}

	return true
}
