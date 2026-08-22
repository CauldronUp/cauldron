package openapi

import (
	"net/url"
	"strings"
)

// BasePath returns the path prefix a description's own servers imply.
//
// A description declares its paths relative to its server. Box's says
// https://api.box.com/2.0 and declares /files/{id}; the Recipe, which has to
// carry the whole path a client requests, says /2.0/files/{id}. Compared
// literally those never match, and every route in the Recipe is reported as a
// path the description does not have -- fifty-eight of them for one provider,
// none of it a real disagreement.
//
// That was left to the caller to notice and pass with --base. It is written
// down in the description, so this reads it.
//
// A server URL may be a template: Chargebee's is
// {protocol}://{site}.{environment}:{port}/api/v2. The placeholders are in the
// authority rather than the path, so trimming everything before the first
// slash that follows them recovers /api/v2. A template inside the path itself
// yields nothing, because a prefix that varies per request is not a prefix.
func BasePath(doc *Document) string {
	for _, server := range doc.Servers {
		if base := basePathOf(server.URL); base != "" {
			return base
		}
	}

	return ""
}

func basePathOf(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	path := raw

	// Strip the authority, whether or not the scheme is a template. Parsing
	// handles the ordinary case; the fallback handles {protocol}://host/x,
	// which url.Parse reads as an opaque scheme it cannot split.
	if parsed, err := url.Parse(raw); err == nil && parsed.Host != "" {
		path = parsed.Path
	} else if i := strings.Index(raw, "://"); i >= 0 {
		rest := raw[i+3:]
		if j := strings.Index(rest, "/"); j >= 0 {
			path = rest[j:]
		} else {
			path = ""
		}
	}

	path = strings.TrimSuffix(path, "/")

	if path == "" || !strings.HasPrefix(path, "/") {
		return ""
	}

	// A placeholder in the path means the prefix is not the same for every
	// request, so there is no one prefix to add.
	if strings.ContainsAny(path, "{}") {
		return ""
	}

	return path
}
