package runtime

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A failure can serve zero bytes, when zero bytes is what the provider sends.
//
// An error with no message falls back to generic wording, which is right almost
// always: a Recipe that forgot to write one should not silently serve nothing.
// But six Recipes here recorded that their provider genuinely answers an empty
// body and could not reproduce it -- Raygun's 404, PropelAuth's 401, the
// auto-generated 405s on Backblaze, Nile and Tigris, and Snipcart's credential
// failures, which carry no Content-Type either.
//
// Silence is a real answer and a hostile one. A client calling .json() on
// nothing throws exactly as it does on prose, and a Recipe that met it should be
// able to say so rather than serve a helpful sentence the provider never sent.
func TestAnErrorCanServeNoBodyAtAll(t *testing.T) {
	r, err := recipe.Open("firehydrant")
	if err != nil {
		t.Fatalf("open firehydrant: %v", err)
	}

	silent := *r
	silent.Errors = map[string]recipe.Error{}

	for name, e := range r.Errors {
		silent.Errors[name] = e
	}

	silent.Errors["authentication_error"] = recipe.Error{Status: 401, Empty: true}

	if err := silent.Validate(); err != nil {
		t.Fatalf("the constructed Recipe does not validate: %v", err)
	}

	s, err := New(&silent, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	// No header at all, so the absent verdict is reached. FireHydrant names a
	// rejected_error and no absent_error, so a wrong token would raise its own
	// sentence rather than the one under test.
	req := httptest.NewRequest(http.MethodGet, "/v1/incidents", nil)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status %d, want 401", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	if len(body) != 0 {
		t.Errorf("body is %d bytes, want none: %q", len(body), body)
	}

	// And no Content-Type invented for a body that is not there. A provider
	// answering nothing usually does not announce a type for it, and Snipcart
	// is on record that its own zero-byte failures carry none.
	if got := w.Header().Get("Content-Type"); got != "" {
		t.Errorf("Content-Type: %q was announced for an empty body", got)
	}
}

// An error that says nothing still says something, unless it asks not to.
//
// This is the compatibility half. Every Recipe that declares a status and no
// message keeps the generic wording it has always had -- the new field is opt-in
// precisely so that a forgotten message stays a visible mistake rather than
// becoming silent output.
func TestAnErrorWithNoMessageStillGetsTheGenericWording(t *testing.T) {
	r, err := recipe.Open("firehydrant")
	if err != nil {
		t.Fatalf("open firehydrant: %v", err)
	}

	quiet := *r
	quiet.Errors = map[string]recipe.Error{}

	for name, e := range r.Errors {
		quiet.Errors[name] = e
	}

	quiet.Errors["authentication_error"] = recipe.Error{Status: 401}

	s, err := New(&quiet, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	// No header at all, so the absent verdict is reached. FireHydrant names a
	// rejected_error and no absent_error, so a wrong token would raise its own
	// sentence rather than the one under test.
	req := httptest.NewRequest(http.MethodGet, "/v1/incidents", nil)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)
	if len(body) == 0 {
		t.Error("an error declaring no message served nothing, where it should still fall back to generic wording")
	}
}

// Declaring both is refused, because an empty body has no room for a message.
func TestAnErrorCannotBeBothEmptyAndSpoken(t *testing.T) {
	r, err := recipe.Open("firehydrant")
	if err != nil {
		t.Fatalf("open firehydrant: %v", err)
	}

	contradictory := *r
	contradictory.Errors = map[string]recipe.Error{}

	for name, e := range r.Errors {
		contradictory.Errors[name] = e
	}

	contradictory.Errors["authentication_error"] = recipe.Error{
		Status:  401,
		Empty:   true,
		Message: "but here is a sentence anyway",
	}

	if err := contradictory.Validate(); err == nil {
		t.Fatal("an error declaring both empty and a message validated")
	}
}
