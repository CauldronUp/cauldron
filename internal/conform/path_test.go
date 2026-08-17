package conform

import "testing"

// Dropbox names a field ".tag", where the leading dot is part of the name
// rather than a separator. Without an escape there is no way to assert on it,
// and the emulator was sending it correctly all along while every case that
// mentioned it failed.
func TestSplitFieldPathHonoursAnEscapedDot(t *testing.T) {
	cases := []struct {
		path string
		want []string
	}{
		{"data.id", []string{"data", "id"}},
		{`\.tag`, []string{".tag"}},
		{`entries[0].\.tag`, []string{"entries[0]", ".tag"}},
		{`a.\.b.c`, []string{"a", ".b", "c"}},
		{"plain", []string{"plain"}},
	}

	for _, c := range cases {
		got := splitFieldPath(c.path)

		if len(got) != len(c.want) {
			t.Errorf("splitFieldPath(%q) = %v, want %v", c.path, got, c.want)
			continue
		}

		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("splitFieldPath(%q) = %v, want %v", c.path, got, c.want)
				break
			}
		}
	}
}

func TestLookupFindsAFieldNamedWithALeadingDot(t *testing.T) {
	document := map[string]any{
		".tag": "file",
		"entries": []any{
			map[string]any{".tag": "folder", "name": "Documents"},
		},
	}

	value, ok := lookup(document, `\.tag`)
	if !ok || value != "file" {
		t.Errorf("lookup of the escaped tag = %v, %v", value, ok)
	}

	nested, ok := lookup(document, `entries[0].\.tag`)
	if !ok || nested != "folder" {
		t.Errorf("lookup of the nested escaped tag = %v, %v", nested, ok)
	}

	// An unescaped dot still separates, so nothing about ordinary paths
	// changes.
	if _, ok := lookup(document, ".tag"); ok {
		t.Error("an unescaped leading dot should not resolve")
	}
}

// A string "0" is not the number 0, and telling them apart is the whole point
// of several Recipes: Vonage reports a successful send with the string "0" and
// Docusign counts with strings. Comparing only the rendered form let those
// cases pass whichever the emulator sent.
func TestScalarsCompareByKindAsWellAsValue(t *testing.T) {
	document := map[string]any{
		"stringZero": "0",
		"numberZero": float64(0),
		"count":      "4",
	}

	if failures := compare("stringZero", "0", document); len(failures) != 0 {
		t.Errorf("a string against a string should match: %v", failures)
	}

	// The bug this test exists for: without a kind check these pass.
	if failures := compare("stringZero", float64(0), document); len(failures) == 0 {
		t.Error("the number 0 should not match the string \"0\"")
	}

	if failures := compare("numberZero", "0", document); len(failures) == 0 {
		t.Error("the string \"0\" should not match the number 0")
	}

	if failures := compare("count", 4, document); len(failures) == 0 {
		t.Error("the number 4 should not match the string \"4\"")
	}

	// Integer and float are the same kind, because YAML and JSON disagree
	// about which one a literal is and no provider distinguishes them.
	if failures := compare("numberZero", 0, document); len(failures) != 0 {
		t.Errorf("an int should match a float of the same value: %v", failures)
	}
}
