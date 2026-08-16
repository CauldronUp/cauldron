package clock

import (
	"sync"
	"testing"
	"time"
)

func TestNewStartsAtTheEpoch(t *testing.T) {
	c := New()

	if !c.Now().Equal(Epoch) {
		t.Errorf("Now() = %v, want %v", c.Now(), Epoch)
	}
}

func TestAdvanceMovesTimeForward(t *testing.T) {
	c := New()

	c.Advance(24 * time.Hour)

	if got, want := c.Now(), Epoch.Add(24*time.Hour); !got.Equal(want) {
		t.Errorf("Now() = %v, want %v", got, want)
	}
}

func TestAdvanceRefusesToGoBackwards(t *testing.T) {
	c := New()

	c.Advance(-48 * time.Hour)

	if !c.Now().Equal(Epoch) {
		t.Errorf("a negative advance must be ignored; Now() = %v", c.Now())
	}
}

func TestResetReturnsToTheEpoch(t *testing.T) {
	c := New()
	c.Advance(90 * 24 * time.Hour)

	c.Reset()

	if !c.Now().Equal(Epoch) {
		t.Errorf("Now() = %v, want the epoch", c.Now())
	}
}

func TestClockIsSafeForConcurrentUse(t *testing.T) {
	c := New()

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			c.Advance(time.Second)
		}()

		go func() {
			defer wg.Done()
			_ = c.Now()
		}()
	}

	wg.Wait()

	if got, want := c.Now(), Epoch.Add(50*time.Second); !got.Equal(want) {
		t.Errorf("Now() = %v, want %v", got, want)
	}
}

func TestTwoClocksAgreeGivenTheSameOperations(t *testing.T) {
	a, b := New(), New()

	for _, d := range []time.Duration{time.Hour, 30 * 24 * time.Hour, 90 * time.Second} {
		a.Advance(d)
		b.Advance(d)
	}

	if !a.Now().Equal(b.Now()) {
		t.Errorf("clocks diverged: %v vs %v", a.Now(), b.Now())
	}
}

func TestParseDuration(t *testing.T) {
	day := 24 * time.Hour

	cases := map[string]time.Duration{
		"30s":   30 * time.Second,
		"15m":   15 * time.Minute,
		"12h":   12 * time.Hour,
		"1h30m": 90 * time.Minute,
		"7d":    7 * day,
		"2w":    14 * day,
		"6mo":   180 * day,
		"1y":    365 * day,
		"1y2w":  365*day + 14*day,
		"500ms": 500 * time.Millisecond,
	}

	for input, want := range cases {
		got, err := ParseDuration(input)
		if err != nil {
			t.Errorf("ParseDuration(%q): %v", input, err)
			continue
		}

		if got != want {
			t.Errorf("ParseDuration(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestParseDurationRejectsNonsense(t *testing.T) {
	for _, input := range []string{"", "soon", "30", "30dd", "d30", "30 dogs"} {
		if _, err := ParseDuration(input); err == nil {
			t.Errorf("ParseDuration(%q) should have failed", input)
		}
	}
}

func TestParseDurationErrorSuggestsUnits(t *testing.T) {
	_, err := ParseDuration("next tuesday")
	if err == nil {
		t.Fatal("expected an error")
	}

	if want := "30s"; !contains(err.Error(), want) {
		t.Errorf("error should suggest valid units, got %q", err)
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (haystack == needle || indexOf(haystack, needle) >= 0)
}

func indexOf(haystack, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}

	return -1
}
