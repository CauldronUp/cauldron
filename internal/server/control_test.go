package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func call(t *testing.T, s *Server, method, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	var reader *strings.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	var req *http.Request
	if reader == nil {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, reader)
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

// Every mutating control endpoint used to answer a GET, which the
// documentation never claimed and nothing enforced.
//
// It matters because a GET to a loopback address is something any web page can
// issue with an image tag, no preflight and no permission asked. A developer
// with a test suite running and a browser open was one visited page away from
// having their sandbox reset, their clock advanced thirty days, or a fault
// armed that made the suite fail for reasons nothing in the repository
// explained.
func TestAMutatingControlEndpointRefusesAGet(t *testing.T) {
	s := mounted(t)

	for _, path := range []string{
		"/_cauldron/stripe/reset",
		"/_cauldron/stripe/seed?fixture=small-shop",
		"/_cauldron/stripe/fault",
		"/_cauldron/stripe/emit",
		"/_cauldron/stripe/subscribe",
		"/_cauldron/clock/advance?duration=30d",
		"/_cauldron/clock/reset",
		"/_cauldron/reset",
		"/_cauldron/restore",
	} {
		rec := call(t, s, http.MethodGet, path, "", nil)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("GET %s = %d, want 405\n%s", path, rec.Code, rec.Body)
		}

		if allow := rec.Header().Get("Allow"); allow != http.MethodPost {
			t.Errorf("GET %s: Allow = %q, want POST", path, allow)
		}
	}

	// Reading is still a GET, and must stay one.
	for _, path := range []string{
		"/_cauldron/status",
		"/_cauldron/stripe/requests",
		"/_cauldron/snapshot",
		"/_cauldron/clock",
	} {
		if rec := call(t, s, http.MethodGet, path, "", nil); rec.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200\n%s", path, rec.Code, rec.Body)
		}
	}
}

// The other way a web page reaches an endpoint without asking the browser
// first: a form posting text/plain. The browser sends it cross-origin with no
// preflight, and a streaming JSON decoder reads the object and ignores the
// trailing "=" a form appends, so the body parses cleanly.
func TestAControlBodyMustBeJson(t *testing.T) {
	s := mounted(t)

	form := call(t, s, http.MethodPost, "/_cauldron/stripe/subscribe",
		`{"url":"http://elsewhere.example/webhooks"}=`,
		map[string]string{"Content-Type": "text/plain"})

	if form.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("a text/plain body = %d, want 415\n%s", form.Code, form.Body)
	}

	sandbox, _ := s.Sandbox("stripe")
	if endpoints := sandbox.Webhooks().Endpoints(); len(endpoints) != 0 {
		t.Errorf("the refused request registered an endpoint anyway: %v", endpoints)
	}

	// The same body as JSON is the supported path and must still work.
	ok := call(t, s, http.MethodPost, "/_cauldron/stripe/subscribe",
		`{"url":"http://localhost:8000/webhooks"}`,
		map[string]string{"Content-Type": "application/json"})

	if ok.Code != http.StatusOK {
		t.Fatalf("a json body = %d, want 200\n%s", ok.Code, ok.Body)
	}
}

// A browser attaches these; a test suite and the Cauldron CLI do not. That
// asymmetry is the whole defence: the caller being refused here is the one that
// cannot choose its own headers.
func TestAControlRequestFromAnotherOriginIsRefused(t *testing.T) {
	s := mounted(t)

	for _, headers := range []map[string]string{
		{"Origin": "https://evil.example"},
		{"Sec-Fetch-Site": "cross-site"},
		{"Sec-Fetch-Site": "same-site"},
	} {
		rec := call(t, s, http.MethodGet, "/_cauldron/status", "", headers)

		if rec.Code != http.StatusForbidden {
			t.Errorf("%v = %d, want 403\n%s", headers, rec.Code, rec.Body)
		}
	}

	// A request from the server's own origin, and one from a client that sends
	// neither header, both proceed.
	for _, headers := range []map[string]string{
		nil,
		{"Sec-Fetch-Site": "same-origin"},
		{"Sec-Fetch-Site": "none"},
		{"Origin": "http://example.com"},
	} {
		rec := call(t, s, http.MethodGet, "/_cauldron/status", "", headers)

		if rec.Code != http.StatusOK {
			t.Errorf("%v = %d, want 200\n%s", headers, rec.Code, rec.Body)
		}
	}
}
