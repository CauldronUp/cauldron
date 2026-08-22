// Failures: choosing which error a Recipe declares and writing it out.

package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/store"
)

// notFound answers for a record that is not there.
//
// The route may name the error rather than take resource_missing, because a
// provider can answer a missing thing more than one way depending on what was
// asked for. An unrecognised path is not that: it is not this route's absence
// to describe, so the override does not reach it.
func (s *Sandbox) notFound(w http.ResponseWriter, err error, resource, id, override string) int {
	if errors.Is(err, store.ErrUnknownResource) {
		return s.writeRecipeError(w, "unknown_route", 404, "unknown_route", "Unrecognised request URL.", resource)
	}

	if override == "" {
		return s.writeRecipeError(
			w, "resource_missing", 404, "resource_missing",
			"No such "+resource+": "+id+".", resource+": "+id,
		)
	}

	// The identifier alone, rather than "resource: id". A route names its own
	// error because this provider describes this absence in its own words, and
	// the composite is the default message's convention rather than something
	// the provider would say. npm's is "version not found: 99.99.99", and
	// "version not found: release: 99.99.99" is not a sentence anybody sends.
	return s.writeRecipeError(
		w, override, 404, override,
		"No such "+resource+": "+id+".", id,
	)
}

// writeRecipeError writes an error, preferring the Recipe's own definition so
// that status codes, codes and headers match what the real provider sends.
// Fallbacks apply only when the Recipe does not define the situation.
//
// A fourth argument is the request-specific detail, such as the parameter that
// was missing. A Recipe reaches it through {detail} in its message.
// builtinFallbacks maps the failures the runtime produces itself onto the
// declared error whose shape they should borrow.
//
// An unrecognised path, a method a route does not take and a body that will
// not parse are all produced by the handler rather than declared by a Recipe,
// and they used to carry their own internal name as both the code and the
// category. No provider has an error type called "unknown_route". Every Stripe
// 404 said type "unknown_route" where the real one says
// "invalid_request_error", and no Recipe declared any of these three, so this
// was live on all 97 of them, on the most commonly exercised error path there
// is. Retry and branching logic keyed on the category took one path locally
// and another in production.
//
// A Recipe may still declare any of these names itself, and that wins.
var builtinFallbacks = map[string]string{
	"unknown_route":      "resource_missing",
	"method_not_allowed": "resource_missing",
	"conflict":           "invalid_request",
	"invalid_request":    "parameter_missing",
}

// resolveBuiltin follows the fallback chain until it reaches a name the Recipe
// declares, and returns the original when it reaches none.
func resolveBuiltin(r *recipe.Recipe, name string) string {
	seen := map[string]bool{}

	for current := name; !seen[current]; {
		if _, declared := r.Errors[current]; declared {
			return current
		}

		seen[current] = true

		next, ok := builtinFallbacks[current]
		if !ok {
			break
		}

		current = next
	}

	return name
}

func (s *Sandbox) writeRecipeError(w http.ResponseWriter, name string, fallback ...any) int {
	status := http.StatusInternalServerError
	code := name
	message := "An error occurred."

	if len(fallback) >= 3 {
		status, _ = fallback[0].(int)
		code, _ = fallback[1].(string)
		message, _ = fallback[2].(string)
	}

	var detail string

	if len(fallback) >= 4 {
		detail, _ = fallback[3].(string)
	}

	// Providers categorise errors more coarsely than they code them, and client
	// libraries switch on the category. Using the name as the type made every
	// error its own category, which no real provider does.
	category := name

	// Extra body properties this particular failure carries.
	var extra map[string]any

	resolved := resolveBuiltin(s.recipe, name)

	// A Recipe that declares this failure by name owns all of it. One that
	// only lends its shape to a built-in lends the category, the code and the
	// extra fields, and not the status: borrowing a 404's status for a 405
	// would change what the failure is, which is a worse lie than the invented
	// category this replaced.
	borrowed := resolved != name

	// The envelope this failure travels in, which a single error may override.
	var style string

	if defined, ok := s.recipe.Errors[resolved]; ok {
		if !borrowed {
			status = defined.Status
		}

		if defined.Code != "" {
			code = defined.Code
		}

		if defined.Type != "" {
			category = defined.Type
		}

		// The Recipe's wording wins. Providers differ here in ways that matter:
		// Stripe names the offending parameter, GitHub says "Not Found" and
		// nothing else. A Recipe that wants the request-specific detail asks
		// for it with {detail}, so the choice is the Recipe author's rather
		// than a hardcoded guess in the handler.
		// The wording travels only when the lender describes the same status.
		// GitHub's 404 really does say "Not Found" and nothing else, and that
		// is worth having; its wording on a 405 would be a sentence about the
		// wrong thing.
		if defined.Message != "" && (!borrowed || defined.Status == status) {
			message = strings.ReplaceAll(defined.Message, "{detail}", detail)
		}

		for key, value := range defined.Headers {
			w.Header().Set(key, value)
		}

		extra = defined.Fields
		style = defined.Style
	}

	if style == "" {
		style = s.recipe.Responses.Error.Style
	}

	// Trello answers with plain text, so a client calling .json() on the
	// response throws rather than reporting the failure. Writing JSON here
	// would hide that.
	if style == "text" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(status)

		fmt.Fprint(w, message)

		return status
	}

	// A bare JSON string, which is valid JSON and is not an object. The npm
	// registry answers one for a version that does not exist on a package
	// that does, while answering an object for a package that does not exist.
	//
	// Text would be the wrong shape for it: a client calling .json() on text
	// throws, and on this it succeeds and hands back a string. Code reading
	// body.error then finds undefined rather than failing, which is the
	// quieter and worse of the two.
	if style == "string" {
		writeJSON(w, status, message)

		return status
	}

	writeJSON(w, status, s.errorBody(category, code, message, status, extra))

	return status
}

// errorBody shapes a failure according to the Recipe's declared error style.
//
// The return type is any rather than a map because Salesforce's failures are a
// bare top-level array with no envelope at all. A client reading .message off
// the response finds undefined, and has to index before it can read anything.
func (s *Sandbox) errorBody(category, code, message string, status int, extra map[string]any) any {
	spec := s.recipe.Responses.Error

	if spec.Style == "string_list" {
		key := spec.Key
		if key == "" {
			key = "errors"
		}

		// Datadog sends the array with bare strings in it. A client looping
		// over the entries and reading .message from each finds undefined on
		// every one, which throws rather than reporting anything.
		return withFields(withFields(map[string]any{key: []any{message}}, spec.Fields), extra)
	}

	if spec.Style != "flat" && spec.Style != "list" {
		// The sub-keys are named the same way the flat style names its
		// top-level ones. Google nests a numeric code beside a string status
		// and calls neither of them "type", so the names have to be
		// declarable rather than assumed. An empty name keeps the old default,
		// which is what every nested Recipe written before this relied on.
		nested := map[string]any{}

		if field := spec.MessageField; field != "-" {
			if field == "" {
				field = "message"
			}

			nested[listName(field)] = messageValue(field, message)
		}

		// Providers disagree about which of these they send. Stripe sends a
		// type and a code, Airtable only a type, Vercel only a code. Declaring
		// the omission keeps the emulator from inventing a field a client
		// could come to depend on and then lose.
		if field := spec.CodeField; field != "-" {
			if field == "" {
				field = "code"
			}

			// The same conversion the flat style applies. PagerDuty's codes
			// really are numbers, and sending "2006" as text meant a client
			// switching on the code never matched.
			nested[field] = codeValue(spec.CodeType, code)
		}

		if field := spec.TypeField; field != "-" {
			if field == "" {
				field = "type"
			}

			nested[field] = category
		}

		if spec.StatusField != "" {
			nested[spec.StatusField] = status
		}

		key := spec.Key
		if key == "" {
			key = "error"
		}

		// Declared fields apply here too. They were dropped for the nested
		// style alone, so a Recipe could add a constant to every error and
		// have it silently ignored.
		return withFields(withFields(map[string]any{key: nested}, spec.Fields), extra)
	}

	body := map[string]any{}

	// "-" means the provider sends no prose at all. Slack's errors are a code
	// and nothing else, and inventing a sentence it never sends is infidelity
	// in the direction that is hardest to notice.
	if field := spec.MessageField; field != "-" {
		if field == "" {
			field = "message"
		}

		// setPath, so a dotted name nests the same way code_field and
		// status_field already do. Front puts all three under _error.
		setPath(body, listName(field), messageValue(field, message))
	}

	// "-" omits the field, exactly as it does for the message above and in the
	// nested style. It used to be honoured only in the nested style here, so a
	// flat Recipe declaring code_field: "-" or type_field: "-" got a literal
	// "-" key in every error body instead of no field at all. Three shipped
	// Recipes were doing that, and every case written about them passed,
	// because a case asserts the fields it names and says nothing about a key
	// nobody thought to look for.
	if spec.CodeField != "" && spec.CodeField != "-" && code != "" {
		setPath(body, spec.CodeField, codeValue(spec.CodeType, code))
	}

	if spec.TypeField != "" && spec.TypeField != "-" {
		setPath(body, spec.TypeField, category)
	}

	if spec.StatusField != "" {
		setPath(body, spec.StatusField, status)
	}

	if spec.Style == "list" {
		key := spec.Key
		if key == "" {
			key = "errors"
		}

		// "-" means there is no envelope: the array is the whole body, which
		// is what Salesforce sends. Reading .message off that response finds
		// undefined, and any declared fields have nowhere to go, so a Recipe
		// that wants both has to choose the shape the provider actually uses.
		if key == "-" {
			return []any{withFields(body, extra)}
		}

		// One request can fail several ways at once, which is why SendGrid and
		// Clerk both send an array. The emulator reports one, in the shape a
		// client already has to loop over.
		//
		// Declared fields sit beside the array rather than inside each entry:
		// Clerk's trace id belongs to the response, not to one failure. A
		// dotted key nests, which QuickBooks needs for Fault.Error.
		envelope := map[string]any{}
		setPath(envelope, key, []any{body})

		return withFields(withFields(envelope, spec.Fields), extra)
	}

	return withFields(withFields(body, spec.Fields), extra)
}

// messageValue wraps the message in an array when the Recipe asks for one.
//
// Mux sends messages, plural: error.message is undefined and
// error.messages[0] is the sentence, which is backwards from every other
// provider in the collection. Writing the sentence into a field called
// "messages" as a bare string would look right in a diff and be wrong on the
// wire, so the Recipe says which it is.
func messageValue(field, message string) any {
	if strings.HasSuffix(field, "[]") {
		return []any{message}
	}

	return message
}

// codeValue renders an error code as the Recipe says the provider sends it.
//
// An undeclared code_type falls back to inferring from the value, which is
// right for Twilio and wrong for Adyen: "000" is a string there, and turning it
// into a number destroys the leading zeros a client is matching on. The
// inference stays the default so that Recipes written before this existed keep
// working, but it is a guess, and a Recipe that knows better says so.
func codeValue(codeType, value string) any {
	if codeType == "string" {
		return value
	}

	return numberOrString(value)
}
