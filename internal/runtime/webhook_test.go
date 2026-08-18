package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

func webhookSandbox(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("stripe")
	if err != nil {
		t.Fatalf("open stripe: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-shop"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	return s
}

// Subscribe takes an address from whoever is talking to the control plane and
// Cauldron then posts record data to it, signed, from inside the developer's
// network. That is an outbound request made on somebody else's say-so, so the
// address is worth checking.
//
// Loopback and private addresses stay allowed on purpose: delivering to the
// application under test on localhost is the documented use and the entire
// point. What is refused is the class no local receiver ever lives in.
func TestSubscribeRefusesAnAddressNoReceiverLivesAt(t *testing.T) {
	s := webhookSandbox(t)

	for _, endpoint := range []string{
		"file:///etc/passwd",
		"gopher://localhost:70/",
		"ftp://example.com/hook",
		// The cloud metadata service. It answers without credentials on most
		// providers and hands back instance role tokens.
		"http://169.254.169.254/latest/meta-data/",
		"http://[fe80::1]/hook",
		"http://",
		"::not a url at all",
	} {
		if err := s.Webhooks().Subscribe(endpoint); err == nil {
			t.Errorf("Subscribe(%q) was accepted", endpoint)
		}
	}

	// The addresses a real receiver lives at, which must keep working.
	for _, endpoint := range []string{
		"http://localhost:8000/webhooks",
		"http://127.0.0.1:3000/hooks/stripe",
		"http://host.docker.internal:8080/webhooks",
		"https://example.com/webhooks",
	} {
		if err := s.Webhooks().Subscribe(endpoint); err != nil {
			t.Errorf("Subscribe(%q) was refused: %v", endpoint, err)
		}
	}
}

// Delivery happens inside the request that triggered it, so the endpoint list
// is a length the control plane lets a caller choose and every write in a
// suite pays for.
func TestSubscribeCapsHowManyEndpointsOneSandboxDeliversTo(t *testing.T) {
	s := webhookSandbox(t)

	for i := 0; i < maxEndpoints; i++ {
		if err := s.Webhooks().Subscribe("http://localhost:8000/hook/" + string(rune('a'+i))); err != nil {
			t.Fatalf("endpoint %d was refused: %v", i, err)
		}
	}

	if err := s.Webhooks().Subscribe("http://localhost:8000/one-too-many"); err == nil {
		t.Errorf("the %dth endpoint was accepted", maxEndpoints+1)
	}
}

// Reset promises a sandbox indistinguishable from a freshly booted one, and
// that is what makes it safe to reuse between tests. A subscription used to
// survive it, so an endpoint registered once kept receiving every record the
// suite created afterwards, through every reset in between.
func TestResetClearsSubscriptions(t *testing.T) {
	s := webhookSandbox(t)

	if err := s.Webhooks().Subscribe("http://localhost:8000/webhooks"); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := s.Reset(); err != nil {
		t.Fatalf("reset: %v", err)
	}

	if endpoints := s.Webhooks().Endpoints(); len(endpoints) != 0 {
		t.Errorf("a reset sandbox still delivers to %v", endpoints)
	}
}

// Endpoints fail by timing out rather than by refusing, and delivery used to
// run in series. Three receivers that had stopped answering cost three
// timeouts, so a suite that registered a handful turned every create into a
// wait measured in tens of seconds and read as flakiness rather than as a
// webhook problem.
func TestDeliveryToSeveralDeadEndpointsCostsOneTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("this one waits for a real timeout")
	}

	s := webhookSandbox(t)

	// A listener that accepts and never answers, so the client waits out its
	// full timeout rather than failing fast on a refused connection.
	quiet := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * deliveryTimeout)
	}))
	defer quiet.Close()

	for i := 0; i < 4; i++ {
		if err := s.Webhooks().Subscribe(quiet.URL + "/hook/" + string(rune('a'+i))); err != nil {
			t.Fatalf("subscribe: %v", err)
		}
	}

	started := time.Now()

	req := httptest.NewRequest(http.MethodPost, "/v1/customers",
		strings.NewReader("email=buyer@cauldron.test"))
	req.Header.Set("Authorization", "Bearer sk_test_cauldron")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusCreated {
		t.Fatalf("create = %d\n%s", rec.Code, rec.Body)
	}

	elapsed := time.Since(started)

	// Four endpoints, one timeout. In series this was four, and the margin is
	// generous enough that a slow machine does not fail it.
	if elapsed > 2*deliveryTimeout {
		t.Errorf("four dead endpoints took %s, which is more than one timeout of %s", elapsed, deliveryTimeout)
	}

	// Every attempt is still recorded, in the order the endpoints were
	// registered, so what a test reads back does not depend on which receiver
	// gave up first.
	deliveries := s.Webhooks().Deliveries()
	if len(deliveries) != 4 {
		t.Fatalf("recorded %d deliveries for 4 endpoints", len(deliveries))
	}

	for i, delivery := range deliveries {
		want := quiet.URL + "/hook/" + string(rune('a'+i))
		if delivery.Endpoint != want {
			t.Errorf("delivery %d went to %s, want %s", i, delivery.Endpoint, want)
		}
	}
}
