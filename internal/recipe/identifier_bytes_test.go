package recipe

import (
	"testing"
	"unicode"
)

// A fixture identifier has to be typeable, greppable and copy-pasteable,
// because that is the whole of what a fixture is for: somebody reads an id out
// of a Recipe and puts it in a request.
//
// Twice while writing Recipes in bulk, a stray non-ASCII character landed in
// the middle of one -- `7c2d5e81-...-9d18c4ف` and `pub_2d7e9ب` -- from an
// editor, a paste, or a model writing the file. Both looked fine at a glance
// and both were caught by reading rather than by anything mechanical. Nothing
// else would have: the Recipe parsed, validated, served and verified, because
// an opaque identifier is opaque and the emulator does not care what is in it.
//
// The failure it causes is downstream and confusing. A developer copies the id
// out of the fixture, the terminal or the browser normalises it differently,
// and the request 404s against a fake that is serving exactly that record.
//
// So: identifiers stay printable ASCII. Prose in comments and record fields is
// unaffected -- Recipes cite provider sentences verbatim and some of those are
// not ASCII -- this is only about the values used to address a record.
func TestFixtureIdentifiersArePrintableASCII(t *testing.T) {
	for _, name := range Bundled() {
		r, err := Open(name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}

		for fixture, records := range r.Fixtures {
			for resource, rows := range records {
				spec, ok := r.Resources[resource]
				if !ok {
					continue
				}

				field := spec.ID.Field
				if field == "" {
					field = "id"
				}

				for i, row := range rows {
					id, ok := row[field].(string)
					if !ok {
						continue
					}

					for _, c := range id {
						if c > unicode.MaxASCII || !unicode.IsPrint(c) {
							t.Errorf("%s fixture %q record %d of %s: identifier %q contains %q, which is not printable ASCII",
								name, fixture, i, resource, id, c)

							break
						}
					}
				}
			}
		}
	}
}
