// Whether a provider's own description still says what a Recipe was written
// against.

package openapi

import (
	"errors"
	"fmt"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The states a Recipe can be in with respect to the description it names.
const (
	// Undeclared is a Recipe that names no machine-readable description.
	//
	// This is most of them, and it is reported rather than skipped so the
	// number of Recipes nothing can check stays visible. A scan that silently
	// omits what it cannot examine reads as a clean bill of health for the
	// whole collection.
	Undeclared = "undeclared"
	// Unrecorded is a description named but never fingerprinted, so there is
	// nothing to compare against yet.
	Unrecorded = "unrecorded"
	// Unchanged is a description that still says what the Recipe claims.
	Unchanged = "unchanged"
	// Moved is the finding: the description changed where the Recipe claims
	// something.
	Moved = "moved"
	// Unsupported is a description in a format this cannot read.
	//
	// Separate from Unreachable because the two are different facts with
	// different remedies. A host answering 503 is a Tuesday; a provider
	// publishing Swagger 2.0 is permanent until somebody converts it. Under
	// one word, a scheduled scan shows the same line every day for a Recipe
	// whose provider is perfectly reachable, and nothing says which it is.
	//
	// The MBTA found it: api-v3.mbta.com publishes Swagger 2.0, answers
	// instantly, and will never once be unreachable.
	Unsupported = "unsupported"
	// Unreachable is a description that could not be fetched or read.
	//
	// Deliberately not drift. A docs host answering 503, a proxy returning an
	// HTML error page, a login redirect -- none of those has said anything
	// about whether the provider changed. Reporting them as drift trains
	// everyone to ignore the scan, and an ignored scan is worse than none,
	// because it also costs the time spent reading it.
	Unreachable = "unreachable"
)

// DriftReport is one Recipe's answer.
type DriftReport struct {
	// Recipe is its name.
	Recipe string
	// Spec is the description it names, if any.
	Spec string
	// Recorded is the fingerprint written into the Recipe.
	Recorded string
	// Now is the fingerprint of what was fetched.
	Now string
	// Seen is the date the recorded fingerprint was taken.
	Seen string
	// Status is one of the constants above.
	Status string
	// Err is why an unreachable description could not be read.
	Err error
}

// Moved reports whether this is the one status that should fail a build.
func (d DriftReport) Moved() bool { return d.Status == Moved }

// Drift fingerprints each Recipe's declared description as it stands now and
// compares it with the fingerprint the Recipe recorded.
//
// The fetch is a parameter rather than a package-level client because the
// interesting cases here are the failures, and a test that has to stand up an
// HTTP server to describe "the host answered with a login page" describes it
// badly.
func Drift(recipes []*recipe.Recipe, fetch func(url string) ([]byte, error)) []DriftReport {
	reports := make([]DriftReport, 0, len(recipes))

	for _, r := range recipes {
		reports = append(reports, driftOf(r, fetch))
	}

	return reports
}

func driftOf(r *recipe.Recipe, fetch func(url string) ([]byte, error)) DriftReport {
	report := DriftReport{
		Recipe:   r.Name,
		Spec:     r.Upstream.Spec,
		Recorded: r.Upstream.SpecHash,
		Seen:     r.Upstream.SpecSeen,
	}

	if report.Spec == "" {
		report.Status = Undeclared

		return report
	}

	raw, err := fetch(report.Spec)
	if err != nil {
		report.Status = Unreachable
		report.Err = err

		return report
	}

	if len(raw) == 0 {
		report.Status = Unreachable
		report.Err = errors.New("the description is empty")

		return report
	}

	doc, err := Parse(raw)
	if err != nil {
		// A document that arrived intact and announced a format this does
		// not read is a fact about the provider rather than about the
		// network, and it stays true tomorrow.
		//
		// This searched the message for the word "Swagger" until Swagger 2.0
		// stopped being unreadable. The search was wrong both ways: it
		// caught a network error that happened to mention Swagger in a URL,
		// and it missed every format that is not Swagger, so a provider
		// publishing RAML was reported unreachable forever.
		var format *FormatError

		if errors.As(err, &format) {
			report.Status = Unsupported
		} else {
			report.Status = Unreachable
		}

		report.Err = fmt.Errorf("the description could not be read: %w", err)

		return report
	}

	// A document that parses but declares no paths at all is an HTML error
	// page or a login redirect that happened to survive the parser, not a
	// description of an API. Treating it as one would fingerprint every route
	// as absent and report the whole collection as drift on the day a docs
	// host goes down.
	if len(doc.Paths) == 0 {
		report.Status = Unreachable
		report.Err = errors.New("the description declares no paths, so it is not a description of this API")

		return report
	}

	// The base comes from the description's own servers, exactly as the check
	// command reads it. Without it, Box's /2.0/files/{id} never matches a
	// description saying /files/{id} beside a server of api.box.com/2.0, every
	// route fingerprints as absent, and the result is a stable value that says
	// nothing -- which would report unchanged forever whatever the provider
	// did to the paths the Recipe actually uses.
	report.Now = Fingerprint(r, doc, BasePath(doc))

	switch {
	case report.Recorded == "":
		report.Status = Unrecorded
	case report.Recorded == report.Now:
		report.Status = Unchanged
	default:
		report.Status = Moved
	}

	return report
}
