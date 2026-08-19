package server

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestNetworkEndpointArmsConditions(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/network",
		`{"latency":"800ms","jitter":"200ms"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var out struct {
		Recipe string `json:"recipe"`
		Armed  string `json:"armed"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if out.Recipe != "stripe" {
		t.Errorf("recipe = %q, want stripe", out.Recipe)
	}

	if out.Armed != "latency 800ms ±200ms" {
		t.Errorf("armed = %q, want a readable description", out.Armed)
	}

	sandbox, _ := s.Sandbox("stripe")
	if len(sandbox.ArmedNetwork()) != 1 {
		t.Errorf("armed %d conditions on the sandbox, want 1", len(sandbox.ArmedNetwork()))
	}
}

func TestNetworkEndpointAcceptsEveryCondition(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/network",
		`{"bandwidth":50,"timeout":"5s","limit":1024,"slice":64,"probability":0.5,"count":3,"path":"/v1/charges","for":"30s"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200: %s", rec.Code, rec.Body.String())
	}

	sandbox, _ := s.Sandbox("stripe")
	armed := sandbox.ArmedNetwork()

	if len(armed) != 1 {
		t.Fatalf("armed %d conditions, want 1", len(armed))
	}

	c := armed[0]

	if c.Bandwidth != 50 || c.Limit != 1024 || c.Slice != 64 || c.Count != 3 {
		t.Errorf("numeric fields did not survive the round trip: %+v", c)
	}

	if c.Timeout.String() != "5s" {
		t.Errorf("timeout = %v, want 5s", c.Timeout)
	}

	if c.Probability != 0.5 {
		t.Errorf("probability = %v, want 0.5", c.Probability)
	}

	if c.Path != "/v1/charges" {
		t.Errorf("path = %q, want /v1/charges", c.Path)
	}

	if c.Until.IsZero() {
		t.Error("a 'for' duration should set an expiry")
	}
}

func TestNetworkEndpointClears(t *testing.T) {
	s := mounted(t)

	do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/network", `{"latency":"1s"}`)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/network", `{"clear":true}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200: %s", rec.Code, rec.Body.String())
	}

	sandbox, _ := s.Sandbox("stripe")
	if len(sandbox.ArmedNetwork()) != 0 {
		t.Error("clearing should leave nothing armed")
	}
}

// Arming nothing is almost always a mistyped flag. Accepting it would leave
// somebody waiting for a failure that is never coming.
func TestNetworkEndpointRejectsEmptyConditions(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/network", `{}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400", rec.Code)
	}

	if body := rec.Body.String(); body == "" {
		t.Error("the refusal should say what is missing")
	}
}

func TestNetworkEndpointRejectsAnUnparseableDuration(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/network", `{"latency":"soon"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400: %s", rec.Code, rec.Body.String())
	}
}

func TestNetworkEndpointRejectsAnImpossibleProbability(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/network",
		`{"reset":true,"probability":25}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400: %s", rec.Code, rec.Body.String())
	}
}

func TestStatusReportsArmedNetworkConditions(t *testing.T) {
	s := mounted(t)

	do(t, s, http.MethodPost, "http://localhost/_cauldron/stripe/network", `{"latency":"1s","bandwidth":10}`)

	rec := do(t, s, http.MethodGet, "http://localhost/_cauldron/status", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}

	var out struct {
		Recipes []struct {
			Recipe  string   `json:"recipe"`
			Network []string `json:"network"`
		} `json:"recipes"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}

	for _, entry := range out.Recipes {
		if entry.Recipe != "stripe" {
			continue
		}

		if len(entry.Network) != 1 {
			t.Fatalf("status reported %d network entries, want 1", len(entry.Network))
		}

		if entry.Network[0] != "latency 1s, bandwidth 10KB/s" {
			t.Errorf("status said %q", entry.Network[0])
		}

		return
	}

	t.Fatal("stripe was missing from the status payload")
}

// The control API is read-only unless the request is a mutation from an
// allowed origin; the network endpoint must be held to the same rule as
// every other one that changes something.
func TestNetworkEndpointIsAMutation(t *testing.T) {
	s := mounted(t)

	rec := do(t, s, http.MethodGet, "http://localhost/_cauldron/stripe/network", "")

	if rec.Code == http.StatusOK {
		t.Error("GET should not arm network conditions")
	}
}
