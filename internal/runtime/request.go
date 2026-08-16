package runtime

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/CauldronUp/cauldron/internal/store"
)

// maxBody caps how much of a request body the sandbox will read. A sandbox is
// not a hardened server, but it should not fall over because a test posted a
// gigabyte by mistake.
const maxBody = 8 << 20 // 8 MiB

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
	var raw map[string]any

	if err := json.Unmarshal(body, &raw); err != nil {
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

	if n, err := strconv.ParseInt(value, 10, 64); err == nil {
		return n
	}

	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	return value
}

// queryInt reads a positive integer query parameter, falling back to a default.
func queryInt(r *http.Request, key string, fallback int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}

	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
	}

	return n
}
