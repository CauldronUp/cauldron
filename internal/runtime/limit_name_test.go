package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Which word a listing reads for its page size has three answers, and the two
// that are not the default are the ones worth a test: a named parameter makes
// every other spelling inert, and "-" makes all of them inert.
//
// Five conformance cases were wrong because nothing said so. Help Scout's,
// Mailchimp's, SES's, Zoom's and Docker Hub's each asked for a page with
// "limit", passed, and were receiving a whole collection -- because until a
// Recipe named its parameter the runtime read "limit" from everybody, so the
// wrong word worked and the right one was never tried.
//
// Each of those is fixed in its own Recipe. This is the behaviour underneath
// them, stated once rather than inferred from five providers that happen to
// exercise it.
func TestWhichWordAListingReadsForItsPageSize(t *testing.T) {
	count := func(t *testing.T, s *Sandbox, path, query string, headers map[string]string, at string) int {
		t.Helper()

		req := httptest.NewRequest(http.MethodGet, path+query, nil)
		for name, value := range headers {
			req.Header.Set(name, value)
		}

		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s%s = %d\n%s", path, query, rec.Code, rec.Body)
		}

		var body any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode: %v", err)
		}

		list, ok := lookupList(body, at)
		if !ok {
			t.Fatalf("GET %s%s: no list at %q in %s", path, query, at, rec.Body)
		}

		return len(list)
	}

	t.Run("a named parameter is the only one read", func(t *testing.T) {
		// Mailchimp calls it count. A request sending limit is sending a word
		// Mailchimp does not accept, and the emulator has to ignore it too --
		// otherwise a client's paging loop works here and does nothing there.
		s := pagedSandbox(t, "mailchimp", "small-account")
		headers := map[string]string{"Authorization": "Bearer cauldron-oauth-token"}
		path := "/3.0/lists/a1b2c3d4e5/members"

		if n := count(t, s, path, "?count=1", headers, "members"); n != 1 {
			t.Errorf("count=1 returned %d members, want 1", n)
		}

		if n := count(t, s, path, "?limit=1", headers, "members"); n != 2 {
			t.Errorf("limit=1 returned %d members, want both: Mailchimp does not read limit", n)
		}
	})

	t.Run("a dash refuses every name", func(t *testing.T) {
		// Rollbar accepts no page size at all. The page is Rollbar's, not the
		// caller's, and an emulator that honoured one would hand back exactly
		// the page that was asked for with a next link Rollbar never sends.
		s := pagedSandbox(t, "rollbar", "small-account")
		headers := map[string]string{"X-Rollbar-Access-Token": "cauldron-rollbar-token"}

		for _, query := range []string{"?limit=1", "?page_size=1", "?per_page=1"} {
			if n := count(t, s, "/api/1/projects", query, headers, "result"); n != 2 {
				t.Errorf("%s returned %d projects, want both: Rollbar takes no size", query, n)
			}
		}
	})

	t.Run("naming nothing reads limit", func(t *testing.T) {
		// The fallback, which is right for the providers that call it limit
		// and is a guess for the rest. Increase is one that does.
		s := pagedSandbox(t, "increase", "small-business")
		headers := map[string]string{"Authorization": "Bearer cauldron_increase_api_key"}

		if n := count(t, s, "/ach_transfers", "?limit=1", headers, "data"); n != 1 {
			t.Errorf("limit=1 returned %d transfers, want 1", n)
		}
	})
}

// lookupList digs a list out of a response body by a single key, or takes the
// body itself when the key is empty.
func lookupList(body any, at string) ([]any, bool) {
	if at == "" {
		list, ok := body.([]any)

		return list, ok
	}

	object, ok := body.(map[string]any)
	if !ok {
		return nil, false
	}

	list, ok := object[at].([]any)

	return list, ok
}
