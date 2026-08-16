package runtime

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CauldronUp/cauldron/internal/store"
)

// listEnvelope is the Stripe-style list shape.
//
// Providers genuinely disagree about this: Stripe wraps in {object, data,
// has_more}, GitHub returns a bare array, Shopify nests under a resource key.
// The Recipe declares which, and the handler honours it — see listBody.
type listEnvelope struct {
	Object  string         `json:"object"`
	Data    []store.Record `json:"data"`
	HasMore bool           `json:"has_more"`
	Next    string         `json:"next_cursor,omitempty"`
}

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
		exchange.Status = s.get(w, matched, vars)
	case "update":
		exchange.Status = s.update(w, r, matched, vars)
	case "delete":
		exchange.Status = s.delete(w, matched, vars)
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

	for field, value := range scopeVars(matched, vars) {
		record[field] = value
	}

	s.applyDefaults(matched.spec.Resource, record)

	if missing := s.missingRequired(matched.spec.Resource, record); len(missing) > 0 {
		return s.writeRecipeError(
			w, "parameter_missing", 400, "parameter_missing",
			"Missing required parameter: "+strings.Join(missing, ", ")+".",
		)
	}

	created, err := s.store.Create(matched.spec.Resource, record)
	if err != nil {
		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", err.Error())
	}

	s.emitFor(matched.spec.Resource, "created", created)

	writeJSON(w, http.StatusOK, created)

	return http.StatusOK
}

func (s *Sandbox) get(w http.ResponseWriter, matched route, vars map[string]string) int {
	record, err := s.store.Get(matched.spec.Resource, vars["id"])
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, vars["id"])
	}

	// A record that exists but belongs to another scope must look absent, not
	// forbidden. Leaking existence across tenants is a real provider bug we do
	// not want to teach an application to rely on.
	if !store.Matches(record, scopeVars(matched, vars)) {
		return s.writeRecipeError(w, "resource_missing", 404, "resource_missing",
			"No such "+matched.spec.Resource+": "+vars["id"]+".")
	}

	writeJSON(w, http.StatusOK, record)

	return http.StatusOK
}

func (s *Sandbox) update(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	changes, err := decodeBody(r)
	if err != nil {
		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", "Request body could not be parsed.")
	}

	existing, err := s.store.Get(matched.spec.Resource, vars["id"])
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, vars["id"])
	}

	if !store.Matches(existing, scopeVars(matched, vars)) {
		return s.writeRecipeError(w, "resource_missing", 404, "resource_missing",
			"No such "+matched.spec.Resource+": "+vars["id"]+".")
	}

	updated, err := s.store.Update(matched.spec.Resource, vars["id"], changes)
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, vars["id"])
	}

	s.emitFor(matched.spec.Resource, "updated", updated)

	writeJSON(w, http.StatusOK, updated)

	return http.StatusOK
}

func (s *Sandbox) delete(w http.ResponseWriter, matched route, vars map[string]string) int {
	record, err := s.store.Get(matched.spec.Resource, vars["id"])
	if err != nil {
		return s.notFound(w, err, matched.spec.Resource, vars["id"])
	}

	if !store.Matches(record, scopeVars(matched, vars)) {
		return s.writeRecipeError(w, "resource_missing", 404, "resource_missing",
			"No such "+matched.spec.Resource+": "+vars["id"]+".")
	}

	if err := s.store.Delete(matched.spec.Resource, vars["id"]); err != nil {
		return s.notFound(w, err, matched.spec.Resource, vars["id"])
	}

	s.emitFor(matched.spec.Resource, "deleted", record)

	writeJSON(w, http.StatusOK, store.Record{
		"id":      vars["id"],
		"object":  matched.spec.Resource,
		"deleted": true,
	})

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

	writeJSON(w, http.StatusOK, s.listBody(page, matched.spec.Resource))

	return http.StatusOK
}

// listBody shapes a page according to the Recipe's declared list style.
func (s *Sandbox) listBody(page store.Page, resource string) any {
	spec := s.recipe.Responses.List

	switch spec.Style {
	case "bare":
		// GitHub and friends return the array itself, with paging in headers.
		// A caller doing json.Unmarshal into a slice must not receive an object.
		return page.Records
	case "wrapped":
		return map[string]any{s.collectionName(resource, spec.Key): page.Records}
	default:
		return listEnvelope{
			Object:  "list",
			Data:    page.Records,
			HasMore: page.HasMore,
			Next:    page.NextCursor,
		}
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
		"No such "+resource+": "+id+".",
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
		user, _, ok := r.BasicAuth()
		if !ok {
			return false
		}

		presented = user
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
func (s *Sandbox) writeRecipeError(w http.ResponseWriter, name string, fallback ...any) int {
	status := http.StatusInternalServerError
	code := name
	message := "An error occurred."

	if len(fallback) == 3 {
		status, _ = fallback[0].(int)
		code, _ = fallback[1].(string)
		message, _ = fallback[2].(string)
	}

	if defined, ok := s.recipe.Errors[name]; ok {
		status = defined.Status

		if defined.Code != "" {
			code = defined.Code
		}

		if defined.Message != "" {
			message = defined.Message
		}

		for key, value := range defined.Headers {
			w.Header().Set(key, value)
		}
	}

	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"type":    name,
			"code":    code,
			"message": message,
		},
	})

	return status
}
