package openapi

import "testing"

// A $ref where a string belongs does not make a document unreadable.
//
// description is specified as a string, and Livepeer's own published document
// puts a $ref there instead -- pointing at another schema's description so the
// wording is written once. That is a reasonable thing to want and it is not
// what the specification says, and a 156KB file was refused outright because
// two fields did it.
//
// This package does nothing with a description. It fingerprints paths, methods,
// response codes and the types of fields a Recipe names, so a description it
// cannot read costs nothing to discard -- and discarding it is much better than
// discarding the document around it.
func TestARefWhereADescriptionBelongsIsSkippedRatherThanFatal(t *testing.T) {
	const document = `
openapi: 3.0.0
info:
  title: Borrowed prose
  version: "1"
components:
  schemas:
    Original:
      type: object
      description: the wording, written once
      properties:
        name:
          type: string
    Borrower:
      type: object
      description:
        $ref: "#/components/schemas/Original/description"
      properties:
        name:
          type: string
paths:
  /things:
    get:
      responses:
        "200":
          description: ok
`

	doc, err := Parse([]byte(document))
	if err != nil {
		t.Fatalf("a description carrying a $ref made the document unreadable: %v", err)
	}

	borrower := doc.Components.Schemas["Borrower"]
	if borrower == nil {
		t.Fatal("the schema with the borrowed description was dropped entirely")
	}

	// The prose is gone and the shape survives, which is the trade.
	if borrower.Description != "" {
		t.Errorf("description = %q, want empty -- a $ref this package cannot expand yields nothing", borrower.Description)
	}

	if borrower.Properties["name"] == nil {
		t.Error("the properties beside the unreadable description were lost with it")
	}
}

// And a $ref where a properties map belongs is skipped the same way.
//
// Livepeer does this too, in three places: a $ref directly under properties,
// pointing at another schema's whole property set so it is not repeated. An
// unexpandable reference yields no fields, which is what an unreadable
// properties block means in any case.
func TestARefWherePropertiesBelongIsSkippedRatherThanFatal(t *testing.T) {
	const document = `
openapi: 3.0.0
info:
  title: Borrowed shape
  version: "1"
components:
  schemas:
    Shared:
      type: object
      properties:
        cid:
          type: string
    Borrower:
      type: object
      properties:
        $ref: "#/components/schemas/Shared/properties"
paths:
  /things:
    get:
      responses:
        "200":
          description: ok
`

	doc, err := Parse([]byte(document))
	if err != nil {
		t.Fatalf("properties carrying a $ref made the document unreadable: %v", err)
	}

	if doc.Components.Schemas["Borrower"] == nil {
		t.Fatal("the schema with the borrowed properties was dropped entirely")
	}

	// The rest of the document still parses, which is the whole point.
	if doc.Components.Schemas["Shared"].Properties["cid"] == nil {
		t.Error("the schema the reference pointed at was lost too")
	}
}

// An ordinary description and an ordinary properties map are unaffected.
//
// The risk in loosening a decoder is loosening it into silence. Both fields
// still carry what they carry when they are written the way the specification
// says.
func TestAnOrdinaryDescriptionAndPropertiesStillDecode(t *testing.T) {
	const document = `
openapi: 3.0.0
info:
  title: Ordinary
  version: "1"
components:
  schemas:
    Thing:
      type: object
      description: a plain sentence
      properties:
        id:
          type: string
paths:
  /things:
    get:
      responses:
        "200":
          description: ok
`

	doc, err := Parse([]byte(document))
	if err != nil {
		t.Fatalf("an ordinary document was refused: %v", err)
	}

	thing := doc.Components.Schemas["Thing"]

	if thing.Description != "a plain sentence" {
		t.Errorf("description = %q, want the sentence it was given", thing.Description)
	}

	if thing.Properties["id"] == nil {
		t.Error("an ordinary properties map came back empty")
	}
}
