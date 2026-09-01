package openapi

import "testing"

// A schema's required is a list of property names. Spotify's published
// description puts a boolean there instead:
//
//	schema:
//	  type: string
//	  format: byte
//	  required: true
//
// which is what required means on a *parameter*, one level up. The document is
// wrong and the wrongness is trivial -- it is one field of one endpoint, and
// nothing anywhere depends on reading it.
//
// Refusing the whole file over it was the worse answer. That is 290KB and every
// path Spotify publishes, unreadable in its entirety because of a boolean in a
// place a boolean does not belong, and the practical effect was that drift could
// not watch Spotify at all.
//
// So a boolean here is read as no required list, which is what it conveys: a
// schema saying required: true names no properties as required.
func TestABooleanWhereRequiredNamesPropertiesIsIgnored(t *testing.T) {
	doc, err := Parse([]byte(`
openapi: 3.0.0
info: {title: Example, version: "1"}
paths:
  /upload:
    put:
      requestBody:
        content:
          image/jpeg:
            schema:
              type: string
              format: byte
              required: true
      responses:
        "202": {description: uploaded}
`))
	if err != nil {
		t.Fatalf("a boolean required made the whole description unreadable: %v", err)
	}

	body := doc.Paths["/upload"].Put.RequestBody
	if body == nil {
		t.Fatal("the request body was lost")
	}

	schema := body.Content["image/jpeg"].Schema
	if schema == nil {
		t.Fatal("the schema was lost")
	}

	if got := string(schema.Type); got != "string" {
		t.Errorf("type = %q, want string -- the rest of the schema should survive", got)
	}

	if len(schema.Required) != 0 {
		t.Errorf("required = %v, and a boolean names no properties", schema.Required)
	}
}

// The ordinary form still works, because that is the one that carries meaning.
func TestRequiredStillReadsAListOfNames(t *testing.T) {
	doc, err := Parse([]byte(`
openapi: 3.0.0
info: {title: Example, version: "1"}
paths:
  /things:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                required: [id, name]
                properties:
                  id: {type: string}
                  name: {type: string}
`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	schema := doc.Paths["/things"].Get.Responses["200"].Content["application/json"].Schema
	if got, want := len(schema.Required), 2; got != want {
		t.Fatalf("required has %d names, want %d", got, want)
	}

	if schema.Required[0] != "id" || schema.Required[1] != "name" {
		t.Errorf("required = %v, want [id name]", schema.Required)
	}
}
