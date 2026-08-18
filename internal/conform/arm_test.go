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
	}, nil)

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
	}, nil)

	if report.Passed() != 0 {
		t.Fatal("a case whose fault could not be installed must not pass, since the status it expects can arrive by accident")
	}
}

var errNoSuchError = errNoSuch{}

type errNoSuch struct{}

func (errNoSuch) Error() string { return "no such error" }

// A webhook claim has to be able to fail, which is not obvious: the block sits
// beside the response assertions and a case whose only claim is about a
// webhook passes trivially if the check is skipped.
func TestAWebhookClaimCanFail(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	emitted := []Delivery{{
		Event:   "payment.created",
		Payload: map[string]any{"data": map[string]any{"object": map[string]any{"amount": 2500}}},
	}}

	run := func(name string, expect recipe.WebhookExpectation) Result {
		t.Helper()

		r := &recipe.Recipe{
			Name:     "example",
			Webhooks: recipe.Webhooks{Events: []string{"payment.created"}},
			Conformance: []recipe.Case{{
				Name:   name,
				Expect: recipe.Expectation{Status: http.StatusOK, Webhook: &expect},
			}},
		}

		// Nothing recorded before the request, everything after, which is the
		// arrangement a real emitting request produces.
		before := true

		return Run(r, handler, "", nil, nil, func() []Delivery {
			if before {
				before = false

				return nil
			}

			return emitted
		}).Results[0]
	}

	if got := run("right event", recipe.WebhookExpectation{Event: "payment.created"}); !got.Passed() {
		t.Errorf("a correct claim failed: %v", got.Failures)
	}

	if got := run("wrong event", recipe.WebhookExpectation{Event: "payment.updated"}); got.Passed() {
		t.Error("a claim naming the wrong event passed")
	}

	if got := run("wrong value", recipe.WebhookExpectation{
		Body: map[string]any{"data.object.amount": 99},
	}); got.Passed() {
		t.Error("a claim about a value the payload does not carry passed")
	}

	// The half that catches an internal field name leaking into a payload,
	// which is the failure this whole mechanism exists for.
	if got := run("present but claimed absent", recipe.WebhookExpectation{
		Absent: []string{"data.object.amount"},
	}); got.Passed() {
		t.Error("a claim that a present field is absent passed")
	}

	if got := run("nothing emitted", recipe.WebhookExpectation{None: true}); got.Passed() {
		t.Error("a claim that nothing was emitted passed while something was")
	}
}
