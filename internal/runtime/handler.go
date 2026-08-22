package runtime

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/store"
)

// Providers genuinely disagree about list shapes: Stripe wraps in
// {object, data, has_more, url}, GitHub returns a bare array, Shopify nests
// under a resource key. The Recipe declares which, and the handler honours it.
// See listBody.

// ServeHTTP makes a Sandbox an http.Handler.
func (s *Sandbox) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	exchange := Exchange{
		At:     s.clock.Now(),
		Method: r.Method,
		Path:   r.URL.Path,
	}

	defer func() {
		s.log.record(exchange)
	}()

	// The network is consulted before the provider is, because that is the
	// order reality uses: a connection that never completes never reaches an
	// application that could have rate limited it.
	if conditions, ok := s.network.next(r.URL.Path); ok {
		exchange.Network = describeConditions(conditions)

		if conditions.Fatal() {
			// No response is coming, so do not build one. hangUp reports false
			// only when the server does not support hijacking, in which case
			// falling through to a normal response is better than pretending.
			if hangUp(w, conditions.Timeout) {
				exchange.Status = 0

				return
			}
		}

		if delay := s.network.delay(conditions); delay > 0 {
			time.Sleep(delay)
		}

		if degraded := newDegradedWriter(w, conditions); conditions.Bandwidth > 0 || conditions.Limit > 0 || conditions.Slice > 0 {
			defer func() {
				if degraded.truncated() {
					exchange.Network += " (truncated)"
				}
			}()

			w = degraded
		}
	}

	// Faults are evaluated before anything else the provider does, including
	// auth. A rate limit that only fires on well-formed authenticated requests
	// would not be a faithful rate limit.
	if name, ok := s.faults.next(r.URL.Path); ok {
		exchange.Fault = name
		exchange.Status = s.writeRecipeError(w, name)

		return
	}

	if !s.authorised(r) {
		exchange.Status = s.writeRecipeError(w, "authentication_error", 401, "authentication_required", "Invalid API key provided.")

		return
	}

	// A header the provider insists on is part of the contract. Forgetting
	// Notion-Version is the classic Notion integration bug, and a fake that
	// waves it through lets code ship that fails on its first real call.
	if header, name, ok := s.missingHeader(r); !ok {
		exchange.Status = s.writeRecipeError(w, name, 400, "missing_header",
			"The "+header+" header is required.", header)

		return
	}

	matched, vars, ok := s.router.matchSelecting(r.Method, r.URL.Path, graphQLQuery(r), r.Header)
	if !ok {
		if allowed := s.router.allowedMethods(r.URL.Path); len(allowed) > 0 {
			w.Header().Set("Allow", strings.Join(allowed, ", "))
			exchange.Status = s.writeRecipeError(w, "method_not_allowed", 405, "method_not_allowed", "This method is not supported on this path.", r.Method+" "+r.URL.Path)

			return
		}

		exchange.Status = s.writeRecipeError(w, "unknown_route", 404, "unknown_route", "Unrecognised request URL.", r.URL.Path)

		return
	}

	exchange.Resource = matched.spec.Resource
	exchange.Op = matched.spec.Operation

	// A route that exists only to fail. Jira's old search endpoint answers 410
	// Gone to a path thousands of integrations still call, and the distinction
	// between that and a 404 is the whole message: the path was right, and it
	// is not coming back. Routing it as an ordinary unknown path would make an
	// emulator that quietly disagrees with the provider about which failure
	// this is.
	if name := matched.spec.Error; name != "" {
		exchange.Status = s.writeRecipeError(w, name)

		return
	}

	switch matched.spec.Operation {
	case "create":
		exchange.Status = s.create(w, r, matched, vars)
	case "get":
		exchange.Status = s.get(w, r, matched, vars)
	case "update":
		exchange.Status = s.update(w, r, matched, vars)
	case "delete":
		exchange.Status = s.delete(w, r, matched, vars)
	case "list":
		exchange.Status = s.list(w, r, matched, vars)
	default:
		exchange.Status = s.writeRecipeError(w, "unsupported_operation", 500, "unsupported_operation", "This operation is not implemented.")
	}
}

func (s *Sandbox) create(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	record, err := decodeBody(r)
	if err != nil {
		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", "Request body could not be parsed.")
	}

	record = s.flatten(matched.spec.Resource, record)

	for field, value := range s.scopeVars(matched, vars) {
		record[field] = value
	}

	s.applyDefaults(matched.spec.Resource, record)

	if missing := s.missingRequired(matched.spec.Resource, record); len(missing) > 0 {
		return s.writeRecipeError(
			w, "parameter_missing", 400, "parameter_missing",
			"Missing required parameter: "+strings.Join(missing, ", ")+".",
			strings.Join(missing, ", "),
		)
	}

	created, err := s.store.Create(matched.spec.Resource, record)
	if err != nil {
		// A caller naming an identifier that already exists is a conflict, and
		// every provider says so with its own status. A Recipe declaring a
		// conflict gets its own shape; the fallback is the 409 the situation
		// deserves rather than the 400 an internal message used to leak
		// through.
		if errors.Is(err, store.ErrConflict) {
			return s.writeRecipeError(w, "conflict", http.StatusConflict, "conflict", "A record with that identifier already exists.")
		}

		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", err.Error())
	}

	s.emitFor(matched.spec.Resource, "created", created)

	s.writeRouteHeaders(w, matched, created)

	return s.writeRecord(w, matched, created)
}

// writeRecord answers with a record, honouring the route's declared status and
// empty_body.
//
// Both used to be read on creates alone, so a Recipe could declare that its
// provider answers an update with 204 and nothing at all, and be quietly
// ignored. Salesforce does exactly that on PATCH and DELETE, and a client
// calling .json() on the real response throws — which is precisely the bug an
// emulator that helpfully returned the updated record would hide.
func (s *Sandbox) writeRecord(w http.ResponseWriter, matched route, record store.Record) int {
	status := matched.spec.Status
	if status == 0 {
		status = http.StatusOK
	}

	if matched.spec.EmptyBody {
		w.WriteHeader(status)

		return status
	}

	record = trim(matched.spec, record)

	body := s.resourceBody(s.recipe.EnvelopeFor(matched.spec), matched.spec.IDAs, matched.spec.Resource, record)

	// Constants this route adds beside the record, which a create needs as
	// much as a listing does. Neon answers a branch create with the branch,
	// the operations it started, and the connection strings -- and the
	// operations are the only thing in that body that says when the branch
	// can actually be used. A response carrying just the record would be the
	// helpful kind of wrong: it looks finished.
	//
	// Both spellings of the same map are handled, because store.Record is a
	// named map[string]any and a type assertion to the unnamed one does not
	// match it. A route's constants were dropped in silence on every response
	// that is not wrapped: the wrapped path builds a plain map and worked,
	// which is why Neon's create was fine and it was Intercom's delete
	// receipt -- the only other route declaring any -- that found this.
	if len(matched.spec.Fields) > 0 {
		switch object := body.(type) {
		case map[string]any:
			body = withFields(object, matched.spec.Fields)
		case store.Record:
			body = withFields(object, matched.spec.Fields)
		}
	}

	writeJSON(w, status, body)

	return status
}

// identifier resolves the record id for a request. Most providers put it in
// the path; RPC-shaped ones like Slack put it in the query string or the body,
// which the Recipe declares with id_from.
func (s *Sandbox) identifier(matched route, r *http.Request, vars map[string]string) (string, bool) {
	return s.resolve(matched, s.rawIdentifier(matched, r, vars))
}

// resolve turns the value a route addresses a record by into the record's own
// identifier, for the routes that address it by something else.
//
// A value that matches nothing is returned unchanged, so the ordinary
// not-found path reports it. That is the right answer for the failure this
// mostly models: a receipt handle from an earlier receive is stale, and SQS
// answers a delete with it by refusing.
func (s *Sandbox) resolve(matched route, value string) (string, bool) {
	field := matched.spec.LookupBy
	if field == "" || value == "" {
		return value, true
	}

	page, err := s.store.ListWhere(matched.spec.Resource, map[string]any{field: value}, "", 0)
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

// nestedValue walks a dotted path to a leaf in a decoded body.
func nestedValue(body store.Record, path string) (any, bool) {
	var current any = map[string]any(body)

	for _, segment := range strings.Split(path, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		current, ok = object[segment]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

// trim reduces a record to the fields a route answers with. The trim happens
// before shaping, so a kept field still nests where the Recipe says it does.
//
// Jira's create hands back an id and a key and nothing else, so anything
// reading created.fields.summary gets undefined. That is the shape this was
// written for, and for a long time it was the only one it served: the trim
// lived inside the create path, so a Recipe declaring returns on a get, an
// update or a list was declaring nothing. The validator accepted it, the
// emulator ignored it, and the Recipe read as though a smaller body had been
// described when the full record was still going out.
//
// A listing is where that matters most. Braze's /campaigns/list hands back
// five properties and /campaigns/details hands back a different, larger set,
// so code reading channels off a list entry gets undefined against the real
// API and a populated array against an emulator that ignored the trim. GitHub
// does the same thing between its issue list and its issue detail, and so do
// enough providers that a listing being a summary should be the assumption
// rather than the surprise.
func trim(spec recipe.Route, record store.Record) store.Record {
	if len(spec.Returns) == 0 {
		return record
	}

	kept := make(store.Record, len(spec.Returns))

	for _, field := range spec.Returns {
		if value, ok := record[field]; ok {
			kept[field] = value
		}
	}

	return kept
}

func (s *Sandbox) get(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	id, addressable := s.identifier(matched, r, vars)
	if !addressable {
		return s.notFound(w, store.ErrNotFound, matched.spec.Resource, id)
	}

	record, err := s.store.GetBy(matched.spec.Resource, id, s.recipe.Resources[matched.spec.Resource].Alias)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id)
	}

	// A record that exists but belongs to another scope must look absent, not
	// forbidden. Leaking existence across tenants is a real provider bug we do
	// not want to teach an application to rely on.
	if !store.Matches(record, s.scopeVars(matched, vars)) {
		return s.writeRecipeError(w, "resource_missing", 404, "resource_missing",
			"No such "+matched.spec.Resource+": "+id+".",
			matched.spec.Resource+": "+id)
	}

	writeJSON(w, http.StatusOK, s.resourceBody(s.recipe.EnvelopeFor(matched.spec), matched.spec.IDAs, matched.spec.Resource, trim(matched.spec, record)))

	return http.StatusOK
}

func (s *Sandbox) update(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	id, addressable := s.identifier(matched, r, vars)
	if !addressable {
		return s.notFound(w, store.ErrNotFound, matched.spec.Resource, id)
	}

	changes, err := decodeBody(r)
	if err != nil {
		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", "Request body could not be parsed.")
	}

	changes = s.flatten(matched.spec.Resource, changes)

	existing, err := s.store.Get(matched.spec.Resource, id)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id)
	}

	if !store.Matches(existing, s.scopeVars(matched, vars)) {
		return s.writeRecipeError(w, "resource_missing", 404, "resource_missing",
			"No such "+matched.spec.Resource+": "+id+".",
			matched.spec.Resource+": "+id)
	}

	updated, err := s.store.Update(matched.spec.Resource, id, changes)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id)
	}

	s.emitFor(matched.spec.Resource, "updated", updated)

	return s.writeRecord(w, matched, updated)
}

func (s *Sandbox) delete(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	id, addressable := s.identifier(matched, r, vars)
	if !addressable {
		return s.notFound(w, store.ErrNotFound, matched.spec.Resource, id)
	}

	record, err := s.store.Get(matched.spec.Resource, id)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id)
	}

	if !store.Matches(record, s.scopeVars(matched, vars)) {
		return s.writeRecipeError(w, "resource_missing", 404, "resource_missing",
			"No such "+matched.spec.Resource+": "+id+".",
			matched.spec.Resource+": "+id)
	}

	if err := s.store.Delete(matched.spec.Resource, id); err != nil {
		return s.notFound(w, err, matched.spec.Resource, id)
	}

	s.emitFor(matched.spec.Resource, "deleted", record)

	// Nothing, unless the Recipe says its provider answers with something.
	// This used to fabricate Stripe's receipt for every provider, using keys
	// no Recipe declares, so a client calling .json() on the response worked
	// here and threw against the thirty-odd APIs that answer 204 and an empty
	// body.
	switch matched.spec.DeletedBody {
	case "receipt":
		return s.writeRecord(w, matched, store.Record{
			"id":      id,
			"object":  matched.spec.Resource,
			"deleted": true,
		})
	case "record":
		return s.writeRecord(w, matched, record)
	case "flagged":
		// The receipt without Stripe's discriminator. Only the object key was
		// ever Stripe's, and a provider that names it something else says so
		// with the route's own fields -- Intercom answers {type: contact, id,
		// deleted: true} and declares that type on the route.
		//
		// The resource's constants are deliberately not copied here. Intercom
		// stamps workspace_id on every contact and sends none of it on a
		// delete, so copying them all would put a field on the wire that the
		// provider does not, which is the thing this format exists to avoid.
		return s.writeRecord(w, matched, store.Record{"id": id, "deleted": true})
	case "id":
		// The identifier and nothing else, which is what Cloudflare answers
		// a delete with -- {"result": {"id": ...}} once the envelope is on.
		//
		// Under its own name when the route says so. Datadog answers a
		// monitor delete with deleted_monitor_id, and a client reading that
		// finds nothing if the key is the resource's ordinary one.
		if key := matched.spec.DeletedKey; key != "" {
			return s.writeRaw(w, matched, map[string]any{key: id})
		}

		return s.writeRecord(w, matched, store.Record{"id": id})
	case "empty":
		// An object with nothing in it. Asana answers {"data": {}}, which is
		// not the same as no body: a client calling .json() succeeds here and
		// throws on a 204, and that difference is the whole point of saying
		// which one a provider does.
		return s.writeRecord(w, matched, store.Record{})
	default:
		status := matched.spec.Status
		if status == 0 {
			status = http.StatusNoContent
		}

		w.WriteHeader(status)

		return status
	}
}

func (s *Sandbox) list(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	limit := matched.spec.Pagination.Limit
	if limit <= 0 {
		limit = 10
	}

	// A declared name replaces the defaults rather than joining them. Google
	// does not accept limit, and an emulator that honours both spellings lets
	// the wrong one work locally and fail nowhere until production.
	// Where the parameters travel is the provider's choice, and for a listing
	// reached by POST it is usually the body.
	from := pagingFrom(r, matched.spec.Pagination)

	switch name := matched.spec.Pagination.LimitParam; name {
	case "":
		limit = from.int("limit", limit)
	case "-":
		// The provider accepts no name for the page size, so the declared one
		// is the only one there is.
	default:
		limit = from.int(name, limit)
	}

	// The declared style decides what the position parameter means. It was
	// never read, so 149 shipped routes declaring offset or page numbering were
	// paged as though they were cursor based, which meant not paged at all: the
	// parameter went unrecognised and every request answered with the first
	// page. A client looping until it had collected total_count records never
	// finished, and one asking for page two processed page one twice.
	var (
		page store.Page
		err  error
	)

	where := s.scopeVars(matched, vars)

	// A listing that narrows itself whether or not the caller asked. The
	// filter joins the scope because it does the same job: both decide which
	// records this request is allowed to see.
	for _, f := range matched.spec.Filters {
		value := r.URL.Query().Get(f.Param)
		if value == "" {
			value = f.Default
		}

		// The escape value, for the providers that have one. Without it there
		// is no way to ask a listing for everything, which is worth knowing
		// about a provider and is not something to invent on its behalf.
		if value == "" || (f.All != "" && value == f.All) {
			continue
		}

		if where == nil {
			where = map[string]any{}
		}

		// A parameter value the provider expands into several field values.
		if options, grouped := f.Values[value]; grouped {
			where[f.Field] = options

			continue
		}

		where[f.Field] = value
	}

	switch matched.spec.Pagination.Style {
	case "offset", "page":
		offset := positionOf(from, matched.spec.Pagination, limit)

		page, err = s.store.ListFrom(matched.spec.Resource, where, offset, limit)

		// A Recipe declaring cursor_field alongside offset or page numbering
		// means the pointer to the next page, and for these styles that
		// pointer is a position rather than an identifier. It used to be the
		// last record's id, which was wrong for every provider that pages this
		// way, and wrong in the direction that reads as working.
		if err == nil && page.HasMore {
			page.NextCursor = nextPosition(matched.spec.Pagination, offset, limit)
		}
	default:
		cursor := cursorOf(from, matched.spec.Pagination)

		page, err = s.store.ListWhere(matched.spec.Resource, where, cursor, limit)

		if err != nil && errors.Is(err, store.ErrNotFound) {
			return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", "No such cursor: "+cursor+".")
		}
	}

	if err != nil {
		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", err.Error())
	}

	for i, record := range page.Records {
		page.Records[i] = trim(matched.spec, record)
	}

	// The envelope this route answers with. A provider's listings do not
	// always share one: Clerk's users and sessions are bare arrays and its
	// organisations are wrapped, and a Recipe-wide shape makes one of them
	// wrong.
	list := s.recipe.ListFor(matched.spec)

	body := s.listBody(list, page, matched.spec.Resource, r.URL.Path,
		nextPageURL(r, matched.spec.Pagination, page.NextCursor))

	// What this request actually served, for the providers that report it.
	// These cannot be constants: a field whose whole purpose is to say where
	// you are is worse than absent when it always says the same thing.
	if object, ok := body.(map[string]any); ok {
		if name := list.LimitField; name != "" {
			setPath(object, name, limit)
		}

		if name := list.PageField; name != "" {
			setPath(object, name, servedPage(from, matched.spec.Pagination))
		}

		// How many pages the whole set makes, which is a different quantity
		// from how many records there are. A count of records under a name
		// that says pages is worse than an invented field: the name is real
		// and the number is plausible, so a client looping while
		// page <= totalPages asks for pages that do not exist and reads them
		// as empty results rather than as a mistake.
		if name := list.PagesField; name != "" {
			setPath(object, name, s.countValue(list, pageCount(page.Total, limit)))
		}
	}

	// Other collections travelling in the same body. One endpoint, several
	// arrays: GoCardless answers a request for transactions with booked and
	// pending together, and the same purchase is in one and then the other.
	if len(matched.spec.Beside) > 0 {
		object, ok := body.(map[string]any)
		if !ok {
			return s.writeRecipeError(w, "invalid_request", 500, "invalid_request",
				"a route naming beside resources must answer with an object")
		}

		for _, name := range matched.spec.Beside {
			// Every record, not a page. The cursor belongs to the route's own
			// resource, and a second collection paged by the same cursor
			// would be paged by a pointer into a different list.
			alongside, err := s.store.ListWhere(name, where, "", 0)
			if err != nil {
				return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", err.Error())
			}

			// The route's returns names the route's own resource's fields
			// and means nothing here. What every collection in the body does
			// share is the scope: the partition is in the path, and a
			// provider that puts it there does not repeat it in the body.
			records := s.presentAll(name, alongside.Records)
			for _, field := range matched.spec.Scope {
				for _, record := range records {
					delete(record, field)
				}
			}

			setPath(object, s.collectionName(name, ""), records)
		}
	}

	// Constants this route adds, after the Recipe-wide ones, because on a
	// one-path API the envelope depends on which query was asked.
	if len(matched.spec.Fields) > 0 {
		if object, ok := body.(map[string]any); ok {
			body = withFields(object, matched.spec.Fields)
		}
	}

	// A listing's declared status and headers used to be read by nothing.
	// status and empty_body were once read on creates alone, then extended to
	// the routes that answer with a record, and listings were left out of
	// both -- so a Recipe could say its provider answers a page with 206 and
	// a Next-Range header, and be quietly ignored.
	//
	// Heroku is the provider that makes this matter. It pages with the Range
	// header, answers 206 Partial Content while there is more, and puts the
	// resume point in Next-Range. A 200 there means you have everything, so
	// an emulator that always answered 200 taught a client the one thing it
	// must not believe.
	// The next page as an RFC 5988 Link header, for the providers that
	// advertise it there rather than in the body. Set before the route's own
	// headers so a Recipe naming Link explicitly still wins.
	if list.LinkHeader && page.HasMore && page.NextCursor != "" {
		if next := nextPageURL(r, matched.spec.Pagination, page.NextCursor); next != "" {
			w.Header().Set("Link", "<"+next+`>; rel="next"`)
		}
	}

	s.writeRouteHeaders(w, matched, nil)

	status := matched.spec.Status
	if status == 0 {
		status = http.StatusOK
	}

	writeJSON(w, status, body)

	return status
}

// nextPageURL renders the request that fetches the page after this one.
//
// It is this request with the position parameter moved on, which is what every
// provider's Link header contains: the same query, one page further in. A
// listing whose paging travels in the body has no such URL, so it gets no
// header rather than a misleading one.
func nextPageURL(r *http.Request, spec recipe.Pagination, cursor string) string {
	if spec.In == "body" {
		return ""
	}

	name := spec.CursorParam
	if name == "" {
		switch spec.Style {
		case "page", "offset":
			name = spec.Style
		default:
			name = "cursor"
		}
	}

	if name == "-" {
		return ""
	}

	next := *r.URL

	// The path the caller asked for, which is not always the one the sandbox
	// sees. The multi-provider server mounts each Recipe under its own name
	// and rewrites URL.Path, leaving RequestURI as it arrived, so a header
	// built from the rewritten path sends the client to /repos/... when the
	// server serves /github/repos/... -- a next page that 404s, which is a
	// worse failure than no next page at all.
	if r.RequestURI != "" {
		if asked, err := url.ParseRequestURI(r.RequestURI); err == nil && asked.Path != "" {
			next.Path = asked.Path
		}
	}

	query := next.Query()
	query.Set(name, cursor)
	next.RawQuery = query.Encode()

	// A provider's Link header carries an absolute URL, and a client is
	// expected to request it as it stands rather than to reassemble one.
	next.Scheme = "http"
	next.Host = r.Host

	return next.String()
}

// cursorOf reads the identifier a cursor-paged listing resumes after.
func cursorOf(from paging, spec recipe.Pagination) string {
	if name := spec.CursorParam; name != "" {
		return from.get(name)
	}

	if cursor := from.get("starting_after"); cursor != "" {
		return cursor
	}

	return from.get("cursor")
}

// positionOf turns an offset or page parameter into a record offset.
//
// The two styles differ only in what the number counts. An offset counts
// records and starts at nought; a page counts pages and starts at one, so page
// two begins one page-length in. Getting that boundary wrong by one page is
// the classic way to lose or duplicate a record at every page break, which is
// exactly the bug an emulator is supposed to catch rather than commit.
func positionOf(from paging, spec recipe.Pagination, limit int) int {
	name := spec.CursorParam
	if name == "" {
		name = spec.Style
	}

	raw := from.get(name)
	if raw == "" {
		return 0
	}

	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return 0
	}

	if spec.Style != "page" {
		return n
	}

	first := spec.FirstPageNumber()

	// Anything at or before the provider's own first page is the first page.
	// Providers disagree about whether a number below it is an error, and
	// answering with the first page is the reading that cannot silently skip
	// records.
	if n <= first {
		return 0
	}

	return (n - first) * limit
}

// nextPosition renders the position of the page after this one.
//
// An offset counts records, so the next one starts a page-length further in. A
// page counts pages from one, so the offset has to be turned back into a page
// number before adding to it.
func nextPosition(spec recipe.Pagination, offset, limit int) string {
	if spec.Style == "page" {
		if limit <= 0 {
			limit = 1
		}

		// The page after this one, counted from wherever the provider starts.
		return strconv.Itoa(offset/limit + spec.FirstPageNumber() + 1)
	}

	return strconv.Itoa(offset + limit)
}

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

// writeRouteHeaders sets the response headers a route declares, substituting
// the record's identifier for {id}.
func (s *Sandbox) writeRouteHeaders(w http.ResponseWriter, matched route, record store.Record) {
	for name, value := range matched.spec.Headers {
		// A listing has no single identifier, so it passes nil and whatever
		// {id} is in the value stays as written rather than being
		// substituted with the first record's, which would be a promise
		// about ordering nobody made.
		if id, ok := record["id"].(string); ok {
			value = strings.ReplaceAll(value, "{id}", id)
		}

		w.Header().Set(name, value)
	}
}

// present renames the identifier to the property the provider actually uses.
// The store keeps every record keyed by "id" so fixtures and internal lookups
// stay uniform; only the wire shape changes, which is where it matters.
// idAs renames the identifier on one route alone, for the providers that call
// it one thing there and another everywhere else. Empty means the resource's
// own name, which is the ordinary case.
func (s *Sandbox) present(resource string, record store.Record, idAs string) store.Record {
	spec, ok := s.recipe.Resources[resource]
	if !ok {
		return record
	}

	if idAs != "" {
		spec.ID.Field = idAs
	}

	// "-" is not a rename but a suppression: the identifier is how Cauldron
	// finds the record, not something the provider puts on the wire.
	hidden := spec.ID.Field == "-"
	renamed := spec.ID.Field != "" && spec.ID.Field != "id" && !hidden

	// A provider whose identifier is a JSON number. The store keeps it as a
	// string, because that is the only form every style shares and the only
	// form a path parameter arrives in, so the conversion belongs here at the
	// edge rather than anywhere a lookup happens.
	retyped := spec.ID.Type == "number"

	if !renamed && !hidden && !retyped && !s.nests(spec) {
		return record
	}

	out := make(store.Record, len(record))

	for key, value := range record {
		if key == "id" && hidden {
			continue
		}

		if key == "id" && retyped {
			setPath(out, identifierName(spec), numeric(value))

			continue
		}

		if key == "id" && renamed {
			// setPath, so a dotted name nests. Contentful keeps the
			// identifier at sys.id rather than at the top level.
			setPath(out, spec.ID.Field, value)
			continue
		}

		// A field declared with "in" moves under that sub-object on the wire.
		// HubSpot reads contact.properties.email, and a client written against
		// it finds nothing at the top level.
		//
		// "-" moves it nowhere: the record holds it and the wire never sees
		// it. That is what a route's scope needs, because a partition living
		// in the path is not repeated in the body by most providers, and the
		// only way to say so before this was to name every other field in the
		// route's returns.
		if field, declared := spec.Fields[key]; declared && field.In == "-" {
			continue
		} else if declared && field.In != "" {
			// A dotted name nests twice. Brex puts a card's limit at
			// spend_controls.limit.amount, and treating the name as one
			// literal key produced a flat "spend_controls.limit" key that no
			// provider sends and nothing was checking.
			nested := nestedObject(out, field.In)

			// The wire name, which is the field's own unless it says
			// otherwise. Without that, a resource needing both title.rendered
			// and content.rendered had to name one of them something else, and
			// the name it was given leaked out as title.title_rendered.
			nested[field.WireName(key)] = value

			continue
		}

		// A declared constant at a dotted name nests, the same way a renamed
		// identifier and a field's "in" already do. Attio keeps two of a
		// record's three identifier UUIDs beside the third -- the whole id is
		// an object -- and writing id.workspace_id as one literal key
		// produced a shape no provider sends. That is the fifth mechanism in
		// this file to have made that mistake, and the reason it keeps
		// happening is that a dotted name is only a path in the places
		// somebody remembered to make it one.
		if _, declared := spec.Fields[key]; !declared && recipe.IsPath(key) {
			setPath(out, key, value)

			continue
		}

		// A top-level field can name itself something else on the wire, for
		// the one thing it cannot otherwise say: that the name it wants is
		// already the record's own. GitHub gives an issue a number and an id,
		// a path addresses it by the number, so the record keys on the number
		// and "id" is left for the field carrying GitHub's id.
		if field, declared := spec.Fields[key]; declared && field.As != "" {
			out[field.As] = value

			continue
		}

		out[key] = value
	}

	return out
}

// lookupNested walks a dotted path through an incoming request body.
func lookupNested(record store.Record, path string) (map[string]any, bool) {
	var current any = map[string]any(record)

	for _, segment := range strings.Split(path, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		current, ok = object[segment]
		if !ok {
			return nil, false
		}
	}

	object, ok := current.(map[string]any)

	return object, ok
}

// nestedObject walks a dotted path, creating the objects it passes through,
// and returns the one the field belongs in.
//
// Most providers nest one level, which is a single key. Brex nests twice for a
// card's spending limit, and a dotted name is how a Recipe says so.
func nestedObject(out map[string]any, path string) map[string]any {
	current := out

	for _, segment := range strings.Split(path, ".") {
		name, indexes := splitIndex(segment)

		if len(indexes) == 0 {
			next, _ := current[segment].(map[string]any)
			if next == nil {
				next = map[string]any{}
				current[segment] = next
			}

			current = next

			continue
		}

		// An array with objects in it. RingCentral's to is a list of one
		// recipient, and without this the whole segment became a literal key
		// spelled "to[0]" — a shape no provider sends, produced silently, and
		// invisible to a conformance suite that has no case naming it. That is
		// the third time a key like this has shipped, so the validator now
		// refuses an index the runtime cannot honour rather than trusting the
		// next Recipe author to notice.
		current = descendIndexed(current, name, indexes)
	}

	return current
}

// descendIndexed walks into an array-valued key, growing the array as needed,
// and returns the object at the innermost index.
func descendIndexed(current map[string]any, name string, indexes []int) map[string]any {
	held, _ := current[name].([]any)
	if held == nil {
		held = []any{}
	}

	for depth, index := range indexes {
		for len(held) <= index {
			held = append(held, map[string]any{})
		}

		if depth == len(indexes)-1 {
			object, _ := held[index].(map[string]any)
			if object == nil {
				object = map[string]any{}
				held[index] = object
			}

			current[name] = held

			return object
		}

		inner, _ := held[index].([]any)
		if inner == nil {
			inner = []any{}
		}

		current[name] = held
		held = inner
	}

	return map[string]any{}
}

// splitIndex separates a path segment from any [n] suffixes on it.
func splitIndex(segment string) (string, []int) {
	open := strings.Index(segment, "[")
	if open < 0 {
		return segment, nil
	}

	name := segment[:open]
	rest := segment[open:]

	var indexes []int

	for rest != "" {
		if !strings.HasPrefix(rest, "[") {
			return segment, nil
		}

		close := strings.Index(rest, "]")
		if close < 0 {
			return segment, nil
		}

		n, err := strconv.Atoi(rest[1:close])
		if err != nil || n < 0 {
			return segment, nil
		}

		indexes = append(indexes, n)
		rest = rest[close+1:]
	}

	return name, indexes
}

// nests reports whether any of a resource's fields live under a sub-object.
//
// A field marked "-" is not sent, and flatten below looks for it under a key
// literally spelled that, finds nothing, and moves on. Excluding it here
// changed no behaviour, so it is not excluded: a guard that guards nothing is
// the thing this project keeps telling Recipes not to write.
func (s *Sandbox) nests(spec recipe.Resource) bool {
	for _, field := range spec.Fields {
		// A field that only renames itself still needs present to run, or the
		// rename is declared and never applied -- which is the quiet kind of
		// wrong this whole file keeps finding.
		if field.In != "" || field.As != "" {
			return true
		}
	}

	// A dotted constant nests too, and a resource can have one without any
	// field declaring an "in" at all -- which is exactly Attio's records,
	// where the nesting is entirely in the identifier and the constants
	// beside it.
	for name := range spec.Constants {
		if recipe.IsPath(name) {
			return true
		}
	}

	return false
}

// flatten is the inverse of present for an incoming request: a client sends
// {"properties": {"email": "..."}} and the store keeps fields flat.
func (s *Sandbox) flatten(resource string, record store.Record) store.Record {
	spec, ok := s.recipe.Resources[resource]
	if !ok || !s.nests(spec) {
		return record
	}

	for name, field := range spec.Fields {
		if field.In == "" {
			continue
		}

		nested, isObject := lookupNested(record, field.In)
		if !isObject {
			continue
		}

		if value, present := nested[field.WireName(name)]; present {
			record[name] = value
		}
	}

	// The sub-objects have been unpacked, so drop them rather than storing the
	// same values twice under two shapes.
	for _, field := range spec.Fields {
		if field.In != "" {
			delete(record, strings.Split(field.In, ".")[0])
		}
	}

	return record
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

// presentAll renames identifiers across a page.
func (s *Sandbox) presentAll(resource string, records []store.Record) []store.Record {
	spec, ok := s.recipe.Resources[resource]
	if !ok || ((spec.ID.Field == "" || spec.ID.Field == "id") && spec.ID.Type != "number" && !s.nests(spec)) {
		return records
	}

	out := make([]store.Record, 0, len(records))

	for _, record := range records {
		out = append(out, s.present(resource, record, ""))
	}

	return out
}

// countValue renders a total as the provider sends it. Docusign sends its
// counts as strings, and emitting a number would quietly fix a bug the caller
// has to handle for themselves.
func (s *Sandbox) countValue(spec recipe.ListResponse, total int) any {
	if spec.CountAsString {
		return strconv.Itoa(total)
	}

	return total
}

// resourceBody shapes a single object according to the Recipe's declared style.
// Shopify nests it under the singular resource name; most providers return the
// object itself, and a client written for one shape breaks on the other.
// The envelope is the route's, which is the Recipe's unless the route says
// otherwise: Datadog wraps a created event and not a created monitor, and
// Vercel wraps a domain and neither a project nor a deployment.
func (s *Sandbox) resourceBody(spec recipe.ResourceResponse, idAs, resource string, record store.Record) any {
	record = s.present(resource, record, idAs)

	success := s.recipe.Responses.Success.Fields

	if spec.Style != "wrapped" {
		if len(success) == 0 {
			return record
		}

		// Nowhere to put the envelope fields without colliding with the
		// object's own, so wrap rather than silently overwrite a real field.
		out := make(map[string]any, len(record)+len(success))

		for key, value := range record {
			out[key] = value
		}

		return withFields(out, success)
	}

	// The default is the resource's own singular name, which is what Shopify,
	// Slack, Square and Zendesk all use. When the object is wrapped in a list
	// the default is the plural collection name instead, because a list of one
	// is still a collection: Xero answers with Invoices, not Invoice.
	key := spec.Key
	if key == "" {
		key = resource

		if spec.Array {
			key = s.collectionName(resource, "")
		}
	}

	var payload any = record

	// Xero answers a request for one invoice with a list of one, so client
	// code reads Invoices[0]. An emulator returning the object directly lets
	// code ship that breaks against the real API on the first call.
	if spec.Array {
		payload = []store.Record{record}
	}

	return withFields(map[string]any{key: payload}, success)
}

// withFields stamps a provider's constant envelope fields onto a body. A dotted
// name nests, so response_metadata.next_cursor needs no second mechanism.
func withFields(body map[string]any, fields map[string]any) map[string]any {
	for name, value := range fields {
		setPath(body, name, value)
	}

	return body
}

func setPath(body map[string]any, path string, value any) {
	head, rest, nested := strings.Cut(path, ".")

	if !nested {
		if name, indexes := splitIndex(head); len(indexes) > 0 {
			if object, ok := value.(map[string]any); ok {
				target := descendIndexed(body, name, indexes)
				for key, nestedValue := range object {
					setPath(target, key, nestedValue)
				}

				return
			}
		}

		// A declared constant must not destroy data already in the body.
		// Intercom's pages object carries both a declared type and a computed
		// next cursor, and whichever landed second used to erase the other.
		if object, ok := value.(map[string]any); ok {
			if existing, ok := body[head].(map[string]any); ok {
				for key, nestedValue := range object {
					setPath(existing, key, nestedValue)
				}

				return
			}
		}

		body[head] = value

		return
	}

	// An indexed segment, the same way nestedObject already handles one for a
	// field's "in". Without it a constant declared at data.boards[0].name
	// produced a key literally spelled "boards[0]" -- a shape no provider
	// sends, produced in silence, and the fourth time a key like that has
	// been written here.
	//
	// The asymmetry was the tell: a conformance case could already assert
	// data.boards[0].name and a Recipe could not emit it.
	if name, indexes := splitIndex(head); len(indexes) > 0 {
		setPath(descendIndexed(body, name, indexes), rest, value)

		return
	}

	child, ok := body[head].(map[string]any)
	if !ok {
		child = map[string]any{}
		body[head] = child
	}

	setPath(child, rest, value)
}

// listBody shapes a page according to the declared list style, which is the
// Recipe's unless the route overrode it.
func (s *Sandbox) listBody(spec recipe.ListResponse, page store.Page, resource, path, nextURL string) any {
	page.Records = s.presentAll(resource, page.Records)

	// Chargebee wraps every item under the resource's own name, so a client
	// reads list[0].subscription.id. Anyone indexing straight into the item
	// finds nothing at all.
	var items any = page.Records

	// A collection of identifiers rather than of records. DynamoDB's
	// ListTables sends TableNames as an array of strings and keeps the table
	// object for DescribeTable beside it; SQS's ListQueues does the same with
	// QueueUrls. Emitting objects there hands a client something it will
	// interpolate into a URL, and [object Object] is what it builds.
	if spec.EntryField != "" {
		values := make([]any, 0, len(page.Records))

		for _, record := range page.Records {
			if value, ok := nestedValue(record, spec.EntryField); ok {
				values = append(values, value)
			}
		}

		items = values
	}

	if spec.EntryStyle == "wrapped" {
		wrapped := make([]any, 0, len(page.Records))

		for _, record := range page.Records {
			wrapped = append(wrapped, map[string]any{resource: record})
		}

		items = wrapped
	}

	// A collection of one, sent as the object. Tradier's words: if you have a
	// single order, it will be returned as a JSON obj/dict whereas multiple
	// orders will be returned as an array.
	//
	// This runs after the entry wrapping above, so a Recipe that does both
	// collapses to the wrapped item rather than to the bare record, which is
	// what an XML-shaped API sending one child element would produce.
	if spec.CollapseSingle {
		if list, ok := items.([]any); ok && len(list) == 1 {
			items = list[0]
		}

		if list, ok := items.([]store.Record); ok && len(list) == 1 {
			items = list[0]
		}
	}

	switch spec.Style {
	case "map":
		// Keyed by identifier rather than ordered. Pusher answers with an
		// object of channel names, so looping over it as a list finds
		// nothing, and a channel nobody is on is absent from the object
		// entirely rather than present with a zero. The key is the
		// identifier, so it does not repeat inside the value.
		keyed := map[string]any{}

		for _, record := range page.Records {
			name, _ := record["id"].(string)
			if name == "" {
				continue
			}

			value := map[string]any{}

			for field, held := range record {
				if field == "id" {
					continue
				}

				value[field] = held
			}

			keyed[name] = value
		}

		body := map[string]any{}
		setPath(body, s.collectionName(resource, spec.Key), keyed)

		return withFields(body, spec.Fields)
	case "bare":
		// GitHub and friends return the array itself, with paging in headers.
		// A caller doing json.Unmarshal into a slice must not receive an object.
		return items
	case "wrapped":
		// A dotted key nests, so a collection can sit two levels down.
		// Segment answers with data.sources rather than a top-level array.
		body := map[string]any{}

		// SQS leaves the key out entirely when there is nothing to send, and
		// that is the difference between a consumer that waits on an idle
		// queue and one that throws. Sending an empty array is the helpful
		// kind of wrong: every test passes and the first quiet minute in
		// production does not.
		if !spec.OmitWhenEmpty || len(page.Records) > 0 {
			setPath(body, s.collectionName(resource, spec.Key), items)
		}

		if spec.CursorField != "" && page.NextCursor != "" {
			setPath(body, spec.CursorField, cursorValue(spec, page.NextCursor, nextURL))
		}

		// Salesforce's done is has_more with the sense reversed, and false is
		// its interesting value: a query that matched more rows than it
		// returned says done: false and expects the caller to follow
		// nextRecordsUrl. Sending true would tell a client its partial result
		// set was the whole thing.
		if spec.CompleteField != "" {
			setPath(body, spec.CompleteField, !page.HasMore)
		}

		// Only on the last page, which is the entire point. Google Calendar
		// gives you a page token or a sync token and never both, so code that
		// grabs whichever one it finds first stores the wrong one.
		if spec.FinalField != "" && !page.HasMore {
			setPath(body, spec.FinalField, finalToken(resource, page))
		}

		if spec.HasMoreField != "" {
			setPath(body, spec.HasMoreField, page.HasMore)
		}

		if spec.CountField != "" {
			setPath(body, spec.CountField, s.countValue(spec, page.Total))
		}

		body = withFields(body, s.recipe.Responses.Success.Fields)

		return withFields(body, spec.Fields)
	default:
		body := map[string]any{
			"object":   "list",
			"data":     items,
			"has_more": page.HasMore,
		}

		if spec.URL {
			body["url"] = path
		}

		// A cursor field is opt in because sending one the provider does not
		// send is the more dangerous mistake: code written against it works
		// locally and fails in production, which is the exact failure this
		// project exists to prevent.
		if spec.CursorField != "" && page.NextCursor != "" {
			setPath(body, spec.CursorField, cursorValue(spec, page.NextCursor, nextURL))
		}

		body = withFields(body, s.recipe.Responses.Success.Fields)

		return withFields(body, spec.Fields)
	}
}

// finalToken derives the opaque token a provider hands back when a listing is
// exhausted.
//
// Stable for a given final page, so a suite can assert it rather than only
// matching a shape, and opaque so nothing is tempted to read structure out of
// it. The real ones carry no parseable structure either, and clients that
// guessed otherwise are the reason they are documented as opaque.
func finalToken(resource string, page store.Page) string {
	last := ""

	if n := len(page.Records); n > 0 {
		last, _ = page.Records[n-1]["id"].(string)
	}

	sum := sha256.Sum256([]byte(resource + ":" + last))

	return base64.RawURLEncoding.EncodeToString(sum[:15])
}

// scopeVars extracts the scope filters for a request from its path parameters.
func (s *Sandbox) scopeVars(matched route, vars map[string]string) map[string]any {
	if len(matched.spec.Scope) == 0 {
		return nil
	}

	fields := s.recipe.Resources[matched.spec.Resource].Fields
	out := make(map[string]any, len(matched.spec.Scope))

	for _, name := range matched.spec.Scope {
		out[name] = scoped(fields[name].Type, vars[name])
	}

	return out
}

// scoped converts a path parameter to the type its field declares.
//
// A path is text and a field is not. Help Scout's conversationId, Gorgias's
// ticket_id, Shortcut's epic_id, Documenso's documentId and Hetzner's
// resources[0].id are all integers on the wire, and a record created through
// a scoped route carried the string out of the URL instead. So the same field
// was a number when a fixture seeded it and a string when a create produced
// it, in the same collection, in the same response shape, and only one of the
// two matched what the provider sends.
//
// Matching is done with fmt.Sprint on both sides, so a scope that changes
// type still finds the records it scoped.
//
// The JSON grammar decides, for the same reason it decides in a form body: a
// leading plus or a redundant leading zero is not something a JSON client
// could have sent, and an identifier that begins with one is an identifier
// rather than a number however the field is declared.
func scoped(declared, value string) any {
	switch declared {
	case "integer", "number", "timestamp", "timestamp_ms":
		if !jsonNumeric(value) {
			return value
		}

		if n, err := strconv.ParseInt(value, 10, 64); err == nil {
			return n
		}

		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}

	return value
}

// collectionName resolves the key a wrapped list is nested under: the
// resource's declared collection, then a recipe-wide override, then the
// resource name itself.
func (s *Sandbox) collectionName(resource, override string) string {
	if spec, ok := s.recipe.Resources[resource]; ok && spec.Collection != "" {
		return spec.Collection
	}

	if override != "" {
		return override
	}

	return resource
}

func (s *Sandbox) notFound(w http.ResponseWriter, err error, resource, id string) int {
	if errors.Is(err, store.ErrUnknownResource) {
		return s.writeRecipeError(w, "unknown_route", 404, "unknown_route", "Unrecognised request URL.", resource)
	}

	return s.writeRecipeError(
		w, "resource_missing", 404, "resource_missing",
		"No such "+resource+": "+id+".", resource+": "+id,
	)
}

// emitFor sends the conventional lifecycle event for a resource, when the
// Recipe declares it. Providers emit customer.created after a customer is
// created; a fake that stays silent teaches an application the wrong lesson.
func (s *Sandbox) emitFor(resource, action string, record store.Record) {
	event := resource + "." + action

	if !s.webhooks.Known(event) {
		return
	}

	// Shaped, exactly as an HTTP response is. Without this the payload carried
	// the store's own field names, so the same record from the same sandbox at
	// the same instant had two shapes depending on how you looked at it: the
	// response nested amount under amount_money as the Recipe declares and the
	// webhook did not. A handler written against the emulator read
	// event.data.object.amount and was entirely green here and entirely wrong
	// against the provider.
	//
	// present rather than resourceBody, because the response envelope belongs
	// to the response. A webhook has its own envelope and the record sits
	// inside it unwrapped.
	_, _ = s.webhooks.Emit(event, s.present(resource, record, ""))
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
	case "bearer":
		presented = r.Header.Get("Authorization")
	case "header":
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

// writeRecipeError writes an error, preferring the Recipe's own definition so
// that status codes, codes and headers match what the real provider sends.
// Fallbacks apply only when the Recipe does not define the situation.
//
// A fourth argument is the request-specific detail, such as the parameter that
// was missing. A Recipe reaches it through {detail} in its message.
// builtinFallbacks maps the failures the runtime produces itself onto the
// declared error whose shape they should borrow.
//
// An unrecognised path, a method a route does not take and a body that will
// not parse are all produced by the handler rather than declared by a Recipe,
// and they used to carry their own internal name as both the code and the
// category. No provider has an error type called "unknown_route". Every Stripe
// 404 said type "unknown_route" where the real one says
// "invalid_request_error", and no Recipe declared any of these three, so this
// was live on all 97 of them, on the most commonly exercised error path there
// is. Retry and branching logic keyed on the category took one path locally
// and another in production.
//
// A Recipe may still declare any of these names itself, and that wins.
var builtinFallbacks = map[string]string{
	"unknown_route":      "resource_missing",
	"method_not_allowed": "resource_missing",
	"conflict":           "invalid_request",
	"invalid_request":    "parameter_missing",
}

// resolveBuiltin follows the fallback chain until it reaches a name the Recipe
// declares, and returns the original when it reaches none.
func resolveBuiltin(r *recipe.Recipe, name string) string {
	seen := map[string]bool{}

	for current := name; !seen[current]; {
		if _, declared := r.Errors[current]; declared {
			return current
		}

		seen[current] = true

		next, ok := builtinFallbacks[current]
		if !ok {
			break
		}

		current = next
	}

	return name
}

func (s *Sandbox) writeRecipeError(w http.ResponseWriter, name string, fallback ...any) int {
	status := http.StatusInternalServerError
	code := name
	message := "An error occurred."

	if len(fallback) >= 3 {
		status, _ = fallback[0].(int)
		code, _ = fallback[1].(string)
		message, _ = fallback[2].(string)
	}

	var detail string

	if len(fallback) >= 4 {
		detail, _ = fallback[3].(string)
	}

	// Providers categorise errors more coarsely than they code them, and client
	// libraries switch on the category. Using the name as the type made every
	// error its own category, which no real provider does.
	category := name

	// Extra body properties this particular failure carries.
	var extra map[string]any

	resolved := resolveBuiltin(s.recipe, name)

	// A Recipe that declares this failure by name owns all of it. One that
	// only lends its shape to a built-in lends the category, the code and the
	// extra fields, and not the status: borrowing a 404's status for a 405
	// would change what the failure is, which is a worse lie than the invented
	// category this replaced.
	borrowed := resolved != name

	if defined, ok := s.recipe.Errors[resolved]; ok {
		if !borrowed {
			status = defined.Status
		}

		if defined.Code != "" {
			code = defined.Code
		}

		if defined.Type != "" {
			category = defined.Type
		}

		// The Recipe's wording wins. Providers differ here in ways that matter:
		// Stripe names the offending parameter, GitHub says "Not Found" and
		// nothing else. A Recipe that wants the request-specific detail asks
		// for it with {detail}, so the choice is the Recipe author's rather
		// than a hardcoded guess in the handler.
		// The wording travels only when the lender describes the same status.
		// GitHub's 404 really does say "Not Found" and nothing else, and that
		// is worth having; its wording on a 405 would be a sentence about the
		// wrong thing.
		if defined.Message != "" && (!borrowed || defined.Status == status) {
			message = strings.ReplaceAll(defined.Message, "{detail}", detail)
		}

		for key, value := range defined.Headers {
			w.Header().Set(key, value)
		}

		extra = defined.Fields
	}

	// Trello answers with plain text, so a client calling .json() on the
	// response throws rather than reporting the failure. Writing JSON here
	// would hide that.
	if s.recipe.Responses.Error.Style == "text" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(status)

		fmt.Fprint(w, message)

		return status
	}

	writeJSON(w, status, s.errorBody(category, code, message, status, extra))

	return status
}

// errorBody shapes a failure according to the Recipe's declared error style.
//
// The return type is any rather than a map because Salesforce's failures are a
// bare top-level array with no envelope at all. A client reading .message off
// the response finds undefined, and has to index before it can read anything.
func (s *Sandbox) errorBody(category, code, message string, status int, extra map[string]any) any {
	spec := s.recipe.Responses.Error

	if spec.Style == "string_list" {
		key := spec.Key
		if key == "" {
			key = "errors"
		}

		// Datadog sends the array with bare strings in it. A client looping
		// over the entries and reading .message from each finds undefined on
		// every one, which throws rather than reporting anything.
		return withFields(withFields(map[string]any{key: []any{message}}, spec.Fields), extra)
	}

	if spec.Style != "flat" && spec.Style != "list" {
		// The sub-keys are named the same way the flat style names its
		// top-level ones. Google nests a numeric code beside a string status
		// and calls neither of them "type", so the names have to be
		// declarable rather than assumed. An empty name keeps the old default,
		// which is what every nested Recipe written before this relied on.
		nested := map[string]any{}

		if field := spec.MessageField; field != "-" {
			if field == "" {
				field = "message"
			}

			nested[listName(field)] = messageValue(field, message)
		}

		// Providers disagree about which of these they send. Stripe sends a
		// type and a code, Airtable only a type, Vercel only a code. Declaring
		// the omission keeps the emulator from inventing a field a client
		// could come to depend on and then lose.
		if field := spec.CodeField; field != "-" {
			if field == "" {
				field = "code"
			}

			// The same conversion the flat style applies. PagerDuty's codes
			// really are numbers, and sending "2006" as text meant a client
			// switching on the code never matched.
			nested[field] = codeValue(spec.CodeType, code)
		}

		if field := spec.TypeField; field != "-" {
			if field == "" {
				field = "type"
			}

			nested[field] = category
		}

		if spec.StatusField != "" {
			nested[spec.StatusField] = status
		}

		key := spec.Key
		if key == "" {
			key = "error"
		}

		// Declared fields apply here too. They were dropped for the nested
		// style alone, so a Recipe could add a constant to every error and
		// have it silently ignored.
		return withFields(withFields(map[string]any{key: nested}, spec.Fields), extra)
	}

	body := map[string]any{}

	// "-" means the provider sends no prose at all. Slack's errors are a code
	// and nothing else, and inventing a sentence it never sends is infidelity
	// in the direction that is hardest to notice.
	if field := spec.MessageField; field != "-" {
		if field == "" {
			field = "message"
		}

		// setPath, so a dotted name nests the same way code_field and
		// status_field already do. Front puts all three under _error.
		setPath(body, listName(field), messageValue(field, message))
	}

	// "-" omits the field, exactly as it does for the message above and in the
	// nested style. It used to be honoured only in the nested style here, so a
	// flat Recipe declaring code_field: "-" or type_field: "-" got a literal
	// "-" key in every error body instead of no field at all. Three shipped
	// Recipes were doing that, and every case written about them passed,
	// because a case asserts the fields it names and says nothing about a key
	// nobody thought to look for.
	if spec.CodeField != "" && spec.CodeField != "-" && code != "" {
		setPath(body, spec.CodeField, codeValue(spec.CodeType, code))
	}

	if spec.TypeField != "" && spec.TypeField != "-" {
		setPath(body, spec.TypeField, category)
	}

	if spec.StatusField != "" {
		setPath(body, spec.StatusField, status)
	}

	if spec.Style == "list" {
		key := spec.Key
		if key == "" {
			key = "errors"
		}

		// "-" means there is no envelope: the array is the whole body, which
		// is what Salesforce sends. Reading .message off that response finds
		// undefined, and any declared fields have nowhere to go, so a Recipe
		// that wants both has to choose the shape the provider actually uses.
		if key == "-" {
			return []any{withFields(body, extra)}
		}

		// One request can fail several ways at once, which is why SendGrid and
		// Clerk both send an array. The emulator reports one, in the shape a
		// client already has to loop over.
		//
		// Declared fields sit beside the array rather than inside each entry:
		// Clerk's trace id belongs to the response, not to one failure. A
		// dotted key nests, which QuickBooks needs for Fault.Error.
		envelope := map[string]any{}
		setPath(envelope, key, []any{body})

		return withFields(withFields(envelope, spec.Fields), extra)
	}

	return withFields(withFields(body, spec.Fields), extra)
}

// listName strips the array marker from a field name.
func listName(field string) string {
	return strings.TrimSuffix(field, "[]")
}

// messageValue wraps the message in an array when the Recipe asks for one.
//
// Mux sends messages, plural: error.message is undefined and
// error.messages[0] is the sentence, which is backwards from every other
// provider in the collection. Writing the sentence into a field called
// "messages" as a bare string would look right in a diff and be wrong on the
// wire, so the Recipe says which it is.
func messageValue(field, message string) any {
	if strings.HasSuffix(field, "[]") {
		return []any{message}
	}

	return message
}

// codeValue renders an error code as the Recipe says the provider sends it.
//
// An undeclared code_type falls back to inferring from the value, which is
// right for Twilio and wrong for Adyen: "000" is a string there, and turning it
// into a number destroys the leading zeros a client is matching on. The
// inference stays the default so that Recipes written before this existed keep
// working, but it is a guess, and a Recipe that knows better says so.
func codeValue(codeType, value string) any {
	if codeType == "string" {
		return value
	}

	return numberOrString(value)
}

// numberOrString keeps an all-digit code a number, because Twilio's error codes
// are integers and a client comparing against 20404 must not be handed "20404".
func numberOrString(value string) any {
	if n, err := strconv.Atoi(value); err == nil {
		return n
	}

	return value
}

// servedPage is the page number this request asked for, counted the way the
// provider counts.
func servedPage(from paging, spec recipe.Pagination) int {
	first := spec.FirstPageNumber()

	name := spec.CursorParam
	if name == "" {
		name = spec.Style
	}

	raw := from.get(name)
	if raw == "" {
		return first
	}

	n, err := strconv.Atoi(raw)
	if err != nil || n < first {
		return first
	}

	return n
}

// pageCount is how many pages a total makes at this page size.
//
// An empty set is nought pages. Providers differ about whether it is nought or
// one, and nought is the reading that stops a loop rather than sending it
// after a page with nothing in it.
func pageCount(total, limit int) int {
	if total <= 0 || limit <= 0 {
		return 0
	}

	return (total + limit - 1) / limit
}

// writeRaw answers with a body exactly as given, past the identifier renaming
// and the resource envelope that writeRecord applies.
//
// A delete receipt keyed on its own name is not the resource, so presenting it
// as one would rename the very key the route just named.
func (s *Sandbox) writeRaw(w http.ResponseWriter, matched route, body map[string]any) int {
	status := matched.spec.Status
	if status == 0 {
		status = http.StatusOK
	}

	// The identifier is retyped the same way present would have retyped it.
	// A provider whose ids are JSON numbers sends a number here too, and a
	// receipt that disagrees with every other route about the type is its
	// own small lie.
	if spec, ok := s.recipe.Resources[matched.spec.Resource]; ok && spec.ID.Type == "number" {
		for key, value := range body {
			body[key] = numeric(value)
		}
	}

	if len(matched.spec.Fields) > 0 {
		body = withFields(body, matched.spec.Fields)
	}

	writeJSON(w, status, body)

	return status
}

// cursorValue renders the cursor as the Recipe says the provider sends it.
//
// A token by default, and the whole next-page URL for the providers whose
// pointer is one. Getting this wrong is not cosmetic: a client written against
// a token concatenates it onto a base URL, which is the exact mistake several
// of these Recipes describe in a comment and then taught anyway.
//
// The URL is empty when paging travels in the body, because there is no such
// URL to render; the token is the honest answer then.
func cursorValue(spec recipe.ListResponse, cursor, nextURL string) string {
	if spec.CursorURL == "" || nextURL == "" {
		return cursor
	}

	if spec.CursorURL == "path" {
		// Salesforce joins its nextRecordsUrl to the instance URL it
		// authenticated against, so the path is the whole of what it sends
		// and an absolute address would be joined to that.
		if at := strings.Index(nextURL, "://"); at >= 0 {
			if slash := strings.IndexByte(nextURL[at+3:], '/'); slash >= 0 {
				return nextURL[at+3+slash:]
			}
		}
	}

	return nextURL
}
