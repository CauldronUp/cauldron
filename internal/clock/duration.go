package clock

import (
	"fmt"
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

	var total time.Duration

	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("cannot understand duration %q", input)
		}

		unit := match[2]

		switch unit {
		case "ms":
			total += time.Duration(value) * time.Millisecond
		case "s":
			total += time.Duration(value) * time.Second
		case "m":
			total += time.Duration(value) * time.Minute
		case "h":
			total += time.Duration(value) * time.Hour
		case "d":
			total += time.Duration(value) * 24 * time.Hour
		case "w":
			total += time.Duration(value) * 7 * 24 * time.Hour
		case "mo":
			total += time.Duration(value) * 30 * 24 * time.Hour
		case "y":
			total += time.Duration(value) * 365 * 24 * time.Hour
		default:
			return 0, fmt.Errorf("unknown unit %q in duration %q", unit, input)
		}
	}

	return total, nil
}
