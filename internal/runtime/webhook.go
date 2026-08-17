package runtime

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CauldronUp/cauldron/internal/clock"
	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/store"
)

// Delivery is one webhook the sandbox produced.
type Delivery struct {
	ID      string
	Event   string
	At      time.Time
	Payload map[string]any
	// Signature is the value sent in the Recipe's signing header.
	Signature string
	// Endpoint is where it was sent, empty if nothing was listening.
	Endpoint string
	// Status is the response code the endpoint returned, zero if undelivered.
	Status int
	// Attempt counts delivery attempts, starting at 1.
	Attempt int
	Error   string
}

// webhookQueue records and delivers webhooks.
//
// Deliveries are recorded even when no endpoint is registered. Being able to
// assert "this event would have fired" without standing up a listener is most
// of what makes webhook behaviour testable at all.
type webhookQueue struct {
	mu sync.Mutex

	recipe    *recipe.Recipe
	clock     *clock.Clock
	endpoints []string
	client    *http.Client
	sent      []Delivery
	seq       int
}

func newWebhookQueue(r *recipe.Recipe, c *clock.Clock) *webhookQueue {
	return &webhookQueue{
		recipe: r,
		clock:  c,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Subscribe registers an endpoint to receive webhooks.
func (q *webhookQueue) Subscribe(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, existing := range q.endpoints {
		if existing == url {
			return
		}
	}

	q.endpoints = append(q.endpoints, url)
}

// Endpoints returns the registered endpoints.
func (q *webhookQueue) Endpoints() []string {
	q.mu.Lock()
	defer q.mu.Unlock()

	out := make([]string, len(q.endpoints))
	copy(out, q.endpoints)

	return out
}

// Deliveries returns every webhook produced, oldest first.
func (q *webhookQueue) Deliveries() []Delivery {
	q.mu.Lock()
	defer q.mu.Unlock()

	out := make([]Delivery, len(q.sent))
	copy(out, q.sent)

	return out
}

func (q *webhookQueue) reset() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.sent = nil
	q.seq = 0
}

// Known reports whether the Recipe declares an event.
func (q *webhookQueue) Known(event string) bool {
	for _, candidate := range q.recipe.Webhooks.Events {
		if candidate == event {
			return true
		}
	}

	return false
}

// Events returns the declared events.
func (q *webhookQueue) Events() []string {
	out := make([]string, len(q.recipe.Webhooks.Events))
	copy(out, q.recipe.Webhooks.Events)
	sort.Strings(out)

	return out
}

// payload builds the envelope the Recipe declares, or Cauldron's default.
//
// The default is Stripe's shape, which was the only shape that existed while
// Stripe was the only Recipe. Keeping it as the fallback means every Recipe
// written before envelopes were declarable behaves exactly as it did; it is
// not a claim that the provider in question sends that shape, and a Recipe
// that knows its provider's envelope should say so.
func (q *webhookQueue) payload(event, id string, at time.Time, data store.Record) map[string]any {
	template := q.recipe.Webhooks.Payload

	if len(template) == 0 {
		return map[string]any{
			"id":      id,
			"type":    event,
			"created": at.Unix(),
			"data": map[string]any{
				"object": map[string]any(data),
			},
		}
	}

	filled, _ := expand(template, substitutions{
		event: event, id: id, at: at, record: map[string]any(data),
	}).(map[string]any)

	return filled
}

// substitutions are the values a payload template can refer to.
type substitutions struct {
	event  string
	id     string
	at     time.Time
	record map[string]any
}

// objectKey marks the place a record's own fields are merged into. Adyen needs
// it: its notification item carries the payment's fields beside eventCode and
// success rather than nested under them, so splicing is the only faithful
// rendering.
const objectKey = "{object}"

// expand walks a payload template, substituting placeholders.
func expand(node any, with substitutions) any {
	switch typed := node.(type) {
	case map[string]any:
		out := map[string]any{}

		for key, value := range typed {
			if key == objectKey {
				for name, field := range with.record {
					out[name] = field
				}

				continue
			}

			out[key] = expand(value, with)
		}

		return out

	case []any:
		out := make([]any, len(typed))
		for i, value := range typed {
			out[i] = expand(value, with)
		}

		return out

	case string:
		return expandString(typed, with)
	}

	return node
}

// expandString substitutes into one scalar.
//
// A value of exactly "{object}" becomes the record itself rather than text, and
// "{created}" alone becomes a number, because a provider that sends a Unix
// timestamp sends a number and a client parsing it as a string would be
// exercising a shape nobody serves.
func expandString(value string, with substitutions) any {
	switch value {
	case objectKey:
		return with.record
	case "{created}":
		return with.at.Unix()
	}

	replaced := strings.NewReplacer(
		"{event}", with.event,
		"{id}", with.id,
		"{created}", strconv.FormatInt(with.at.Unix(), 10),
		"{created_iso}", with.at.UTC().Format(time.RFC3339),
	).Replace(value)

	return replaced
}

// Emit builds, signs and delivers a webhook.
//
// Unknown events are refused. A typo in an event name should fail loudly here
// rather than producing a webhook no real provider would ever send.
func (q *webhookQueue) Emit(event string, data store.Record) (Delivery, error) {
	if !q.Known(event) {
		return Delivery{}, fmt.Errorf("recipe %s does not emit %q (available: %s)", q.recipe.Name, event, strings.Join(q.Events(), ", "))
	}

	q.mu.Lock()

	q.seq++

	id := fmt.Sprintf("evt_%011d", q.seq)
	at := q.clock.Now()

	delivery := Delivery{
		ID:      id,
		Event:   event,
		At:      at,
		Payload: q.payload(event, id, at, data),
	}

	endpoints := make([]string, len(q.endpoints))
	copy(endpoints, q.endpoints)

	q.mu.Unlock()

	body, err := json.Marshal(delivery.Payload)
	if err != nil {
		return Delivery{}, err
	}

	delivery.Signature = q.sign(body, delivery.At)

	if len(endpoints) == 0 {
		q.record(delivery)

		return delivery, nil
	}

	for _, endpoint := range endpoints {
		attempt := delivery
		attempt.Endpoint = endpoint
		attempt.Attempt = 1

		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			attempt.Error = err.Error()
			q.record(attempt)

			continue
		}

		req.Header.Set("Content-Type", "application/json")

		if header := q.recipe.Webhooks.Signing.Header; header != "" && attempt.Signature != "" {
			req.Header.Set(header, attempt.Signature)
		}

		resp, err := q.client.Do(req)
		if err != nil {
			attempt.Error = err.Error()
			q.record(attempt)

			continue
		}

		attempt.Status = resp.StatusCode
		_ = resp.Body.Close()

		q.record(attempt)
	}

	return delivery, nil
}

func (q *webhookQueue) record(d Delivery) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.sent = append(q.sent, d)
}

// sign produces the signature header value.
//
// The format mirrors Stripe's — a timestamp and a versioned HMAC over
// "timestamp.body" — because applications verify signatures with the
// provider's own SDK, and a signature that does not round-trip through that
// SDK is worse than no signature at all.
func (q *webhookQueue) sign(body []byte, at time.Time) string {
	signing := q.recipe.Webhooks.Signing

	if signing.Scheme != "hmac-sha256" || signing.Secret == "" {
		return ""
	}

	timestamp := at.Unix()
	payload := fmt.Sprintf("%d.%s", timestamp, body)

	mac := hmac.New(sha256.New, []byte(signing.Secret))
	mac.Write([]byte(payload))

	return fmt.Sprintf("t=%d,v1=%s", timestamp, hex.EncodeToString(mac.Sum(nil)))
}
