package runtime

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
)

// A multipart body is read rather than guessed at.
//
// The request decoder handled JSON and url-encoded forms and fell through to
// guessing at anything else, so a multipart upload -- which is most of what
// several providers do -- arrived as nothing. Onfido's document and photo
// uploads are the larger half of its API. Coveralls takes a coverage report
// that way and nothing else. Dropbox Sign takes the document being signed.
//
// A file part becomes its filename, not its bytes. What a Recipe can check
// about an upload is that it arrived, under the right field name, with the name
// and type the provider cares about -- and several refuse a file sent as an
// ordinary field. The content is not a claim anybody can verify against a
// provider.
func TestAMultipartBodyIsReadIntoTheSameRecordShape(t *testing.T) {
	body, contentType := build(t, []part{
		{name: "report", filename: "coverage.json", contentType: "application/json", text: `{"covered":9}`},
		{name: "flag_name", text: "unit"},
	})

	got, err := decodeMultipart(body, contentType)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got["report"] != "coverage.json" {
		t.Errorf("the file part read as %v, want its filename", got["report"])
	}

	if got["report_content_type"] != "application/json" {
		t.Errorf("the file part's own type read as %v", got["report_content_type"])
	}

	if got["flag_name"] != "unit" {
		t.Errorf("the ordinary field read as %v, want its value", got["flag_name"])
	}
}

// A part carrying JSON is decoded, so metadata beside a file reads the way a
// JSON body would. Zenvia sends exactly that: one request with a JSON part and
// a CSV part together, which is why parts are a list rather than a flag.
func TestAJSONPartBesideAFileIsDecoded(t *testing.T) {
	body, contentType := build(t, []part{
		{name: "metadata", contentType: "application/json", text: `{"encoding":"utf-8","rows":3}`},
		{name: "contacts", filename: "contacts.csv", contentType: "text/csv", text: "a,b,c"},
	})

	got, err := decodeMultipart(body, contentType)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got["metadata.encoding"] != "utf-8" {
		t.Errorf("the JSON part's field read as %v", got["metadata.encoding"])
	}

	if got["contacts"] != "contacts.csv" {
		t.Errorf("the file beside it read as %v", got["contacts"])
	}
}

// And the decoder reaches it through the ordinary path, on the content type,
// rather than needing a caller to know what shape arrived.
func TestTheBodyDecoderDispatchesOnTheMultipartContentType(t *testing.T) {
	body, contentType := build(t, []part{{name: "flag_name", text: "unit"}})

	req := httptest.NewRequest(http.MethodPost, "/v1/jobs", bytes.NewReader(body))
	req.Header.Set("Content-Type", contentType)

	got, err := decodeBody(req)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got["flag_name"] != "unit" {
		t.Errorf("decodeBody read %v from a multipart body", got)
	}
}

// A body claiming to be multipart with no boundary is an error rather than an
// empty record, because silently reading nothing is how a broken upload passes.
func TestAMultipartBodyWithNoBoundaryIsAnError(t *testing.T) {
	if _, err := decodeMultipart([]byte("whatever"), "multipart/form-data"); err == nil {
		t.Error("a multipart body with no boundary decoded without complaint")
	}
}

type part struct {
	name        string
	filename    string
	contentType string
	text        string
}

func build(t *testing.T, parts []part) ([]byte, string) {
	t.Helper()

	var body bytes.Buffer

	form := multipart.NewWriter(&body)

	for _, p := range parts {
		headers := textproto.MIMEHeader{}

		disposition := `form-data; name="` + p.name + `"`
		if p.filename != "" {
			disposition += `; filename="` + p.filename + `"`
		}

		headers.Set("Content-Disposition", disposition)

		if p.contentType != "" {
			headers.Set("Content-Type", p.contentType)
		}

		w, err := form.CreatePart(headers)
		if err != nil {
			t.Fatalf("create part: %v", err)
		}

		if _, err := w.Write([]byte(p.text)); err != nil {
			t.Fatalf("write part: %v", err)
		}
	}

	if err := form.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	return body.Bytes(), form.FormDataContentType()
}
