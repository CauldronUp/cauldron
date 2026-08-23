package runtime

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/store"
)

// Emit returned the delivery it had built before anything was sent, and
// deliver sets the outcome on a copy. So the value a caller read reported an
// endpoint of "" and a status of 0 whether the receiver verified the
// signature and answered 200, refused it, threw a 500, or was never
// listening.
//
// That is the failure this tool exists to prevent, made by the tool: somebody
// wiring up signature verification gets the same output for working and
// broken code, and `cauldron emit` printed a cheerful line either way.
func TestEmitReportsWhatEachSubscriberSaid(t *testing.T) {
	sink := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer sink.Close()

	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Webhooks().Subscribe(sink.URL); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	delivery, attempts, err := s.Webhooks().Emit("customer.created", store.Record{"id": "cus_1"})
	if err != nil {
		t.Fatalf("emit: %v", err)
	}

	if delivery.ID == "" {
		t.Error("the delivery carries no id, which is what names it afterwards")
	}

	if len(attempts) != 1 {
		t.Fatalf("got %d attempts, want one per subscriber", len(attempts))
	}

	if attempts[0].Status != http.StatusTeapot {
		t.Errorf("status = %d, want %d: the receiver's answer, not a zero", attempts[0].Status, http.StatusTeapot)
	}

	if attempts[0].Endpoint != sink.URL {
		t.Errorf("endpoint = %q, want %q", attempts[0].Endpoint, sink.URL)
	}
}

// An endpoint nobody is listening on is the case a developer most wants told
// about, and it was reported identically to a success.
func TestEmitReportsAnEndpointThatRefuses(t *testing.T) {
	sink := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	dead := sink.URL
	sink.Close()

	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Webhooks().Subscribe(dead); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	_, attempts, err := s.Webhooks().Emit("customer.created", store.Record{"id": "cus_1"})
	if err != nil {
		t.Fatalf("emit: %v", err)
	}

	if len(attempts) != 1 {
		t.Fatalf("got %d attempts, want one", len(attempts))
	}

	if attempts[0].Error == "" {
		t.Error("delivering to a closed port reported no error")
	}
}

// Nothing attempted is no attempts, rather than one carrying a zero status
// that reads as somebody having answered.
func TestEmitWithNoSubscribersReportsNoAttempts(t *testing.T) {
	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	delivery, attempts, err := s.Webhooks().Emit("customer.created", store.Record{"id": "cus_1"})
	if err != nil {
		t.Fatalf("emit: %v", err)
	}

	if len(attempts) != 0 {
		t.Errorf("got %d attempts with nobody subscribed", len(attempts))
	}

	// Still recorded, because the event happened whether or not anybody was
	// there to receive it.
	if got := s.Webhooks().Deliveries(); len(got) != 1 || got[0].ID != delivery.ID {
		t.Errorf("the delivery was not recorded: %v", got)
	}
}
