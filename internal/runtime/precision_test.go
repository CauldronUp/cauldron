package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A JSON request body was decoded with json.Unmarshal, which turns every
// number into a float64. That holds 53 bits of integer, so anything an API
// uses a 64-bit id or cursor for came back changed: a Discord snowflake, a
// Stripe cursor, a ledger amount in minor units above 2^53.
//
// 578730123365711993 went in and 578730123365712000 came out, with a 200 and
// a body that looks entirely reasonable.
//
// It was inconsistent as well, which is worse. A YAML fixture decodes to an
// int and stays exact, so the same field was right when seeded and wrong when
// written, in the same collection -- and a case asserting anything but that
// field passes either way.
func TestALongIntegerSurvivesBeingWritten(t *testing.T) {
	const exact = "578730123365711993"

	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/customers",
		strings.NewReader(`{"email":"ada@example.com","metadata":{"cursor":`+exact+`}}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk_test_cauldron")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("create = %d\n%s", rec.Code, rec.Body)
	}

	// The bytes, not a decoded value: what a client reads off the wire is the
	// thing that was wrong, and decoding it here with the wrong settings
	// would lose the digits a second time and hide it again.
	if body := rec.Body.String(); !strings.Contains(body, exact) {
		t.Errorf("the response does not carry %s:\n%s", exact, body)
	}
}
