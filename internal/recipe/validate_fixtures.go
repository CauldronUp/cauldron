// Fixture and conformance-case rules: the seeded data and the evidence.

package recipe

import (
	"fmt"
	"regexp"
	"strings"
)

// validateFixtures checks the seeded records against the resources they fill.
//
// Extracted from Validate, which was 991 lines of sequential rules in one
// function. Each of these appends through the same add, so the order and the
// wording of every problem are unchanged.
func (r *Recipe) validateFixtures(add func(string, ...any)) {
	for _, fixtureName := range sortedKeys(r.Fixtures) {
		for _, resourceName := range sortedKeys(r.Fixtures[fixtureName]) {
			resource, ok := r.Resources[resourceName]
			if !ok {
				add("fixture %q seeds unknown resource %q", fixtureName, resourceName)
				continue
			}

			// A route that reads its identifier from the credentials answers
			// with the one record there is. Two of them and the runtime picks
			// the first, which is the seeding order, which is not a decision
			// anybody made.
			if len(r.Fixtures[fixtureName][resourceName]) > 1 {
				for _, route := range r.Routes {
					if route.Resource == resourceName && route.IDFrom == "auth" {
						add("fixture %q seeds %d records of %q, and %s %s answers about whoever asked, so there is nothing to choose between them with",
							fixtureName, len(r.Fixtures[fixtureName][resourceName]), resourceName, route.Method, route.Path)

						break
					}
				}
			}

			// A list field declared as one and then seeded with a scalar would
			// serve the scalar and pass every case written about it, so the
			// declaration has to mean something.
			//
			// An explicit null is allowed, because several providers send one
			// and it is a distinct state from both an absent field and an
			// empty array. AssemblyAI sends utterances: null on a transcript
			// that succeeded without speaker labels, so a rule that refused
			// null would block a true claim.
			for i, record := range r.Fixtures[fixtureName][resourceName] {
				// The identifier, which is not among the declared fields and
				// so is checked on its own.
				idField := resource.ID.Field
				if idField == "" || idField == "-" {
					idField = "id"
				}

				// An identifier the provider never echoes cannot disagree
				// with anything a client can see, so its shape is free.
				// Cohere keys embeddings e1 and e2 to find them again and
				// emits neither.
				if seeded, present := record[idField]; present && resource.ID.Field != "-" {
					if why, ok := seededID(resource.ID, seeded); !ok {
						add("fixture %q record %d of %q seeds an id that %s, so seeded records and created ones would not have the same shape",
							fixtureName, i, resourceName, why)
					}
				}

				for _, field := range sortedKeys(record) {
					// A key the resource does not declare is dropped on the
					// way in, so a fixture can set it, a reader can believe
					// it, and the emulator never sends it. That silence is
					// the problem: it fooled a mutation into passing, which
					// means it could fool a conformance case the same way.
					//
					// id is always allowed, and so is whatever the Recipe
					// renamed it to, because the identifier is not declared
					// among the fields.
					if _, declared := resource.Fields[field]; !declared {
						if _, constant := resource.Constants[field]; !constant &&
							field != "id" && field != resource.ID.Field {
							add("fixture %q record %d of %q sets %q, which is not a field on that resource, so it would be dropped in silence",
								fixtureName, i, resourceName, field)
						}
					}

					if resource.Fields[field].Type != "list" || record[field] == nil {
						if seeded, ok := seededAs(resource.Fields[field].Type, record[field]); !ok {
							add("fixture %q record %d of %q declares %q as %s and seeds %s, and nothing coerces it, so %s is what goes on the wire",
								fixtureName, i, resourceName, field, resource.Fields[field].Type, seeded, seeded)
						}

						continue
					}

					if _, isList := record[field].([]any); !isList {
						add("fixture %q record %d of %q sets list field %q to something that is neither a sequence nor null",
							fixtureName, i, resourceName, field)
					}
				}
			}
		}
	}
}

// validateCases checks each conformance case against the Recipe it belongs to.
//
// Extracted from Validate, which was 991 lines of sequential rules in one
// function. Each of these appends through the same add, so the order and the
// wording of every problem are unchanged.
func (r *Recipe) validateCases(add func(string, ...any)) {
	for i, c := range r.Conformance {
		where := fmt.Sprintf("conformance case %d", i+1)
		if c.Name != "" {
			where = fmt.Sprintf("conformance case %q", c.Name)
		} else {
			add("%s: name is required", where)
		}

		if c.Source == "" {
			add("%s: source is required, so a reader can check the claim against the provider", where)
		}

		if c.Verified != "" && !datePattern.MatchString(c.Verified) {
			add("%s: verified %q must be a date like 2026-08-15", where, c.Verified)
		}

		if c.Request.Method == "" {
			add("%s: request.method is required", where)
		} else if c.Request.Method != strings.ToUpper(c.Request.Method) {
			add("%s: request.method must be upper case", where)
		}

		if !strings.HasPrefix(c.Request.Path, "/") {
			add("%s: request.path must start with /", where)
		}

		if len(c.Request.Form) > 0 && len(c.Request.JSON) > 0 {
			add("%s: a request sends form or json, not both", where)
		}

		if c.Fixture != "" {
			if _, ok := r.Fixtures[c.Fixture]; !ok {
				add("%s: unknown fixture %q", where, c.Fixture)
			}
		}

		if c.Expect.Status == 0 {
			add("%s: expect.status is required", where)
		}

		for field, pattern := range c.Expect.Matches {
			if _, err := regexp.Compile(pattern); err != nil {
				add("%s: expect.matches[%s] is not a valid regular expression: %v", where, field, err)
			}
		}

		for header, pattern := range c.Expect.HeaderMatches {
			if _, err := regexp.Compile(pattern); err != nil {
				add("%s: expect.header_matches[%s] is not a valid regular expression: %v", where, header, err)
			}
		}

		if c.Expect.BodyMatches != "" {
			if _, err := regexp.Compile(c.Expect.BodyMatches); err != nil {
				add("%s: expect.body_matches is not a valid regular expression: %v", where, err)
			}
		}

		// An absence is a claim too: "another project's issues are not visible"
		// is exactly the kind of thing worth asserting, and it has no positive
		// half. Leaving it out of this check rejected real cases.
		// A case whose every body assertion repeats a value it sent proves
		// nothing about the emulator: it asserts its own request came back.
		// Four of these were written and only caught by deliberately breaking
		// the Recipe and noticing the case still passed, so the check is here
		// rather than in somebody's memory.
		//
		// One non-echoed claim is enough to clear it, because then something
		// the emulator decided is under test. matches and absent count, since
		// neither can be satisfied by echoing.
		if w := c.Expect.Webhook; w != nil {
			if w.None && (w.Event != "" || len(w.Body) > 0 || len(w.Matches) > 0 || len(w.Absent) > 0) {
				add("%s: expect.webhook.none claims nothing was emitted, so there is nothing for event, body, matches or absent to look at", where)
			}

			if !w.None && w.Event == "" && len(w.Body) == 0 && len(w.Matches) == 0 && len(w.Absent) == 0 {
				add("%s: expect.webhook asserts nothing, which is not the same as claiming none was emitted", where)
			}

			if w.Event != "" && !contains(r.Webhooks.Events, w.Event) {
				add("%s: expects webhook %q, which the Recipe does not declare", where, w.Event)
			}
		}

		if echoesOnly(c) {
			add("%s: every assertion repeats a value the request sent, so the case passes whatever the emulator does; assert something it decides, or add a matches or absent claim", where)
		}

		// A webhook claim counts. It is a claim about what the request caused,
		// which is evidence of exactly the kind this rule exists to demand.
		if c.Expect.Status < 400 && !c.Expect.NoBody && c.Expect.BodyMatches == "" && c.Expect.Webhook == nil &&
			len(c.Expect.Body) == 0 && len(c.Expect.Matches) == 0 &&
			len(c.Expect.Headers) == 0 && len(c.Expect.HeaderMatches) == 0 && len(c.Expect.Absent) == 0 &&
			len(c.Expect.AbsentHeaders) == 0 {
			add("%s: a case that asserts nothing about the response is not evidence of anything", where)
		}
	}
}
