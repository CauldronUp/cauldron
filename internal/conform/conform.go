// Package conform runs a Recipe's conformance cases against the emulator.
//
// The uncomfortable question about any fake is "how do you know it behaves like
// the real thing?", and the honest answer for most of them is "you don't". A
// conformance suite is Cauldron's attempt at a better answer: each case is a
// checkable claim with a citation, and the report separates what was observed
// against the real API from what was only read in the documentation.
//
// Cases run in process against the server's own handler. That is deliberate:
// verification must work in CI, on a fork, and with no credentials, or nobody
// will run it and the evidence goes stale.
package conform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// Result is the outcome of one case.
type Result struct {
	Case     recipe.Case
	Failures []string
}

// Passed reports whether every claim in the case held.
func (r Result) Passed() bool { return len(r.Failures) == 0 }

// Report is the outcome for one Recipe.
type Report struct {
	Recipe  string
	Version string
	Results []Result
}

// Passed counts the cases that held.
func (r Report) Passed() int {
	var n int

	for _, result := range r.Results {
		if result.Passed() {
			n++
		}
	}

	return n
}

// Failed returns the cases that did not hold.
func (r Report) Failed() []Result {
	var out []Result

	for _, result := range r.Results {
		if !result.Passed() {
			out = append(out, result)
		}
	}

	return out
}

// Provenance summarises where the evidence came from and how fresh it is.
func (r Report) Provenance() (observed, documented int, latest string) {
	for _, result := range r.Results {
		if result.Case.Verified == "" {
			documented++
			continue
		}

		observed++

		if result.Case.Verified > latest {
			latest = result.Case.Verified
		}
	}

	return observed, documented, latest
}

// Seeder loads a fixture between cases. The server supplies this; conform does
// not reach into the runtime itself.
type Seeder func(fixture string) error

// Run executes every case in a Recipe against handler, mounted at prefix.
func Run(r *recipe.Recipe, handler http.Handler, prefix string, seed Seeder) Report {
	report := Report{Recipe: r.Name, Version: r.Version}

	for _, c := range r.Conformance {
		report.Results = append(report.Results, runCase(c, handler, prefix, seed))
	}

	return report
}

func runCase(c recipe.Case, handler http.Handler, prefix string, seed Seeder) Result {
	result := Result{Case: c}

	if c.Fixture != "" && seed != nil {
		if err := seed(c.Fixture); err != nil {
			result.Failures = append(result.Failures, fmt.Sprintf("could not seed %q: %v", c.Fixture, err))
			return result
		}
	}

	req, err := buildRequest(c.Request, prefix)
	if err != nil {
		result.Failures = append(result.Failures, err.Error())
		return result
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	response := rec.Result()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		result.Failures = append(result.Failures, fmt.Sprintf("could not read the response: %v", err))
		return result
	}

	result.Failures = check(c.Expect, response, body)

	return result
}

func buildRequest(r recipe.Request, prefix string) (*http.Request, error) {
	target := strings.TrimRight(prefix, "/") + r.Path

	if len(r.Query) > 0 {
		values := url.Values{}

		for _, key := range sortedKeys(r.Query) {
			values.Set(key, r.Query[key])
		}

		target += "?" + values.Encode()
	}

	var (
		req *http.Request
		err error
	)

	switch {
	case len(r.Form) > 0:
		values := url.Values{}

		for _, key := range sortedKeys(r.Form) {
			values.Set(key, r.Form[key])
		}

		req, err = http.NewRequest(r.Method, "http://sandbox"+target, strings.NewReader(values.Encode()))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	case len(r.JSON) > 0:
		encoded, marshalErr := json.Marshal(r.JSON)
		if marshalErr != nil {
			return nil, marshalErr
		}

		req, err = http.NewRequest(r.Method, "http://sandbox"+target, strings.NewReader(string(encoded)))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	default:
		req, err = http.NewRequest(r.Method, "http://sandbox"+target, nil)
	}

	if err != nil {
		return nil, err
	}

	for _, key := range sortedKeys(r.Headers) {
		req.Header.Set(key, r.Headers[key])
	}

	return req, nil
}

func check(expect recipe.Expectation, response *http.Response, body []byte) []string {
	var failures []string

	if expect.Status != 0 && response.StatusCode != expect.Status {
		failures = append(failures, fmt.Sprintf("status: want %d, got %d", expect.Status, response.StatusCode))
	}

	for _, name := range sortedKeys(expect.Headers) {
		want := expect.Headers[name]

		got := response.Header.Get(name)
		if !strings.Contains(got, want) {
			failures = append(failures, fmt.Sprintf("header %s: want %q, got %q", name, want, got))
		}
	}

	for _, name := range sortedKeys(expect.HeaderMatches) {
		pattern := expect.HeaderMatches[name]

		got := response.Header.Get(name)
		if got == "" {
			failures = append(failures, fmt.Sprintf("header %s: expected a value matching %s, but the header is absent", name, pattern))
			continue
		}

		compiled, err := regexp.Compile(pattern)
		if err != nil {
			failures = append(failures, fmt.Sprintf("header %s: %v", name, err))
			continue
		}

		if !compiled.MatchString(got) {
			failures = append(failures, fmt.Sprintf("header %s: %q does not match %s", name, got, pattern))
		}
	}

	if len(expect.Body) == 0 && len(expect.Matches) == 0 && len(expect.Absent) == 0 {
		return failures
	}

	// An empty body is a real answer, not a parse failure: SendGrid accepts a
	// send with 202 and nothing else. Every field is absent, so absence
	// assertions hold and any positive claim fails with a useful reason.
	if len(bytes.TrimSpace(body)) == 0 {
		for _, path := range sortedKeys(expect.Body) {
			failures = append(failures, fmt.Sprintf("%s: the response has no body", path))
		}

		for _, path := range sortedKeys(expect.Matches) {
			failures = append(failures, fmt.Sprintf("%s: the response has no body", path))
		}

		return failures
	}

	var decoded any

	if err := json.Unmarshal(body, &decoded); err != nil {
		return append(failures, fmt.Sprintf("response is not JSON: %v", err))
	}

	for _, path := range sortedKeys(expect.Body) {
		failures = append(failures, compare(path, expect.Body[path], decoded)...)
	}

	for _, path := range sortedKeys(expect.Matches) {
		pattern := expect.Matches[path]

		value, ok := lookup(decoded, path)
		if !ok {
			failures = append(failures, fmt.Sprintf("%s: expected a value matching %s, but the field is absent", path, pattern))
			continue
		}

		compiled, err := regexp.Compile(pattern)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", path, err))
			continue
		}

		if text := text(value); !compiled.MatchString(text) {
			failures = append(failures, fmt.Sprintf("%s: %q does not match %s", path, text, pattern))
		}
	}

	for _, path := range expect.Absent {
		if _, ok := lookup(decoded, path); ok {
			failures = append(failures, fmt.Sprintf("%s: the provider does not send this field, but the emulator did", path))
		}
	}

	return failures
}

// compare walks an expected value against the response. Maps are compared as
// subsets so a case asserts what it claims and stays silent about the rest;
// lists are compared element by element, because a case that says "two items"
// is making a claim about the length.
func compare(path string, want, document any) []string {
	got, ok := lookup(document, path)
	if !ok {
		return []string{fmt.Sprintf("%s: missing", path)}
	}

	return compareValue(path, want, got)
}

func compareValue(path string, want, got any) []string {
	switch expected := want.(type) {
	case map[string]any:
		actual, ok := got.(map[string]any)
		if !ok {
			return []string{fmt.Sprintf("%s: want an object, got %s", path, describe(got))}
		}

		var failures []string

		for _, key := range sortedKeys(expected) {
			value, present := actual[key]
			if !present {
				failures = append(failures, fmt.Sprintf("%s.%s: missing", path, key))
				continue
			}

			failures = append(failures, compareValue(path+"."+key, expected[key], value)...)
		}

		return failures

	case []any:
		actual, ok := got.([]any)
		if !ok {
			return []string{fmt.Sprintf("%s: want a list, got %s", path, describe(got))}
		}

		if len(actual) != len(expected) {
			return []string{fmt.Sprintf("%s: want %d item(s), got %d", path, len(expected), len(actual))}
		}

		var failures []string

		for i := range expected {
			failures = append(failures, compareValue(fmt.Sprintf("%s[%d]", path, i), expected[i], actual[i])...)
		}

		return failures

	default:
		// YAML and JSON disagree about number types (int64 against float64),
		// and a case author should not have to care. Compare the rendered form.
		if text(want) != text(got) {
			return []string{fmt.Sprintf("%s: want %s, got %s", path, text(want), text(got))}
		}

		return nil
	}
}

// splitFieldPath splits a dotted path into segments, honouring a backslash
// escape so a field name may contain a dot of its own. Dropbox names a field
// ".tag", where the leading dot is part of the name rather than a separator,
// and without an escape there is no way to assert on it at all.
func splitFieldPath(path string) []string {
	var (
		segments []string
		current  strings.Builder
		escaped  bool
	)

	for _, r := range path {
		switch {
		case escaped:
			current.WriteRune(r)

			escaped = false
		// 0x5C is a backslash, written numerically so the escape character
		// itself does not need escaping here. It escapes the next character,
		// letting a dot be part of a field name rather than a separator.
		case r == 0x5C:
			escaped = true
		case r == '.':
			segments = append(segments, current.String())

			current.Reset()
		default:
			current.WriteRune(r)
		}
	}

	return append(segments, current.String())
}

// lookup resolves a dotted path with optional [n] indexes, e.g. data[0].id.
func lookup(document any, path string) (any, bool) {
	current := document

	for _, segment := range splitFieldPath(path) {
		name, indexes := splitIndexes(segment)

		if name != "" {
			object, ok := current.(map[string]any)
			if !ok {
				return nil, false
			}

			current, ok = object[name]
			if !ok {
				return nil, false
			}
		}

		for _, index := range indexes {
			list, ok := current.([]any)
			if !ok || index < 0 || index >= len(list) {
				return nil, false
			}

			current = list[index]
		}
	}

	return current, true
}

var indexPattern = regexp.MustCompile(`\[(\d+)\]`)

func splitIndexes(segment string) (string, []int) {
	name := segment
	if at := strings.Index(segment, "["); at >= 0 {
		name = segment[:at]
	}

	var indexes []int

	for _, match := range indexPattern.FindAllStringSubmatch(segment, -1) {
		n, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}

		indexes = append(indexes, n)
	}

	return name, indexes
}

func text(value any) string {
	switch typed := value.(type) {
	case nil:
		return "null"
	case string:
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case bool:
		return strconv.FormatBool(typed)
	default:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprint(typed)
		}

		return string(encoded)
	}
}

func describe(value any) string {
	switch value.(type) {
	case map[string]any:
		return "an object"
	case []any:
		return "a list"
	case nil:
		return "null"
	default:
		return text(value)
	}
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
