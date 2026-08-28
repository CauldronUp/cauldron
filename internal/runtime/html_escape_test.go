package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Angle brackets and ampersands go on the wire as themselves.
//
// Go's encoding/json escapes the three characters that matter inside an HTML
// document unless it is told not to, because the standard library's default
// assumes the result may be embedded in a page. No provider in this collection
// does that, and it is not a difference a conformance case comparing decoded
// values can see: both forms decode to the same string, so every case here
// passed while the bytes on the wire were not the bytes the provider sends.
//
// Open-Meteo is what found it. A misspelled variable name answers with the
// service's internal Swift generic type signature, angle brackets and all, and
// the real API sends a literal 0x3c where Cauldron sent an escape. Anything
// reading the raw body rather than a parsed one -- a snapshot test, a regex
// over the response, a signature over a webhook payload -- disagreed with
// production for a reason nothing in the Recipe could express.
//
// It was never only the new Recipe. Discourse ships a topic whose fancy_title
// is an HTML entity, which is the whole point of that field, and it was going
// out escaped. Akeneo's HAL links carry two query parameters joined by an
// ampersand, so the paging URL a client is meant to follow was not the one
// Akeneo sends. This test uses a shipped Recipe on purpose.
func TestAngleBracketsAndAmpersandsAreNotEscaped(t *testing.T) {
	r, err := recipe.Open("discourse")
	if err != nil {
		t.Fatalf("open discourse: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-forum"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/t/101.json", nil)
	req.Header.Set("Api-Key", "cauldron_discourse_key")
	req.Header.Set("Api-Username", "ada")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()

	if strings.Contains(body, `\u0026`) {
		t.Errorf("the body escapes its ampersand: %.300q", body)
	}

	if !strings.Contains(body, "Ada&rsquo;s notes on the engine") {
		t.Errorf("the entity is not on the wire as itself: %.300q", body)
	}
}
