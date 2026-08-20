package runtime

import (
	"sync"
	"time"
)

// Exchange is one request the sandbox handled, and what it answered.
//
// The log exists because the most common question during an integration is not
// "what does the API return?" but "what did my code actually send?" — which is
// exactly what a provider's dashboard will not show you for a local request.
type Exchange struct {
	Seq      int
	At       time.Time
	Method   string
	Path     string
	Status   int
	Fault    string
	Resource string
	Op       string
	// Network describes any degraded network conditions this request met, so
	// a baffling timing in the log has a visible cause.
	Network string
}

// requestLog is a bounded ring of recent exchanges.
type requestLog struct {
	mu   sync.Mutex
	all  []Exchange
	seq  int
	size int
}

func newRequestLog(size int) *requestLog {
	if size <= 0 {
		size = 500
	}

	return &requestLog{size: size}
}

func (l *requestLog) record(e Exchange) Exchange {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.seq++
	e.Seq = l.seq

	l.all = append(l.all, e)

	if len(l.all) > l.size {
		l.all = l.all[len(l.all)-l.size:]
	}

	return e
}

// Entries returns the most recent exchanges, oldest first.
func (l *requestLog) Entries(limit int) []Exchange {
	l.mu.Lock()
	defer l.mu.Unlock()

	if limit <= 0 || limit > len(l.all) {
		limit = len(l.all)
	}

	out := make([]Exchange, limit)
	copy(out, l.all[len(l.all)-limit:])

	return out
}

func (l *requestLog) reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.all = nil
	l.seq = 0
}
