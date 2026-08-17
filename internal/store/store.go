// Package store holds the state of a sandbox: the records each Recipe has
// created, in the order they were created.
//
// Two properties matter more than anything else here.
//
// Determinism: identifiers come from a seeded generator, never from crypto/rand
// or the clock. The same fixture and seed produce byte-identical IDs on every
// machine, which is what lets a CI failure be reproduced on a laptop.
//
// Ordering: providers return lists in a stable order and applications page
// through them. A map alone cannot do that, so insertion order is tracked
// explicitly.
package store

import (
	"errors"
	"fmt"
	"sync"
)

// Record is a single object — a customer, a payment intent, an order.
type Record map[string]any

// Clone returns a copy so callers cannot mutate stored state by accident.
func (r Record) Clone() Record {
	out := make(Record, len(r))

	for key, value := range r {
		out[key] = value
	}

	return out
}

// ErrNotFound is returned when an identifier does not exist.
var ErrNotFound = errors.New("not found")

// ErrUnknownResource is returned when a resource type was never declared by the
// Recipe. Creating one on the fly would let a typo invent an endpoint.
var ErrUnknownResource = errors.New("unknown resource")

// collection is one resource type's records plus their insertion order.
type collection struct {
	order   []string
	records map[string]Record
}

// Store is the mutable state of one Recipe instance.
type Store struct {
	mu sync.RWMutex

	collections map[string]*collection
	ids         *Generator
}

// New returns an empty store whose identifiers derive from seed.
func New(seed int64) *Store {
	return &Store{
		collections: map[string]*collection{},
		ids:         NewGenerator(seed),
	}
}

// Declare registers a resource type and how its identifiers are shaped. It must
// be called before records of that type can be created.
func (s *Store) Declare(resource, prefix string, length int) {
	s.DeclareStyle(resource, "prefixed", prefix, length)
}

// DeclareStyle registers a resource type with an explicit identifier style.
// Clock gives the identifier generator the sandbox clock, for styles that mint
// time-based identifiers.
func (s *Store) Clock(now func() int64) {
	s.ids.Now(now)
}

func (s *Store) DeclareStyle(resource, style, prefix string, length int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.collections[resource]; !ok {
		s.collections[resource] = &collection{records: map[string]Record{}}
	}

	s.ids.DeclareStyle(resource, style, prefix, length)
}

// Declared reports whether a resource type is known.
func (s *Store) Declared(resource string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.collections[resource]

	return ok
}

// Create stores a new record, minting an identifier unless the record already
// carries one — fixtures pin their IDs so that seeded data is quotable in docs
// and tests.
func (s *Store) Create(resource string, record Record) (Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.collections[resource]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownResource, resource)
	}

	stored := record.Clone()

	id, _ := stored["id"].(string)
	if id == "" {
		id = s.ids.Next(resource)
		stored["id"] = id
	}

	if _, exists := c.records[id]; !exists {
		c.order = append(c.order, id)
	}

	c.records[id] = stored

	return stored.Clone(), nil
}

// Get returns one record by identifier.
func (s *Store) Get(resource, id string) (Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.collections[resource]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownResource, resource)
	}

	record, ok := c.records[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s %s", ErrNotFound, resource, id)
	}

	return record.Clone(), nil
}

// Update merges changes into an existing record. Identifiers are immutable.
func (s *Store) Update(resource, id string, changes Record) (Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.collections[resource]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownResource, resource)
	}

	record, ok := c.records[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s %s", ErrNotFound, resource, id)
	}

	updated := record.Clone()

	for key, value := range changes {
		if key == "id" {
			continue
		}

		updated[key] = value
	}

	c.records[id] = updated

	return updated.Clone(), nil
}

// Delete removes a record.
func (s *Store) Delete(resource, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.collections[resource]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownResource, resource)
	}

	if _, ok := c.records[id]; !ok {
		return fmt.Errorf("%w: %s %s", ErrNotFound, resource, id)
	}

	delete(c.records, id)

	for i, candidate := range c.order {
		if candidate == id {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}

	return nil
}

// Page is a slice of a collection, in insertion order.
type Page struct {
	Records []Record
	// HasMore reports whether further records exist after this page.
	HasMore bool
	// NextCursor is the identifier to start the following page after.
	NextCursor string
}

// List returns up to limit records, starting after the given cursor. An empty
// cursor starts at the beginning.
func (s *Store) List(resource, after string, limit int) (Page, error) {
	return s.ListWhere(resource, nil, after, limit)
}

// ListWhere returns records whose fields match every entry in where.
//
// This is what makes a scoped route work: /repos/{owner}/{repo}/issues must
// only ever see that repository's issues, and paging has to happen after the
// filter or a page could come back empty while claiming there is more.
func (s *Store) ListWhere(resource string, where map[string]any, after string, limit int) (Page, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.collections[resource]
	if !ok {
		return Page{}, fmt.Errorf("%w: %s", ErrUnknownResource, resource)
	}

	if limit <= 0 {
		limit = 10
	}

	// Filter first, then page, so a cursor always refers to a position in the
	// filtered view rather than the whole collection.
	matching := make([]string, 0, len(c.order))

	for _, id := range c.order {
		if Matches(c.records[id], where) {
			matching = append(matching, id)
		}
	}

	start := 0

	if after != "" {
		found := false

		for i, id := range matching {
			if id == after {
				start = i + 1
				found = true

				break
			}
		}

		if !found {
			return Page{}, fmt.Errorf("%w: cursor %s", ErrNotFound, after)
		}
	}

	page := Page{Records: []Record{}}

	for i := start; i < len(matching) && len(page.Records) < limit; i++ {
		page.Records = append(page.Records, c.records[matching[i]].Clone())
	}

	end := start + len(page.Records)

	if end < len(matching) {
		page.HasMore = true

		if len(page.Records) > 0 {
			page.NextCursor = matching[end-1]
		}
	}

	return page, nil
}

// Matches reports whether a record satisfies every filter.
//
// Comparison is on the string form, because a path parameter arrives as text
// while the stored value may be a number. A scoped route would otherwise
// silently return nothing whenever the scope field happened to be numeric.
func Matches(record Record, where map[string]any) bool {
	for field, want := range where {
		if fmt.Sprint(record[field]) != fmt.Sprint(want) {
			return false
		}
	}

	return true
}

// Count returns how many records a resource holds.
func (s *Store) Count(resource string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.collections[resource]
	if !ok {
		return 0
	}

	return len(c.order)
}

// Resources returns the declared resource types.
func (s *Store) Resources() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]string, 0, len(s.collections))

	for name := range s.collections {
		out = append(out, name)
	}

	return out
}

// Reset empties every collection and rewinds the identifier generator, so a
// reset sandbox mints exactly the same IDs as a fresh one.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, c := range s.collections {
		c.order = nil
		c.records = map[string]Record{}
	}

	s.ids.Reset()
}

// Export is a serialisable copy of a store's state.
//
// Records are kept as an ordered list per resource rather than a map, because
// providers return lists in a stable order and a restored sandbox that pages
// differently from the original is not the same sandbox.
type Export struct {
	Resources map[string][]Record `json:"resources"`
	Order     map[string][]string `json:"order"`
	// Drawn counts identifiers minted per resource. The random source cannot
	// be serialised, so restoring replays that many draws to put the generator
	// back where it was. Without this, a restored sandbox would re-mint
	// identifiers it has already used.
	Drawn map[string]int `json:"drawn"`
}

// Export captures the store's state.
func (s *Store) Export() Export {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := Export{
		Resources: map[string][]Record{},
		Order:     map[string][]string{},
		Drawn:     s.ids.Drawn(),
	}

	for name, c := range s.collections {
		records := make([]Record, 0, len(c.order))

		for _, id := range c.order {
			records = append(records, c.records[id].Clone())
		}

		out.Resources[name] = records
		out.Order[name] = append([]string(nil), c.order...)
	}

	return out
}

// Import replaces the store's state.
//
// Resource types that the Recipe never declared are refused: a snapshot taken
// against a different Recipe would otherwise restore into a sandbox that
// cannot serve it.
func (s *Store) Import(in Export) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name := range in.Resources {
		if _, ok := s.collections[name]; !ok {
			return fmt.Errorf("%w: snapshot contains %s", ErrUnknownResource, name)
		}
	}

	for name, c := range s.collections {
		c.order = nil
		c.records = map[string]Record{}

		records, ok := in.Resources[name]
		if !ok {
			continue
		}

		for _, record := range records {
			id, _ := record["id"].(string)
			if id == "" {
				continue
			}

			c.records[id] = record.Clone()
			c.order = append(c.order, id)
		}

		// Prefer the recorded order when present: it is authoritative, and JSON
		// round-trips of the record list could otherwise drift.
		if order, ok := in.Order[name]; ok && len(order) == len(c.order) {
			c.order = append([]string(nil), order...)
		}
	}

	s.ids.Restore(in.Drawn)

	return nil
}
