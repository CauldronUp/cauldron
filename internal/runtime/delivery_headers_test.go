package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A conformance case reads the recorded delivery, which is not the thing a
// handler sees. The header has to be on the request as well, and the two are
// set in different places -- one when the delivery is built and one when it
// is posted -- so a change to either could leave a case green and a
// subscriber without the timestamp its verifier needs.
func TestTheSignedTimestampReachesASubscriber(t *testing.T) {
	seen := make(chan http.Header, 2)

	sink := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen <- r.Header.Clone()
	}))
	defer sink.Close()

	r, err := recipe.Open("slack")
	if err != nil {
		t.Fatalf("open slack: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Webhooks().Subscribe(sink.URL); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/conversations.create",
		strings.NewReader(`{"name":"incident-response"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer xoxb-cauldron")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("create = %d\n%s", rec.Code, rec.Body)
	}

	select {
	case headers := <-seen:
		if got := headers.Get("X-Slack-Signature"); !strings.HasPrefix(got, "v0=") {
			t.Errorf("X-Slack-Signature = %q, want a v0= value", got)
		}

		// The half a verifier cannot do without, and the half that used to be
		// missing: the signature covers "v0:<ts>:<body>" and nothing else
		// told the subscriber what <ts> was.
		if got := headers.Get("X-Slack-Request-Timestamp"); got == "" {
			t.Error("the subscriber got a signature over a timestamp and no timestamp")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("nothing was delivered")
	}
}
