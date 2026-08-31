package openapi

import (
	"errors"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

func driftRecipeWith(spec, hash string) *recipe.Recipe {
	r := driftRecipe()
	r.Upstream = recipe.Upstream{
		API:      "v1",
		Spec:     spec,
		SpecHash: hash,
		SpecSeen: "2026-08-30",
	}

	return r
}

func serving(body string) func(string) ([]byte, error) {
	return func(string) ([]byte, error) { return []byte(body), nil }
}

func onlyReport(t *testing.T, r *recipe.Recipe, fetch func(string) ([]byte, error)) DriftReport {
	t.Helper()

	reports := Drift([]*recipe.Recipe{r}, fetch)
	if len(reports) != 1 {
		t.Fatalf("want one report, got %d", len(reports))
	}

	return reports[0]
}

// Most Recipes here have no machine-readable description to check against --
// the provider never published one, or published something that is not
// OpenAPI. That is a fact about the collection rather than a failure, and it
// is reported rather than skipped, so the number of Recipes nothing can check
// stays visible instead of quietly rounding to zero.
func TestARecipeWithNoDescriptionIsUndeclaredRatherThanFailing(t *testing.T) {
	report := onlyReport(t, driftRecipe(), func(string) ([]byte, error) {
		t.Error("a Recipe with no spec was fetched anyway")

		return nil, nil
	})

	if report.Status != Undeclared {
		t.Errorf("status is %q, want %q", report.Status, Undeclared)
	}

	if report.Err != nil {
		t.Errorf("a Recipe with no description reported an error: %v", report.Err)
	}
}

func TestADescriptionThatStillSaysTheSameThingIsUnchanged(t *testing.T) {
	doc, err := Parse([]byte(driftSpec))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	recorded := Fingerprint(driftRecipe(), doc, "")

	report := onlyReport(t, driftRecipeWith("https://example.test/openapi.json", recorded), serving(driftSpec))

	if report.Status != Unchanged {
		t.Errorf("status is %q, want %q (recorded %q, now %q)", report.Status, Unchanged, report.Recorded, report.Now)
	}
}

func TestAClaimedFieldChangingUpstreamIsDrift(t *testing.T) {
	doc, err := Parse([]byte(driftSpec))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	recorded := Fingerprint(driftRecipe(), doc, "")

	moved := strings.Replace(driftSpec,
		"                  name: {type: string}",
		"                  title: {type: string}", 1)

	report := onlyReport(t, driftRecipeWith("https://example.test/openapi.json", recorded), serving(moved))

	if report.Status != Moved {
		t.Errorf("status is %q, want %q", report.Status, Moved)
	}

	if report.Now == report.Recorded {
		t.Error("the report claims drift while both fingerprints are the same")
	}
}

// The distinction the whole scan rests on. A docs host answering 503, or a
// proxy returning an HTML error page, has told us nothing about whether the
// provider changed anything. Calling that drift trains everyone to ignore the
// scan, which costs more than not running it.
func TestAHostThatCannotBeReachedIsNotDrift(t *testing.T) {
	for _, fetch := range []func(string) ([]byte, error){
		func(string) ([]byte, error) { return nil, errors.New("503 Service Unavailable") },
		// A body that arrives and is not a description: an HTML error page,
		// a login redirect, a JSON blob that is not OpenAPI.
		serving("<html><body>Service Unavailable</body></html>"),
		serving(""),
	} {
		report := onlyReport(t, driftRecipeWith("https://example.test/openapi.json", "0123456789abcdef"), fetch)

		if report.Status == Moved {
			t.Errorf("an unreachable or unreadable description was reported as drift: %+v", report)
		}

		if report.Status != Unreachable {
			t.Errorf("status is %q, want %q", report.Status, Unreachable)
		}

		if report.Err == nil {
			t.Error("an unreachable description reported no reason")
		}
	}
}

// A Recipe that names a description but has never fingerprinted it cannot be
// said to have drifted. It reports what the fingerprint is now, so the value
// can be recorded, and says plainly that nothing was compared.
func TestADescriptionWithNoRecordedFingerprintIsUnrecorded(t *testing.T) {
	report := onlyReport(t, driftRecipeWith("https://example.test/openapi.json", ""), serving(driftSpec))

	if report.Status != Unrecorded {
		t.Errorf("status is %q, want %q", report.Status, Unrecorded)
	}

	if report.Now == "" {
		t.Error("an unrecorded Recipe did not report the fingerprint to record")
	}

	if report.Recorded != "" {
		t.Errorf("an unrecorded Recipe reported a recorded fingerprint of %q", report.Recorded)
	}
}

// Drift is the only status that should ever fail a build, so nothing else may
// wear it.
func TestOnlyAMovedFingerprintCountsAsDrift(t *testing.T) {
	doc, err := Parse([]byte(driftSpec))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	recorded := Fingerprint(driftRecipe(), doc, "")

	for _, r := range []*recipe.Recipe{
		driftRecipe(),
		driftRecipeWith("https://example.test/openapi.json", ""),
		driftRecipeWith("https://example.test/openapi.json", recorded),
	} {
		if report := onlyReport(t, r, serving(driftSpec)); report.Moved() {
			t.Errorf("%q was counted as drift", report.Status)
		}
	}

	moved := strings.Replace(driftSpec, "  /v1/things:", "  /v1/items:", 1)

	if report := onlyReport(t, driftRecipeWith("https://example.test/openapi.json", recorded), serving(moved)); !report.Moved() {
		t.Error("a moved fingerprint was not counted as drift")
	}
}

// A description in a format this cannot read is a permanent fact about the
// provider, and a host that answered 503 is a Tuesday. Both leave the Recipe
// unchecked, and reporting them under one word means a scheduled scan shows
// the same line every day for a Recipe whose provider is perfectly reachable,
// with nothing to say which of the two it is.
//
// The MBTA found it: api-v3.mbta.com answers instantly and will never once
// be unreachable, and it published the one description here that could not
// be read. Swagger 2.0 is now rewritten and read, so the example has to be
// a format with no reader -- which is the better test anyway, because the
// old one stopped being unreadable and took its own premise with it.
func TestAFormatThatCannotBeReadIsNotAHostThatCannotBeReached(t *testing.T) {
	raml := "#%RAML 1.0\ntitle: Old\nbaseUri: https://example.test/v1\n"

	report := onlyReport(t, driftRecipeWith("https://example.test/api.raml", "0123456789abcdef"), serving(raml))

	if report.Status != Unsupported {
		t.Errorf("status is %q, want %q", report.Status, Unsupported)
	}

	if report.Err == nil {
		t.Error("an unreadable format reported no reason")
	}

	if report.Moved() {
		t.Error("an unreadable format was counted as drift")
	}
}

// A description declares its paths relative to its own server, and a Recipe
// carries the whole path a client requests. Box is the example the check
// command already documents: the description says /files/{id} beside a server
// of https://api.box.com/2.0, while the Recipe says /2.0/files/{id}.
//
// Read literally those never match, so every route fingerprints as absent --
// a value that is stable, and says nothing. Worse than a wrong answer: a
// confident one about nothing, which would report unchanged forever no matter
// what the provider did to the paths the Recipe actually uses.
func TestTheBasePathComesFromTheDescriptionsOwnServers(t *testing.T) {
	prefixed := `
openapi: 3.0.0
info: {title: Things, version: "1"}
servers:
  - url: https://api.example.test/2.0
paths:
  /things:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  name: {type: string}
`

	r := driftRecipeWith("https://example.test/openapi.json", "")
	r.Routes = []recipe.Route{
		{Method: "GET", Path: "/2.0/things", Resource: "thing", Operation: "list"},
	}

	report := onlyReport(t, r, serving(prefixed))

	if report.Status != Unrecorded {
		t.Fatalf("status is %q, want %q", report.Status, Unrecorded)
	}

	// The same Recipe against a description that has genuinely dropped the
	// path must not fingerprint the same, which is what would happen if both
	// were reading as absent.
	gone := strings.Replace(prefixed, "  /things:", "  /gone:", 1)

	away := onlyReport(t, r, serving(gone))

	if report.Now == away.Now {
		t.Errorf("a found path and a missing one fingerprinted alike as %q, so the base path was not applied", report.Now)
	}
}
