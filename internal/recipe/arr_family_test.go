// Sonarr, Radarr, Prowlarr, Lidarr and Readarr: what the family claims, and
// what keeps the claim true.

package recipe

import "testing"

// These five Recipes exist as a set, and the set carries a claim no single
// Recipe can: that five forks of one codebase -- commonly installed on the same
// machine, sharing a credential header, a path prefix, a paged envelope and an
// API shape -- diverge in specific places, and that the divergence is invisible
// to a client until it produces wrong answers rather than errors.
//
// The sharpest form of it is the eventType enum. None of the five documents
// pins explicit values, so the integer a caller sends is a member's declaration
// position, and only position nought agrees across all five:
//
//	pos | sonarr                 | radarr                 | prowlarr       | lidarr                | readarr
//	 0  | unknown                | unknown                | unknown        | unknown               | unknown
//	 1  | grabbed                | grabbed                | releaseGrabbed | grabbed               | grabbed
//	 2  | seriesFolderImported   | downloadFolderImported | indexerQuery   | artistFolderImported  | bookFileImported
//	 3  | downloadFolderImported | downloadFailed         | indexerRss     | trackFileImported     | downloadFailed
//	 4  | downloadFailed         | movieFileDeleted       | indexerAuth    | downloadFailed        | bookFileDeleted
//
// A claim like that rots the first time somebody edits one Recipe and not the
// others. So it is checked here rather than left to five documents agreeing by
// habit.

// arrFamily is the set, in the order the tables above use.
var arrFamily = []string{"sonarr", "radarr", "prowlarr", "lidarr", "readarr"}

// arrPrefix is each fork's own path prefix. Two of the five are on v3 and three
// on v1, which is one more thing that looks shared and is not.
var arrPrefix = map[string]string{
	"sonarr":   "/api/v3",
	"radarr":   "/api/v3",
	"prowlarr": "/api/v1",
	"lidarr":   "/api/v1",
	"readarr":  "/api/v1",
}

func openArr(t *testing.T, name string) *Recipe {
	t.Helper()

	r, err := Open(name)
	if err != nil {
		t.Fatalf("open %s: %v", name, err)
	}

	return r
}

// The half that must stay the same, because the resemblance is what makes the
// difference dangerous. If these ever diverge, pointing a client at the wrong
// one stops being silent and the set stops being worth having.
func TestTheArrFamilyAgreesWhereTheProductsDo(t *testing.T) {
	t.Parallel()

	for _, name := range arrFamily {
		r := openArr(t, name)

		if r.Auth.Header != "X-Api-Key" {
			t.Errorf("%s: the credential header is %q, want X-Api-Key", name, r.Auth.Header)
		}

		if got := r.Responses.List.Key; got != "records" {
			t.Errorf("%s: the listing key is %q, want records", name, got)
		}

		for _, field := range []struct{ what, got, want string }{
			{"count", r.Responses.List.CountField, "totalRecords"},
			{"page", r.Responses.List.PageField, "page"},
			{"limit", r.Responses.List.LimitField, "pageSize"},
		} {
			if field.got != field.want {
				t.Errorf("%s: the %s field is %q, want %q", name, field.what, field.got, field.want)
			}
		}

		// And the paths -- under each fork's own version segment, which is
		// itself a divergence: Sonarr and Radarr serve /api/v3, and Prowlarr,
		// Lidarr and Readarr serve /api/v1. Same product family, same
		// generated shape, two different numbers in the URL, and no
		// relationship at all to info.version, which is "3.0.0" for the two
		// on v3 and "1.0.0" for the three on v1.
		for _, path := range []string{"/history", "/system/status"} {
			full := arrPrefix[name] + path

			if !hasPath(r, full) {
				t.Errorf("%s no longer serves %s", name, full)
			}
		}
	}
}

// The half that must stay different, because it is the finding.
//
// Each fork's scope field is its own, and no sibling accepts another's
// spelling: ASP.NET Core model binding drops a query key with no matching
// property rather than refusing it, so the wrong one is a filter that returns
// everything instead of an error.
func TestTheArrFamilyDivergesWhereTheProductsDo(t *testing.T) {
	t.Parallel()

	scopes := map[string][]string{
		"sonarr":   {"seriesId", "episodeId"},
		"radarr":   {"movieId"},
		"prowlarr": {"indexerId"},
		"lidarr":   {"artistId", "albumId", "trackId"},
		"readarr":  {"authorId", "bookId"},
	}

	for name, own := range scopes {
		r := openArr(t, name)
		fields := r.Resources["history"].Fields

		for _, field := range own {
			if _, ok := fields[field]; !ok {
				t.Errorf("%s's history no longer carries %s, which is part of the family's claim", name, field)
			}
		}

		// And it must carry none of the others'.
		for other, theirs := range scopes {
			if other == name {
				continue
			}

			for _, field := range theirs {
				if contains(own, field) {
					continue
				}

				if _, ok := fields[field]; ok {
					t.Errorf("%s's history carries %s, which is %s's field", name, field, other)
				}
			}
		}
	}
}

// Only appName tells the five apart. instanceName looks like a second
// discriminator and is not: it is settable, and an operator running several may
// rename them to anything.
func TestOnlyAppNameSeparatesTheArrFamily(t *testing.T) {
	t.Parallel()

	seen := map[string]string{}

	for _, name := range arrFamily {
		r := openArr(t, name)

		fixture := "library"
		if name == "prowlarr" {
			fixture = "indexers"
		}

		appName, _ := fixtureField(r, fixture, "status", "appName").(string)
		if appName == "" {
			t.Errorf("%s's status fixture sets no appName, and it is the only discriminator", name)
			continue
		}

		if previous, ok := seen[appName]; ok {
			t.Errorf("%s and %s both report appName %q, so nothing tells them apart", previous, name, appName)
		}

		seen[appName] = name
	}

	if len(seen) != len(arrFamily) {
		t.Errorf("%d distinct appNames across %d forks", len(seen), len(arrFamily))
	}
}

// The event names must not collide across the family in the places the enums
// diverge, and must agree in the one place they do.
//
// The Recipes cannot encode the integers -- none of the five documents pins
// explicit values, so the mapping is inferred from declaration order -- but
// they serve the names, and the names are what a reader of these Recipes
// compares.
func TestTheArrEventNamesDoNotAgree(t *testing.T) {
	t.Parallel()

	events := map[string][]string{}

	for _, name := range arrFamily {
		fixture := "library"
		if name == "prowlarr" {
			fixture = "indexers"
		}

		events[name] = fixtureValues(openArr(t, name), fixture, "history", "eventType")

		if len(events[name]) == 0 {
			t.Fatalf("%s's fixture no longer sets eventType, so the family demonstrates nothing", name)
		}
	}

	// grabbed is position one on four of the five, and Prowlarr renames it.
	for _, name := range []string{"sonarr", "radarr", "lidarr", "readarr"} {
		if !contains(events[name], "grabbed") {
			t.Errorf("%s's fixture no longer shows `grabbed`, which four of the five share at position one", name)
		}
	}

	if contains(events["prowlarr"], "grabbed") {
		t.Error("prowlarr's fixture shows `grabbed`, and prowlarr calls that event releaseGrabbed")
	}

	if !contains(events["prowlarr"], "releaseGrabbed") {
		t.Error("prowlarr's fixture no longer shows `releaseGrabbed`, which is its own name for the shared event")
	}

	// And downloadFailed, the event with no consistent integer across the
	// family -- 4 on Sonarr and Lidarr, 3 on Radarr and Readarr, absent from
	// Prowlarr. Its presence in a fixture is what makes the table checkable.
	failing := 0

	for _, name := range []string{"sonarr", "radarr", "lidarr", "readarr"} {
		if contains(events[name], "downloadFailed") {
			failing++
		}
	}

	if failing < 2 {
		t.Errorf("only %d of the four download managers show downloadFailed; the family's sharpest claim is that its integer differs", failing)
	}

	if contains(events["prowlarr"], "downloadFailed") {
		t.Error("prowlarr's fixture shows downloadFailed, and prowlarr's enum has no such member")
	}
}

func hasPath(r *Recipe, path string) bool {
	for _, route := range r.Routes {
		if route.Path == path {
			return true
		}
	}

	return false
}

func fixtureField(r *Recipe, fixture, resource, field string) any {
	for _, record := range r.Fixtures[fixture][resource] {
		if value, ok := record[field]; ok {
			return value
		}
	}

	return nil
}

func fixtureValues(r *Recipe, fixture, resource, field string) []string {
	var out []string

	for _, record := range r.Fixtures[fixture][resource] {
		if value, ok := record[field].(string); ok {
			out = append(out, value)
		}
	}

	return out
}
