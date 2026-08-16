package store

import (
	"math/rand"
	"strings"
	"sync"
)

// alphabet is the character set provider identifiers are drawn from. Digits and
// both cases, matching what Stripe, Shopify and most others actually emit —
// applications occasionally validate the shape, so the fakes have to match it.
const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// shape describes how one resource's identifiers look.
type shape struct {
	prefix string
	length int
}

// Generator mints identifiers that look like the provider's own but are fully
// determined by the seed.
//
// Using math/rand here is deliberate and is not a security weakness: these are
// fixture identifiers inside a sandbox, and reproducibility is the entire
// point. Anything unpredictable would make CI failures impossible to replay.
type Generator struct {
	mu     sync.Mutex
	seed   int64
	rng    *rand.Rand
	shapes map[string]shape
}

// NewGenerator returns a generator for the given seed.
func NewGenerator(seed int64) *Generator {
	return &Generator{
		seed:   seed,
		rng:    rand.New(rand.NewSource(seed)),
		shapes: map[string]shape{},
	}
}

// Declare registers the identifier shape for a resource.
func (g *Generator) Declare(resource, prefix string, length int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if length <= 0 {
		length = 14
	}

	g.shapes[resource] = shape{prefix: prefix, length: length}
}

// Next mints the next identifier for a resource.
func (g *Generator) Next(resource string) string {
	g.mu.Lock()
	defer g.mu.Unlock()

	s, ok := g.shapes[resource]
	if !ok {
		s = shape{prefix: resource + "_", length: 14}
	}

	var b strings.Builder

	b.WriteString(s.prefix)

	for i := 0; i < s.length; i++ {
		b.WriteByte(alphabet[g.rng.Intn(len(alphabet))])
	}

	return b.String()
}

// Reset rewinds the generator to its initial state, so a reset sandbox
// reproduces the same identifiers as a fresh one.
func (g *Generator) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.rng = rand.New(rand.NewSource(g.seed))
}
