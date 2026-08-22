package openapi

import "testing"

// Chargebee's listing answers {"list": [{"subscription": {...}}], ...}, so the
// collection's items are objects holding a key rather than the subscription
// itself. Reading an item as the resource reported every field the Recipe
// declares as one the description does not: eleven findings with the shape of
// a Recipe that invented its whole model, against a Recipe that had the
// nesting exactly right and said so with entry_style.
const wrappedEntries = `
openapi: 3.0.0
info: {title: Wrapped, version: "1"}
paths:
  /subscriptions:
    get:
      responses:
        "200":
          content:
            application/json:
              schema:
                type: object
                properties:
                  list:
                    type: array
                    items:
                      type: object
                      properties:
                        subscription:
                          $ref: "#/components/schemas/Subscription"
                  next_offset: {type: string}
components:
  schemas:
    Subscription:
      type: object
      properties:
        id: {type: string}
        currency_code: {type: string}
        cancelled_at: {type: integer}
`

func TestAWrappedEntryIsReachedThroughItsKey(t *testing.T) {
	doc, err := Parse([]byte(wrappedEntries))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	envelope := doc.Paths["/subscriptions"].Get.Responses["200"].Content["application/json"].Schema

	items := collectionSchema(doc, envelope, "list", "")
	if items == nil {
		t.Fatal("the collection was not found")
	}

	// Without descending again this is the wrapper, whose only property is
	// the resource's name.
	if _, ok := items.Properties["subscription"]; !ok {
		t.Fatalf("expected the wrapper, got %v", items.Properties)
	}

	entry := doc.Resolve(descend(doc, items, "subscription"))
	if entry == nil {
		t.Fatal("the entry was not reached through its key")
	}

	for _, name := range []string{"id", "currency_code", "cancelled_at"} {
		if _, ok := entry.Properties[name]; !ok {
			t.Errorf("%s is declared and was not found on the entry", name)
		}
	}
}
