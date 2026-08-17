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
