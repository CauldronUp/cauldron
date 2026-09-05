package conform

import (
	"io"
	"mime"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// A case can send a multipart body, which several providers take and nothing
// else here could describe.
//
// Onfido's document and photo uploads are the larger half of its API. Coveralls
// takes a coverage report that way and nothing else. Dropbox Sign takes the
// document being signed. Zenvia is the sharpest and is why parts are a list
// rather than a flag: one request carrying a JSON part and a CSV part together,
// which a boolean could say nothing about.
func TestACaseCanSendAMultipartBody(t *testing.T) {
	req, err := buildRequest(recipe.Request{
		Method: "POST",
		Path:   "/v1/jobs",
		Multipart: []recipe.Part{
			{Name: "metadata", ContentType: "application/json", JSON: map[string]any{"encoding": "utf-8"}},
			{Name: "contacts", Filename: "contacts.csv", ContentType: "text/csv", Text: "a,b,c"},
		},
	}, "")
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	contentType := req.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		t.Fatalf("Content-Type = %q, want multipart/form-data", contentType)
	}

	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Fatalf("parse content type: %v", err)
	}

	// The boundary is generated rather than written into a Recipe, because a
	// fixture naming one would name something no client sends twice.
	if params["boundary"] == "" {
		t.Fatal("no boundary on a multipart request")
	}

	reader := multipart.NewReader(req.Body, params["boundary"])

	seen := map[string]string{}
	files := map[string]string{}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatalf("read part: %v", err)
		}

		body, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("read part body: %v", err)
		}

		seen[part.FormName()] = string(body)

		if name := part.FileName(); name != "" {
			files[part.FormName()] = name
		}
	}

	if seen["metadata"] != `{"encoding":"utf-8"}` {
		t.Errorf("the JSON part went out as %q", seen["metadata"])
	}

	if seen["contacts"] != "a,b,c" {
		t.Errorf("the file part went out as %q", seen["contacts"])
	}

	if files["contacts"] != "contacts.csv" {
		t.Errorf("the file part's name went out as %q", files["contacts"])
	}

	// And the part that is not a file is not marked as one, because several
	// providers refuse a field sent as a file and the reverse.
	if _, marked := files["metadata"]; marked {
		t.Error("an ordinary field went out marked as a file")
	}
}
