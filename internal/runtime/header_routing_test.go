package runtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CauldronUp/cauldron/internal/recipe"
)

// The AWS JSON protocol is every operation on POST /, with the operation named
// in X-Amz-Target. Without a way to route on that, the AWS Recipes encoded the
// operation in the path instead -- /ListSecrets, /tables, /queues -- and served
// URLs AWS does not have. A client can be written entirely against those
// paths, pass every case, and be entirely wrong.

func secretsManager(t *testing.T) *Sandbox {
	t.Helper()

	r, err := recipe.Open("secretsmanager")
	if err != nil {
		t.Fatalf("open secretsmanager: %v", err)
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

func awsCall(t *testing.T, s *Sandbox, target, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization",
		"AWS4-HMAC-SHA256 Credential=CAULDRONKEY/20260115/eu-west-2/secretsmanager/aws4_request, "+
			"SignedHeaders=host;x-amz-date, "+
			"Signature=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	req.Header.Set("Content-Type", "application/x-amz-json-1.1")

	if target != "" {
		req.Header.Set("X-Amz-Target", target)
	}

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	return rec
}

func TestTheHeaderPicksTheOperation(t *testing.T) {
	s := secretsManager(t)

	list := awsCall(t, s, "secretsmanager.ListSecrets", `{}`)
	if list.Code != http.StatusOK {
		t.Fatalf("ListSecrets: status %d: %s", list.Code, list.Body.String())
	}

	if !strings.Contains(list.Body.String(), "SecretList") {
		t.Errorf("ListSecrets did not answer with a listing: %s", list.Body.String())
	}

	// A different operation, on the same method and the same path.
	describe := awsCall(t, s, "secretsmanager.DescribeSecret", `{"SecretId":"prod/db"}`)
	if describe.Code != http.StatusOK {
		t.Fatalf("DescribeSecret: status %d: %s", describe.Code, describe.Body.String())
	}

	if strings.Contains(describe.Body.String(), "SecretList") {
		t.Errorf("DescribeSecret was answered by the listing route: %s", describe.Body.String())
	}
}

func TestWithoutTheHeaderThereIsNoOperation(t *testing.T) {
	// A 405 would be worse than a 404: it tells a client to change the
	// method, which was never the problem.
	rec := awsCall(t, secretsManager(t), "", `{}`)
	if rec.Code != http.StatusNotFound {
		t.Errorf("a request with no target answered %d, want 404", rec.Code)
	}
}

func TestAnUnmodelledOperationIsNotAnswered(t *testing.T) {
	// DeleteSecret is a real operation and not one modelled here. Answering
	// it from the nearest route would be a confident wrong answer.
	rec := awsCall(t, secretsManager(t), "secretsmanager.DeleteSecret", `{}`)
	if rec.Code != http.StatusNotFound {
		t.Errorf("an unmodelled operation answered %d, want 404", rec.Code)
	}
}

func TestTheTargetIsComparedExactly(t *testing.T) {
	// A substring comparison would route this to ListSecrets and answer it
	// with somebody else's data. The mutation that made it one passed every
	// other case in the Recipe.
	rec := awsCall(t, secretsManager(t), "secretsmanager.ListSecretsAndMore", `{}`)
	if rec.Code != http.StatusNotFound {
		t.Errorf("a near-miss target answered %d, want 404", rec.Code)
	}
}
