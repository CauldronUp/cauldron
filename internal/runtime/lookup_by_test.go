package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// SQS deletes a message by the receipt handle from a receive, and a receipt
// handle is deliberately not a message id: it is issued per receive, two
// consumers holding two handles for the same message is normal, and a handle
// from an earlier receive is stale.
//
// Looked up as though it were an id it finds nothing, so every delete failed
// and the only thing the Recipe could show about DeleteMessage was the way it
// fails.

func sqsSandbox(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("sqs")
	if err != nil {
		t.Fatalf("open sqs: %v", err)
	}

	s, err := New(r, Options{Seed: 1})
	if err != nil {
		t.Fatalf("new sandbox: %v", err)
	}

	if err := s.Seed("small-account"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	return s
}

func sqsCall(t *testing.T, s *Sandbox, target, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization",
		"AWS4-HMAC-SHA256 Credential=CAULDRONKEY/20260115/eu-west-2/sqs/aws4_request, "+
			"SignedHeaders=host;x-amz-date, "+
			"Signature=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")
	req.Header.Set("X-Amz-Target", target)

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

func messageCount(t *testing.T, s *Sandbox) int {
	t.Helper()

	rec := sqsCall(t, s, "AmazonSQS.ReceiveMessage", `{}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("receive: status %d", rec.Code)
	}

	var out struct {
		Messages []map[string]any `json:"Messages"`
	}

	decodeInto(t, rec, &out)

	return len(out.Messages)
}

func TestALiveHandleDeletesTheMessageItCameFrom(t *testing.T) {
	s := sqsSandbox(t)

	before := messageCount(t, s)

	rec := sqsCall(t, s, "AmazonSQS.DeleteMessage",
		`{"QueueUrl":"https://sqs.eu-west-2.amazonaws.com/000000000000/orders",`+
			`"ReceiptHandle":"AQEBcauldron0000000000000000000000000002"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("delete with a live handle: status %d: %s", rec.Code, rec.Body.String())
	}

	if after := messageCount(t, s); after != before-1 {
		t.Errorf("delete removed %d messages, want 1", before-after)
	}
}

func TestAStaleHandleDeletesNothing(t *testing.T) {
	// A handle from an earlier receive. The consumer took longer than the
	// visibility timeout, the message is already back on the queue, and the
	// delete does nothing -- which is the failure the Recipe exists to show.
	s := sqsSandbox(t)

	before := messageCount(t, s)

	rec := sqsCall(t, s, "AmazonSQS.DeleteMessage",
		`{"QueueUrl":"https://sqs.eu-west-2.amazonaws.com/000000000000/orders",`+
			`"ReceiptHandle":"AQEBcauldron0000000000000000000000000099"}`)
	if rec.Code == http.StatusOK {
		t.Error("a stale handle deleted something")
	}

	if after := messageCount(t, s); after != before {
		t.Errorf("a stale handle removed %d messages", before-after)
	}
}

func TestTheRightMessageIsTheOneDeleted(t *testing.T) {
	// The handle names one message, and deleting by a field rather than an id
	// must not take whichever record happened to come first.
	s := sqsSandbox(t)

	rec := sqsCall(t, s, "AmazonSQS.DeleteMessage",
		`{"QueueUrl":"https://sqs.eu-west-2.amazonaws.com/000000000000/orders",`+
			`"ReceiptHandle":"AQEBcauldron0000000000000000000000000002"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("delete: status %d", rec.Code)
	}

	var out struct {
		Messages []map[string]any `json:"Messages"`
	}

	decodeInto(t, sqsCall(t, s, "AmazonSQS.ReceiveMessage", `{}`), &out)

	for _, message := range out.Messages {
		if message["MessageId"] == "22222222-2222-4222-8222-222222222222" {
			t.Error("the message the handle named is still there")
		}
	}

	if len(out.Messages) != 3 {
		t.Errorf("expected the other three to survive, got %d", len(out.Messages))
	}
}

func TestAnIdentifierIsNotAHandle(t *testing.T) {
	// The direction that was wrong. A message id in the ReceiptHandle field
	// used to answer 200 and delete the message, because a value matching no
	// handle fell back to an identifier lookup -- and a message id is one.
	//
	// The emulator was teaching what the Recipe's own header warns against:
	// that a handle and an id are interchangeable.
	s := sqsSandbox(t)

	before := messageCount(t, s)

	rec := sqsCall(t, s, "AmazonSQS.DeleteMessage",
		`{"QueueUrl":"https://sqs.eu-west-2.amazonaws.com/000000000000/orders",`+
			`"ReceiptHandle":"11111111-1111-4111-8111-111111111111"}`)
	if rec.Code == http.StatusOK {
		t.Error("a message id was accepted as a receipt handle")
	}

	if after := messageCount(t, s); after != before {
		t.Errorf("a message id deleted %d messages", before-after)
	}
}
