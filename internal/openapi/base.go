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
// A later server may supply the prefix, but only if it is the same host.
//
// The first server is the default, so its path is the answer when it has one.
// When it has none, a second server on the same host carrying a version is
// describing the same API twice -- api.example.com and api.example.com/v3 --
// and taking the version is what a caller means.
//
// A second server on a *different* host is not that. Printify found the
// difference: its description gained https://enterprise.printful.com/api/pfy/public,
// a different company's domain, listed after its own api.printify.com which has
// no path. Falling through to it prefixed every declared path, so both of that
// Recipe's routes read as absent and its fingerprint moved -- for a change to a
// server list that has nothing to do with the orders API the Recipe describes.
//
// That is the failure this package exists to prevent: drift reported against
// something the provider did not change to anything the Recipe claims. A scan
// that cries wolf gets switched off, and then the change that mattered arrives
// unannounced.
func BasePath(doc *Document) string {
	if len(doc.Servers) == 0 {
		return ""
	}

	first := doc.Servers[0].URL

	if base := basePathOf(first); base != "" {
		return base
	}

	for _, server := range doc.Servers[1:] {
		if sameHost(first, server.URL) {
			if base := basePathOf(server.URL); base != "" {
				return base
			}
		}
	}

	return ""
}

// sameHost reports whether two server URLs name the same authority, so a
// prefix can be borrowed from one for the other. hostOf lives in discover.go
// and reads only http and https, which is the conservative answer here: a
// templated scheme yields no host, so nothing is borrowed for it.
func sameHost(a, b string) bool {
	ha, hb := hostOf(a), hostOf(b)

	return ha != "" && ha == hb
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
