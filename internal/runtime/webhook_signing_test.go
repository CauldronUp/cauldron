package runtime

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The README says webhook signing works. Nothing checked it.
//
// Signing is the one webhook feature whose value depends entirely on being
// exactly right. An application verifies the signature with the provider's
// own SDK, so a signature that is the wrong length, over the wrong bytes, or
// keyed on the wrong secret does not fail here -- it fails in the one place
// the fake exists to keep it from failing, and until then it makes the
// application's verification code look correct.
//
// Every other webhook test asserts that delivery happened, that subscriptions
// are capped, that a reset clears them, that the payload has the right shape.
// None of them looked at the signature at all, so "Working" covered a value
// nothing had ever recomputed.
func TestAWebhookSignatureIsAnHMACOverTimestampAndBody(t *testing.T) {
	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	if r.Webhooks.Signing.Scheme != "hmac-sha256" {
		t.Fatalf("stripe no longer signs with hmac-sha256, which this test is about: %q", r.Webhooks.Signing.Scheme)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	body := []byte(`{"id":"evt_1","type":"customer.created"}`)
	at := s.clock.Now()

	signed := s.webhooks.sign(body, at)
	if signed == "" {
		t.Fatal("a Recipe declaring hmac-sha256 and a secret produced no signature")
	}

	// t=<unix>,v1=<hex>, which is the shape an SDK parses before it verifies.
	parts := strings.Split(signed, ",")
	if len(parts) != 2 || !strings.HasPrefix(parts[0], "t=") || !strings.HasPrefix(parts[1], "v1=") {
		t.Fatalf("signature is not t=,v1= shaped: %q", signed)
	}

	stamp, err := strconv.ParseInt(strings.TrimPrefix(parts[0], "t="), 10, 64)
	if err != nil {
		t.Fatalf("timestamp is not a number: %v", err)
	}

	if stamp != at.Unix() {
		t.Errorf("signature is stamped %d and the sandbox clock says %d", stamp, at.Unix())
	}

	// Recomputed here rather than compared against what sign produced, so the
	// test fails if the payload, the key or the algorithm changes -- not only
	// if the formatting does.
	mac := hmac.New(sha256.New, []byte(r.Webhooks.Signing.Secret))
	mac.Write([]byte(fmt.Sprintf("%d.%s", stamp, body)))
	want := hex.EncodeToString(mac.Sum(nil))

	if got := strings.TrimPrefix(parts[1], "v1="); got != want {
		t.Errorf("signature is %s and an HMAC over \"timestamp.body\" with the declared secret is %s", got, want)
	}
}

// A different body must not sign the same, which is the property the whole
// mechanism rests on and the one a constant would satisfy the other test with.
func TestADifferentBodySignsDifferently(t *testing.T) {
	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	at := s.clock.Now()

	first := s.webhooks.sign([]byte(`{"id":"evt_1"}`), at)
	second := s.webhooks.sign([]byte(`{"id":"evt_2"}`), at)

	if first == second {
		t.Errorf("two different payloads signed identically at the same instant: %s", first)
	}
}

// A Recipe that declares no signing sends none, rather than a signature over
// an empty secret that would verify against nothing and look like one.
func TestARecipeWithoutSigningSendsNoSignature(t *testing.T) {
	r, err := recipe.Open("github")
	if err != nil {
		t.Fatalf("open github: %v", err)
	}

	if r.Webhooks.Signing.Scheme == "hmac-sha256" && r.Webhooks.Signing.Secret != "" {
		t.Skip("github now declares signing, so this test needs a Recipe that does not")
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if signed := s.webhooks.sign([]byte(`{"id":1}`), s.clock.Now()); signed != "" {
		t.Errorf("a Recipe declaring no signing produced %q", signed)
	}
}
