package runtime

import (
	"sync"
	"time"

	"github.com/CauldronUp/cauldron/internal/clock"
)

// Fault is an armed failure mode. It is what `cauldron fault stripe --error
// rate_limit --for 30s` installs.
//
// Failure injection is the reason this project exists as much as the happy
// path is: a provider's own sandbox lets you make things succeed, and the paths
// that actually page people at 3am are the ones you cannot stage there.
type Fault struct {
	// Error names an entry in the Recipe's errors table.
	Error string
	// Until is the sandbox time the fault expires. Zero means it never does.
	Until time.Time
	// Count limits how many requests the fault affects. Zero means unlimited.
	Count int
	// Every makes the fault intermittent: fail one request in N. Zero or one
	// means every request fails.
	Every int
	// Path, when set, restricts the fault to routes whose path contains it.
	Path string

	seen int
}

// faultSet holds the faults armed against one sandbox.
type faultSet struct {
	mu     sync.Mutex
	faults []*Fault
	clock  *clock.Clock
}

func newFaultSet(c *clock.Clock) *faultSet {
	return &faultSet{clock: c}
}

// Arm installs a fault.
func (f *faultSet) Arm(fault Fault) {
	f.mu.Lock()
	defer f.mu.Unlock()

	copied := fault
	f.faults = append(f.faults, &copied)
}

// Clear removes every armed fault.
func (f *faultSet) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.faults = nil
}

// Armed returns the currently armed faults.
func (f *faultSet) Armed() []Fault {
	f.mu.Lock()
	defer f.mu.Unlock()

	out := make([]Fault, 0, len(f.faults))

	for _, fault := range f.faults {
		out = append(out, *fault)
	}

	return out
}

// next returns the error name to inject for a request, if any.
func (f *faultSet) next(path string) (string, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := f.clock.Now()

	// Drop expired faults first so an expired fault never fires and never
	// consumes a slot from the ones behind it.
	live := f.faults[:0]

	for _, fault := range f.faults {
		if !fault.Until.IsZero() && !now.Before(fault.Until) {
			continue
		}

		if fault.Count > 0 && fault.seen >= fault.Count {
			continue
		}

		live = append(live, fault)
	}

	f.faults = live

	for _, fault := range f.faults {
		if fault.Path != "" && !contains(path, fault.Path) {
			continue
		}

		fault.seen++

		if fault.Every > 1 && fault.seen%fault.Every != 0 {
			continue
		}

		return fault.Error, true
	}

	return "", false
}

func contains(haystack, needle string) bool {
	if needle == "" {
		return true
	}

	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}

	return false
}
