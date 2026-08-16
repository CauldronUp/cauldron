package server

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/CauldronUp/cauldron/internal/clock"
	"github.com/CauldronUp/cauldron/internal/runtime"
)

// control implements the sandbox control API.
//
//	GET  /_cauldron/status
//	GET  /_cauldron/{recipe}/requests
//	POST /_cauldron/{recipe}/seed?fixture=small-shop
//	POST /_cauldron/{recipe}/reset
//	POST /_cauldron/{recipe}/fault      {"error":"rate_limit","count":3}
//	POST /_cauldron/{recipe}/emit       {"event":"customer.created","data":{}}
//	POST /_cauldron/{recipe}/subscribe  {"url":"http://localhost:8000/webhooks"}
//	POST /_cauldron/clock/advance       {"duration":"30d"}
//	POST /_cauldron/reset
func (s *Server) control(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, controlPrefix)
	path = strings.Trim(path, "/")

	if path == "" || path == "status" {
		s.status(w)
		return
	}

	head, rest, _ := strings.Cut(path, "/")

	switch head {
	case "clock":
		s.clockControl(w, r, rest)
		return
	case "reset":
		s.resetAll(w)
		return
	case "snapshot":
		writeJSON(w, http.StatusOK, s.Snapshot())
		return
	case "restore":
		s.restore(w, r)
		return
	}

	sandbox, ok := s.Sandbox(head)
	if !ok {
		s.writeError(w, http.StatusNotFound, "No recipe "+head+" is running.")
		return
	}

	switch rest {
	case "requests":
		s.requests(w, sandbox)
	case "seed":
		s.seed(w, r, sandbox)
	case "reset":
		s.reset(w, sandbox)
	case "fault":
		s.fault(w, r, sandbox)
	case "emit":
		s.emit(w, r, sandbox)
	case "subscribe":
		s.subscribe(w, r, sandbox)
	default:
		s.writeError(w, http.StatusNotFound, "Unknown control endpoint: "+rest+".")
	}
}

func (s *Server) status(w http.ResponseWriter) {
	type entry struct {
		Recipe   string   `json:"recipe"`
		Version  string   `json:"version"`
		Fixture  string   `json:"fixture,omitempty"`
		Requests int      `json:"requests"`
		Faults   int      `json:"armed_faults"`
		Webhooks int      `json:"webhooks_sent"`
		Errors   []string `json:"injectable_errors"`
	}

	out := struct {
		Time    string  `json:"time"`
		Recipes []entry `json:"recipes"`
	}{Time: s.clock.Now().Format("2006-01-02T15:04:05Z")}

	for _, name := range s.Names() {
		sandbox, ok := s.Sandbox(name)
		if !ok {
			continue
		}

		out.Recipes = append(out.Recipes, entry{
			Recipe:   sandbox.Name(),
			Version:  sandbox.Recipe().Version,
			Fixture:  sandbox.Fixture(),
			Requests: len(sandbox.Exchanges(0)),
			Faults:   len(sandbox.ArmedFaults()),
			Webhooks: len(sandbox.Webhooks().Deliveries()),
			Errors:   sandbox.Errors(),
		})
	}

	writeJSON(w, http.StatusOK, out)
}

func (s *Server) requests(w http.ResponseWriter, sandbox *runtime.Sandbox) {
	writeJSON(w, http.StatusOK, map[string]any{
		"recipe":   sandbox.Name(),
		"requests": sandbox.Exchanges(0),
	})
}

func (s *Server) seed(w http.ResponseWriter, r *http.Request, sandbox *runtime.Sandbox) {
	fixture := r.URL.Query().Get("fixture")

	if fixture == "" {
		body, err := recordFrom(r)
		if err == nil {
			fixture, _ = body["fixture"].(string)
		}
	}

	if fixture == "" {
		s.writeError(w, http.StatusBadRequest, "A fixture name is required. Available: "+strings.Join(sandbox.Fixtures(), ", ")+".")
		return
	}

	if err := sandbox.Seed(fixture); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"recipe": sandbox.Name(), "fixture": fixture, "seeded": true})
}

func (s *Server) reset(w http.ResponseWriter, sandbox *runtime.Sandbox) {
	if err := sandbox.Reset(); err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"recipe": sandbox.Name(), "reset": true})
}

// restore applies a snapshot posted as the request body.
func (s *Server) restore(w http.ResponseWriter, r *http.Request) {
	var snapshot Snapshot

	if err := json.NewDecoder(r.Body).Decode(&snapshot); err != nil {
		s.writeError(w, http.StatusBadRequest, "That is not a Cauldron snapshot.")
		return
	}

	if err := s.Restore(snapshot); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	restored := make([]string, 0, len(snapshot.Recipes))
	for name := range snapshot.Recipes {
		restored = append(restored, name)
	}

	sort.Strings(restored)

	writeJSON(w, http.StatusOK, map[string]any{"restored": restored, "now": s.clock.Now().Format("2006-01-02T15:04:05Z")})
}

func (s *Server) resetAll(w http.ResponseWriter) {
	for _, name := range s.Names() {
		sandbox, ok := s.Sandbox(name)
		if !ok {
			continue
		}

		if err := sandbox.Reset(); err != nil {
			s.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	s.clock.Reset()

	writeJSON(w, http.StatusOK, map[string]any{"reset": s.Names()})
}

func (s *Server) fault(w http.ResponseWriter, r *http.Request, sandbox *runtime.Sandbox) {
	body, err := recordFrom(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Body could not be parsed.")
		return
	}

	name, _ := body["error"].(string)
	if name == "" {
		s.writeError(w, http.StatusBadRequest, "An error name is required. Available: "+strings.Join(sandbox.Errors(), ", ")+".")
		return
	}

	fault := runtime.Fault{
		Error: name,
		Count: intFrom(body["count"]),
		Every: intFrom(body["every"]),
	}

	if path, ok := body["path"].(string); ok {
		fault.Path = path
	}

	if raw, ok := body["for"].(string); ok && raw != "" {
		d, err := clock.ParseDuration(raw)
		if err != nil {
			s.writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		fault.Until = s.clock.Now().Add(d)
	}

	if err := sandbox.Arm(fault); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"recipe": sandbox.Name(), "armed": name})
}

func (s *Server) emit(w http.ResponseWriter, r *http.Request, sandbox *runtime.Sandbox) {
	body, err := recordFrom(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Body could not be parsed.")
		return
	}

	event, _ := body["event"].(string)
	if event == "" {
		s.writeError(w, http.StatusBadRequest, "An event name is required.")
		return
	}

	data := map[string]any{}
	if raw, ok := body["data"].(map[string]any); ok {
		data = raw
	}

	delivery, err := sandbox.Webhooks().Emit(event, data)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"recipe":   sandbox.Name(),
		"event":    delivery.Event,
		"id":       delivery.ID,
		"endpoint": delivery.Endpoint,
		"status":   delivery.Status,
	})
}

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request, sandbox *runtime.Sandbox) {
	body, err := recordFrom(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Body could not be parsed.")
		return
	}

	url, _ := body["url"].(string)
	if url == "" {
		s.writeError(w, http.StatusBadRequest, "A url is required.")
		return
	}

	sandbox.Webhooks().Subscribe(url)

	writeJSON(w, http.StatusOK, map[string]any{"recipe": sandbox.Name(), "endpoints": sandbox.Webhooks().Endpoints()})
}

func (s *Server) clockControl(w http.ResponseWriter, r *http.Request, action string) {
	switch action {
	case "advance":
		body, err := recordFrom(r)
		if err != nil {
			s.writeError(w, http.StatusBadRequest, "Body could not be parsed.")
			return
		}

		raw, _ := body["duration"].(string)
		if raw == "" {
			raw = r.URL.Query().Get("duration")
		}

		d, err := clock.ParseDuration(raw)
		if err != nil {
			s.writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		now := s.clock.Advance(d)

		writeJSON(w, http.StatusOK, map[string]any{"now": now.Format("2006-01-02T15:04:05Z"), "advanced": raw})
	case "reset":
		s.clock.Reset()
		writeJSON(w, http.StatusOK, map[string]any{"now": s.clock.Now().Format("2006-01-02T15:04:05Z")})
	default:
		writeJSON(w, http.StatusOK, map[string]any{"now": s.clock.Now().Format("2006-01-02T15:04:05Z")})
	}
}

func intFrom(value any) int {
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	}

	return 0
}
