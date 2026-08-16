package store

import (
	"errors"
	"strings"
	"testing"
)

func newCustomerStore(t *testing.T, seed int64) *Store {
	t.Helper()

	s := New(seed)
	s.Declare("customer", "cus_", 14)

	return s
}

func TestCreateMintsAnIdentifier(t *testing.T) {
	s := newCustomerStore(t, 1)

	record, err := s.Create("customer", Record{"email": "ada@example.com"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	id, _ := record["id"].(string)

	if !strings.HasPrefix(id, "cus_") {
		t.Errorf("id = %q, want a cus_ prefix", id)
	}

	if len(id) != len("cus_")+14 {
		t.Errorf("id = %q, want %d characters after the prefix", id, 14)
	}
}

func TestCreateKeepsAnExplicitIdentifier(t *testing.T) {
	s := newCustomerStore(t, 1)

	record, _ := s.Create("customer", Record{"id": "cus_pinned", "email": "grace@example.com"})

	if record["id"] != "cus_pinned" {
		t.Errorf("id = %v, want the pinned value — fixtures depend on this", record["id"])
	}
}

func TestIdentifiersAreDeterministicAcrossStores(t *testing.T) {
	var first []string

	for run := 0; run < 5; run++ {
		s := newCustomerStore(t, 42)

		var ids []string

		for i := 0; i < 5; i++ {
			record, _ := s.Create("customer", Record{})
			ids = append(ids, record["id"].(string))
		}

		if first == nil {
			first = ids
			continue
		}

		for i := range ids {
			if ids[i] != first[i] {
				t.Fatalf("run %d produced %s at index %d, first run produced %s", run, ids[i], i, first[i])
			}
		}
	}
}

func TestDifferentSeedsProduceDifferentIdentifiers(t *testing.T) {
	a, _ := newCustomerStore(t, 1).Create("customer", Record{})
	b, _ := newCustomerStore(t, 2).Create("customer", Record{})

	if a["id"] == b["id"] {
		t.Errorf("different seeds produced the same id %v", a["id"])
	}
}

func TestResetRewindsIdentifiers(t *testing.T) {
	s := newCustomerStore(t, 7)

	firstRecord, _ := s.Create("customer", Record{})
	first := firstRecord["id"]

	s.Reset()

	againRecord, _ := s.Create("customer", Record{})

	if againRecord["id"] != first {
		t.Errorf("after reset id = %v, want %v — a reset sandbox must match a fresh one", againRecord["id"], first)
	}

	if s.Count("customer") != 1 {
		t.Errorf("Count = %d, want 1", s.Count("customer"))
	}
}

func TestUnknownResourceIsRejected(t *testing.T) {
	s := newCustomerStore(t, 1)

	if _, err := s.Create("invoice", Record{}); !errors.Is(err, ErrUnknownResource) {
		t.Errorf("err = %v, want ErrUnknownResource — a typo must not invent an endpoint", err)
	}
}

func TestGetReturnsACopy(t *testing.T) {
	s := newCustomerStore(t, 1)
	created, _ := s.Create("customer", Record{"email": "ada@example.com"})

	fetched, err := s.Get("customer", created["id"].(string))
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	fetched["email"] = "tampered@example.com"

	again, _ := s.Get("customer", created["id"].(string))

	if again["email"] != "ada@example.com" {
		t.Errorf("stored state was mutated through a returned record: %v", again["email"])
	}
}

func TestGetMissingRecord(t *testing.T) {
	s := newCustomerStore(t, 1)

	if _, err := s.Get("customer", "cus_nope"); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestUpdateMergesAndProtectsTheIdentifier(t *testing.T) {
	s := newCustomerStore(t, 1)
	created, _ := s.Create("customer", Record{"email": "ada@example.com", "name": "Ada"})
	id := created["id"].(string)

	updated, err := s.Update("customer", id, Record{"name": "Ada Lovelace", "id": "cus_hijack"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if updated["id"] != id {
		t.Errorf("identifiers must be immutable; got %v", updated["id"])
	}

	if updated["name"] != "Ada Lovelace" {
		t.Errorf("name = %v", updated["name"])
	}

	if updated["email"] != "ada@example.com" {
		t.Errorf("unspecified fields must survive an update; got %v", updated["email"])
	}
}

func TestDelete(t *testing.T) {
	s := newCustomerStore(t, 1)
	created, _ := s.Create("customer", Record{})
	id := created["id"].(string)

	if err := s.Delete("customer", id); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if s.Count("customer") != 0 {
		t.Errorf("Count = %d, want 0", s.Count("customer"))
	}

	if err := s.Delete("customer", id); !errors.Is(err, ErrNotFound) {
		t.Errorf("second delete err = %v, want ErrNotFound", err)
	}
}

func TestListPreservesInsertionOrder(t *testing.T) {
	s := newCustomerStore(t, 1)

	var ids []string

	for i := 0; i < 5; i++ {
		record, _ := s.Create("customer", Record{"n": i})
		ids = append(ids, record["id"].(string))
	}

	page, err := s.List("customer", "", 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(page.Records) != 5 {
		t.Fatalf("got %d records, want 5", len(page.Records))
	}

	for i, record := range page.Records {
		if record["id"] != ids[i] {
			t.Errorf("index %d = %v, want %v", i, record["id"], ids[i])
		}
	}
}

func TestListPagesWithACursor(t *testing.T) {
	s := newCustomerStore(t, 1)

	for i := 0; i < 7; i++ {
		_, _ = s.Create("customer", Record{"n": i})
	}

	first, err := s.List("customer", "", 3)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(first.Records) != 3 || !first.HasMore {
		t.Fatalf("first page: %d records, HasMore=%v", len(first.Records), first.HasMore)
	}

	second, err := s.List("customer", first.NextCursor, 3)
	if err != nil {
		t.Fatalf("List page 2: %v", err)
	}

	if len(second.Records) != 3 || !second.HasMore {
		t.Fatalf("second page: %d records, HasMore=%v", len(second.Records), second.HasMore)
	}

	third, err := s.List("customer", second.NextCursor, 3)
	if err != nil {
		t.Fatalf("List page 3: %v", err)
	}

	if len(third.Records) != 1 {
		t.Errorf("third page: %d records, want 1", len(third.Records))
	}

	if third.HasMore {
		t.Error("the final page must not claim there is more")
	}

	// No record may appear on two pages.
	seen := map[any]bool{}

	for _, page := range []Page{first, second, third} {
		for _, record := range page.Records {
			if seen[record["id"]] {
				t.Errorf("record %v appeared on more than one page", record["id"])
			}

			seen[record["id"]] = true
		}
	}
}

func TestListRejectsAnUnknownCursor(t *testing.T) {
	s := newCustomerStore(t, 1)
	_, _ = s.Create("customer", Record{})

	if _, err := s.List("customer", "cus_nonsense", 10); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestDeletedRecordsLeaveTheOrdering(t *testing.T) {
	s := newCustomerStore(t, 1)

	a, _ := s.Create("customer", Record{"n": 1})
	b, _ := s.Create("customer", Record{"n": 2})
	c, _ := s.Create("customer", Record{"n": 3})

	if err := s.Delete("customer", b["id"].(string)); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	page, _ := s.List("customer", "", 10)

	if len(page.Records) != 2 {
		t.Fatalf("got %d records, want 2", len(page.Records))
	}

	if page.Records[0]["id"] != a["id"] || page.Records[1]["id"] != c["id"] {
		t.Errorf("ordering broke after a delete: %v", page.Records)
	}
}

func TestExportAndImportRoundTrip(t *testing.T) {
	s := newCustomerStore(t, 3)

	var ids []string

	for i := 0; i < 4; i++ {
		record, _ := s.Create("customer", Record{"n": i})
		ids = append(ids, record["id"].(string))
	}

	snapshot := s.Export()

	restored := newCustomerStore(t, 3)
	if err := restored.Import(snapshot); err != nil {
		t.Fatalf("Import: %v", err)
	}

	page, err := restored.List("customer", "", 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(page.Records) != 4 {
		t.Fatalf("restored %d records, want 4", len(page.Records))
	}

	for i, record := range page.Records {
		if record["id"] != ids[i] {
			t.Errorf("order changed at %d: %v, want %v", i, record["id"], ids[i])
		}
	}
}

// A restored sandbox must not re-mint identifiers it has already handed out.
func TestImportContinuesIdentifiersRatherThanRepeatingThem(t *testing.T) {
	s := newCustomerStore(t, 3)

	for i := 0; i < 4; i++ {
		_, _ = s.Create("customer", Record{})
	}

	next, _ := s.Create("customer", Record{})
	expected := next["id"]

	// Re-take the snapshot from before that fifth record.
	original := newCustomerStore(t, 3)
	for i := 0; i < 4; i++ {
		_, _ = original.Create("customer", Record{})
	}

	restored := newCustomerStore(t, 3)
	if err := restored.Import(original.Export()); err != nil {
		t.Fatalf("Import: %v", err)
	}

	after, _ := restored.Create("customer", Record{})

	if after["id"] != expected {
		t.Errorf("after restore the next id was %v, want %v", after["id"], expected)
	}
}

func TestImportRefusesAForeignSnapshot(t *testing.T) {
	s := newCustomerStore(t, 1)

	err := s.Import(Export{Resources: map[string][]Record{"invoice": {{"id": "in_1"}}}})

	if !errors.Is(err, ErrUnknownResource) {
		t.Errorf("err = %v, want ErrUnknownResource; a snapshot from another recipe must not restore", err)
	}
}

func TestImportClearsRecordsTheSnapshotDoesNotHave(t *testing.T) {
	s := newCustomerStore(t, 1)
	_, _ = s.Create("customer", Record{"n": 1})

	if err := s.Import(Export{Resources: map[string][]Record{}, Drawn: map[string]int{}}); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if s.Count("customer") != 0 {
		t.Errorf("restoring an empty snapshot left %d records", s.Count("customer"))
	}
}
