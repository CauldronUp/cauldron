package clock

import (
	"testing"
	"time"
)

// A duration is an int64 of nanoseconds and runs out at about 292 years, and
// the components were accumulated with no check. "1000y" -- the natural way
// to say "for the rest of this suite" -- came back as -1488191h, and nothing
// noticed.
//
// What that produced is the outcome this tool exists to prevent. `fault
// --for 1000y` computed an expiry in the past, the fault was culled on the
// first request that looked at it, and the control plane had already answered
// {"armed":"rate_limit"}. A resilience test written to exercise a 429 ran
// green against a provider that never failed.
func TestADurationTooLongToRepresentIsRefused(t *testing.T) {
	for _, input := range []string{"300y", "1000y", "9999mo", "99999999d", "292y1y"} {
		d, err := ParseDuration(input)
		if err == nil {
			t.Errorf("ParseDuration(%q) = %v, want it refused", input, d)
		}
	}

	// The boundary is where the type ends, not somewhere arbitrary short of
	// it. 292 years is representable and is still accepted.
	for _, c := range []struct {
		input string
		want  time.Duration
	}{
		{"292y", 292 * 365 * 24 * time.Hour},
		{"30d", 30 * 24 * time.Hour},
		{"6mo", 6 * 30 * 24 * time.Hour},
		{"1y2w", 365*24*time.Hour + 2*7*24*time.Hour},
		{"90m", 90 * time.Minute},
	} {
		got, err := ParseDuration(c.input)
		if err != nil {
			t.Errorf("ParseDuration(%q) = %v", c.input, err)
			continue
		}

		if got != c.want {
			t.Errorf("ParseDuration(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

// Whatever it returns, it is never negative -- which is the property the
// callers rely on when they add it to now.
func TestAParsedDurationIsNeverNegative(t *testing.T) {
	for _, input := range []string{"1ms", "292y", "9999999d", "500y", "1y1y1y1y1y"} {
		d, err := ParseDuration(input)
		if err != nil {
			continue
		}

		if d < 0 {
			t.Errorf("ParseDuration(%q) = %v, which is in the past", input, d)
		}
	}
}
