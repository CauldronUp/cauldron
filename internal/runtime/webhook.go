package runtime

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
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
	Payload any
	// Signature is the value sent in the Recipe's signing header.
	Signature string
	// Headers are the other headers a delivery carries, beyond the signature
	// and the content type. A signature over a timestamp is unverifiable
	// without the timestamp beside it.
	Headers map[string]string
	// SignatureHeader is the header that value travels in.
	//
	// Recorded whether or not anything is listening, because it is a claim the
	// Recipe makes rather than a property of a delivery: seventy-four Recipes
	// name a header and no case could assert one, since the name was only ever
	// applied at send time and a conformance case has no subscriber.
	//
	// It is the first thing a webhook handler reads. A Recipe naming the wrong
	// one produces a handler that looks for a header which never arrives, and
	// signature verification that fails only against the real provider.
	SignatureHeader string
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

// maxEndpoints caps how many places one sandbox will deliver to.
//
// Delivery happens inside the request that triggered it, so an uncapped list
// is a way to make every write in a test suite slow, and the subscribe
// endpoint is reachable by anything that can talk to the control plane. Nobody
// testing webhooks needs more receivers than this.
const maxEndpoints = 16

// deliveryTimeout bounds one attempt. Short, because these are local receivers
// and the caller is a test waiting on a create.
const deliveryTimeout = 2 * time.Second

func newWebhookQueue(r *recipe.Recipe, c *clock.Clock) *webhookQueue {
	return &webhookQueue{
		recipe: r,
		clock:  c,
		client: &http.Client{
			Timeout: deliveryTimeout,
			// A webhook receiver has no business redirecting, and following
			// one turns an endpoint somebody vetted into an address they never
			// saw. The payload carries record data and a valid signature, so
			// where it lands is worth being strict about.
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// Subscribe registers an endpoint to receive webhooks.
//
// The address is checked because this is an outbound request Cauldron makes on
// somebody else's say-so, carrying record data and a signature the receiver
// will treat as genuine.
//
// Loopback and private addresses are deliberately allowed: delivering to the
// application under test on localhost is the entire point, and refusing it
// would break the documented use for a threat this cannot fix on its own. What
// is refused is the class no local receiver ever lives in — a non-HTTP scheme,
// and the link-local range that cloud metadata services sit on.
func (q *webhookQueue) Subscribe(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("webhook endpoint is not a URL: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("webhook endpoint must be http or https, not %q", parsed.Scheme)
	}

	if parsed.Host == "" {
		return fmt.Errorf("webhook endpoint has no host: %s", raw)
	}

	if linkLocal(parsed.Hostname()) {
		return fmt.Errorf("refusing a webhook endpoint on a link-local address: %s", parsed.Host)
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	for _, existing := range q.endpoints {
		if existing == raw {
			return nil
		}
	}

	if len(q.endpoints) >= maxEndpoints {
		return fmt.Errorf("a sandbox delivers to at most %d endpoints", maxEndpoints)
	}

	q.endpoints = append(q.endpoints, raw)

	return nil
}

// linkLocal reports whether a host is in the range cloud metadata services
// live on. 169.254.169.254 is the address worth naming: it answers without
// credentials on most providers and hands back instance role tokens.
func linkLocal(host string) bool {
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	return ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
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
	// Subscriptions too. Leaving them made a reset sandbox distinguishable
	// from a fresh one in the way that matters most: an endpoint registered
	// once kept receiving every record the suite created afterwards, through
	// every reset between tests.
	q.endpoints = nil
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
func (q *webhookQueue) payload(event, id string, at time.Time, data store.Record) any {
	template := q.recipe.Webhooks.Payload

	if template == nil {
		return map[string]any{
			"id":      id,
			"type":    event,
			"created": at.Unix(),
			"data": map[string]any{
				"object": map[string]any(data),
			},
		}
	}

	resource := q.resourceFor(event)

	return expand(template, substitutions{
		event: event, id: id, at: at, record: map[string]any(data),
		resource: resource, action: actionOf(event, resource),
		recordID: q.identifierOf(resource, data),
	})
}

// substitutions are the values a payload template can refer to.
type substitutions struct {
	event  string
	id     string
	at     time.Time
	record map[string]any
	// resource is the name of the thing the event is about, which some
	// envelopes key on. Square nests the record under it twice over.
	resource string
	// action is the rest of the event name, which some envelopes carry
	// separately from the resource: Asana sends {resource: {resource_type},
	// action} and QuickBooks sends {name, operation}, so the one string this
	// project calls an event is two fields to them.
	action string
	// recordID is the changed record's own identifier, which is not the
	// delivery's. Mollie's webhook body is that id and nothing else, on
	// purpose: the notification says something happened and you fetch the
	// payment to find out what.
	//
	// Resolved against the resource rather than read from a field called id,
	// because 121 resources in this collection call it something else --
	// Asana's is gid, Adyen's is pspReference, Calendly's is a uri -- and
	// reading "id" off those gives the empty string, which renders as a
	// present field with nothing in it rather than as an error.
	recordID string
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

		// The record first, then the template's own keys, rather than both in
		// whatever order the map yields.
		//
		// Go randomises map iteration, so a template naming a key the record
		// also carries produced one payload or the other from run to run --
		// Slack's envelope has a literal type beside the merged event, and a
		// record with a type field would have won or lost the coin toss. An
		// explicit key is the Recipe author saying what the provider sends
		// there, so it takes precedence, and it does so every time.
		for key := range typed {
			if key != objectKey {
				continue
			}

			for name, field := range with.record {
				out[name] = field
			}
		}

		for key, value := range typed {
			if key == objectKey {
				continue
			}

			// The key as well as the value. Square's data.object holds the
			// record under a key that is the resource's own name -- a payment
			// arrives at data.object.payment and a customer at
			// data.object.customer -- so a Recipe-wide envelope can only
			// describe it if the key can vary.
			out[expandKey(key, with)] = expand(value, with)
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

// expandKey substitutes into a template's key.
//
// Only the plain replacements. A key is a name rather than a value, so the
// forms that return a non-string -- {object}, {created} -- have no meaning
// here, and {object} as a key is the merge marker and is handled before this.
func expandKey(key string, with substitutions) string {
	replaced, _ := expandString(key, with).(string)
	if replaced == "" {
		return key
	}

	return replaced
}

// expandString substitutes into one scalar.
//
// A value of exactly "{object}" becomes the record itself rather than text, and
// "{created}" alone becomes a number, because a provider that sends a Unix
// timestamp sends a number and a client parsing it as a string would be
// exercising a shape nobody serves.
func expandString(value string, with substitutions) any {
	// A single field of the record, kept at whatever type it already has.
	//
	// Splicing the whole record is the common case and was the only one, which
	// left a provider that renames or prefixes the fields it sends
	// undescribable: Freshdesk wraps its payload in freshdesk_webhook and
	// prefixes every key with ticket_, so the record's own names appear
	// nowhere in it. A template that could only merge had to either send the
	// wrong key names or send Stripe's envelope instead.
	if name, found := strings.CutPrefix(value, "{record."); found {
		if path, closed := strings.CutSuffix(name, "}"); closed && !strings.Contains(path, "{") {
			return lookupIn(with.record, path)
		}
	}

	switch value {
	case objectKey:
		return with.record
	case "{created}":
		return with.at.Unix()
	case "{created_ms}":
		// Milliseconds, as a number. The Recipe format already treats this as
		// a distinction that matters -- it has both timestamp and timestamp_ms
		// field types, because a factor of a thousand puts a date in 1970 or
		// in the year 55000 -- and a payload envelope needs the same choice.
		// Zoom's event_ts is milliseconds.
		return with.at.UnixMilli()
	}

	replaced := strings.NewReplacer(
		"{event}", with.event,
		"{resource}", with.resource,
		"{action}", with.action,
		"{id}", with.id,
		"{record_id}", with.recordID,
		"{created}", strconv.FormatInt(with.at.Unix(), 10),
		"{created_ms}", strconv.FormatInt(with.at.UnixMilli(), 10),
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

	if delivery.Signature != "" {
		delivery.SignatureHeader = q.recipe.Webhooks.Signing.Header

		if header := q.recipe.Webhooks.Signing.TimestampHeader; header != "" {
			delivery.Headers = map[string]string{header: strconv.FormatInt(at.Unix(), 10)}
		}
	}

	if len(endpoints) == 0 {
		q.record(delivery)

		return delivery, nil
	}

	// Concurrently, because delivery happens inside the request that triggered
	// it and endpoints fail by timing out rather than by refusing. In series,
	// three receivers that have stopped answering cost three timeouts, and a
	// suite that registers a handful of them turns every create into a
	// quarter-minute wait. Fanning out costs one timeout however many there
	// are.
	//
	// The results are collected in order and recorded afterwards, so what a
	// test reads back does not depend on which receiver answered first.
	attempts := make([]Delivery, len(endpoints))

	var wg sync.WaitGroup

	for i, endpoint := range endpoints {
		wg.Add(1)

		go func(i int, endpoint string) {
			defer wg.Done()

			attempts[i] = q.deliver(delivery, endpoint, body)
		}(i, endpoint)
	}

	wg.Wait()

	for _, attempt := range attempts {
		q.record(attempt)
	}

	return delivery, nil
}

// deliver posts one webhook and reports what happened, without touching shared
// state: the caller records the results in endpoint order once they are all in.
func (q *webhookQueue) deliver(delivery Delivery, endpoint string, body []byte) Delivery {
	attempt := delivery
	attempt.Endpoint = endpoint
	attempt.Attempt = 1

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		attempt.Error = err.Error()

		return attempt
	}

	req.Header.Set("Content-Type", "application/json")

	if header := q.recipe.Webhooks.Signing.Header; header != "" && attempt.Signature != "" {
		req.Header.Set(header, attempt.Signature)
	}

	for name, value := range attempt.Headers {
		req.Header.Set(name, value)
	}

	resp, err := q.client.Do(req)
	if err != nil {
		attempt.Error = err.Error()

		return attempt
	}

	attempt.Status = resp.StatusCode
	_ = resp.Body.Close()

	return attempt
}

func (q *webhookQueue) record(d Delivery) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.sent = append(q.sent, d)
}

// sign produces the signature header value.
//
// The shape comes from the Recipe, because applications verify signatures
// with the provider's own SDK and a signature that does not round-trip
// through that SDK is worse than no signature at all.
//
// That sentence was here before, above code that gave every Recipe Stripe's
// shape. Seventy-four declare hmac-sha256 under sixty-eight different header
// names, and a Stripe-shaped value under X-Shopify-Hmac-Sha256 fails
// Shopify's verifier on the first call.
func (q *webhookQueue) sign(body []byte, at time.Time) string {
	signing := q.recipe.Webhooks.Signing

	if signing.Scheme != "hmac-sha256" || signing.Secret == "" {
		return ""
	}

	timestamp := at.Unix()

	// What is signed differs as much as how it is wrapped. Slack prefixes a
	// version and the timestamp; GitHub and Shopify sign the body alone.
	signed := string(body)

	switch signing.Format {
	case "v0-hex":
		signed = fmt.Sprintf("v0:%d:%s", timestamp, body)
	case "prefixed-hex", "base64", "hex":
	default:
		signed = fmt.Sprintf("%d.%s", timestamp, body)
	}

	mac := hmac.New(sha256.New, []byte(signing.Secret))
	mac.Write([]byte(signed))
	sum := mac.Sum(nil)

	switch signing.Format {
	case "hex":
		return hex.EncodeToString(sum)
	case "prefixed-hex":
		return "sha256=" + hex.EncodeToString(sum)
	case "base64":
		return base64.StdEncoding.EncodeToString(sum)
	case "v0-hex":
		return "v0=" + hex.EncodeToString(sum)
	default:
		return fmt.Sprintf("t=%d,v1=%s", timestamp, hex.EncodeToString(sum))
	}
}

// resourceFor names the resource an event is about.
//
// Taken from the route that declares it where one does, and from the
// convention otherwise, so the manual `cauldron emit` path resolves it the
// same way a write does rather than needing to be told.
func (q *webhookQueue) resourceFor(event string) string {
	for _, route := range q.recipe.Routes {
		if route.Emits == event {
			return route.Resource
		}

		for _, conditional := range route.EmitsWhen {
			if conditional.Event == event {
				return route.Resource
			}
		}
	}

	name, _, found := strings.Cut(event, ".")
	if !found {
		return ""
	}

	if _, known := q.recipe.Resources[name]; known {
		return name
	}

	return ""
}

// lookupIn reads a dotted path out of a record, returning nil when any step of
// it is missing.
//
// nil rather than the empty string, because a field a provider sends as null
// and a field it does not send are different things, and a template naming a
// field the resource does not declare should look like the second.
func lookupIn(record map[string]any, path string) any {
	var current any = record

	for _, step := range strings.Split(path, ".") {
		object, isObject := current.(map[string]any)
		if !isObject {
			return nil
		}

		value, present := object[step]
		if !present {
			return nil
		}

		current = value
	}

	return current
}

// Fields is the payload as an object, or nil when the provider sends
// something else.
//
// Most envelopes are objects and reading one field out of them should not
// need a type assertion at every call site. The ones that are not -- an array
// of batched events -- have nothing to return here, and nil says so.
func (d Delivery) Fields() map[string]any {
	fields, _ := d.Payload.(map[string]any)

	return fields
}

// actionOf is what is left of an event name once the resource is taken out of
// it.
//
// Not the part after the last dot, which is what this was and which only
// worked for providers that put the resource first. Pipedrive does not:
// its events are added.deal and updated.person, so the suffix is the resource
// and the prefix is the action. Asana's task.added and QuickBooks'
// Customer.Create are the other way round, and Dropbox separates with an
// underscore rather than a dot.
//
// Removing the resource handles all of them without needing to know which
// order a provider chose, and it is the more honest definition anyway: the
// action is the part that is not the thing.
//
// Empty when the resource is not in the name, which is the honest answer
// rather than a guess. Handing back the whole event would put a plausible
// value somewhere the provider sends something else.
func actionOf(event, resource string) string {
	// Case-insensitively, because a provider's event names and a Recipe's
	// resource names do not have to agree on it: Xero declares INVOICE.CREATE
	// against a resource called invoice, and Documenso DOCUMENT_CREATED
	// against document. Matching exactly returned the empty string for both,
	// which is the failure this whole substitution exists to avoid -- a
	// present field with nothing in it.
	at := strings.Index(strings.ToLower(event), strings.ToLower(resource))
	if resource == "" || at < 0 {
		return ""
	}

	rest := event[:at] + event[at+len(resource):]

	return strings.Trim(rest, "._:-/ ")
}

// identifierOf reads a record's own identifier, under whatever name the
// resource gives it.
func (q *webhookQueue) identifierOf(resource string, record store.Record) string {
	name := "id"

	if spec, known := q.recipe.Resources[resource]; known && spec.ID.Field != "" {
		name = spec.ID.Field
	}

	switch value := record[name].(type) {
	case string:
		return value
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", value)
	}
}
