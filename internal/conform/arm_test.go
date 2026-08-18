package conform

import (
	"net/http"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The Armer is called for every case, not only the ones arming something,
// because that is what clears a fault the previous case installed.
//
// Without it a declared failure leaks into every case after the one that armed
// it, and the symptom is a suite that passes case by case and fails in order,
// which is among the least pleasant things to debug. The count of one that the
// caller arms with is a second guard against the same thing; this test pins the
// first, since a guard nothing checks is a guard somebody deletes.
func TestEveryCaseClearsWhateverThePreviousOneArmed(t *testing.T) {
	r := &recipe.Recipe{
		Name: "example",
		Conformance: []recipe.Case{
			{Name: "arms", Arm: "teapot", Expect: recipe.Expectation{Status: http.StatusTeapot}},
			{Name: "arms nothing", Expect: recipe.Expectation{Status: http.StatusOK}},
			{Name: "arms nothing either", Expect: recipe.Expectation{Status: http.StatusOK}},
		},
	}

	// A handler that fails while a fault is installed, which is the behaviour
	// the runtime has and this package deliberately knows nothing about.
	armed := ""

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if armed != "" {
			w.WriteHeader(http.StatusTeapot)

			return
		}

		w.WriteHeader(http.StatusOK)
	})

	var asked []string

	report := Run(r, handler, "", nil, func(name string) error {
		asked = append(asked, name)
		armed = name

		return nil
	})

	// Three calls, not one. The two empty names are the clearing.
	if len(asked) != 3 {
		t.Fatalf("the armer was called %d times for 3 cases: %q", len(asked), asked)
	}

	if asked[1] != "" || asked[2] != "" {
		t.Errorf("a case arming nothing must still clear, got %q", asked)
	}

	for _, result := range report.Results {
		if len(result.Failures) > 0 {
			t.Errorf("%s: %v", result.Case.Name, result.Failures)
		}
	}
}

// An Armer that refuses is a failure of the case, not a panic and not a silent
// pass. A Recipe naming an error the runtime will not install has to say so
// somewhere, and the case is where it is visible.
func TestACaseThatCannotBeArmedFails(t *testing.T) {
	r := &recipe.Recipe{
		Name: "example",
		Conformance: []recipe.Case{
			{Name: "arms something unknown", Arm: "nonesuch", Expect: recipe.Expectation{Status: http.StatusTeapot}},
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	report := Run(r, handler, "", nil, func(string) error {
		return errNoSuchError
	})

	if report.Passed() != 0 {
		t.Fatal("a case whose fault could not be installed must not pass, since the status it expects can arrive by accident")
	}
}

var errNoSuchError = errNoSuch{}

type errNoSuch struct{}

func (errNoSuch) Error() string { return "no such error" }
