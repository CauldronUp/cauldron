package openapi

import (
	"strings"
	"testing"
)

// A description written as JSON is read as YAML, because YAML 1.2 is a
// superset of JSON and one parser is simpler than two. yaml.v3 implements
// YAML 1.1, which is not: it rejects the surrogate pairs JSON uses for
// anything outside the basic plane.
//
// One emoji in an example is enough. Twilio's description carries a thumbs
// up in a message body, and the whole file -- a hundred and twenty-one
// paths -- could not be read at all. The failure arrives as a parse error
// halfway through a 1.8MB file, which reads like a corrupt download rather
// than an emoji.

const jsonWithSurrogatePair = `{
  "openapi": "3.0.0",
  "info": {"title": "Emoji", "version": "1"},
  "paths": {
    "/v1/messages": {
      "post": {
        "responses": {
          "201": {
            "description": "sent",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "body": {"type": "string", "example": "Hello! \ud83d\udc4d"}
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`

func TestADescriptionCarryingAnEmojiEscapeIsRead(t *testing.T) {
	doc, err := Parse([]byte(jsonWithSurrogatePair))
	if err != nil {
		t.Fatalf("a valid JSON description was refused: %v", err)
	}

	if len(doc.Paths) != 1 {
		t.Fatalf("read %d paths, want 1", len(doc.Paths))
	}

	if _, ok := doc.Paths["/v1/messages"]; !ok {
		t.Errorf("the path did not survive the round trip: %v", doc.Paths)
	}
}

func TestBrokenInputStillSaysWhatIsWrong(t *testing.T) {
	// The fallback must not turn every failure into a JSON error. A file that
	// is neither is still reported as unreadable.
	_, err := Parse([]byte("this: [is: not: valid"))
	if err == nil {
		t.Fatal("expected a parse failure")
	}

	if !strings.Contains(err.Error(), "not a readable OpenAPI description") {
		t.Errorf("got %q", err)
	}
}
