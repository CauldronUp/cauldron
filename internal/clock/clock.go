// Package clock provides the only source of time inside a Cauldron sandbox.
//
// Recipes never read the wall clock. Determinism is enforced here, at the
// boundary, rather than requested politely in documentation: if a Recipe could
// call time.Now() then the same fixture would produce different data on
// different machines, and a test that passes on a laptop would fail in CI on a
// Tuesday.
//
// The clock is also movable, which is what makes `cauldron clock advance 30d`
// possible — ageing a subscription into dunning without waiting a month.
package clock

import (
	"sync"
	"time"
)

// Clock is a movable, deterministic source of time.
type Clock struct {
	mu  sync.RWMutex
	now time.Time
}

// Epoch is the default start time for a sandbox: 2026-01-01T00:00:00Z.
//
// A fixed, round epoch means fixture data reads the same in every environment
// and every screenshot, and makes offsets easy to reason about by eye.
var Epoch = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

// New returns a clock started at the epoch.
func New() *Clock {
	return &Clock{now: Epoch}
}

// At returns a clock started at a specific instant.
func At(t time.Time) *Clock {
	return &Clock{now: t.UTC()}
}

// Now returns the sandbox's current time.
func (c *Clock) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.now
}

// Unix returns the sandbox time as a Unix timestamp, which is the form most
// provider APIs actually put on the wire.
func (c *Clock) Unix() int64 {
	return c.Now().Unix()
}

// Advance moves the clock forward. Negative durations are rejected: time
// running backwards inside a sandbox would produce data that no real provider
// could ever emit, which is a trap rather than a feature.
func (c *Clock) Advance(d time.Duration) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	if d > 0 {
		c.now = c.now.Add(d)
	}

	return c.now
}

// Set moves the clock to an absolute instant.
func (c *Clock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = t.UTC()
}

// Reset returns the clock to the epoch.
func (c *Clock) Reset() {
	c.Set(Epoch)
}
