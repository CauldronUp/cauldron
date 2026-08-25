package runtime

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	matched, vars, ok := s.router.matchSelecting(r.Method, r.URL.Path, graphQLQuery(r), rawBody(r), r.URL.Query(), r.Header)
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
	record = s.declaredOnly(matched.spec.Resource, record)

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

	s.emitFor(matched.spec.Resource, "created", matched.spec.Emits, created)

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

func (s *Sandbox) get(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	id, addressable := s.identifier(matched, r, vars)
	if !addressable {
		return s.notFound(w, store.ErrNotFound, matched.spec.Resource, id, matched.spec.NotFound)
	}

	if s.misshapen(matched.spec.Resource, id) {
		return s.writeRecipeError(w, "id_malformed", 400, "id_malformed",
			"The identifier is not in the expected format: "+id+".", id)
	}

	record, err := s.store.GetBy(matched.spec.Resource, id, s.recipe.Resources[matched.spec.Resource].Alias)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id, matched.spec.NotFound)
	}

	// A record that exists but belongs to another scope must look absent, not
	// forbidden. Leaking existence across tenants is a real provider bug we do
	// not want to teach an application to rely on.
	if !store.Matches(record, s.scopeVars(matched, vars)) {
		return s.notFound(w, store.ErrNotFound, matched.spec.Resource, id, matched.spec.NotFound)
	}

	// Through writeRecord rather than building the body here, which is what
	// this used to do.
	//
	// The duplication was the bug: writeRecord applies the route's declared
	// constants and this did not, so fields: on a get route was accepted by
	// the validator, written into a Recipe, and then silently dropped. No
	// shipped Recipe declared any -- which is why it went unnoticed -- and
	// Lemon Squeezy is the first to need them, because JSON:API puts a type
	// beside every single record and that is a constant per route.
	//
	// It is the same shape as the fix a few lines above, where a route's
	// constants were dropped on every response that is not wrapped. Removing
	// the second copy is better than patching it.
	return s.writeRecord(w, matched, record)
}

func (s *Sandbox) update(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	id, addressable := s.identifier(matched, r, vars)
	if !addressable {
		return s.notFound(w, store.ErrNotFound, matched.spec.Resource, id, matched.spec.NotFound)
	}

	if s.misshapen(matched.spec.Resource, id) {
		return s.writeRecipeError(w, "id_malformed", 400, "id_malformed",
			"The identifier is not in the expected format: "+id+".", id)
	}

	changes, err := decodeBody(r)
	if err != nil {
		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", "Request body could not be parsed.")
	}

	changes = s.flatten(matched.spec.Resource, changes)
	changes = s.declaredOnly(matched.spec.Resource, changes)

	existing, err := s.store.Get(matched.spec.Resource, id)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id, matched.spec.NotFound)
	}

	if !store.Matches(existing, s.scopeVars(matched, vars)) {
		return s.notFound(w, store.ErrNotFound, matched.spec.Resource, id, matched.spec.NotFound)
	}

	// The optimistic lock, for the providers that keep one. A write that does
	// not say which version it is replacing, or says the wrong one, is
	// refused -- and refusing it here is the whole point: code that skips the
	// check passes every local test, because a test suite is the one place
	// where nothing else is writing to the same record.
	if raise, status, ok := s.versionCheck(matched.spec.Resource, existing, changes); !ok {
		return s.writeRecipeError(w, raise, status, "conflict",
			"The version sent is not the current version of this record.")
	}

	updated, err := s.store.Update(matched.spec.Resource, id, changes)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id, matched.spec.NotFound)
	}

	s.emitFor(matched.spec.Resource, "updated", matched.spec.Emits, updated)
	s.emitChanges(matched.spec, existing, updated)

	return s.writeRecord(w, matched, updated)
}

func (s *Sandbox) delete(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	id, addressable := s.identifier(matched, r, vars)
	if !addressable {
		return s.notFound(w, store.ErrNotFound, matched.spec.Resource, id, matched.spec.NotFound)
	}

	if s.misshapen(matched.spec.Resource, id) {
		return s.writeRecipeError(w, "id_malformed", 400, "id_malformed",
			"The identifier is not in the expected format: "+id+".", id)
	}

	record, err := s.store.Get(matched.spec.Resource, id)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id, matched.spec.NotFound)
	}

	if !store.Matches(record, s.scopeVars(matched, vars)) {
		return s.notFound(w, store.ErrNotFound, matched.spec.Resource, id, matched.spec.NotFound)
	}

	if err := s.store.Delete(matched.spec.Resource, id); err != nil {
		return s.notFound(w, err, matched.spec.Resource, id, matched.spec.NotFound)
	}

	s.emitFor(matched.spec.Resource, "deleted", matched.spec.Emits, record)

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

	// And the provider's own ceiling, for the ones that have one. A capped
	// page size is answered with less rather than refused, which is the
	// quieter of the two: a loop that stops when it receives fewer records
	// than it asked for stops on the first page, and the shop looks empty
	// rather than broken.
	if max := matched.spec.Pagination.MaxLimit; max > 0 && limit > max {
		// Unless the provider refuses instead, which is the other half of the
		// same rule and the louder half. A refusal is a bug the caller finds
		// on the first oversized request; a trim is one they find when a
		// report comes back short.
		if raise := matched.spec.Pagination.OverLimit; raise != "" {
			return s.writeRecipeError(w, raise, 400, "invalid_request",
				"The page size asked for is larger than this endpoint serves.")
		}

		limit = max
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

	body := s.listBody(list, page, limit, matched.spec.Resource, r.URL.Path,
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
	if list.LinkHeader {
		var links []string

		if page.HasMore && page.NextCursor != "" {
			if next := nextPageURL(r, matched.spec.Pagination, page.NextCursor); next != "" {
				links = append(links, "<"+next+`>; rel="next"`)
			}
		}

		// The last page still carries a Link header, and this emulator sent
		// none. Checked against api.github.com on 2026-08-22: page two of two
		// answers `rel="prev"` and no next.
		//
		// The difference is the whole reason a client reads this header. One
		// that stops when the header is missing never stops against GitHub,
		// because the header is there -- it is rel="next" that is not. The
		// Recipe asserting the header's absence was teaching exactly the loop
		// its own comment said the header exists to prevent.
		if list.PrevLink {
			if prev := prevPageURL(r, matched.spec.Pagination, limit); prev != "" {
				links = append(links, "<"+prev+`>; rel="prev"`)
			}
		}

		if len(links) > 0 {
			w.Header().Set("Link", strings.Join(links, ", "))
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

// versionCheck enforces a resource's optimistic lock and advances it.
//
// Three outcomes, because providers distinguish three: a write carrying the
// current version is allowed and moves the number on; a write carrying a
// different one is a conflict; and a write carrying none at all is a separate
// refusal, since a client that retries on the first should not retry on the
// second.
//
// A resource that declares no version field is unaffected, which is every
// Recipe written before this one.
func (s *Sandbox) versionCheck(resource string, existing store.Record, changes map[string]any) (string, int, bool) {
	spec, ok := s.recipe.Resources[resource]
	if !ok || spec.VersionField == "" {
		return "", 0, true
	}

	field := spec.VersionField

	sent, present := changes[field]
	if !present {
		if spec.VersionMissing == "" {
			return "", 0, true
		}

		return spec.VersionMissing, http.StatusBadRequest, false
	}

	current, known := asNumber(existing[field])
	asked, valid := asNumber(sent)

	if !known || !valid || current != asked {
		return spec.VersionConflict, http.StatusConflict, false
	}

	// The write is accepted, and the number moves. A caller that replays the
	// same body a second time is refused by the branch above, which is what
	// makes the lock worth having.
	changes[field] = int(current) + 1

	return "", 0, true
}

// asNumber reads a version however it happened to arrive. A fixture holds an
// int, and a request body holds a json.Number -- because decodeJSON uses one
// deliberately, to keep identifiers above 2^53 exact -- so a comparison that
// understood only one of the two would call every write stale.
//
// That is not a hypothetical. The first version of this compared the two
// directly and refused a write quoting exactly the right number, which is the
// failure mode an optimistic lock has to get right before it is worth having:
// a lock that never opens is not safer than no lock, it is just broken in a
// direction nobody tests for.
func asNumber(value any) (float64, bool) {
	switch typed := value.(type) {
	case json.Number:
		n, err := typed.Float64()

		return n, err == nil
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case string:
		n, err := strconv.ParseFloat(typed, 64)

		return n, err == nil
	}

	return 0, false
}
