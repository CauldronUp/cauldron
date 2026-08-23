package runtime

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A signature has two halves that can be wrong independently: how the digest
// is wrapped, and what string it was taken over. A conformance case can pin
// the first with a pattern -- sha256= and hex, or bare base64 -- and cannot
// see the second at all, because every shape of the right length looks alike.
//
// So the base string is checked here, against an HMAC computed the long way.
// Signing the body where Slack signs "v0:<ts>:<body>" produces a perfectly
// well-formed v0= value that Slack's verifier rejects, and no pattern in any
// Recipe would notice.
func TestWhatEachFormatSigns(t *testing.T) {
	const secret = "cauldron-signing-secret"

	body := []byte(`{"hello":"world"}`)
	at := time.Unix(1767225600, 0).UTC()

	digest := func(over string) []byte {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(over))

		return mac.Sum(nil)
	}

	for _, c := range []struct {
		format string
		want   string
	}{
		// The body alone.
		{"prefixed-hex", "sha256=" + hex.EncodeToString(digest(string(body)))},
		{"base64", base64.StdEncoding.EncodeToString(digest(string(body)))},

		// The body with something in front of it.
		{"slack", "v0=" + hex.EncodeToString(digest(fmt.Sprintf("v0:%d:%s", at.Unix(), body)))},
		{"stripe", fmt.Sprintf("t=%d,v1=%s", at.Unix(), hex.EncodeToString(digest(fmt.Sprintf("%d.%s", at.Unix(), body))))},

		// Unset is stripe, so a Recipe nobody has looked at keeps the shape it
		// had rather than changing under it.
		{"", fmt.Sprintf("t=%d,v1=%s", at.Unix(), hex.EncodeToString(digest(fmt.Sprintf("%d.%s", at.Unix(), body))))},
	} {
		q := &webhookQueue{recipe: &recipe.Recipe{
			Webhooks: recipe.Webhooks{
				Signing: recipe.Signing{Scheme: "hmac-sha256", Secret: secret, Format: c.format},
			},
		}}

		if got := q.sign(body, at); got != c.want {
			t.Errorf("format %q signed to %q, want %q", c.format, got, c.want)
		}
	}
}

// Nothing to sign with is nothing to send, rather than a digest of the empty
// string, which would look like a signature and verify as nothing.
func TestAnUnsignedRecipeSendsNoSignature(t *testing.T) {
	for _, signing := range []recipe.Signing{
		{Scheme: "none"},
		{Scheme: "hmac-sha256"},
	} {
		q := &webhookQueue{recipe: &recipe.Recipe{Webhooks: recipe.Webhooks{Signing: signing}}}

		if got := q.sign([]byte("{}"), time.Unix(0, 0)); got != "" {
			t.Errorf("scheme %q secret %q signed to %q, want nothing", signing.Scheme, signing.Secret, got)
		}
	}
}
