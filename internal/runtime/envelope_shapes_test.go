package runtime

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// SendGrid answers a send with 202, no body, and the message id in a header,
// and reports failures as an array because one request can be wrong twice.

func sendgridSandbox(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("sendgrid")
	if err != nil {
		t.Fatalf("open sendgrid: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	return s
}

func sendgridRequest(t *testing.T, s *Sandbox, method, target, body string, authorised bool) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}

	if authorised {
		req.Header.Set("Authorization", "Bearer SG.cauldron")
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

func TestAnEmptyBodyResponseIsActuallyEmpty(t *testing.T) {
	s := sendgridSandbox(t)

	rec := sendgridRequest(t, s, http.MethodPost, "/v3/mail/send",
		`{"subject":"Your receipt","from":"billing@example.com"}`, true)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202\n%s", rec.Code, rec.Body)
	}

	// A client calling .json() on this throws, and that is correct. An
	// emulator returning a helpful object would hide the bug until production.
	if body := strings.TrimSpace(rec.Body.String()); body != "" {
		t.Errorf("body = %q, want nothing at all", body)
	}
}

func TestTheGeneratedIdentifierReachesTheDeclaredHeader(t *testing.T) {
	s := sendgridSandbox(t)

	first := sendgridRequest(t, s, http.MethodPost, "/v3/mail/send",
		`{"subject":"One","from":"billing@example.com"}`, true)
	second := sendgridRequest(t, s, http.MethodPost, "/v3/mail/send",
		`{"subject":"Two","from":"billing@example.com"}`, true)

	opaque := regexp.MustCompile(`^[A-Za-z0-9]{22}$`)

	firstID := first.Header().Get("X-Message-Id")
	secondID := second.Header().Get("X-Message-Id")

	if !opaque.MatchString(firstID) {
		t.Errorf("X-Message-Id = %q, want an opaque 22 character id", firstID)
	}

	// Each send must be separately correlatable, which a constant header value
	// would quietly prevent.
	if firstID == secondID {
		t.Errorf("both sends reported %q", firstID)
	}
}

func TestAnOpaqueIdentifierCarriesNoPrefix(t *testing.T) {
	s := sendgridSandbox(t)

	rec := sendgridRequest(t, s, http.MethodPost, "/v3/mail/send",
		`{"subject":"Prefix check","from":"billing@example.com"}`, true)

	id := rec.Header().Get("X-Message-Id")

	if strings.Contains(id, "_") || strings.Contains(id, "-") {
		t.Errorf("id = %q, want no prefix separator", id)
	}
}

func TestAFailureIsReportedAsAnArray(t *testing.T) {
	s := sendgridSandbox(t)

	rec := sendgridRequest(t, s, http.MethodPost, "/v3/mail/send",
		`{"from":"billing@example.com"}`, true)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}

	body := decode(t, rec)

	errs, _ := body["errors"].([]any)
	if len(errs) == 0 {
		t.Fatalf("body = %v, want an errors array", body)
	}

	first, _ := errs[0].(map[string]any)
	if message, _ := first["message"].(string); !strings.Contains(message, "subject") {
		t.Errorf("message = %v, want it to name the missing field", first["message"])
	}

	// Code written to loop over errors must not also find a bare message, or
	// the two shapes drift apart in the client.
	for _, absent := range []string{"error", "message"} {
		if _, present := body[absent]; present {
			t.Errorf("SendGrid sends no top-level %q", absent)
		}
	}
}

func TestAuthFailsBeforeAnythingIsSent(t *testing.T) {
	s := sendgridSandbox(t)

	rec := sendgridRequest(t, s, http.MethodPost, "/v3/mail/send",
		`{"subject":"Should not send","from":"billing@example.com"}`, false)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}

	if _, present := decode(t, rec)["errors"]; !present {
		t.Error("an auth failure should use the same array shape as any other")
	}
}
