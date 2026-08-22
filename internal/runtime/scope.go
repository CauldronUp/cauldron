// Scope: which records a request is allowed to see.

package runtime

import (
	"strconv"
)

// scopeVars extracts the scope filters for a request from its path parameters.
func (s *Sandbox) scopeVars(matched route, vars map[string]string) map[string]any {
	if len(matched.spec.Scope) == 0 {
		return nil
	}

	fields := s.recipe.Resources[matched.spec.Resource].Fields
	out := make(map[string]any, len(matched.spec.Scope))

	for _, name := range matched.spec.Scope {
		out[name] = scoped(fields[name].Type, vars[name])
	}

	return out
}

// scoped converts a path parameter to the type its field declares.
//
// A path is text and a field is not. Help Scout's conversationId, Gorgias's
// ticket_id, Shortcut's epic_id, Documenso's documentId and Hetzner's
// resources[0].id are all integers on the wire, and a record created through
// a scoped route carried the string out of the URL instead. So the same field
// was a number when a fixture seeded it and a string when a create produced
// it, in the same collection, in the same response shape, and only one of the
// two matched what the provider sends.
//
// Matching is done with fmt.Sprint on both sides, so a scope that changes
// type still finds the records it scoped.
//
// The JSON grammar decides, for the same reason it decides in a form body: a
// leading plus or a redundant leading zero is not something a JSON client
// could have sent, and an identifier that begins with one is an identifier
// rather than a number however the field is declared.
func scoped(declared, value string) any {
	switch declared {
	case "integer", "number", "timestamp", "timestamp_ms":
		if !jsonNumeric(value) {
			return value
		}

		if n, err := strconv.ParseInt(value, 10, 64); err == nil {
			return n
		}

		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}

	return value
}
