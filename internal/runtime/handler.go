package runtime

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CauldronUp/cauldron/internal/store"
)

// listEnvelope is the shape list responses take.
//
// This currently mirrors Stripe's ({object, data, has_more}). Providers differ
// here — Shopify nests under a resource key, GitHub returns a bare array — so
// the envelope will need to move into the Recipe format before a second
// provider can be modelled faithfully. Recorded as a known limitation rather
// than left to be discovered.
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
		exchange.Status = s.create(w, r, matched)
	case "get":
		exchange.Status = s.get(w, matched, vars)
	case "update":
		exchange.Status = s.update(w, r, matched, vars)
	case "delete":
		exchange.Status = s.delete(w, matched, vars)
	case "list":
		exchange.Status = s.list(w, r, matched)
	default:
		exchange.Status = s.writeRecipeError(w, "unsupported_operation", 500, "unsupported_operation", "This operation is not implemented.")
	}
}

func (s *Sandbox) create(w http.ResponseWriter, r *http.Request, matched route) int {
	record, err := decodeBody(r)
	if err != nil {
		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", "Request body could not be parsed.")
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

	writeJSON(w, http.StatusOK, record)

	return http.StatusOK
}

func (s *Sandbox) update(w http.ResponseWriter, r *http.Request, matched route, vars map[string]string) int {
	changes, err := decodeBody(r)
	if err != nil {
		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", "Request body could not be parsed.")
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

func (s *Sandbox) list(w http.ResponseWriter, r *http.Request, matched route) int {
	limit := matched.spec.Pagination.Limit
	if limit <= 0 {
		limit = 10
	}

	limit = queryInt(r, "limit", limit)

	cursor := r.URL.Query().Get("starting_after")
	if cursor == "" {
		cursor = r.URL.Query().Get("cursor")
	}

	page, err := s.store.List(matched.spec.Resource, cursor, limit)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", "No such cursor: "+cursor+".")
		}

		return s.writeRecipeError(w, "invalid_request", 400, "invalid_request", err.Error())
	}

	writeJSON(w, http.StatusOK, listEnvelope{
		Object:  "list",
		Data:    page.Records,
		HasMore: page.HasMore,
		Next:    page.NextCursor,
	})

	return http.StatusOK
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
