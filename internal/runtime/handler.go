package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	// Faults are evaluated before anything else, including auth. A rate limit
	// that only fires on well-formed authenticated requests would not be a
	// faithful rate limit.
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

	matched, vars, ok := s.router.match(r.Method, r.URL.Path)
	if !ok {
		if allowed := s.router.allowedMethods(r.URL.Path); len(allowed) > 0 {
			w.Header().Set("Allow", strings.Join(allowed, ", "))
			exchange.Status = s.writeRecipeError(w, "method_not_allowed", 405, "method_not_allowed", "This method is not supported on this path.")

			return
		}

		exchange.Status = s.writeRecipeError(w, "unknown_route", 404, "unknown_route", "Unrecognised request URL.")

		return
	}

	exchange.Resource = matched.spec.Resource
	exchange.Op = matched.spec.Operation

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

	for field, value := range scopeVars(matched, vars) {
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
		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", err.Error())
	}

	s.emitFor(matched.spec.Resource, "created", created)

	status := matched.spec.Status
	if status == 0 {
		status = http.StatusOK
	}

	s.writeRouteHeaders(w, matched, created)

	if matched.spec.EmptyBody {
		w.WriteHeader(status)
		return status
	}

	writeJSON(w, status, s.resourceBody(matched.spec.Resource, created))

	return status
}

// identifier resolves the record id for a request. Most providers put it in
// the path; RPC-shaped ones like Slack put it in the query string or the body,
// which the Recipe declares with id_from.
func identifier(matched route, r *http.Request, vars map[string]string) string {
	if matched.spec.IDFrom == "" {
		return vars["id"]
	}

	source, name, _ := strings.Cut(matched.spec.IDFrom, ":")

	switch source {
	case "query":
		return r.URL.Query().Get(name)
	case "body":
		body, err := decodeBody(r)
		if err != nil {
			return ""
		}

		if value, ok := body[name]; ok {
			return fmt.Sprint(value)
		}
	}

	return ""
}

func (s *Sandbox) get(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	id := identifier(matched, r, vars)

	record, err := s.store.Get(matched.spec.Resource, id)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id)
	}

	// A record that exists but belongs to another scope must look absent, not
	// forbidden. Leaking existence across tenants is a real provider bug we do
	// not want to teach an application to rely on.
	if !store.Matches(record, scopeVars(matched, vars)) {
		return s.writeRecipeError(w, "resource_missing", 404, "resource_missing",
			"No such "+matched.spec.Resource+": "+id+".",
			matched.spec.Resource+": "+id)
	}

	writeJSON(w, http.StatusOK, s.resourceBody(matched.spec.Resource, record))

	return http.StatusOK
}

func (s *Sandbox) update(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	id := identifier(matched, r, vars)

	changes, err := decodeBody(r)
	if err != nil {
		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", "Request body could not be parsed.")
	}

	changes = s.flatten(matched.spec.Resource, changes)

	existing, err := s.store.Get(matched.spec.Resource, id)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id)
	}

	if !store.Matches(existing, scopeVars(matched, vars)) {
		return s.writeRecipeError(w, "resource_missing", 404, "resource_missing",
			"No such "+matched.spec.Resource+": "+id+".",
			matched.spec.Resource+": "+id)
	}

	updated, err := s.store.Update(matched.spec.Resource, id, changes)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id)
	}

	s.emitFor(matched.spec.Resource, "updated", updated)

	writeJSON(w, http.StatusOK, s.resourceBody(matched.spec.Resource, updated))

	return http.StatusOK
}

func (s *Sandbox) delete(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	id := identifier(matched, r, vars)

	record, err := s.store.Get(matched.spec.Resource, id)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, id)
	}

	if !store.Matches(record, scopeVars(matched, vars)) {
		return s.writeRecipeError(w, "resource_missing", 404, "resource_missing",
			"No such "+matched.spec.Resource+": "+id+".",
			matched.spec.Resource+": "+id)
	}

	if err := s.store.Delete(matched.spec.Resource, id); err != nil {
		return s.notFound(w, err, matched.spec.Resource, id)
	}

	s.emitFor(matched.spec.Resource, "deleted", record)

	writeJSON(w, http.StatusOK, s.resourceBody(matched.spec.Resource, store.Record{
		"id":      id,
		"object":  matched.spec.Resource,
		"deleted": true,
	}))

	return http.StatusOK
}

func (s *Sandbox) list(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	limit := matched.spec.Pagination.Limit
	if limit <= 0 {
		limit = 10
	}

	limit = queryInt(r, "limit", limit)

	cursor := r.URL.Query().Get("starting_after")
	if cursor == "" {
		cursor = r.URL.Query().Get("cursor")
	}

	page, err := s.store.ListWhere(matched.spec.Resource, scopeVars(matched, vars), cursor, limit)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", "No such cursor: "+cursor+".")
		}

		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", err.Error())
	}

	writeJSON(w, http.StatusOK, s.listBody(page, matched.spec.Resource, r.URL.Path))

	return http.StatusOK
}

// missingHeader reports the first required header a request does not carry.
func (s *Sandbox) missingHeader(r *http.Request) (header, errorName string, ok bool) {
	for name, raise := range s.recipe.RequiredHeaders {
		if r.Header.Get(name) != "" {
			continue
		}

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
		if id, ok := record["id"].(string); ok {
			value = strings.ReplaceAll(value, "{id}", id)
		}

		w.Header().Set(name, value)
	}
}

// present renames the identifier to the property the provider actually uses.
// The store keeps every record keyed by "id" so fixtures and internal lookups
// stay uniform; only the wire shape changes, which is where it matters.
func (s *Sandbox) present(resource string, record store.Record) store.Record {
	spec, ok := s.recipe.Resources[resource]
	if !ok {
		return record
	}

	renamed := spec.ID.Field != "" && spec.ID.Field != "id"

	if !renamed && !s.nests(spec) {
		return record
	}

	out := make(store.Record, len(record))

	for key, value := range record {
		if key == "id" && renamed {
			// setPath, so a dotted name nests. Contentful keeps the
			// identifier at sys.id rather than at the top level.
			setPath(out, spec.ID.Field, value)
			continue
		}

		// A field declared with "in" moves under that sub-object on the wire.
		// HubSpot reads contact.properties.email, and a client written against
		// it finds nothing at the top level.
		if field, declared := spec.Fields[key]; declared && field.In != "" {
			nested, _ := out[field.In].(map[string]any)
			if nested == nil {
				nested = map[string]any{}
				out[field.In] = nested
			}

			nested[key] = value

			continue
		}

		out[key] = value
	}

	return out
}

// nests reports whether any of a resource's fields live under a sub-object.
func (s *Sandbox) nests(spec recipe.Resource) bool {
	for _, field := range spec.Fields {
		if field.In != "" {
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

		nested, isObject := record[field.In].(map[string]any)
		if !isObject {
			continue
		}

		if value, present := nested[name]; present {
			record[name] = value
		}
	}

	// The sub-objects have been unpacked, so drop them rather than storing the
	// same values twice under two shapes.
	for _, field := range spec.Fields {
		if field.In != "" {
			delete(record, field.In)
		}
	}

	return record
}

// presentAll renames identifiers across a page.
func (s *Sandbox) presentAll(resource string, records []store.Record) []store.Record {
	spec, ok := s.recipe.Resources[resource]
	if !ok || ((spec.ID.Field == "" || spec.ID.Field == "id") && !s.nests(spec)) {
		return records
	}

	out := make([]store.Record, 0, len(records))

	for _, record := range records {
		out = append(out, s.present(resource, record))
	}

	return out
}

// resourceBody shapes a single object according to the Recipe's declared style.
// Shopify nests it under the singular resource name; most providers return the
// object itself, and a client written for one shape breaks on the other.
func (s *Sandbox) resourceBody(resource string, record store.Record) any {
	record = s.present(resource, record)

	success := s.recipe.Responses.Success.Fields

	if s.recipe.Responses.Resource.Style != "wrapped" {
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

	spec := s.recipe.Responses.Resource

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

	child, ok := body[head].(map[string]any)
	if !ok {
		child = map[string]any{}
		body[head] = child
	}

	setPath(child, rest, value)
}

// listBody shapes a page according to the Recipe's declared list style.
func (s *Sandbox) listBody(page store.Page, resource, path string) any {
	spec := s.recipe.Responses.List
	page.Records = s.presentAll(resource, page.Records)

	switch spec.Style {
	case "bare":
		// GitHub and friends return the array itself, with paging in headers.
		// A caller doing json.Unmarshal into a slice must not receive an object.
		return page.Records
	case "wrapped":
		body := map[string]any{s.collectionName(resource, spec.Key): page.Records}

		if spec.CursorField != "" && page.NextCursor != "" {
			setPath(body, spec.CursorField, page.NextCursor)
		}

		if spec.HasMoreField != "" {
			setPath(body, spec.HasMoreField, page.HasMore)
		}

		if spec.CountField != "" {
			setPath(body, spec.CountField, page.Total)
		}

		body = withFields(body, s.recipe.Responses.Success.Fields)

		return withFields(body, spec.Fields)
	default:
		body := map[string]any{
			"object":   "list",
			"data":     page.Records,
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
			setPath(body, spec.CursorField, page.NextCursor)
		}

		body = withFields(body, s.recipe.Responses.Success.Fields)

		return withFields(body, spec.Fields)
	}
}

// scopeVars extracts the scope filters for a request from its path parameters.
func scopeVars(matched route, vars map[string]string) map[string]any {
	if len(matched.spec.Scope) == 0 {
		return nil
	}

	out := make(map[string]any, len(matched.spec.Scope))

	for _, name := range matched.spec.Scope {
		out[name] = vars[name]
	}

	return out
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
		return s.writeRecipeError(w, "unknown_route", 404, "unknown_route", "Unrecognised request URL.")
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

	_, _ = s.webhooks.Emit(event, record)
}

// authorised checks the credential according to the Recipe's auth scheme.
func (s *Sandbox) authorised(r *http.Request) bool {
	auth := s.recipe.Auth

	// A Recipe that declares no keys accepts anything, so an author can model
	// routes first and tighten auth later.
	if auth.Scheme == "" || auth.Scheme == "none" || len(auth.Keys) == 0 {
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
		return true
	}

	if auth.Prefix != "" {
		if !strings.HasPrefix(presented, auth.Prefix) {
			return false
		}

		presented = strings.TrimPrefix(presented, auth.Prefix)
	}

	presented = strings.TrimSpace(presented)

	for _, key := range auth.Keys {
		if presented == key {
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

	if defined, ok := s.recipe.Errors[name]; ok {
		status = defined.Status

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
		if defined.Message != "" {
			message = strings.ReplaceAll(defined.Message, "{detail}", detail)
		}

		for key, value := range defined.Headers {
			w.Header().Set(key, value)
		}

		extra = defined.Fields
	}

	writeJSON(w, status, s.errorBody(category, code, message, status, extra))

	return status
}

// errorBody shapes a failure according to the Recipe's declared error style.
func (s *Sandbox) errorBody(category, code, message string, status int, extra map[string]any) map[string]any {
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
		nested := map[string]any{"message": message}

		// Providers disagree about which of these they send. Stripe sends a
		// type and a code, Airtable only a type, Vercel only a code. Declaring
		// the omission keeps the emulator from inventing a field a client
		// could come to depend on and then lose.
		if spec.CodeField != "-" {
			nested["code"] = code
		}

		if spec.TypeField != "-" {
			nested["type"] = category
		}

		key := spec.Key
		if key == "" {
			key = "error"
		}

		return withFields(map[string]any{key: nested}, extra)
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
		setPath(body, field, message)
	}

	if spec.CodeField != "" && code != "" {
		setPath(body, spec.CodeField, numberOrString(code))
	}

	if spec.TypeField != "" {
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

		// One request can fail several ways at once, which is why SendGrid and
		// Clerk both send an array. The emulator reports one, in the shape a
		// client already has to loop over.
		//
		// Declared fields sit beside the array rather than inside each entry:
		// Clerk's trace id belongs to the response, not to one failure.
		return withFields(withFields(map[string]any{key: []any{body}}, spec.Fields), extra)
	}

	return withFields(withFields(body, spec.Fields), extra)
}

// numberOrString keeps an all-digit code a number, because Twilio's error codes
// are integers and a client comparing against 20404 must not be handed "20404".
func numberOrString(value string) any {
	if n, err := strconv.Atoi(value); err == nil {
		return n
	}

	return value
}
