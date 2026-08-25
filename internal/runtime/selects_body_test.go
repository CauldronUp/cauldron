package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Some providers answer one path with two shapes depending on what the request
// asked for, and nothing outside the body tells them apart.
//
// Gemini is the reason this exists. A blocked prompt is documented as a 200
// where "the prompt was blocked and no candidates are returned", and a
// permitted one carries candidates and no block reason. Same path, same
// method, same headers, same query string -- so selects, which is fed the
// GraphQL query field and nothing else, could not choose between them, and the
// backlog recorded the whole provider as unservable for it.
//
// The two halves are asserted together because either alone would pass against
// a router that ignored the body: a fallback answering everything looks right
// for the permitted prompt and wrong only for the blocked one.
func TestARouteCanBeChosenByAWordAnywhereInTheBody(t *testing.T) {
	r, err := recipe.Open("gemini")
	if err != nil {
		t.Fatalf("open gemini: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("project"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	ask := func(prompt string) map[string]any {
		t.Helper()

		body := `{"contents":[{"role":"user","parts":[{"text":"` + prompt + `"}]}]}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost,
			"/v1beta/models/gemini-2.5-flash:generateContent", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-goog-api-key", "AIzaCauldronGeminiFixtureKey00000000000")
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200: %s", prompt, rec.Code, rec.Body.String())
		}

		var decoded map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
			t.Fatalf("%s: decode: %v", prompt, err)
		}

		return decoded
	}

	// The marker is a word in a prompt, nested two objects and an array deep,
	// which is the whole point: nothing at the top level of the body says
	// which answer this is.
	refused := ask("a request that gets blocked")

	feedback, ok := refused["promptFeedback"].(map[string]any)
	if !ok {
		t.Fatalf("the blocked prompt carried no promptFeedback: %v", refused)
	}

	if feedback["blockReason"] != "SAFETY" {
		t.Errorf("blockReason = %v, want SAFETY", feedback["blockReason"])
	}

	if _, present := refused["candidates"]; present {
		t.Error("the blocked prompt carried candidates, and Gemini's own words are that none are returned")
	}

	// And the other half, which a router ignoring the body would also answer
	// correctly -- which is why both are here.
	allowed := ask("at what temperature does a kettle boil")

	candidates, ok := allowed["candidates"].([]any)
	if !ok || len(candidates) == 0 {
		t.Fatalf("the permitted prompt carried no candidates: %v", allowed)
	}

	if feedback, ok := allowed["promptFeedback"].(map[string]any); ok {
		if _, blocked := feedback["blockReason"]; blocked {
			t.Error("the permitted prompt carried a blockReason")
		}
	}
}
