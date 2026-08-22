// Reaching the schema a Recipe is talking about, through whatever the
// description wrapped it in.

package openapi

import (
	"sort"
	"strconv"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// collectionSchema digs the array of objects out of a list envelope.
func collectionSchema(doc *Document, envelope *Schema, keys ...string) *Schema {
	if envelope.Type == "array" && envelope.Items != nil {
		return doc.Resolve(envelope.Items)
	}

	for _, key := range keys {
		if key == "" {
			continue
		}

		if resolved := descend(doc, envelope, key); resolved != nil {
			if resolved.Type == "array" && resolved.Items != nil {
				return doc.Resolve(resolved.Items)
			}

			return resolved
		}
	}

	// Nothing named. Fall back to the first array-of-objects property, which
	// is what a description with one collection in its envelope looks like.
	for _, property := range doc.Properties(envelope) {
		resolved := doc.Resolve(property.Schema)
		if resolved != nil && resolved.Type == "array" && resolved.Items != nil {
			return doc.Resolve(resolved.Items)
		}
	}

	return nil
}

// statusBlock reads a range like 4XX and answers the hundreds it covers.
//
// "default" is deliberately not a range. It means every status the operation
// did not name, which is a description saying it has stopped enumerating
// rather than one saying a particular failure is expected, and treating it as
// coverage would silence this check on every description that writes it.
func statusBlock(code string) (int, bool) {
	if len(code) != 3 {
		return 0, false
	}

	if !strings.EqualFold(code[1:], "XX") {
		return 0, false
	}

	n, err := strconv.Atoi(code[:1])
	if err != nil || n < 1 || n > 5 {
		return 0, false
	}

	return n, true
}

// descend walks a dotted key through a schema's properties, resolving
// references as it goes, the same way the runtime nests a dotted name.
func descend(doc *Document, schema *Schema, key string) *Schema {
	current := schema

	for _, segment := range strings.Split(key, ".") {
		current = doc.Resolve(current)
		if current == nil {
			return nil
		}

		// Properties rather than the map, because a description may compose
		// a response out of allOf members and declare nothing on the object
		// itself. Neon describes a branch that way -- the 200 has two
		// members, one carrying the branch and one carrying an annotation --
		// and reading the map walked into an object that looked empty, gave
		// up, and left every field compared against the envelope instead.
		var next *Schema

		for _, property := range doc.Properties(current) {
			if property.Name == segment {
				next = property.Schema
				break
			}
		}

		if next == nil {
			return nil
		}

		current = next
	}

	return doc.Resolve(current)
}

// resourceSchema descends a response envelope to the object one resource
// actually arrives in.
//
// Only listings were unwrapped, so a Recipe whose single-resource responses
// are wrapped had every field of every resource reported as undeclared: the
// comparison ran against the envelope, whose properties are things like
// result, status and time, and no resource has fields called those. Cloudflare
// wraps under result, Xero wraps under a plural and puts one object inside an
// array, and Qdrant wraps under result and then again under a name that
// differs per endpoint.
//
// A report that cries wolf is one nobody reads the second time, which is the
// same reason deletes are skipped.
// spec is the route's envelope, which is the Recipe's unless the route says
// otherwise. Reading the Recipe's here reported every field of a resource
// whose route wraps differently: Vercel wraps a single domain and wraps
// neither a project nor a deployment, so the domain fetch declares its own,
// and comparing it against the Recipe's found no wrapper and gave up.
func resourceSchema(doc *Document, r *recipe.Recipe, spec recipe.ResourceResponse, envelope *Schema, name string) *Schema {
	if spec.Style != "wrapped" {
		return envelope
	}

	key := spec.Key
	if key == "" {
		key = name

		// A list of one is still a collection, so the runtime defaults an
		// array-wrapped resource to the collection name rather than the
		// resource name: Xero answers a single invoice with Invoices, not
		// invoice. Reading the resource name here found no such key, fell
		// through to the envelope, and reported 61 disagreements against a
		// Recipe that was right about all of them.
		if spec.Array {
			if spec, ok := r.Resources[name]; ok && spec.Collection != "" {
				key = spec.Collection
			}
		}
	}

	found := descend(doc, envelope, key)
	if found == nil {
		// The description does not nest where the Recipe says it does. That
		// is worth reporting one day and is not this function's to report,
		// and comparing against the envelope instead would turn it into one
		// finding per field.
		return envelope
	}

	if found.Type == "array" && found.Items != nil {
		return doc.Resolve(found.Items)
	}

	return found
}

func sortedFieldNames(fields map[string]recipe.Field) []string {
	out := make([]string, 0, len(fields))
	for name := range fields {
		out = append(out, name)
	}

	sort.Strings(out)

	return out
}

func sortedErrorNames(errors map[string]recipe.Error) []string {
	out := make([]string, 0, len(errors))
	for name := range errors {
		out = append(out, name)
	}

	sort.Strings(out)

	return out
}

// declaredNames collects every property name a schema or any of its
// alternatives declares.
//
// Lenient on purpose: this decides whether to report a field as undeclared,
// and a report that is wrong costs more than one that is missing. A name in
// one branch of a union is a name the description declares.
func declaredNames(doc *Document, schema *Schema, seen map[*Schema]bool) map[string]bool {
	names := map[string]bool{}

	schema = doc.Resolve(schema)
	if schema == nil || seen[schema] {
		return names
	}

	seen[schema] = true

	for name := range schema.Properties {
		names[name] = true
	}

	for _, group := range [][]*Schema{schema.AllOf, schema.OneOf, schema.AnyOf} {
		for _, member := range group {
			for name := range declaredNames(doc, member, seen) {
				names[name] = true
			}
		}
	}

	return names
}
