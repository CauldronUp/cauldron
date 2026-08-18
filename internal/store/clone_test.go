package store

import (
	"testing"
)

// Clone promises callers cannot mutate stored state by accident, and for a
// year it copied one level deep, so every nested map and slice a record held
// was the same object the store held.
//
// The consequence is worse than aliasing. A read handed a caller a live
// pointer into the store, so a handler shaping a response and a snapshot being
// encoded at the same time were reading and writing one map, which is a data
// race across a package boundary rather than a tidiness problem.
func TestCloneCopiesNestedValuesTo(t *testing.T) {
	s := New(1)
	s.Declare("issue", "", 0)

	original := Record{
		"id": "1",
		"fields": map[string]any{
			"summary": "original",
			"labels":  []any{"one", "two"},
			"status": map[string]any{
				"statusCategory": map[string]any{"key": "new"},
			},
		},
		"links": []any{
			map[string]any{"rel": "self"},
		},
	}

	if _, err := s.Create("issue", original); err != nil {
		t.Fatalf("create: %v", err)
	}

	held, err := s.Get("issue", "1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	// Everything a caller could reasonably reach through.
	fields, _ := held["fields"].(map[string]any)
	fields["summary"] = "mutated"
	fields["labels"].([]any)[0] = "mutated"
	fields["status"].(map[string]any)["statusCategory"].(map[string]any)["key"] = "mutated"
	held["links"].([]any)[0].(map[string]any)["rel"] = "mutated"

	fresh, err := s.Get("issue", "1")
	if err != nil {
		t.Fatalf("second get: %v", err)
	}

	freshFields, _ := fresh["fields"].(map[string]any)

	if freshFields["summary"] != "original" {
		t.Errorf("a nested map was mutated through a returned record: %v", freshFields["summary"])
	}

	if freshFields["labels"].([]any)[0] != "one" {
		t.Errorf("a nested slice was mutated through a returned record: %v", freshFields["labels"])
	}

	category := freshFields["status"].(map[string]any)["statusCategory"].(map[string]any)
	if category["key"] != "new" {
		t.Errorf("a twice-nested map was mutated through a returned record: %v", category)
	}

	if fresh["links"].([]any)[0].(map[string]any)["rel"] != "self" {
		t.Errorf("a map inside a slice was mutated through a returned record: %v", fresh["links"])
	}

	// And the record the caller handed in is not a way back in either. Create
	// stores a clone, so writing to the original afterwards changes nothing.
	original["fields"].(map[string]any)["summary"] = "mutated by the caller"

	after, _ := s.Get("issue", "1")
	if after["fields"].(map[string]any)["summary"] != "original" {
		t.Error("the record passed to Create is still connected to the stored one")
	}
}
