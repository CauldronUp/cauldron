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
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.collections[resource]
	if !ok {
		return Page{}, fmt.Errorf("%w: %s", ErrUnknownResource, resource)
	}

	if limit <= 0 {
		limit = 10
	}

	start := 0

	if after != "" {
		found := false

		for i, id := range c.order {
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

	for i := start; i < len(c.order) && len(page.Records) < limit; i++ {
		page.Records = append(page.Records, c.records[c.order[i]].Clone())
	}

	end := start + len(page.Records)

	if end < len(c.order) {
		page.HasMore = true

		if len(page.Records) > 0 {
			page.NextCursor = c.order[end-1]
		}
	}

	return page, nil
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
