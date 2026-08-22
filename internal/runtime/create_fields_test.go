package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A create does not always answer with just the record. Neon answers a branch
// create with the branch, the operations it started, and the connection
// strings -- and the operations are the only thing in that body that says when
// the branch can actually be used.
//
// A response carrying only the record would be the helpful kind of wrong: it
// looks finished. The branch comes back in state "init", which Neon defines as
// "being created but is not available for querying", and without the
// operations beside it there is nothing to wait on.

func createBranch(t *testing.T) map[string]any {
	t.Helper()

	r, err := recipe.Open("neon")
	if err != nil {
		t.Fatalf("open neon: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-project"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost,
		"/projects/cool-darkness-12345678/branches",
		strings.NewReader(`{"name":"feature/x"}`))
	req.Header.Set("Authorization", "Bearer neon_api_key_cauldron")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}

	var out map[string]any

	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}

	return out
}

func TestACreateCarriesTheRouteConstants(t *testing.T) {
	out := createBranch(t)

	for _, key := range []string{"branch", "operations", "connection_uris"} {
		if _, present := out[key]; !present {
			t.Errorf("the create response is missing %s", key)
		}
	}
}

func TestTheCreatedRecordIsStillThere(t *testing.T) {
	// The constants go beside the record, not instead of it.
	out := createBranch(t)

	branch, ok := out["branch"].(map[string]any)
	if !ok {
		t.Fatalf("the branch is not an object: %T", out["branch"])
	}

	if branch["name"] != "feature/x" {
		t.Errorf("the created branch was replaced by the constants: %v", branch)
	}

	// And it is not ready, which is the point of the operations beside it.
	if branch["current_state"] != "init" {
		t.Errorf("a created branch reported state %v, want init", branch["current_state"])
	}
}

func TestTheWorkItStartedIsUnfinished(t *testing.T) {
	// If the operations came back finished there would be nothing to wait
	// for, and the Recipe would be teaching that a 201 means done.
	out := createBranch(t)

	operations, ok := out["operations"].([]any)
	if !ok || len(operations) == 0 {
		t.Fatalf("no operations: %v", out["operations"])
	}

	for _, entry := range operations {
		operation, _ := entry.(map[string]any)
		if operation["status"] == "finished" {
			t.Errorf("a create answered with finished work: %v", operation)
		}
	}
}
