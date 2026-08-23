package runtime

import (
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

	// The body only says which event this is when the envelope happens to
	// carry one, and plenty do not: Freshdesk's payload is the ticket's
	// fields under freshdesk_webhook and names no event anywhere. Reading
	// payload["type"] was reading the default envelope, so this test broke
	// the moment Freshdesk described its own -- while the thing it exists to
	// check, that the route emits the name it declares, was still true.
	//
	// The sink proves a delivery arrived. Which event it was comes off the
	// delivery record, where it does not depend on the envelope.
	sink := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
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
	case <-delivered:
	case <-time.After(2 * time.Second):
		t.Fatal("nothing was delivered, which is what the convention did on its own")
	}

	deliveries := s.Webhooks().Deliveries()
	if len(deliveries) == 0 {
		t.Fatal("nothing was recorded, so there is no event to check")
	}

	if name := deliveries[len(deliveries)-1].Event; name != "ticket_create" {
		t.Errorf("emitted %q, want ticket_create", name)
	}
}
