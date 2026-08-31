// Telling a format this cannot read from a host that did not answer.

package openapi

import (
	"fmt"
	"strings"
)

// FormatError is a description that arrived intact and is not something this
// package reads.
//
// The distinction it carries is the one a scheduled scan lives or dies by. A
// host answering 503 is a Tuesday and will probably be fine tomorrow. A
// provider publishing RAML is permanent until somebody writes a reader, and a
// scan that reports the two under one word shows the same line every morning
// for a Recipe whose provider is perfectly reachable, with nothing to say
// which of the two it is. Everyone stops reading the scan, and an ignored scan
// is worse than no scan because it also costs the time spent reading it.
//
// This was a substring search for the word "Swagger" until Swagger 2.0 stopped
// being unreadable, and the search was wrong in both directions: it caught a
// network error that happened to mention Swagger in a URL, and it missed every
// format that is not Swagger. A provider publishing RAML or AsyncAPI was
// reported unreachable forever.
type FormatError struct {
	// Format is what the document said it was, in the provider's own words
	// where it gave any.
	Format string
	// Err is the underlying reason, where there is one beyond the format
	// itself.
	Err error
}

func (e *FormatError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("this is %s and it could not be read as OpenAPI 3: %v", e.Format, e.Err)
	}

	return fmt.Sprintf("this is %s, which is not a format this reads", e.Format)
}

func (e *FormatError) Unwrap() error { return e.Err }

// detectForeignFormat names the format a document announces, where that format
// is one this package knows it does not read.
//
// Only formats that announce themselves are named. Guessing at an unmarked
// document would put a permanent verdict on something that might be a login
// page, and "this provider publishes a format we cannot read" is not a thing
// to say on a guess -- it reads as a fact about the provider.
func detectForeignFormat(raw []byte) string {
	head := raw
	if len(head) > 4096 {
		head = head[:4096]
	}

	text := string(head)

	for _, candidate := range []struct {
		marker string
		format string
	}{
		// RAML announces itself on the first line and nowhere else.
		{"#%RAML", "RAML"},
		{`"asyncapi"`, "AsyncAPI"},
		{"asyncapi:", "AsyncAPI"},
		{"<wsdl:", "WSDL"},
		{"<definitions", "WSDL"},
		{`"$schema": "https://json-schema.org`, "a JSON Schema"},
	} {
		if strings.Contains(text, candidate.marker) {
			return candidate.format
		}
	}

	return ""
}
