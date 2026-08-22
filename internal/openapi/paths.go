// Matching a Recipe's paths against a description's templates.

package openapi

import (
	"regexp"
	"strings"
)

// pathIndex matches a Recipe's paths against a description's templates.
type pathIndex struct {
	base      string
	templates []pathTemplate
}

type pathTemplate struct {
	template string
	pattern  *regexp.Regexp
}

type pathMatch struct {
	template string
}

func indexPaths(doc *Document, base string) pathIndex {
	index := pathIndex{base: strings.TrimRight(base, "/")}

	for _, path := range doc.SortedPaths() {
		index.templates = append(index.templates, pathTemplate{
			template: path,
			pattern:  templatePattern(index.base + path),
		})
	}

	return index
}

// templatePattern turns /v1/customers/{id} into a pattern that matches a path
// with any parameter names in it, since a Recipe names its parameters for
// itself and a description names them for its own reasons.
func templatePattern(path string) *regexp.Regexp {
	var b strings.Builder

	b.WriteString("^")

	for _, segment := range strings.Split(strings.Trim(path, "/"), "/") {
		b.WriteString("/")

		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			b.WriteString(`\{[^/}]*\}`)

			continue
		}

		// A segment that is partly literal and partly a parameter, which is
		// how Shopify writes /orders/{id}.json.
		var piece strings.Builder

		for _, part := range splitTemplate(segment) {
			if strings.HasPrefix(part, "{") {
				piece.WriteString(`\{[^/}]*\}`)

				continue
			}

			piece.WriteString(regexp.QuoteMeta(part))
		}

		b.WriteString(piece.String())
	}

	b.WriteString("/?$")

	return regexp.MustCompile(b.String())
}

// splitTemplate breaks a segment into literal and {parameter} pieces.
func splitTemplate(segment string) []string {
	var (
		out     []string
		current strings.Builder
	)

	for i := 0; i < len(segment); i++ {
		if segment[i] != '{' {
			current.WriteByte(segment[i])

			continue
		}

		close := strings.IndexByte(segment[i:], '}')
		if close < 0 {
			current.WriteByte(segment[i])

			continue
		}

		if current.Len() > 0 {
			out = append(out, current.String())
			current.Reset()
		}

		out = append(out, segment[i:i+close+1])
		i += close
	}

	if current.Len() > 0 {
		out = append(out, current.String())
	}

	return out
}

// find locates the template a Recipe's path is declared under, preferring one
// that declares the method being asked about.
//
// Two templates can differ only in what the parameter is called, and a
// description is free to split one path across both. Vercel does: the get on a
// deployment is declared under {idOrUrl} and the delete on the same path under
// {id}. Both match, and taking the first reported the delete as a method
// Vercel does not declare -- on a path where it is declared, one entry further
// down.
//
// The method is a preference rather than a requirement, so a path declared
// with none of the methods a Recipe uses still matches and still reports the
// method as undeclared, which is the finding that should survive.
func (p pathIndex) find(path, method string, doc *Document) (pathMatch, bool) {
	normalised := "/" + strings.Trim(path, "/")

	var first *pathMatch

	for _, candidate := range p.templates {
		if !candidate.pattern.MatchString(normalised) {
			continue
		}

		match := pathMatch{template: candidate.template}

		if operationFor(doc.Paths[candidate.template], method) != nil {
			return match, true
		}

		if first == nil {
			found := match
			first = &found
		}
	}

	if first != nil {
		return *first, true
	}

	return pathMatch{}, false
}

func operationFor(item PathItem, method string) *Operation {
	for _, mo := range item.Operations() {
		if mo.Method == strings.ToUpper(method) {
			return mo.Operation
		}
	}

	return nil
}

func methodsOf(item PathItem) []string {
	out := make([]string, 0, 7)
	for _, mo := range item.Operations() {
		out = append(out, mo.Method)
	}

	return out
}
