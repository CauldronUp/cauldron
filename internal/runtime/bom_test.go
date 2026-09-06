package runtime

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Some providers put a byte-order mark in front of their JSON.
//
// Authorize.Net does, on every response. Recorded live from its public sandbox
// on 2026-09-06: the first three bytes are EF BB BF and the fourth is {.
// Python's json.loads raises "Unexpected UTF-8 BOM (decode using utf-8-sig)"
// on the decoded text, and Go's encoding/json refuses it too.
//
// A fake that leaves it out is wrong in the direction that costs the most: the
// client works locally and throws on its first real response, which is the
// failure this project exists to prevent happening inside the tool meant to
// prevent it. So it is served, opt-in, by the Recipes whose providers send one.
func TestABodyOrderMarkIsServedWhenDeclared(t *testing.T) {
	body := map[string]any{"messages": map[string]any{"resultCode": "Error"}}

	for _, tt := range []struct {
		name string
		bom  bool
		want bool
	}{
		{"declared", true, true},
		{"not declared", false, false},
	} {
		w := httptest.NewRecorder()

		s := &Sandbox{recipe: &recipe.Recipe{Responses: recipe.Responses{BOM: tt.bom}}}
		s.writeJSON(w, http.StatusOK, body)

		raw := w.Body.Bytes()
		got := bytes.HasPrefix(raw, []byte{0xEF, 0xBB, 0xBF})

		if got != tt.want {
			t.Errorf("%s: leading BOM = %v, want %v (first bytes %x)", tt.name, got, tt.want, raw[:4])
		}

		// And whatever else is true, the rest of it still has to be the
		// document: a mark in front of something unparseable helps nobody.
		if err := json.Unmarshal(bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF}), &map[string]any{}); err != nil {
			t.Errorf("%s: the body after the mark does not parse: %v", tt.name, err)
		}
	}
}

// And the mark is what breaks a strict decoder, which is the whole point.
func TestAStrictDecoderRefusesTheMarkedBody(t *testing.T) {
	w := httptest.NewRecorder()

	s := &Sandbox{recipe: &recipe.Recipe{Responses: recipe.Responses{BOM: true}}}
	s.writeJSON(w, http.StatusOK, map[string]any{"ok": true})

	var out map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &out); err == nil {
		t.Error("encoding/json accepted a body with a leading BOM; the Recipe would be teaching that it is harmless")
	}
}
