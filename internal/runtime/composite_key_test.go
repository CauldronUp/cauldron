package runtime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A natural key is only as unique as the thing that owns it, which is what a
// scope is for.
//
// AfterShip keys a tracking by carrier and number together, because carriers
// mint tracking numbers in their own namespaces and they collide: the same
// number really does exist under usps and under fedex, on two different
// parcels. lookup_by resolved the number across the whole collection and then
// let the scope filter reject what it found, so one carrier answered and the
// other 404'd -- and which one worked depended on the order the store returned
// records in rather than on the request.
//
// The failure is worse than a 404. A provider that answers about the wrong
// parcel is the reason the key is composite in the first place.
func TestANaturalKeyIsResolvedWithinTheRoutesScope(t *testing.T) {
	r, err := recipe.Open("aftership")
	if err != nil {
		t.Fatalf("open aftership: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// One tracking number, two carriers, two parcels.
	const number = "9400111899223197428490"

	for slug, want := range map[string]string{
		"usps":  "Cast iron cauldron",
		"fedex": "Brass ladle",
	} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v4/trackings/"+slug+"/"+number, nil)
		req.Header.Set("aftership-api-key", "cauldron_aftership_key_00000000")
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200 (body %s)", slug, rec.Code, rec.Body.String())
		}

		var body struct {
			Data struct {
				Tracking struct {
					Slug  string `json:"slug"`
					Title string `json:"title"`
				} `json:"tracking"`
			} `json:"data"`
		}

		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: decode: %v", slug, err)
		}

		if got := body.Data.Tracking.Title; got != want {
			t.Errorf("%s/%s answered with %q, want %q -- the number was resolved outside the scope", slug, number, got, want)
		}

		if got := body.Data.Tracking.Slug; got != slug {
			t.Errorf("%s/%s answered with a %s tracking", slug, number, got)
		}
	}

	// And a carrier that has no such number is still a miss rather than a
	// guess at the nearest one.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v4/trackings/dhl/"+number, nil)
	req.Header.Set("aftership-api-key", "cauldron_aftership_key_00000000")
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("dhl/%s: status = %d, want 404", number, rec.Code)
	}
}

// A dotted envelope key nests, and this was the one of four places that did
// not. The list envelope nests, the error envelope nests, and the constants in
// withFields nest; a single record's wrapper built its map with the key
// straight in, so a Recipe wrapping under data.tracking got one key with a dot
// in its name -- a shape no provider sends.
//
// It failed quietly. The Recipe validated, the sandbox answered 200, and only
// a conformance case asserting the nested path noticed.
func TestASingleRecordsEnvelopeKeyNestsOnDots(t *testing.T) {
	r, err := recipe.Open("aftership")
	if err != nil {
		t.Fatalf("open aftership: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v4/trackings/usps/9400111899223197428490", nil)
	req.Header.Set("aftership-api-key", "cauldron_aftership_key_00000000")
	s.ServeHTTP(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if _, literal := body["data.tracking"]; literal {
		t.Fatal(`the body carries a key literally called "data.tracking"`)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("data = %#v, want an object", body["data"])
	}

	if _, ok := data["tracking"].(map[string]any); !ok {
		t.Errorf("data.tracking = %#v, want the record nested under it", data["tracking"])
	}
}
