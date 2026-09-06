// Sonarr and Radarr: what the pair claims, and what keeps the claim true.

package recipe

import "testing"

// These two Recipes exist as a pair, and the pair carries a claim no single
// Recipe can: that two APIs which look identical -- same fork, same framework,
// same /api/v3 prefix, same info.version string, same credential schemes, 134
// shared paths -- diverge in specific places, and that the divergence is
// invisible to a client until it produces wrong answers rather than errors.
//
// A claim like that rots the first time somebody edits one Recipe and not the
// other. So it is checked here rather than left to two documents agreeing by
// habit.

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
// one stops being silent and the pair stops being worth having.
func TestSonarrAndRadarrAgreeWhereTheProductsDo(t *testing.T) {
	t.Parallel()

	sonarr, radarr := openArr(t, "sonarr"), openArr(t, "radarr")

	if sonarr.Auth.Header != radarr.Auth.Header {
		t.Errorf("the credential header differs: %q and %q", sonarr.Auth.Header, radarr.Auth.Header)
	}

	if sonarr.Auth.Header != "X-Api-Key" {
		t.Errorf("the credential header is %q, want X-Api-Key", sonarr.Auth.Header)
	}

	if got, want := sonarr.Responses.List.Key, radarr.Responses.List.Key; got != want {
		t.Errorf("the listing key differs: %q and %q", got, want)
	}

	for _, field := range []struct{ name, a, b string }{
		{"count", sonarr.Responses.List.CountField, radarr.Responses.List.CountField},
		{"page", sonarr.Responses.List.PageField, radarr.Responses.List.PageField},
		{"limit", sonarr.Responses.List.LimitField, radarr.Responses.List.LimitField},
	} {
		if field.a != field.b {
			t.Errorf("the %s field differs: %q and %q", field.name, field.a, field.b)
		}
	}

	// And the paths. A client's URL builder is the thing that would notice a
	// difference here, and it does not notice anything else.
	for _, path := range []string{"/api/v3/history", "/api/v3/system/status"} {
		if !hasPath(sonarr, path) {
			t.Errorf("sonarr no longer serves %s", path)
		}

		if !hasPath(radarr, path) {
			t.Errorf("radarr no longer serves %s", path)
		}
	}
}

// The half that must stay different, because it is the finding.
func TestSonarrAndRadarrDivergeWhereTheProductsDo(t *testing.T) {
	t.Parallel()

	sonarr, radarr := openArr(t, "sonarr"), openArr(t, "radarr")

	// The scope field. seriesId here, movieId there, and neither API refuses
	// the other's spelling -- ASP.NET Core model binding drops a query key
	// with no matching property, so the wrong one is a filter that returns
	// everything rather than an error.
	if _, ok := sonarr.Resources["history"].Fields["seriesId"]; !ok {
		t.Error("sonarr's history no longer carries seriesId, which is half of the pair's claim")
	}

	if _, ok := radarr.Resources["history"].Fields["movieId"]; !ok {
		t.Error("radarr's history no longer carries movieId, which is half of the pair's claim")
	}

	if _, ok := sonarr.Resources["history"].Fields["movieId"]; ok {
		t.Error("sonarr's history carries movieId, which is Radarr's field")
	}

	if _, ok := radarr.Resources["history"].Fields["seriesId"]; ok {
		t.Error("radarr's history carries seriesId, which is Sonarr's field")
	}

	// A series has episodes and a movie is the whole thing, so one side has a
	// second identifier and the other cannot be narrowed below the title.
	if _, ok := sonarr.Resources["history"].Fields["episodeId"]; !ok {
		t.Error("sonarr's history no longer carries episodeId")
	}

	if _, ok := radarr.Resources["history"].Fields["episodeId"]; ok {
		t.Error("radarr's history carries episodeId, and a movie has no episodes")
	}

	// The filter parameter, which is what a client actually sends.
	if got := filterParam(sonarr, "/api/v3/history", "seriesId"); got != "seriesIds" {
		t.Errorf("sonarr filters series by %q, want seriesIds", got)
	}

	if got := filterParam(radarr, "/api/v3/history", "movieId"); got != "movieIds" {
		t.Errorf("radarr filters movies by %q, want movieIds", got)
	}

	// And the one field that tells an operator which of the two they reached.
	sonarrName := fixtureField(sonarr, "library", "status", "appName")
	radarrName := fixtureField(radarr, "library", "status", "appName")

	if sonarrName != "Sonarr" || radarrName != "Radarr" {
		t.Errorf("appName says %v and %v, want Sonarr and Radarr -- it is the only field that discriminates", sonarrName, radarrName)
	}
}

// The eventType enums are the same length with different members, and the
// filter of that name takes integers. Positions two through six mean different
// things, so ?eventType=3 asks for a successful import on Sonarr and a failed
// download on Radarr.
//
// The Recipes cannot encode the integers -- neither document pins explicit
// values, so the mapping is inferred from declaration order -- but they can and
// do serve the names, and the names must not collide.
func TestTheArrEventTypeNamesDoNotAgree(t *testing.T) {
	t.Parallel()

	sonarr, radarr := openArr(t, "sonarr"), openArr(t, "radarr")

	sonarrEvents := fixtureValues(sonarr, "library", "history", "eventType")
	radarrEvents := fixtureValues(radarr, "library", "history", "eventType")

	if len(sonarrEvents) == 0 || len(radarrEvents) == 0 {
		t.Fatal("one of the fixtures no longer sets eventType, so the pair demonstrates nothing")
	}

	// downloadFolderImported is on both, at different positions in their
	// enums. That it appears in both fixtures is the point: the name is
	// shared and the integer behind it is not.
	shared := false

	for _, event := range sonarrEvents {
		for _, other := range radarrEvents {
			if event == other && event == "downloadFolderImported" {
				shared = true
			}
		}
	}

	if !shared {
		t.Error("neither fixture pair still shows downloadFolderImported, which is the name whose integer differs between the two")
	}

	// And each side has to carry a name the other's enum does not have, or
	// the divergence is not visible in what the emulator serves.
	if !contains(sonarrEvents, "grabbed") || !contains(radarrEvents, "grabbed") {
		t.Error("grabbed is position one on both and should appear on both")
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

func filterParam(r *Recipe, path, field string) string {
	for _, route := range r.Routes {
		if route.Path != path {
			continue
		}

		for _, filter := range route.Filters {
			if filter.Field == field {
				return filter.Param
			}
		}
	}

	return ""
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
