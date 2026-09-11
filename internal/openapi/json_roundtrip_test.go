package openapi

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// The document below is reduced from Todoist's published description, which is
// 1.2MB of valid JSON. Two things in it matter, and the failure needs both.
//
// The `summary` carries an astral character as a JSON surrogate pair. yaml.v3
// implements YAML 1.1, which is not a superset of JSON, and refuses that escape
// with "found invalid Unicode character escape code" -- so the raw parse fails
// and the JSON path takes over. That path exists for exactly this: one
// thumbs-up in one example of Twilio's description used to make all 121 of its
// paths unreadable.
//
// And the `source` string begins with a newline. yaml.v3 emits a
// leading-newline literal block with an explicit indentation indicator --
// `source: |4` -- and its own parser then reads that indicator as a larger
// indent than the lines it wrote. The block ends early, the next line is parsed
// as a mapping, and `-H "Authorization: Bearer ..."` becomes a mapping key. The
// error surfaces as "mapping values are not allowed in this context" hundreds
// of lines from anything a reader could act on.
//
// So the JSON path -- JSON to value to YAML *text* to value -- threw away any
// description that was valid JSON and tripped that emitter quirk. Nine strings
// in Todoist's description begin with a newline, and `cauldron drift` reported
// the whole provider as unreadable.
const leadingNewlineJSON = `{
  "openapi": "3.0.0",
  "info": {"title": "T", "version": "1"},
  "paths": {
    "/api/v1/uploads": {
      "post": {
        "summary": "Upload \uD83D\uDC4D",
        "x-codeSamples": [
          {
            "lang": "curl",
            "source": "\n$ curl https://api.todoist.com/api/v1/uploads \\\n       -H \"Authorization: Bearer 0123456789abcdef\" \\\n       -F file=@/path/to/file.pdf\n"
          }
        ],
        "responses": {"200": {"description": "ok"}}
      }
    }
  }
}`

// First, prove both halves of the quirk are real and still present in the
// pinned yaml.v3. Without this the test below could pass for the wrong reason
// -- because the emitter was fixed upstream rather than because this package
// stopped relying on it.
func TestTheYAMLTextRoundTripStillLosesALeadingNewlineBlock(t *testing.T) {
	raw := []byte(leadingNewlineJSON)

	var direct any
	if err := yaml.Unmarshal(raw, &direct); err == nil {
		t.Skip("yaml.v3 now reads JSON surrogate escapes; the JSON path is no longer reached")
	}

	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		t.Fatalf("the fixture is not valid JSON: %v", err)
	}

	text, err := yaml.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var back any
	if err := yaml.Unmarshal(text, &back); err == nil {
		t.Skip("yaml.v3 now round-trips this; the JSON path no longer needs to avoid text")
	}
}

// And then: a description that is valid JSON is read, whatever its strings
// contain.
func TestADescriptionThatIsValidJSONIsReadWhateverItsStringsContain(t *testing.T) {
	parsed, err := Parse([]byte(leadingNewlineJSON))
	if err != nil {
		t.Fatalf("a valid JSON description was refused: %v", err)
	}

	if _, ok := parsed.Paths["/api/v1/uploads"]; !ok {
		t.Errorf("the path is missing; got %v", keysOfPaths(parsed))
	}
}

// The surrogate-pair case on its own still works, so the fix does not trade one
// provider for another.
func TestADescriptionWithAnAstralEscapeIsStillRead(t *testing.T) {
	const doc = `{
  "openapi": "3.0.0",
  "info": {"title": "Things", "version": "1"},
  "paths": {
    "/v1/things": {
      "get": {
        "summary": "Nice 👍",
        "responses": {"200": {"description": "ok"}}
      }
    }
  }
}`

	parsed, err := Parse([]byte(doc))
	if err != nil {
		t.Fatalf("a description with an astral escape was refused: %v", err)
	}

	if _, ok := parsed.Paths["/v1/things"]; !ok {
		t.Errorf("the path is missing; got %v", keysOfPaths(parsed))
	}
}

// A file that is neither JSON nor YAML still reports as unreadable rather than
// being quietly accepted.
func TestSomethingThatIsNeitherJSONNorYAMLIsStillRefused(t *testing.T) {
	_, err := Parse([]byte("\x00\x01 this: is: not: anything\x00"))
	if err == nil {
		t.Fatal("a file that is neither JSON nor YAML was accepted")
	}

	if strings.Contains(err.Error(), "panic") {
		t.Errorf("unreadable input produced a panic-shaped error: %v", err)
	}
}

func keysOfPaths(d *Document) []string {
	out := make([]string, 0, len(d.Paths))
	for k := range d.Paths {
		out = append(out, k)
	}

	return out
}
