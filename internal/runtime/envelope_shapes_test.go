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

// An error code that is all digits is not automatically a number. Twilio's
// 20404 really is an integer and a client comparing against "20404" never
// matches; Adyen's "000" really is a string and turning it into a number
// destroys the leading zeros a client is matching on. Inferring from the shape
// of the literal gets one of those two right, which is why a Recipe can say.
func TestAnErrorCodeIsSentAsTheRecipeDeclares(t *testing.T) {
	code := func(t *testing.T, name, path, header, key string) any {
		t.Helper()

		r, err := recipe.Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		sandbox, err := New(r, Options{Seed: 1})
		if err != nil {
			t.Fatalf("new sandbox: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set(header, "not-a-key")

		rec := httptest.NewRecorder()
		sandbox.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("%s status = %d, want 401\n%s", name, rec.Code, rec.Body)
		}

		return decode(t, rec)[key]
	}

	// Adyen declares code_type: string, so the leading zeros survive.
	if got := code(t, "adyen", "/v71/paymentMethods", "X-API-Key", "errorCode"); got != "000" {
		t.Errorf("adyen errorCode = %#v, want the string \"000\"", got)
	}

	// Twilio declares nothing, so the inference stands and the code stays a
	// number, which is what its clients switch on.
	if got := code(t, "twilio", "/2010-04-01/Accounts/ACcauldron00000000000000000000000/Messages.json",
		"Authorization", "code"); got != float64(20003) {
		t.Errorf("twilio code = %#v, want the number 20003", got)
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

// Datadog sends its errors as an array of bare strings rather than an array of
// objects. A client looping over the entries and reading .message from each,
// which is what almost every other provider needs, finds undefined on every
// one and throws instead of reporting anything.
func TestErrorsCanBeAnArrayOfBareStrings(t *testing.T) {
	r, err := recipe.Open("datadog")
	if err != nil {
		t.Fatalf("open datadog: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-org"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/monitor/99999999", nil)
	req.Header.Set("DD-API-KEY", "cauldron-api-key")
	req.Header.Set("DD-APPLICATION-KEY", "cauldron-app-key")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404\n%s", rec.Code, rec.Body)
	}

	body := decode(t, rec)

	errs, _ := body["errors"].([]any)
	if len(errs) != 1 {
		t.Fatalf("errors = %v, want one entry", body["errors"])
	}

	message, isString := errs[0].(string)
	if !isString {
		t.Fatalf("errors[0] is %T, want a bare string", errs[0])
	}

	if !strings.Contains(message, "Monitor not found") {
		t.Errorf("errors[0] = %q", message)
	}
}

// A route's declared status and empty_body used to be read on creates alone, so
// a Recipe could say its provider answers an update with 204 and nothing at all
// and be quietly ignored. Salesforce does exactly that on PATCH and DELETE, and
// a client calling .json() on the real response throws — the bug an emulator
// that helpfully returned the updated record would hide.
func TestAnUpdateAndDeleteHonourTheDeclaredStatusAndEmptyBody(t *testing.T) {
	r, err := recipe.Open("salesforce")
	if err != nil {
		t.Fatalf("open salesforce: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-org"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	call := func(method, target, body string) *httptest.ResponseRecorder {
		t.Helper()

		var req *http.Request
		if body != "" {
			req = httptest.NewRequest(method, target, strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req = httptest.NewRequest(method, target, nil)
		}

		req.Header.Set("Authorization", "Bearer 00Dcauldron!AQEAQCauldronAccessToken")

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		return rec
	}

	const account = "/services/data/v60.0/sobjects/Account/001cauldron000001A"

	updated := call(http.MethodPatch, account, `{"Industry":"Software"}`)
	if updated.Code != http.StatusNoContent {
		t.Errorf("update status = %d, want 204\n%s", updated.Code, updated.Body)
	}

	if body := strings.TrimSpace(updated.Body.String()); body != "" {
		t.Errorf("an update should answer with nothing at all, got %q", body)
	}

	// The update still happened, which is the half a 204 makes impossible to
	// see from the response.
	after := decode(t, call(http.MethodGet, account, ""))
	if after["Industry"] != "Software" {
		t.Errorf("Industry = %v, want the update to have landed", after["Industry"])
	}

	removed := call(http.MethodDelete, account, "")
	if removed.Code != http.StatusNoContent {
		t.Errorf("delete status = %d, want 204", removed.Code)
	}

	if body := strings.TrimSpace(removed.Body.String()); body != "" {
		t.Errorf("a delete should answer with nothing at all, got %q", body)
	}
}
