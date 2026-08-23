package clock

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// durationPattern matches a single component of a human duration: 30d, 2w, 6mo.
var durationPattern = regexp.MustCompile(`(\d+)\s*(mo|ms|[smhdwy])`)

// ParseDuration understands the durations people actually type at a CLI —
// "30d", "6mo", "1y2w" — as well as everything time.ParseDuration accepts.
//
// Go's own parser stops at hours, which is useless for the cases this exists
// to serve: nobody ages a subscription into dunning in units of hours.
//
// Months and years are approximated as 30 and 365 days. That is deliberate: a
// sandbox needs a predictable offset, not a calendar.
func ParseDuration(input string) (time.Duration, error) {
	trimmed := strings.ToLower(strings.TrimSpace(input))

	if trimmed == "" {
		return 0, fmt.Errorf("empty duration")
	}

	// Anything Go already understands wins, so "90m" and "1h30m" behave
	// exactly as a Go developer expects.
	if d, err := time.ParseDuration(trimmed); err == nil {
		return d, nil
	}

	matches := durationPattern.FindAllStringSubmatch(trimmed, -1)
	if len(matches) == 0 {
		return 0, fmt.Errorf("cannot understand duration %q. Try 30s, 15m, 12h, 7d, 2w, 6mo or 1y", input)
	}

	// Reject trailing junk: "30dd" or "30d!" should fail loudly rather than
	// silently becoming 30 days.
	consumed := 0
	for _, match := range matches {
		consumed += len(match[0])
	}

	if consumed != len(strings.ReplaceAll(trimmed, " ", "")) {
		return 0, fmt.Errorf("cannot understand duration %q. Try 30s, 15m, 12h, 7d, 2w, 6mo or 1y", input)
	}

	units := map[string]time.Duration{
		"ms": time.Millisecond,
		"s":  time.Second,
		"m":  time.Minute,
		"h":  time.Hour,
		"d":  24 * time.Hour,
		"w":  7 * 24 * time.Hour,
		"mo": 30 * 24 * time.Hour,
		"y":  365 * 24 * time.Hour,
	}

	var total time.Duration

	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("cannot understand duration %q", input)
		}

		unit, known := units[match[2]]
		if !known {
			return 0, fmt.Errorf("unknown unit %q in duration %q", match[2], input)
		}

		total, err = add(total, value, unit)
		if err != nil {
			return 0, fmt.Errorf("duration %q is longer than this can represent, which is about 292 years", input)
		}
	}

	return total, nil
}

// add accumulates one component, refusing to wrap.
//
// A duration is an int64 of nanoseconds, so it runs out at about 292 years and
// "1000y" -- the natural way to say "for the rest of this suite" -- came back
// negative. Nothing noticed: `fault --for 1000y` computed an expiry in the
// past, the fault was culled on the very first request, and the control plane
// had already answered {"armed":"rate_limit"}. The resilience test that was
// supposed to exercise a 429 ran green against a provider that never failed.
//
// A false green is the one outcome this tool exists to prevent, so an
// unrepresentable duration is refused rather than wrapped.
func add(total time.Duration, value int, unit time.Duration) (time.Duration, error) {
	if value < 0 {
		return 0, fmt.Errorf("negative component")
	}

	if value != 0 && int64(value) > int64(math.MaxInt64-total)/int64(unit) {
		return 0, fmt.Errorf("overflow")
	}

	return total + time.Duration(value)*unit, nil
}
