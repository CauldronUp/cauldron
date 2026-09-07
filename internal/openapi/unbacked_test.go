package openapi

import (
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A moved fingerprint says that something in the claim list changed and
// nothing about what, and the cheapest response to a line nobody can read is
// to re-record it. That reflex is how a drift scan becomes a weekly rubber
// stamp, so the report carries the half a reader can act on: the claims this
// Recipe makes that the description no longer backs.
func TestAMovedReportNamesTheClaimsTheDescriptionDoesNotBack(t *testing.T) {
	doc, err := Parse([]byte(driftSpec))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	// Two fields, because one field vanishing from a one-field Recipe is
	// indistinguishable from this parser never opening an envelope, and the
	// rule below prefers silence in that case. Here id survives, so the
	// reading demonstrably got inside the response and name's absence is a
	// fact about the description.
	r := driftRecipe()
	r.Resources["thing"] = recipe.Resource{Fields: map[string]recipe.Field{
		"id": {Type: "string"}, "name": {Type: "string"},
	}}

	recorded := Fingerprint(r, doc, "")

	// The field the Recipe names is gone from the success schema, and so is
	// the status it declares an error for.
	moved := strings.Replace(driftSpec,
		"                  name: {type: string}",
		"                  title: {type: string}", 1)
	moved = strings.Replace(moved, `        "404": {description: gone}`, "", 1)

	report := onlyReport(t, withUpstream(r, "https://example.test/openapi.json", recorded), serving(moved))

	if report.Status != Moved {
		t.Fatalf("status is %q, want %q", report.Status, Moved)
	}

	joined := strings.Join(report.Gaps, "\n")

	if !strings.Contains(joined, "name") {
		t.Errorf("the report does not say the schema stopped declaring name:\n%s", joined)
	}

	if !strings.Contains(joined, "404") {
		t.Errorf("the report does not say nothing answers 404 any more:\n%s", joined)
	}
}

// A description that still backs every claim reports no gaps, so the presence
// of a gap means something rather than being decoration on every line.
func TestADescriptionThatBacksEveryClaimHasNoGaps(t *testing.T) {
	doc, err := Parse([]byte(driftSpec))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if gaps := Unbacked(driftRecipe(), doc, ""); len(gaps) != 0 {
		t.Errorf("a description backing every claim reported gaps: %v", gaps)
	}
}

// The two questions are separate, and this is the case that proves it: a
// Recipe pinned to one of several reference files a provider splits its API
// across has had a route missing from that document since the day it was
// written. Its fingerprint has never moved. Reading "unchanged" as "the
// description covers this Recipe" is the mistake -- Basiq and Customer.io are
// both in that position in the shipped collection.
func TestAGapCanExistWhileTheFingerprintHasNeverMoved(t *testing.T) {
	r := driftRecipe()
	r.Routes = append(r.Routes, recipeRouteNotInTheDocument())

	doc, err := Parse([]byte(driftSpec))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	recorded := Fingerprint(r, doc, "")

	report := onlyReport(t, withUpstream(r, "https://example.test/openapi.json", recorded), serving(driftSpec))

	if report.Status != Unchanged {
		t.Fatalf("status is %q, want %q", report.Status, Unchanged)
	}

	if len(report.Gaps) == 0 {
		t.Fatal("a route the description has never declared reported no gap")
	}

	if !strings.Contains(strings.Join(report.Gaps, "\n"), "/v1/elsewhere") {
		t.Errorf("the gap does not name the missing route: %v", report.Gaps)
	}
}

func recipeRouteNotInTheDocument() recipe.Route {
	return recipe.Route{Method: "GET", Path: "/v1/elsewhere", Resource: "thing", Operation: "list"}
}

func withUpstream(r *recipe.Recipe, spec, hash string) *recipe.Recipe {
	r.Upstream = recipe.Upstream{
		API:      "v1",
		Spec:     spec,
		SpecHash: hash,
		SpecSeen: "2026-08-30",
	}

	return r
}

// The envelope case, which would otherwise drown the report. Document.Properties
// stops at the success schema's top level, so a provider wrapping its records in
// {"data": {...}} declares exactly one property and every field the Recipe names
// reads as absent. None of those absences is a fact about the provider, and nine
// shipped Recipes are in that position -- Vercel alone would print forty lines a
// week saying that this parser did not open an envelope.
//
// The claims are still made, so a Recipe adding a field still moves its
// fingerprint. Only the gaps are suppressed.
func TestAWrappedResponseDoesNotReportEveryFieldAsAGap(t *testing.T) {
	const wrapped = `
openapi: 3.0.0
info: {title: Things, version: "1"}
paths:
  /v1/things:
    get:
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  data: {type: object}
        "404": {description: gone}
`

	doc, err := Parse([]byte(wrapped))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if gaps := Unbacked(driftRecipe(), doc, ""); len(gaps) != 0 {
		t.Errorf("an unopened envelope reported its fields as gaps: %v", gaps)
	}

	// And the claim is still recorded, so the fingerprint still moves when
	// the Recipe's own field list changes.
	before := Fingerprint(driftRecipe(), doc, "")

	other := driftRecipe()
	other.Resources["thing"] = recipe.Resource{Fields: map[string]recipe.Field{
		"name": {Type: "string"}, "extra": {Type: "string"},
	}}

	if Fingerprint(other, doc, "") == before {
		t.Error("suppressing the gaps also stopped the fingerprint moving")
	}
}

// One field missing where the others are present is the real finding, and it
// has to survive the suppression above.
func TestOneMissingFieldAmongPresentOnesIsStillAGap(t *testing.T) {
	partial := strings.Replace(driftSpec,
		"                  name: {type: string}",
		"                  other: {type: string}", 1)

	doc, err := Parse([]byte(partial))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	r := driftRecipe()
	r.Resources["thing"] = recipe.Resource{Fields: map[string]recipe.Field{
		"name": {Type: "string"}, "id": {Type: "string"},
	}}

	gaps := Unbacked(r, doc, "")

	if len(gaps) != 1 || !strings.Contains(gaps[0], "name") {
		t.Errorf("want one gap naming name, got %v", gaps)
	}
}
