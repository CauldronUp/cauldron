package runtime

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The runtime emits resource.action and nothing else, so a provider whose
// events are named otherwise had them declared and never fired.
//
// Freshdesk calls it ticket_create. The convention would have produced
// ticket.created -- an event Freshdesk does not send and this Recipe does not
// declare -- so the change fired nothing at all, silently, and a handler
// waiting for one waited for ever with no error to show for it.
//
// 438 of the 482 events declared across this collection were in that state,
// and every one of the names is the provider's own: Bitbucket's repo:push,
// Zoom's meeting.started, Recurly's new_subscription_notification.
func TestARouteEmitsTheEventItNames(t *testing.T) {
	delivered := make(chan string, 4)

	sink := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err == nil {
			if name, ok := payload["type"].(string); ok {
				delivered <- name

				return
			}
		}

		delivered <- string(body)
	}))
	defer sink.Close()

	r, err := recipe.Open("freshdesk")
	if err != nil {
		t.Fatalf("open freshdesk: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Webhooks().Subscribe(sink.URL); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v2/tickets",
		strings.NewReader(`{"subject":"Kettle is broken","description":"It will not boil","email":"ada@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("cauldron-freshdesk-key", "X")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d\n%s", rec.Code, rec.Body)
	}

	select {
	case name := <-delivered:
		if name != "ticket_create" {
			t.Errorf("emitted %q, want ticket_create", name)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("nothing was delivered, which is what the convention did on its own")
	}
}
