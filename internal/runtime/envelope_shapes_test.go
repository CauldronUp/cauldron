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

// "-" means the provider sends no such field. That was honoured for the message
// everywhere and for the code and type only in the nested style, so a flat
// Recipe saying type_field: "-" got a literal "-" key in every error body.
//
// Three shipped Recipes were doing it and every case about them passed, because
// a case asserts the fields it names and stays silent about a key nobody
// thought to look for. It was found by reading an actual response rather than
// by the suite, which is the argument for doing that occasionally.
func TestTheOmissionMarkerNeverBecomesAFieldName(t *testing.T) {
	for _, provider := range []struct {
		recipe string
		path   string
		header string
		value  string
	}{
		{"woocommerce", "/wp-json/wc/v3/orders/999", "Authorization", "Basic Y2tfY2F1bGRyb246Y3NfY2F1bGRyb25jb25zdW1lcnNlY3JldA=="},
		{"wordpress", "/wp-json/wp/v2/posts/999", "Authorization", "Basic Y2F1bGRyb246Y2F1bGRyb24gYXBwIHBhc3Mgd29yZCBoZXJl"},
		{"digitalocean", "/v2/droplets/999", "Authorization", "Bearer dop_v1_cauldron"},
	} {
		r, err := recipe.Open(provider.recipe)
		if err != nil {
			t.Fatalf("open %s: %v", provider.recipe, err)
		}

		sandbox, err := New(r, Options{Seed: 1})
		if err != nil {
			t.Fatalf("new sandbox: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, provider.path, nil)
		req.Header.Set(provider.header, provider.value)

		rec := httptest.NewRecorder()
		sandbox.ServeHTTP(rec, req)

		body := decode(t, rec)

		if _, present := body["-"]; present {
			t.Errorf("%s sends a literal %q key: %v", provider.recipe, "-", body)
		}
	}
}

// SQS answers an idle queue with 200 and no Messages key at all. Sending an
// empty array instead is the helpful kind of wrong: every test passes, and the
// first quiet minute in production throws inside
// `for (const m of response.Messages)`.
func TestACollectionKeyCanBeOmittedWhenEmpty(t *testing.T) {
	r, err := recipe.Open("sqs")
	if err != nil {
		t.Fatalf("open sqs: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	const signed = "AWS4-HMAC-SHA256 Credential=CAULDRONKEY/20260115/eu-west-2/sqs/aws4_request, " +
		"SignedHeaders=host;x-amz-date, Signature=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	receive := func() map[string]any {
		t.Helper()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Authorization", signed)

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200\n%s", rec.Code, rec.Body)
		}

		return decode(t, rec)
	}

	// Nothing to give is a success that says nothing.
	if _, present := receive()["Messages"]; present {
		t.Error("an idle queue should omit the key rather than send an empty array")
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// And omitting it when empty only means something if it is there when not.
	messages, present := receive()["Messages"].([]any)
	if !present || len(messages) == 0 {
		t.Error("a busy queue should send the key")
	}
}

// AWS signs every request, so there is no fixed credential to compare. The
// shape is what can be checked, and that catches the failure that actually
// happens: credentials not configured, or the header wired up like every other
// API in this collection. A wrongly signed request is accepted, and the Recipe
// header says so rather than leaving it to be found.
func TestACredentialCanBeCheckedByShape(t *testing.T) {
	r, err := recipe.Open("sqs")
	if err != nil {
		t.Fatalf("open sqs: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	call := func(authorization string) int {
		t.Helper()

		req := httptest.NewRequest(http.MethodGet, "/queues", nil)
		if authorization != "" {
			req.Header.Set("Authorization", authorization)
		}

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		return rec.Code
	}

	signed := "AWS4-HMAC-SHA256 Credential=CAULDRONKEY/20260115/eu-west-2/sqs/aws4_request, " +
		"SignedHeaders=host;x-amz-date, Signature=" + strings.Repeat("ab", 32)

	if code := call(signed); code != http.StatusOK {
		t.Errorf("a well-formed signature should be accepted, got %d", code)
	}

	// The two failures that actually happen.
	if code := call(""); code != http.StatusForbidden {
		t.Errorf("no credentials at all should be refused, got %d", code)
	}

	if code := call("Bearer cauldron-sqs-key"); code != http.StatusForbidden {
		t.Errorf("a bearer token is not a signature, got %d", code)
	}

	// Truncated, so the shape is wrong even though the algorithm is right.
	if code := call("AWS4-HMAC-SHA256 Credential=CAULDRONKEY/20260115/eu-west-2/sqs/aws4_request"); code != http.StatusForbidden {
		t.Errorf("a half-formed signature should be refused, got %d", code)
	}
}

<<<<<<< HEAD
// Mux sends messages, plural, as an array: error.message is undefined and
// error.messages[0] is the sentence, which is backwards from every other
// provider in the collection. Writing the sentence into a field called
// "messages" as a bare string would look right in a diff and be wrong on the
// wire, so the Recipe says which it is with a [] marker.
//
// The nested style also used to drop responses.error.fields entirely, so a
// Recipe could declare a constant on every error and have it ignored. Both are
// checked here because both are about a declaration meaning what it says.
func TestAnErrorMessageCanBeAnArray(t *testing.T) {
	r, err := recipe.Open("mux")
	if err != nil {
		t.Fatalf("open mux: %v", err)
=======
// Pusher answers with an object keyed by channel name rather than an array, so
// looping over it as a list finds nothing. A channel nobody is on is absent
// from the object entirely rather than present with a zero, which is the same
// omit-when-empty idea one level down.
func TestACollectionCanBeKeyedRatherThanOrdered(t *testing.T) {
	r, err := recipe.Open("pusher")
	if err != nil {
		t.Fatalf("open pusher: %v", err)
>>>>>>> b936046 (Add a Pusher Recipe, a keyed collection, and an assertable text body)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

<<<<<<< HEAD
	req := httptest.NewRequest(http.MethodGet, "/video/v1/assets/doesnotexist", nil)
	req.SetBasicAuth("cauldron-token-id", "cauldron-mux-secret-key")
=======
	if err := s.Seed("small-app"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/apps/cauldron/channels?auth_key=cauldronappkey", nil)
>>>>>>> b936046 (Add a Pusher Recipe, a keyed collection, and an assertable text body)

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

<<<<<<< HEAD
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404\n%s", rec.Code, rec.Body)
	}

	failure, _ := decode(t, rec)["error"].(map[string]any)

	messages, isList := failure["messages"].([]any)
	if !isList || len(messages) != 1 {
		t.Fatalf("messages = %#v, want a one-element array", failure["messages"])
	}

	if text, ok := messages[0].(string); !ok || !strings.Contains(text, "Not Found") {
		t.Errorf("messages[0] = %#v, want the sentence", messages[0])
	}

	// The singular must not be there beside it, or a client finds both shapes
	// and picks whichever it was written for.
	if _, present := failure["message"]; present {
		t.Error("Mux sends no error.message, only error.messages")
=======
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200\n%s", rec.Code, rec.Body)
	}

	channels, isObject := decode(t, rec)["channels"].(map[string]any)
	if !isObject {
		t.Fatalf("channels is %T, want an object keyed by name", decode(t, rec)["channels"])
	}

	presence, _ := channels["presence-orders"].(map[string]any)
	if presence == nil {
		t.Fatalf("channels = %v, want the channel under its own name", channels)
	}

	// The key is the identifier, so it must not repeat inside the value.
	if _, repeated := presence["id"]; repeated {
		t.Error("the channel name is the key, so it should not also be a field")
	}

	if presence["user_count"] != float64(4) {
		t.Errorf("user_count = %v", presence["user_count"])
	}

	// A channel nobody is on is not in the object at all.
	if _, present := channels["quiet"]; present {
		t.Error("an unoccupied channel should be absent rather than zeroed")
>>>>>>> b936046 (Add a Pusher Recipe, a keyed collection, and an assertable text body)
	}
}
