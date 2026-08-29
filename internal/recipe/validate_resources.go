// Resource rules: what a provider has, and whether this Recipe can address it.

package recipe

import (
	"strings"
)

// validateResources checks the resources a Recipe declares.
//
// Extracted from Validate, which was 991 lines of sequential rules in one
// function. Each of these appends through the same add, so the order and the
// wording of every problem are unchanged.
func (r *Recipe) validateResources(add func(string, ...any)) {
	if len(r.Resources) == 0 {
		add("at least one resource is required")
	}

	for _, name := range sortedKeys(r.Resources) {
		resource := r.Resources[name]

		if !contains(validIDStyles, resource.ID.Style) {
			add("resource %q has id.style %q, which must be one of %s", name, resource.ID.Style, strings.Join(validIDStyles[1:], ", "))
		}

		if resource.ID.CarriedBy != "" {
			if _, declared := resource.Fields[resource.ID.CarriedBy]; !declared {
				add("resource %q says its identifier is carried by %q, which is not a field on it", name, resource.ID.CarriedBy)
			}

			// Carried by something means not sent under its own name. A
			// resource that sends both is not describing this shape, and the
			// declaration would read as though it were.
			if resource.ID.Field != "-" {
				add("resource %q says its identifier is carried by %q and also sends it as %q; carried_by is for the providers that send it in one place only",
					name, resource.ID.CarriedBy, identifierFieldName(resource.ID.Field))
			}
		}

		if !contains(validIDTypes, resource.ID.Type) {
			add("resource %q has id.type %q, which must be one of %s", name, resource.ID.Type, strings.Join(validIDTypes[1:], ", "))
		}

		// A number on the wire has to be a number in the store too, or the
		// response disagrees with the declaration and the disagreement shows
		// up as a string where a number was promised. The styles that mint
		// something with letters in it cannot be numbers, and saying so here
		// is cheaper than finding out from a response.
		if resource.ID.Type == "number" && !contains([]string{"numeric", "digits"}, resource.ID.Style) {
			add("resource %q says its identifier is a number and mints it with the %s style, which does not produce one",
				name, styleName(resource.ID.Style))
		}

		// digits exists for identifiers that are numeric and must not be
		// parsed as numbers, because they exceed what a JavaScript number can
		// hold. Declaring one of those a number is asking for the rounding bug
		// the style was added to prevent.
		//
		// The length is what decides that, not the style. A double holds every
		// integer up to 2^53, which is sixteen figures, so fifteen is safe
		// with room to spare. Datadog's monitor ids are eight and really are
		// numbers -- its own description gives an integer for the id and for
		// the deleted_monitor_id a delete answers with -- and refusing that
		// forced the Recipe to send strings where Datadog sends numbers,
		// which is the disagreement this rule exists to catch, pointed the
		// other way.
		if resource.ID.Type == "number" && resource.ID.Style == "digits" && resource.ID.Length > safeDigits {
			add("resource %q mints a %d-digit numeric string and declares it a number, which is the rounding bug the digits style exists to avoid: a double holds %d figures",
				name, resource.ID.Length, safeDigits)
		}

		// A prefix is required for the prefixed style precisely because
		// applications prefix-match identifiers. Providers whose ids carry no
		// prefix say so by declaring the opaque style, rather than leaving it
		// ambiguous whether the omission was deliberate.
		if (resource.ID.Style == "" || resource.ID.Style == "prefixed") && resource.ID.Prefix == "" {
			add("resource %q must declare an id.prefix, or use id.style opaque if the provider's identifiers have none", name)
		}

		// A resource with no fields is almost always an unfinished one, and
		// the exception is narrow enough to name: a receipt. Knock's trigger
		// answers with a workflow_run_id and nothing else, and describing
		// that as a resource with one invented field would put a property on
		// the wire the provider never sends. A receipt is created and never
		// read back, so if anything reads this resource it is not one.
		if len(resource.Fields) == 0 && len(resource.Constants) == 0 && !onlyCreated(r, name) {
			add("resource %q has no fields", name)
		}

		for _, field := range sortedKeys(resource.Fields) {
			spec := resource.Fields[field]

			// A top-level field is named by its key, so an "as" that repeats
			// the key says nothing. One that differs is the only way to put a
			// name on the wire the record cannot use for itself, and there is
			// a real case: GitHub gives an issue a number and an id, the
			// number is what a path addresses it by, and the record therefore
			// has to key on the number -- which leaves "id" unavailable for
			// the field that carries GitHub's id.
			//
			// The rule was refusing every "as" here to catch the redundant
			// kind. It catches the redundant kind now.
			if spec.As != "" && spec.In == "" && spec.As == field {
				add("resource %q field %q sets as to its own name, which is what a top-level field is called anyway", name, field)
			}

			// Two declarations cannot claim one wire key. The comment above
			// works through the case where a record keys on something other
			// than "id" and a field takes "id" for itself -- GitHub's issue
			// number and id are the example -- and that reasoning only holds
			// while the identifier is not also being emitted under the same
			// name. When it is, the record has said twice what goes at that
			// key, the runtime picks one silently, and no conformance case can
			// see the loser: the wire is identical either way.
			//
			// The National Weather Service is the Recipe that found it. A
			// point is addressed at /points/39.7456,-97.0892 and its body
			// carries the whole URL under "id", so the Recipe keys on the
			// coordinate pair, hides the minted identifier, and renames a
			// field onto "id". Un-hiding the minted one -- which is a real
			// mistake, and the kind an edit makes -- changed nothing at all in
			// the response, which is exactly the shape of a claim nothing
			// checks.
			//
			// Only an explicit "as" counts. A field whose own name matches the
			// identifier's key is how a Recipe declares the identifier's type
			// and lines it up against an OpenAPI schema -- id.field: name
			// beside a declared name field, which three fixtures in
			// internal/openapi and a runtime test all do deliberately. "as"
			// has no such reading: it exists to move a field to a key it does
			// not already have, and that key is taken.
			//
			// A field carrying an "in" is emitted under that parent rather
			// than beside the identifier, so it is not competing either.
			//
			// Nothing else needs excluding. A hidden identifier is compared as
			// the literal "-", which no field is sent as, and a dotted one as
			// "sys.id", which no top-level field is called -- so both fall out
			// of the comparison rather than needing a guard, and a guard that
			// never changes an answer is a claim no test can check.
			if spec.As != "" && spec.In == "" {
				id := resource.ID.Field
				if id == "" {
					id = "id"
				}

				if spec.As == id {
					add("resource %q field %q is sent as %q and the identifier is emitted under %q too, so two declarations claim one key and one of them is silently dropped; hide the identifier with id.field \"-\" if the provider does not send it",
						name, field, id, id)
				}
			}

			// A field that never reaches the wire cannot be renamed on the
			// way there, and cannot be null there either.
			if spec.In == "-" {
				if spec.As != "" {
					add("resource %q field %q is not sent and also declares as, so it names a wire field that does not exist", name, field)
				}

				if spec.NullWhenUnset {
					add("resource %q field %q is not sent and also declares null_when_unset, which is a claim about a key nobody receives", name, field)
				}
			}

			if spec.NullWhenUnset && spec.Default != nil {
				add("resource %q field %q is null when unset and also has a default, so it is never unset", name, field)
			}

			if spec.NullWhenUnset && spec.Required {
				add("resource %q field %q is null when unset and also required, so it is never unset", name, field)
			}

			// Every segment of an in: path has to be something the runtime can
			// walk. A segment it cannot parse became a literal key: a field
			// nested under "to[0]" was emitted under a property actually
			// spelled "to[0]", which is a shape no provider sends, produced in
			// silence, and invisible to a conformance suite with no case
			// naming it. That has now happened three times in different
			// disguises, so it is a rule rather than a habit.
			for _, segment := range splitNonEmpty(spec.In) {
				// "-" is the whole value rather than a path: the field is not
				// sent, so there is no segment to be well formed.
				if spec.In == "-" {
					continue
				}

				if plainSegment.MatchString(segment) || indexedSegment.MatchString(segment) {
					continue
				}

				add("resource %q field %q nests under %q, and the segment %q is not a name or a name with an array index, so it would be sent as a property literally spelled that",
					name, field, spec.In, segment)
			}

			// A field named parent_thing nested under parent emits
			// parent.parent_thing, which no provider sends. It happened
			// thirty-one times across six Recipes before anything checked,
			// and every conformance case about them passed because they
			// asserted the shape the emulator produced. A provider that
			// really does repeat the parent says so with an explicit as.
			// A field named the same as its parent emits parent.parent, which
			// is usually an accident: the author meant the field to be the
			// parent object itself and wrote in: rather than leaving it off.
			// Four Recipes were written with that shape in one sitting.
			//
			// It is not always an accident. ClickUp really does send a status
			// object with a status inside it, so an explicit as: says the
			// provider means it, exactly as the sibling rule below allows.
			if spec.In != "" && spec.As == "" && strings.EqualFold(field, spec.In) {
				add("resource %q field %q nests under %q and would be sent as %s.%s; drop the in: if the field is the parent object, or set as: %s if the provider really sends that",
					name, field, spec.In, spec.In, field, field)
			}

			if spec.In != "" && spec.As == "" && strings.HasPrefix(strings.ToLower(field), strings.ToLower(spec.In)+"_") {
				add("resource %q field %q nests under %q and would be sent as %s.%s, which repeats the parent; set as: %s, or set as explicitly if the provider really sends that",
					name, field, spec.In, spec.In, field, field[len(spec.In)+1:])
			}

			if spec.Type == "list" && spec.Default != nil {
				if _, ok := spec.Default.([]any); !ok {
					add("resource %q field %q is a list, but its default is not a sequence", name, field)
				}
			}

			if !contains(validFieldTypes, spec.Type) {
				add("resource %q field %q has type %q, which must be one of %s",
					name, field, spec.Type, strings.Join(validFieldTypes[1:], ", "))
			}

			// Only time fields are ever filled in automatically, so declaring
			// stamped on anything else does nothing and reads as though it
			// does something.
			if spec.Stamped != nil && !contains(timeFieldTypes, spec.Type) {
				add("resource %q field %q declares stamped, which only applies to %s",
					name, field, strings.Join(timeFieldTypes, ", "))
			}
		}
	}
}
