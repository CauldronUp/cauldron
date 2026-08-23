package runtime

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/store"
)

// maxBody caps how much of a request body the sandbox will read. A sandbox is
// not a hardened server, but it should not fall over because a test posted a
// gigabyte by mistake.
const maxBody = 8 << 20 // 8 MiB

// graphQLQuery reads the query out of a GraphQL request so routing can tell
// several routes on one path apart.
//
// It reads the body and puts it back, because the handlers downstream read it
// again, and it answers with an empty string for anything that is not a
// GraphQL-shaped POST. That is the common case by far: 158 of the 159 Recipes
// here route on the path alone and must keep doing so.
func graphQLQuery(r *http.Request) string {
	if r.Method != http.MethodPost || r.Body == nil {
		return ""
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		return ""
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	var envelope struct {
		Query string `json:"query"`
	}

	if err := json.Unmarshal(body, &envelope); err != nil {
		return ""
	}

	return envelope.Query
}

// decodeBody reads a request body as either JSON or form encoding.
//
// Supporting both is a fidelity requirement, not a convenience: Stripe's
// official SDKs post application/x-www-form-urlencoded, so a fake that only
// spoke JSON would reject the very clients it exists to serve.
func decodeBody(r *http.Request) (store.Record, error) {
	if r.Body == nil {
		return store.Record{}, nil
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		return nil, err
	}

	// Put it back. A route whose identifier comes from the body reads it once
	// to find the id and again to apply the changes, and a drained body would
	// make the second read silently return nothing.
	r.Body = io.NopCloser(bytes.NewReader(body))

	if len(body) == 0 {
		return store.Record{}, nil
	}

	contentType := r.Header.Get("Content-Type")

	switch {
	case strings.Contains(contentType, "application/json"):
		return decodeJSON(body)
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		return decodeForm(body)
	default:
		// No usable content type. Try JSON, then form — a client that sends a
		// body without a type is still trying to say something.
		if record, err := decodeJSON(body); err == nil {
			return record, nil
		}

		return decodeForm(body)
	}
}

func decodeJSON(body []byte) (store.Record, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))

	// The same reason jsonBody and conform decode this way, and this was the
	// one path that did not: json.Unmarshal turns every number into a
	// float64, which holds 53 bits of integer. A Discord snowflake, a Stripe
	// cursor, a ledger amount in minor units above 2^53 -- anything posted
	// through a create or an update came back changed.
	//
	// 578730123365711993 in, 578730123365712000 out, with a 200 and a
	// plausible body. Inconsistently, too: a YAML fixture decodes to an int
	// and stays exact, so the same field was right when seeded and wrong when
	// written, in the same collection.
	//
	// conform's own comment on this says it best -- the tool that exists to
	// catch an emulator sending the wrong number was itself unable to see the
	// right one.
	decoder.UseNumber()

	var raw map[string]any

	if err := decoder.Decode(&raw); err != nil {
		return nil, err
	}

	return store.Record(raw), nil
}

// decodeForm parses form encoding, including the bracketed nesting providers
// use for sub-objects: metadata[order_id]=123.
func decodeForm(body []byte) (store.Record, error) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	record := store.Record{}

	for key, list := range values {
		if len(list) == 0 {
			continue
		}

		value := coerce(list[0])

		open := strings.Index(key, "[")
		if open <= 0 || !strings.HasSuffix(key, "]") {
			record[key] = value
			continue
		}

		parent := key[:open]
		child := strings.TrimSuffix(key[open+1:], "]")

		nested, ok := record[parent].(map[string]any)
		if !ok {
			nested = map[string]any{}
			record[parent] = nested
		}

		nested[child] = value
	}

	return record, nil
}

// coerce turns form strings into the types a JSON client would have sent, so
// that a form-encoded request and a JSON request produce identical records.
func coerce(value string) any {
	switch value {
	case "true":
		return true
	case "false":
		return false
	}

	if !jsonNumeric(value) {
		return value
	}

	if n, err := strconv.ParseInt(value, 10, 64); err == nil {
		return n
	}

	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	return value
}

// jsonNumeric reports whether value is a JSON number literal.
//
// strconv is more generous than JSON in three ways, and each one turns a
// string a client posted into a number the provider never sends.
//
// It accepts a leading plus, so "+15017122661" becomes the integer
// 15017122661 and the phone number loses the one character that made it
// E.164 rather than a number somebody has to guess the country of. Every
// provider taking a destination in a form body takes it in that format.
//
// It accepts redundant leading zeros, so an account code of "007" comes back
// as 7 and a sort code beginning 04 comes back five digits long.
//
// And it accepts NaN, Inf, Infinity and +Inf, case-insensitively. Storing one
// meant the record could never be encoded again: encoding/json refuses it, so
// every later read of that collection answered with an empty body until
// something reset the sandbox.
//
// None of the three is a value a JSON client could have sent, and matching a
// JSON client is the whole job of coerce. json.Valid applies the real
// grammar. The first-byte check is what separates a number from the other
// things that grammar accepts, since "null" and a quoted string are both
// valid JSON and neither is a number.
func jsonNumeric(value string) bool {
	if value == "" {
		return false
	}

	if c := value[0]; c != '-' && (c < '0' || c > '9') {
		return false
	}

	return json.Valid([]byte(value))
}

// paging reads a pagination parameter from wherever the provider carries it.
//
// A listing reached by POST usually takes its paging in the JSON body, and
// reading the query string for it finds nothing. Nothing is the dangerous
// answer here: the limit falls back to the route's default, the default is
// larger than any fixture, so the first response holds the whole collection
// and reports no next page. The paging loop a client wrote runs once, takes
// neither branch, and passes.
type paging struct {
	query url.Values
	body  map[string]any
}

func pagingFrom(r *http.Request, spec recipe.Pagination) paging {
	if spec.In != "body" {
		return paging{query: r.URL.Query()}
	}

	return paging{body: jsonBody(r)}
}

// get returns a parameter as text, whichever JSON type carried it. A dotted
// name nests: Plaid keeps count and offset under options.
func (p paging) get(name string) string {
	if name == "" {
		return ""
	}

	if p.body == nil {
		return p.query.Get(name)
	}

	var current any = p.body

	for _, segment := range strings.Split(name, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return ""
		}

		current, ok = object[segment]
		if !ok {
			return ""
		}
	}

	switch value := current.(type) {
	case string:
		return value

	case json.Number:
		return value.String()

	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)

	case bool:
		return strconv.FormatBool(value)
	}

	return ""
}

// int reads a parameter as a positive count, falling back when it is absent
// or unusable.
func (p paging) int(name string, fallback int) int {
	raw := p.get(name)
	if raw == "" {
		return fallback
	}

	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
	}

	return n
}

// jsonBody parses the request body and puts it back, so that everything
// downstream still reads a full body.
func jsonBody(r *http.Request) map[string]any {
	if r.Body == nil {
		return map[string]any{}
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		return map[string]any{}
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	if len(body) == 0 {
		return map[string]any{}
	}

	decoder := json.NewDecoder(bytes.NewReader(body))

	// The same reason conform decodes this way: a cursor can be a long
	// number, and float64 loses the end of one.
	decoder.UseNumber()

	var out map[string]any

	if err := decoder.Decode(&out); err != nil {
		return map[string]any{}
	}

	return out
}
