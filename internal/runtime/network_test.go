package runtime

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// An undegraded sandbox must stay fast. If arming nothing still cost
// milliseconds, every existing test would slow down for no reason.
func TestNetworkIsFreeWhenNothingIsArmed(t *testing.T) {
	s := stripe(t)

	started := time.Now()
	rec := call(t, s, http.MethodGet, "/v1/customers", "")
	elapsed := time.Since(started)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}

	if elapsed > 100*time.Millisecond {
		t.Errorf("an undegraded request took %v; nothing should be sleeping", elapsed)
	}
}

func TestLatencyDelaysTheResponse(t *testing.T) {
	s := stripe(t)

	if err := s.Degrade(Conditions{Latency: 150 * time.Millisecond}); err != nil {
		t.Fatalf("degrade: %v", err)
	}

	started := time.Now()
	rec := call(t, s, http.MethodGet, "/v1/customers", "")
	elapsed := time.Since(started)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200; latency must delay a response, not replace it", rec.Code)
	}

	if elapsed < 150*time.Millisecond {
		t.Errorf("responded in %v, want at least 150ms", elapsed)
	}
}

func TestJitterStaysWithinItsBounds(t *testing.T) {
	s := stripe(t)

	conditions := Conditions{Latency: 100 * time.Millisecond, Jitter: 40 * time.Millisecond}

	for i := 0; i < 25; i++ {
		d := s.network.delay(conditions)

		if d < 60*time.Millisecond || d > 140*time.Millisecond {
			t.Fatalf("delay %v is outside 100ms ± 40ms", d)
		}
	}
}

// Jitter larger than the latency would otherwise produce a negative delay, and
// time.Sleep on a negative duration returns instantly rather than erroring, so
// the bug would be silent.
func TestJitterNeverProducesANegativeDelay(t *testing.T) {
	s := stripe(t)

	conditions := Conditions{Latency: 10 * time.Millisecond, Jitter: 500 * time.Millisecond}

	for i := 0; i < 50; i++ {
		if d := s.network.delay(conditions); d < 0 {
			t.Fatalf("got a negative delay: %v", d)
		}
	}
}

func TestBandwidthThrottlesTheBody(t *testing.T) {
	s := stripe(t)

	// 8KB at 32KB/s is a quarter of a second: slow enough to measure, quick
	// enough not to drag the suite.
	if err := s.Degrade(Conditions{Bandwidth: 32}); err != nil {
		t.Fatalf("degrade: %v", err)
	}

	body := strings.Repeat("x", 8*1024)
	rec := httptest.NewRecorder()
	writer := newDegradedWriter(rec, Conditions{Bandwidth: 32, Slice: 1024})

	started := time.Now()

	if _, err := writer.Write([]byte(body)); err != nil {
		t.Fatalf("write: %v", err)
	}

	elapsed := time.Since(started)

	if rec.Body.Len() != len(body) {
		t.Errorf("wrote %d bytes, want %d; throttling must not lose data", rec.Body.Len(), len(body))
	}

	if elapsed < 150*time.Millisecond {
		t.Errorf("8KB at 32KB/s took %v, which is not throttled", elapsed)
	}
}

func TestLimitTruncatesTheBodyAndSaysSo(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := newDegradedWriter(rec, Conditions{Limit: 64})

	written, err := writer.Write([]byte(strings.Repeat("y", 4096)))
	if err != nil {
		t.Fatalf("write: %v", err)
	}

	// The handler is told everything was written, because there is nothing
	// useful it could do about a connection that has already gone.
	if written != 4096 {
		t.Errorf("reported %d bytes written to the handler, want 4096", written)
	}

	if rec.Body.Len() != 64 {
		t.Errorf("delivered %d bytes, want exactly the 64 byte limit", rec.Body.Len())
	}

	if !writer.truncated() {
		t.Error("the writer should report the response as truncated")
	}
}

func TestLimitAcrossSeveralWritesStillStopsAtTheBudget(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := newDegradedWriter(rec, Conditions{Limit: 10})

	for i := 0; i < 5; i++ {
		if _, err := writer.Write([]byte("abcd")); err != nil {
			t.Fatalf("write %d: %v", i, err)
		}
	}

	if rec.Body.Len() != 10 {
		t.Errorf("delivered %d bytes across five writes, want 10", rec.Body.Len())
	}
}

func TestSliceSplitsTheBodyIntoChunks(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := newDegradedWriter(rec, Conditions{Slice: 16})

	body := strings.Repeat("z", 100)

	if _, err := writer.Write([]byte(body)); err != nil {
		t.Fatalf("write: %v", err)
	}

	if rec.Body.String() != body {
		t.Error("slicing changed the bytes; it must only change how they arrive")
	}

	if rec.Flushed != true {
		t.Error("sliced writes must flush, or Go buffers them back into one and the client sees nothing different")
	}
}

func TestConditionsExpire(t *testing.T) {
	s := stripe(t)

	if err := s.Degrade(Conditions{
		Latency: time.Second,
		Until:   s.Clock().Now().Add(30 * time.Second),
	}); err != nil {
		t.Fatalf("degrade: %v", err)
	}

	if _, ok := s.network.next("/v1/customers"); !ok {
		t.Fatal("the conditions should be live before they expire")
	}

	s.Clock().Advance(time.Minute)

	if _, ok := s.network.next("/v1/customers"); ok {
		t.Error("the conditions should have expired")
	}
}

func TestCountLimitsHowManyRequestsAreAffected(t *testing.T) {
	s := stripe(t)

	if err := s.Degrade(Conditions{Latency: time.Millisecond, Count: 2}); err != nil {
		t.Fatalf("degrade: %v", err)
	}

	for i := 1; i <= 2; i++ {
		if _, ok := s.network.next("/v1/customers"); !ok {
			t.Fatalf("request %d should have been affected", i)
		}
	}

	if _, ok := s.network.next("/v1/customers"); ok {
		t.Error("a spent condition should stop firing")
	}
}

func TestPathRestrictsWhichRoutesAreAffected(t *testing.T) {
	s := stripe(t)

	if err := s.Degrade(Conditions{Latency: time.Millisecond, Path: "charges"}); err != nil {
		t.Fatalf("degrade: %v", err)
	}

	if _, ok := s.network.next("/v1/customers"); ok {
		t.Error("a path-scoped condition should leave other paths alone")
	}

	if _, ok := s.network.next("/v1/charges"); !ok {
		t.Error("the matching path should be affected")
	}
}

// Chaos that cannot be replayed is a bug report nobody can act on, so the
// probability draw comes from the sandbox seed rather than the wall clock.
func TestProbabilityIsReproducibleFromTheSeed(t *testing.T) {
	run := func() []bool {
		s := stripe(t)

		if err := s.Degrade(Conditions{Latency: time.Millisecond, Probability: 0.5}); err != nil {
			t.Fatalf("degrade: %v", err)
		}

		var hits []bool
		for i := 0; i < 20; i++ {
			_, ok := s.network.next("/v1/customers")
			hits = append(hits, ok)
		}

		return hits
	}

	first, second := run(), run()

	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("two runs with the same seed disagreed at request %d", i)
		}
	}

	var affected int

	for _, hit := range first {
		if hit {
			affected++
		}
	}

	if affected == 0 || affected == len(first) {
		t.Errorf("a probability of 0.5 affected %d of %d requests, which is not a coin flip", affected, len(first))
	}
}

func TestDegradingNothingIsRejected(t *testing.T) {
	s := stripe(t)

	if err := s.Degrade(Conditions{}); err == nil {
		t.Error("arming empty conditions should be an error, not a silent no-op")
	}
}

func TestAnImpossibleProbabilityIsRejected(t *testing.T) {
	s := stripe(t)

	for _, p := range []float64{-0.5, 1.5} {
		if err := s.Degrade(Conditions{Latency: time.Second, Probability: p}); err == nil {
			t.Errorf("probability %v should be rejected", p)
		}
	}
}

func TestClearRemovesEverything(t *testing.T) {
	s := stripe(t)

	if err := s.Degrade(Conditions{Latency: time.Second}); err != nil {
		t.Fatalf("degrade: %v", err)
	}

	if len(s.ArmedNetwork()) != 1 {
		t.Fatalf("armed %d conditions, want 1", len(s.ArmedNetwork()))
	}

	s.ClearNetwork()

	if len(s.ArmedNetwork()) != 0 {
		t.Error("clearing should leave nothing armed")
	}
}

// A reset or a timeout must produce no response at all. Answering 200 with an
// empty body would be a completely different thing for the client under test,
// and the one an emulator is most tempted to do.
func TestResetEndsTheRequestWithNoResponse(t *testing.T) {
	s := stripe(t)

	if err := s.Degrade(Conditions{Reset: true}); err != nil {
		t.Fatalf("degrade: %v", err)
	}

	server := httptest.NewServer(s)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/v1/customers", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer sk_test_cauldron")

	resp, err := server.Client().Do(req)
	if err == nil {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		t.Fatalf("got HTTP %d with body %q, want a severed connection", resp.StatusCode, body)
	}
}

func TestTimeoutHoldsTheConnectionThenSeversIt(t *testing.T) {
	s := stripe(t)

	if err := s.Degrade(Conditions{Timeout: 120 * time.Millisecond}); err != nil {
		t.Fatalf("degrade: %v", err)
	}

	server := httptest.NewServer(s)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/v1/customers", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer sk_test_cauldron")

	started := time.Now()
	resp, err := server.Client().Do(req)
	elapsed := time.Since(started)

	if err == nil {
		resp.Body.Close()
		t.Fatalf("got HTTP %d, want a connection that hangs then dies", resp.StatusCode)
	}

	if elapsed < 100*time.Millisecond {
		t.Errorf("the connection died after %v; a timeout should hold it open first", elapsed)
	}
}

// The network is consulted before the provider's own failures, because that is
// the order reality uses: a connection that never completes never reaches an
// application that could have rate limited it.
func TestNetworkIsAppliedBeforeFaults(t *testing.T) {
	s := stripe(t)

	if err := s.Arm(Fault{Error: "rate_limit"}); err != nil {
		t.Fatalf("arm fault: %v", err)
	}

	if err := s.Degrade(Conditions{Reset: true}); err != nil {
		t.Fatalf("degrade: %v", err)
	}

	server := httptest.NewServer(s)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/v1/customers", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer sk_test_cauldron")

	if resp, err := server.Client().Do(req); err == nil {
		resp.Body.Close()
		t.Fatalf("got HTTP %d; the severed connection should win over the rate limit", resp.StatusCode)
	}
}

func TestDegradedRequestsAreVisibleInTheLog(t *testing.T) {
	s := stripe(t)

	if err := s.Degrade(Conditions{Latency: 10 * time.Millisecond, Bandwidth: 64}); err != nil {
		t.Fatalf("degrade: %v", err)
	}

	call(t, s, http.MethodGet, "/v1/customers", "")

	entries := s.Exchanges(1)
	if len(entries) != 1 {
		t.Fatalf("logged %d exchanges, want 1", len(entries))
	}

	if entries[0].Network == "" {
		t.Error("a degraded request should say so in the log, or the timing has no visible cause")
	}
}

func TestDescribeReadsLikeSomethingAPersonWrote(t *testing.T) {
	cases := []struct {
		conditions Conditions
		want       string
	}{
		{Conditions{Latency: 800 * time.Millisecond}, "latency 800ms"},
		{Conditions{Latency: time.Second, Jitter: 200 * time.Millisecond}, "latency 1s ±200ms"},
		{Conditions{Bandwidth: 50}, "bandwidth 50KB/s"},
		{Conditions{Reset: true}, "reset"},
		{Conditions{Timeout: 30 * time.Second}, "timeout 30s"},
		{Conditions{Limit: 1024}, "limit 1024B"},
		{Conditions{Slice: 64}, "slice 64B"},
		{Conditions{Reset: true, Probability: 0.25}, "reset @ 25% of requests"},
		{Conditions{Latency: time.Second, Path: "/v1/charges"}, "latency 1s on paths containing /v1/charges"},
		{Conditions{Latency: time.Second, Bandwidth: 10}, "latency 1s, bandwidth 10KB/s"},
	}

	for _, tc := range cases {
		if got := Describe(tc.conditions); got != tc.want {
			t.Errorf("Describe() = %q, want %q", got, tc.want)
		}
	}
}
