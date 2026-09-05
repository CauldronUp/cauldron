package recipe

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// docs/format.md names every key a Recipe may carry, and no key it may not.
//
// The format had no written reference at all. Six hundred Recipes were written
// against the Go source and against each other, and the mistakes that produced
// are the specific kind you would expect: two Recipes put type_field and
// status_field inside an entry in the errors table, where those names belong to
// responses.error and an entry has type and status instead. Both parse as
// nothing and were caught by chance.
//
// A reference that drifts is worse than none, so this pins it both ways. A
// field added to the format and not written down fails here, and so does a key
// this document invents.
func TestTheFormatReferenceNamesEveryKey(t *testing.T) {
	doc, err := os.ReadFile(filepath.Join("..", "..", "docs", "format.md"))
	if err != nil {
		t.Fatalf("read the format reference: %v", err)
	}

	documented := map[string]bool{}

	// Every key the reference names sits in the first column of a table, as
	// `key`, and a row may name two that share a description.
	rows := regexp.MustCompile("(?m)^\\| (`[a-z_]+`(?:, `[a-z_]+`)*) \\|")
	for _, row := range rows.FindAllStringSubmatch(string(doc), -1) {
		for _, name := range strings.Split(row[1], ", ") {
			documented[strings.Trim(name, "`")] = true
		}
	}

	if len(documented) == 0 {
		t.Fatal("no keys found in the reference at all, so this test is checking nothing")
	}

	declared := map[string]bool{}
	seen := map[reflect.Type]bool{}

	var walk func(t reflect.Type)

	walk = func(typ reflect.Type) {
		for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map {
			typ = typ.Elem()
		}

		if typ.Kind() != reflect.Struct || seen[typ] {
			return
		}

		seen[typ] = true

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			if tag := field.Tag.Get("yaml"); tag != "" && tag != "-" {
				if name := strings.Split(tag, ",")[0]; name != "" {
					declared[name] = true
				}
			}

			walk(field.Type)
		}
	}

	walk(reflect.TypeOf(Recipe{}))

	var missing, invented []string

	for name := range declared {
		if !documented[name] {
			missing = append(missing, name)
		}
	}

	for name := range documented {
		if !declared[name] {
			invented = append(invented, name)
		}
	}

	sort.Strings(missing)
	sort.Strings(invented)

	if len(missing) > 0 {
		t.Errorf("the format reference does not name these keys, which a Recipe may carry: %s", strings.Join(missing, ", "))
	}

	if len(invented) > 0 {
		t.Errorf("the format reference names these keys, which no Recipe may carry: %s", strings.Join(invented, ", "))
	}
}
