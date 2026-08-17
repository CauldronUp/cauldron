package runtime

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Slack is the reason these capabilities exist: an RPC-shaped API where the
// identifier is a query parameter or a body field, every response carries an
// ok flag, and failures arrive with HTTP 200. Each of those breaks a different
// assumption a REST-shaped emulator would otherwise bake in permanently.

func slackSandbox(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("slack")
	if err != nil {
		t.Fatalf("open slack: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-workspace"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	return s
}

// slackCall drives the sandbox with a Slack bot token rather than the Stripe
// key the shared helper assumes.
func slackCall(t *testing.T, s *Sandbox, method, target, body string) map[string]any {
	t.Helper()

	return decode(t, slackRequest(t, s, method, target, body))
}

func slackRequest(t *testing.T, s *Sandbox, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, target, reader)
	req.Header.Set("Authorization", "Bearer xoxb-cauldron")

	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

func TestIdentifierFromAQueryParameter(t *testing.T) {
	s := slackSandbox(t)

	body := slackCall(t, s, http.MethodGet, "/api/conversations.info?channel=C0000000002", "")

	channel, _ := body["channel"].(map[string]any)
	if channel == nil {
		t.Fatalf("body = %v", body)
	}

	if channel["name"] != "engineering" {
		t.Errorf("name = %v", channel["name"])
	}
}

func TestIdentifierFromTheRequestBody(t *testing.T) {
	s := slackSandbox(t)

	body := slackCall(t, s, http.MethodPost, "/api/conversations.rename",
		`{"channel":"C0000000001","name":"announcements"}`)

	channel, _ := body["channel"].(map[string]any)
	if channel == nil {
		t.Fatalf("body = %v", body)
	}

	// The body is read once to find the identifier and again to apply the
	// change. A drained body would leave the name untouched and say nothing.
	if channel["name"] != "announcements" {
		t.Errorf("name = %v, want the rename to have applied", channel["name"])
	}

	if channel["id"] != "C0000000001" {
		t.Errorf("id = %v", channel["id"])
	}
}

func TestAFailureKeepsAnHTTP200(t *testing.T) {
	s := slackSandbox(t)

	rec := slackRequest(t, s, http.MethodGet, "/api/conversations.info?channel=C9999999999", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: a client branching on the status code must still see the failure", rec.Code)
	}

	body := decode(t, rec)

	if body["ok"] != false {
		t.Errorf("ok = %v, want false", body["ok"])
	}

	if body["error"] != "channel_not_found" {
		t.Errorf("error = %v", body["error"])
	}

	if _, present := body["message"]; present {
		t.Error("Slack sends no message on an error, so neither should the emulator")
	}
}

func TestEveryResponseCarriesTheSuccessFlag(t *testing.T) {
	s := slackSandbox(t)

	for _, target := range []string{
		"/api/conversations.list",
		"/api/users.list",
		"/api/conversations.info?channel=C0000000001",
	} {
		body := slackCall(t, s, http.MethodGet, target, "")

		if body["ok"] != true {
			t.Errorf("%s: ok = %v, want true", target, body["ok"])
		}
	}
}

func TestANestedCursorIsNestedRatherThanFlattened(t *testing.T) {
	s := slackSandbox(t)

	body := slackCall(t, s, http.MethodGet, "/api/conversations.list?limit=2", "")

	metadata, _ := body["response_metadata"].(map[string]any)
	if metadata == nil {
		t.Fatalf("response_metadata missing from %v", body)
	}

	cursor, _ := metadata["next_cursor"].(string)
	if cursor == "" {
		t.Error("expected a next_cursor inside response_metadata")
	}

	if _, flat := body["next_cursor"]; flat {
		t.Error("the cursor should not also appear at the top level")
	}
}

func TestTimestampIdentifiersComeFromTheSandboxClock(t *testing.T) {
	s := slackSandbox(t)

	first := slackCall(t, s, http.MethodPost, "/api/chat.postMessage",
		`{"channel":"C0000000002","text":"one"}`)
	second := slackCall(t, s, http.MethodPost, "/api/chat.postMessage",
		`{"channel":"C0000000002","text":"two"}`)

	firstTS, _ := first["message"].(map[string]any)["ts"].(string)
	secondTS, _ := second["message"].(map[string]any)["ts"].(string)

	if firstTS == "" || secondTS == "" {
		t.Fatalf("missing ts: %q %q", firstTS, secondTS)
	}

	// Two messages in the same second must still be distinguishable, which is
	// what the counter after the decimal point is for.
	if firstTS == secondTS {
		t.Errorf("both messages got ts %q", firstTS)
	}

	if _, present := first["message"].(map[string]any)["id"]; present {
		t.Error("Slack calls it ts, not id")
	}
}
