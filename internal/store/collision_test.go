package store

import (
	"testing"
)

// A numeric identifier counter that starts at zero, and a fixture that pins
// id "1", are on a collision course: the first record the API creates is
// handed the identifier a seeded record already holds.
//
// Create used to treat that as an in-place replacement. The seeded record was
// destroyed, the collection did not grow, and the response was a 201 carrying
// somebody else's identifier. GitHub and Zendesk both allocate the next unused
// number, so this was not a difference in behaviour but a loss of data.
func TestCreatingARecordCannotDestroyASeededOne(t *testing.T) {
	s := New(1)
	s.DeclareStyle("issue", "numeric", "", 0)

	for _, id := range []string{"1", "2"} {
		if _, err := s.Create("issue", Record{"id": id, "title": "seeded " + id}); err != nil {
			t.Fatalf("seed %s: %v", id, err)
		}
	}

	created, err := s.Create("issue", Record{"title": "from the API"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// The minted identifier must skip past what the fixture already holds.
	if created["id"] == "1" || created["id"] == "2" {
		t.Fatalf("the API was handed a seeded identifier: %v", created["id"])
	}

	seeded, err := s.Get("issue", "1")
	if err != nil {
		t.Fatalf("the seeded record is gone: %v", err)
	}

	if seeded["title"] != "seeded 1" {
		t.Errorf("the seeded record was overwritten: %v", seeded["title"])
	}

	page, err := s.ListWhere("issue", nil, "", 100)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(page.Records) != 3 {
		t.Errorf("the collection holds %d records after seeding 2 and creating 1", len(page.Records))
	}
}

// A create naming an identifier that already exists is a conflict, not an
// instruction to replace. Silently overwriting meant a suite could destroy a
// record and be told it had made one.
func TestCreatingOverAnExistingIdentifierIsRefused(t *testing.T) {
	s := New(1)
	s.DeclareStyle("issue", "numeric", "", 0)

	if _, err := s.Create("issue", Record{"id": "abc", "title": "first"}); err != nil {
		t.Fatalf("first create: %v", err)
	}

	if _, err := s.Create("issue", Record{"id": "abc", "title": "second"}); err == nil {
		t.Fatal("a create over an existing identifier was accepted")
	}

	held, err := s.Get("issue", "abc")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if held["title"] != "first" {
		t.Errorf("the original record was replaced: %v", held["title"])
	}
}
