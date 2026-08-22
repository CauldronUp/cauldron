// Paging: where a page starts, how big it is, and how a caller asks for the
// next one.

package runtime

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// nextPageURL renders the request that fetches the page after this one.
//
// It is this request with the position parameter moved on, which is what every
// provider's Link header contains: the same query, one page further in. A
// listing whose paging travels in the body has no such URL, so it gets no
// header rather than a misleading one.
func nextPageURL(r *http.Request, spec recipe.Pagination, cursor string) string {
	if spec.In == "body" {
		return ""
	}

	name := spec.CursorParam
	if name == "" {
		switch spec.Style {
		case "page", "offset":
			name = spec.Style
		default:
			name = "cursor"
		}
	}

	if name == "-" {
		return ""
	}

	next := *r.URL

	// The path the caller asked for, which is not always the one the sandbox
	// sees. The multi-provider server mounts each Recipe under its own name
	// and rewrites URL.Path, leaving RequestURI as it arrived, so a header
	// built from the rewritten path sends the client to /repos/... when the
	// server serves /github/repos/... -- a next page that 404s, which is a
	// worse failure than no next page at all.
	if r.RequestURI != "" {
		if asked, err := url.ParseRequestURI(r.RequestURI); err == nil && asked.Path != "" {
			next.Path = asked.Path
		}
	}

	query := next.Query()
	query.Set(name, cursor)
	next.RawQuery = query.Encode()

	// A provider's Link header carries an absolute URL, and a client is
	// expected to request it as it stands rather than to reassemble one.
	next.Scheme = "http"
	next.Host = r.Host

	return next.String()
}

// prevPageURL builds the previous page's URL, for the styles that have one.
//
// Only offset and page numbering do. A cursor names a position the caller was
// handed and cannot be arithmetic'd backwards, so a cursor-paged listing gets
// no prev link rather than a guessed one -- which is what the providers
// themselves do.
func prevPageURL(r *http.Request, spec recipe.Pagination, limit int) string {
	switch spec.Style {
	case "offset", "page":
	default:
		return ""
	}

	if spec.In == "body" {
		return ""
	}

	name := spec.CursorParam
	if name == "" {
		name = spec.Style
	}

	if name == "-" {
		return ""
	}

	from := pagingFrom(r, spec)

	position := positionOf(from, spec, limit)
	if position <= 0 {
		return ""
	}

	previous := position - limit
	if previous < 0 {
		previous = 0
	}

	value := strconv.Itoa(previous)
	if spec.Style == "page" {
		value = strconv.Itoa(previous/max(limit, 1) + spec.FirstPageNumber())
	}

	prev := *r.URL

	if r.RequestURI != "" {
		if asked, err := url.ParseRequestURI(r.RequestURI); err == nil && asked.Path != "" {
			prev.Path = asked.Path
		}
	}

	query := prev.Query()
	query.Set(name, value)
	prev.RawQuery = query.Encode()

	prev.Scheme = "http"
	prev.Host = r.Host

	return prev.String()
}

// cursorOf reads the identifier a cursor-paged listing resumes after.
func cursorOf(from paging, spec recipe.Pagination) string {
	if name := spec.CursorParam; name != "" {
		return from.get(name)
	}

	if cursor := from.get("starting_after"); cursor != "" {
		return cursor
	}

	return from.get("cursor")
}

// positionOf turns an offset or page parameter into a record offset.
//
// The two styles differ only in what the number counts. An offset counts
// records and starts at nought; a page counts pages and starts at one, so page
// two begins one page-length in. Getting that boundary wrong by one page is
// the classic way to lose or duplicate a record at every page break, which is
// exactly the bug an emulator is supposed to catch rather than commit.
func positionOf(from paging, spec recipe.Pagination, limit int) int {
	name := spec.CursorParam
	if name == "" {
		name = spec.Style
	}

	raw := from.get(name)
	if raw == "" {
		return 0
	}

	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return 0
	}

	if spec.Style != "page" {
		return n
	}

	first := spec.FirstPageNumber()

	// Anything at or before the provider's own first page is the first page.
	// Providers disagree about whether a number below it is an error, and
	// answering with the first page is the reading that cannot silently skip
	// records.
	if n <= first {
		return 0
	}

	return (n - first) * limit
}

// nextPosition renders the position of the page after this one.
//
// An offset counts records, so the next one starts a page-length further in. A
// page counts pages from one, so the offset has to be turned back into a page
// number before adding to it.
func nextPosition(spec recipe.Pagination, offset, limit int) string {
	if spec.Style == "page" {
		if limit <= 0 {
			limit = 1
		}

		// The page after this one, counted from wherever the provider starts.
		return strconv.Itoa(offset/limit + spec.FirstPageNumber() + 1)
	}

	return strconv.Itoa(offset + limit)
}

// servedPage is the page number this request asked for, counted the way the
// provider counts.
func servedPage(from paging, spec recipe.Pagination) int {
	first := spec.FirstPageNumber()

	name := spec.CursorParam
	if name == "" {
		name = spec.Style
	}

	raw := from.get(name)
	if raw == "" {
		return first
	}

	n, err := strconv.Atoi(raw)
	if err != nil || n < first {
		return first
	}

	return n
}

// pageCount is how many pages a total makes at this page size.
//
// An empty set is nought pages. Providers differ about whether it is nought or
// one, and nought is the reading that stops a loop rather than sending it
// after a page with nothing in it.
func pageCount(total, limit int) int {
	if total <= 0 || limit <= 0 {
		return 0
	}

	return (total + limit - 1) / limit
}

// cursorValue renders the cursor as the Recipe says the provider sends it.
//
// A token by default, and the whole next-page URL for the providers whose
// pointer is one. Getting this wrong is not cosmetic: a client written against
// a token concatenates it onto a base URL, which is the exact mistake several
// of these Recipes describe in a comment and then taught anyway.
//
// The URL is empty when paging travels in the body, because there is no such
// URL to render; the token is the honest answer then.
func cursorValue(spec recipe.ListResponse, cursor, nextURL string) string {
	if spec.CursorURL == "" || nextURL == "" {
		return cursor
	}

	if spec.CursorURL == "path" {
		// Salesforce joins its nextRecordsUrl to the instance URL it
		// authenticated against, so the path is the whole of what it sends
		// and an absolute address would be joined to that.
		if at := strings.Index(nextURL, "://"); at >= 0 {
			if slash := strings.IndexByte(nextURL[at+3:], '/'); slash >= 0 {
				return nextURL[at+3+slash:]
			}
		}
	}

	return nextURL
}
