// Package recipe defines the Recipe format — the description of how Cauldron
// emulates one external dependency.
//
// A Recipe is not a mock. A mock returns a shape; a Recipe models behaviour:
// what resources exist, how they change, what the provider emits afterwards,
// and how it fails. The format is declarative on purpose, so that the majority
// of Recipes are data a contributor can read and review rather than code.
package recipe

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Recipe is a complete emulation description for one provider.
type Recipe struct {
	Name    string `yaml:"recipe"`
	Version string `yaml:"version"`
	// Capability is what kind of thing this provider is, so a hundred Recipes
	// can be found by what they do rather than by whether you remember the
	// company's name.
	//
	// One word from a fixed list, not a free string. The value of a category
	// is that two people reaching for it independently land on the same one,
	// and a free string gives you "payments", "payment", "billing" and "money"
	// within a month. Adding a word is a deliberate change to the list.
	//
	// Deliberately not part of the Recipe's name. Renaming stripe to
	// payments.stripe would break every configuration and every command
	// anybody has already written, to buy a grouping a field gives for
	// nothing.
	Capability string              `yaml:"capability"`
	Upstream   Upstream            `yaml:"upstream"`
	Auth       Auth                `yaml:"auth"`
	Resources  map[string]Resource `yaml:"resources"`
	Routes     []Route             `yaml:"routes"`
	Webhooks   Webhooks            `yaml:"webhooks"`
	Responses  Responses           `yaml:"responses"`
	Errors     map[string]Error    `yaml:"errors"`
	Fixtures   map[string]Fixture  `yaml:"fixtures"`
	// RequiredHeaders are headers a request must carry, mapped to the error
	// name to raise when one is missing. Forgetting Notion-Version is the
	// classic Notion integration bug, and a fake that does not enforce it lets
	// code ship that fails on the first real call.
	RequiredHeaders map[string]RequiredHeader `yaml:"required_headers"`
	// Conformance is the evidence that this Recipe resembles the real provider.
	Conformance []Case `yaml:"conformance"`
}

// assertsName reports whether any case claims something about a field by name.
//
// The comparison is on the last segment of a dotted path, and it is deliberately
// loose: a Recipe declaring data.pagination.next may be asserted as a dotted
// path or as nested maps, and either pins the name down. What it will not
// accept is silence.
func assertsName(cases []Case, field string) bool {
	segments := strings.Split(field, ".")
	leaf := segments[len(segments)-1]

	for _, c := range cases {
		for path := range c.Expect.Matches {
			if mentions(path, leaf) {
				return true
			}
		}

		if mentionedIn(c.Expect.Body, leaf) {
			return true
		}
	}

	return false
}

// mentions reports whether a dotted path names this field at any depth.
func mentions(path, leaf string) bool {
	for _, segment := range strings.Split(path, ".") {
		// Trim any [0] so items[0].next matches next.
		if index := strings.IndexByte(segment, '['); index >= 0 {
			segment = segment[:index]
		}

		if segment == leaf {
			return true
		}
	}

	return false
}

// mentionedIn walks a body assertion looking for the field, at any depth,
// because a nested claim is as good as a dotted one.
func mentionedIn(node any, leaf string) bool {
	switch typed := node.(type) {
	case map[string]any:
		for key, value := range typed {
			if mentions(key, leaf) || mentionedIn(value, leaf) {
				return true
			}
		}
	case []any:
		for _, value := range typed {
			if mentionedIn(value, leaf) {
				return true
			}
		}
	}

	return false
}

// echoesOnly reports whether a case's only claims repeat its own request.
//
// A create that sends {"name": "x"} and asserts {"name": "x"} is testing that
// the request body survived the round trip, which every fake does by
// construction. It reads as evidence and is not, and the failure is invisible
// until somebody breaks the Recipe on purpose and watches the case pass.
func echoesOnly(c Case) bool {
	if len(c.Expect.Body) == 0 {
		return false
	}

	// Either of these puts something under test that cannot be echoed.
	if len(c.Expect.Matches) > 0 || len(c.Expect.Absent) > 0 {
		return false
	}

	sent := map[string]any{}

	for name, value := range c.Request.JSON {
		sent[name] = value
	}

	for name, value := range c.Request.Form {
		sent[name] = value
	}

	if len(sent) == 0 {
		return false
	}

	for name, want := range c.Expect.Body {
		value, present := sent[name]
		if !present {
			return false
		}

		// Form values arrive as text, so compare rendered forms rather than
		// requiring the YAML types to match.
		if !reflect.DeepEqual(value, want) && fmt.Sprint(value) != fmt.Sprint(want) {
			return false
		}
	}

	return true
}

// Case is one checkable claim about the provider's behaviour.
//
// The point of a conformance case is not that the emulator passes it. Any fake
// passes its own tests. The point is provenance: every case cites where the
// expectation came from, and records whether it was observed against the real
// API or only read in the documentation. A developer deciding whether to trust
// this emulator can then read the evidence rather than the marketing.
type Case struct {
	Name string `yaml:"name"`
	// Source cites the provider documentation or transcript the expectation
	// came from. Required: an uncited claim about someone else's API is a
	// guess wearing a test's clothing.
	Source string `yaml:"source"`
	// Verified is the date this case was last checked against the real API,
	// as YYYY-MM-DD. Empty means the expectation was read, not observed, and
	// the report says so rather than quietly counting it as proof.
	Verified string `yaml:"verified"`
	// Fixture is seeded before the case runs. Empty leaves the sandbox as it is,
	// which lets a group of cases build on each other in order.
	Fixture string `yaml:"fixture"`
	// Arm names an entry in the Recipe's errors table to install before this
	// case's request, and only for it.
	//
	// Without this a Recipe's error table is a list of unverified claims.
	// Every failure a conformance suite could reach was one the runtime
	// produces on its own: a 404 for a missing record, a 401 for a bad
	// credential. The interesting entries, the ones describing a declined
	// card or an expired sync token or a rate limit, were declared and never
	// once exercised, so a field could be renamed, a status changed or a
	// nested detail dropped and nothing anywhere would notice.
	//
	// The fault is armed for exactly one request and cleared afterwards, so a
	// case cannot leak a failure into the next one.
	Arm     string      `yaml:"arm"`
	Request Request     `yaml:"request"`
	Expect  Expectation `yaml:"expect"`
}

// Request is the call a conformance case makes.
type Request struct {
	Method  string            `yaml:"method"`
	Path    string            `yaml:"path"`
	Query   map[string]string `yaml:"query"`
	Headers map[string]string `yaml:"headers"`
	// Form sends application/x-www-form-urlencoded, which is what Stripe's own
	// SDKs send. JSON sends a JSON body. A case may set at most one.
	Form map[string]string `yaml:"form"`
	JSON map[string]any    `yaml:"json"`
}

// Expectation is what the provider is claimed to answer.
//
// Body matching is a subset: a case asserts the fields it is making a claim
// about and ignores the rest, so a Recipe can grow a field without invalidating
// every case ever written about it.
type Expectation struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Body    map[string]any    `yaml:"body"`
	// Matches holds dotted field paths to regular expressions, for values that
	// are correct in shape rather than exact, such as generated identifiers.
	Matches map[string]string `yaml:"matches"`
	// HeaderMatches holds response header names to regular expressions, for
	// headers that carry a generated value. A plain `headers` entry compares
	// substrings, which cannot assert that a header is merely present and
	// well-formed.
	HeaderMatches map[string]string `yaml:"header_matches"`
	// AbsentHeaders lists response headers that must not appear.
	//
	// The absence of a header is a claim as real as its presence, and for
	// paging it is the one that terminates the loop: a provider advertises
	// the next page in Link and sends no Link on the last page, so a client
	// that keeps following one until it is gone stops exactly there. An
	// emulator that sent Link on every page would loop forever, and there
	// was no way to write that down.
	AbsentHeaders []string `yaml:"absent_headers"`
	// Absent lists fields that must not appear. Providers are as specific about
	// what they omit as what they send.
	Absent []string `yaml:"absent"`
	// BodyMatches is a regular expression applied to the raw response body,
	// without parsing it.
	//
	// A provider whose failures are plain text has no assertable body
	// otherwise: matches walks a decoded document, so the only thing Trello's
	// text-error case could pin down was its Content-Type. The prose is the
	// part support threads quote and the part a client ends up regex-matching
	// in anger, so it is worth being able to claim.
	BodyMatches string `yaml:"body_matches"`
	// NoBody asserts the response body is empty.
	//
	// This is a positive claim, not the absence of one. Salesforce answers an
	// update with 204 and nothing at all, so a client calling .json() on it
	// throws rather than seeing that the update worked, and an emulator that
	// helpfully returned an object would hide that. An `absent` list cannot
	// express it: absences are vacuously true against an empty body, so a case
	// built from them would pass whatever the emulator sent.
	NoBody bool `yaml:"no_body"`
	// Webhook asserts what the request emitted, which nothing could assert
	// before.
	//
	// Webhook payloads were the largest unverified surface in the project: 85
	// Recipes emit them, the record went in raw rather than shaped, and no
	// case could look. An application's handler written against the emulator
	// could read a field the provider never sends and be entirely green.
	Webhook *WebhookExpectation `yaml:"webhook"`
}

// WebhookExpectation is what a case claims about the webhook its request
// emitted.
//
// The last delivery is the one examined, because a request emits at most one
// event and asserting on "the one this caused" is the only reading that stays
// true as a Recipe grows.
type WebhookExpectation struct {
	// Event is the type the delivery must carry.
	Event string `yaml:"event"`
	// Body asserts dotted paths in the payload, envelope included, so a
	// Recipe declaring its own envelope can pin that too.
	Body map[string]any `yaml:"body"`
	// Matches asserts regular expressions against payload paths.
	Matches map[string]string `yaml:"matches"`
	// Absent names paths the payload must not carry. This is the half that
	// catches an internal field name leaking into a payload.
	Absent []string `yaml:"absent"`
	// None claims the request emitted nothing at all, which is worth being
	// able to say: an event that fires when it should not is as wrong as one
	// that does not fire.
	None bool `yaml:"none"`
}

// Upstream records which real API version this Recipe targets. Without it, a
// Recipe silently rots as the provider moves on.
type Upstream struct {
	API  string `yaml:"api"`
	Docs string `yaml:"docs"`
}

// RequiredHeader is one header a request must carry.
//
// It reads from YAML either as a bare error name, meaning every request needs
// the header, or as a mapping with a methods list, meaning only those methods
// do. The second form exists because Greenhouse only wants On-Behalf-Of on a
// write: reads work without it, so an integration passes every test it has and
// then gets a 403 the first time it tries to change something.
type RequiredHeader struct {
	// Error names the error to raise when the header is missing.
	Error string `yaml:"error"`
	// Methods limits the requirement to those HTTP methods. Empty means all.
	Methods []string `yaml:"methods"`
}

// UnmarshalYAML accepts either a bare error name or the full mapping.
func (h *RequiredHeader) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		return value.Decode(&h.Error)
	}

	type plain RequiredHeader

	return value.Decode((*plain)(h))
}

// Applies reports whether the header is required for this HTTP method.
func (h RequiredHeader) Applies(method string) bool {
	if len(h.Methods) == 0 {
		return true
	}

	for _, m := range h.Methods {
		if strings.EqualFold(m, method) {
			return true
		}
	}

	return false
}

// Auth describes how the provider authenticates callers.
type Auth struct {
	// Scheme is one of: bearer, basic, header, query, none.
	//
	// A query credential travels in the URL, which is worth reproducing
	// precisely because it is a bad idea: URLs end up in access logs, browser
	// history and error reports. Trello and a good deal of older software do
	// it anyway, and an emulator that quietly accepted a header instead would
	// hide the exposure.
	Scheme string `yaml:"scheme"`
	// Header is the header carrying the credential, when scheme is header.
	Header string `yaml:"header"`
	// Param is the query parameter carrying the credential, when the scheme
	// is query.
	Param string `yaml:"param"`
	// Prefix is stripped from the credential before comparison, e.g. "Bearer ".
	Prefix string `yaml:"prefix"`
	// Credential says which half of a basic credential carries the secret:
	// "username" (the default, which is what Twilio does with the account SID)
	// or "password" (Mailgun, whose username is the constant "api"). Checking
	// the wrong half means a bad key is never rejected at all.
	Credential string `yaml:"credential"`
	// Keys are the credentials the emulator accepts. Test keys only — a Recipe
	// must never carry a real secret.
	Keys []string `yaml:"keys"`
	// Pattern accepts any credential matching this regular expression, for
	// schemes where the value is computed per request and cannot be compared
	// against a fixed list.
	//
	// AWS signs every request with SigV4, so the Authorization header is
	// different each time and there is no key to hold. Verifying the signature
	// would mean implementing the algorithm, which is not what this project is
	// for. Checking the shape catches the failure that actually happens —
	// credentials not configured, or the header missing entirely — and the
	// Recipe header has to say plainly that a wrongly signed request is
	// accepted. Silence about that would be worse than the gap.
	Pattern string `yaml:"pattern"`
}

// Responses describes the envelopes a provider wraps its payloads in.
//
// This exists because providers genuinely disagree: Stripe returns
// {object, data, has_more}, GitHub returns a bare array, Shopify nests under a
// resource key. Hardcoding one of them would make every other provider a
// second-class citizen.
type Responses struct {
	List     ListResponse     `yaml:"list"`
	Error    ErrorResponse    `yaml:"error"`
	Resource ResourceResponse `yaml:"resource"`
	Success  SuccessResponse  `yaml:"success"`
}

// SuccessResponse describes what a provider adds to every successful body.
//
// Slack stamps {"ok": true} on everything and its clients check it before
// looking at anything else. A fake that omits it fails at the first line of
// every handler written against the real API.
type SuccessResponse struct {
	Fields map[string]any `yaml:"fields"`
}

// ResourceResponse describes how a single object comes back.
//
// Shopify wraps it under the singular resource name, so a client reads
// body.order.id. Stripe and GitHub return the object itself. Getting this
// wrong is not a cosmetic difference: every field access is one level out.
type ResourceResponse struct {
	// Style is bare (the default) or wrapped.
	Style string `yaml:"style"`
	// Key is the wrapping property name. Empty uses the resource's own name,
	// which is what Shopify does. Cloudflare wraps everything under "result"
	// regardless of what the object is, so the name has to be declarable.
	Key string `yaml:"key"`
	// Array wraps the single object in a list. Xero answers a request for one
	// invoice with {"Invoices": [{...}]}, so client code reads Invoices[0] and
	// anyone expecting an object gets an array with no warning. With Array
	// set, Key defaults to the resource's plural collection name rather than
	// its singular one, because a list of one is still a collection.
	Array bool `yaml:"array"`
}

// ErrorResponse describes the envelope a provider puts failures in.
//
// Stripe nests under "error" with a type and a code; GitHub sends a flat object
// with a message and a documentation link. Code that unwraps one and receives
// the other does not report a helpful failure, it panics.
type ErrorResponse struct {
	// Style is nested (Stripe, the default), flat (GitHub), list (SendGrid,
	// which sends {"errors": [{...}]} because one request can fail several
	// ways at once), string_list (Datadog, which sends the same array with
	// bare strings in it rather than objects) or text (Trello, whose failures
	// are not JSON at all, so a client calling .json() on one throws).
	Style string `yaml:"style"`
	// Key is the property holding the array when the style is list. A dotted
	// name nests, which QuickBooks needs: its failures arrive under
	// Fault.Error rather than at the top level. "-" removes the envelope
	// entirely, which Salesforce needs: its failures are a bare top-level
	// array, so a client reading .message off the response finds undefined and
	// has to index before it can read anything at all.
	Key string `yaml:"key"`
	// MessageField names the property carrying the human-readable message when
	// the style is flat. Empty means "message". Set it to "-" to omit the
	// message entirely, which Slack does: its errors are a code and nothing
	// else, and inventing prose the provider never sends is still infidelity.
	MessageField string `yaml:"message_field"`
	// CodeField names the property carrying the error code in a flat envelope.
	// Twilio sends one and its clients switch on it; GitHub does not send one
	// at all, so this stays empty unless a Recipe claims otherwise. As with
	// MessageField, "-" omits the code, which the nested style needs too:
	// Airtable nests its error but sends only a type and a message.
	CodeField string `yaml:"code_field"`
	// CodeType says whether the code is sent as a number or as a string.
	//
	// Empty infers it from the value: all digits becomes a number, anything
	// else stays text. That inference is right for Twilio, whose codes really
	// are integers, and wrong for Adyen, whose "000" is a string and loses its
	// leading zeros on the way through. Inferring a provider's behaviour from
	// the shape of a literal is a guess, so a Recipe that knows can say, and
	// one that says overrides the guess.
	CodeType string `yaml:"code_type"`
	// StatusField names a property echoing the HTTP status inside the body,
	// which Twilio does.
	StatusField string `yaml:"status_field"`
	// TypeField names a property carrying the error category in a flat
	// envelope. Plaid sends error_type and its clients switch on that before
	// they look at the code, because the category decides whether to retry,
	// re-authenticate or give up. As with CodeField, "-" omits the category,
	// which the nested style needs too: Vercel nests its error and sends only
	// a code and a message.
	TypeField string `yaml:"type_field"`
	// Fields are constants the provider adds to every error, such as GitHub's
	// documentation_url.
	Fields map[string]any `yaml:"fields"`
}

// ListResponse describes how a collection is returned.
type ListResponse struct {
	// Style is one of: envelope (Stripe), bare (GitHub), wrapped (Shopify),
	// map (Pusher, whose channels arrive as an object keyed by channel name
	// rather than an array, so looping over it as a list finds nothing and a
	// channel with no occupants is absent from the object entirely rather
	// than present with a zero).
	// Empty means envelope, which keeps existing Recipes working.
	Style string `yaml:"style"`
	// Key is the wrapping property name, required when style is wrapped. A
	// dotted name nests, so a collection can sit two levels down: Segment
	// answers with data.sources rather than a top-level array.
	Key string `yaml:"key"`
	// URL asks the envelope to echo the request path, which Stripe does.
	URL bool `yaml:"url"`
	// CursorField names a property carrying the next cursor. Most providers do
	// not send one: Stripe expects the caller to pass the last id back as
	// starting_after. Leaving it empty is therefore the faithful default, and
	// setting it is a deliberate claim that the provider really sends it.
	//
	// A dotted name nests, so Slack's response_metadata.next_cursor is
	// expressible without a second mechanism.
	CursorField string `yaml:"cursor_field"`
	// CountField names a property carrying how many records matched in total,
	// which is not the same as how many are on this page. Zendesk sends one and
	// a pagination UI cannot be built without it.
	CountField string `yaml:"count_field"`
	// PagesField names a property carrying how many pages the whole set makes
	// at this page size, which is a different quantity from CountField.
	//
	// Documenso's list envelope is {documents, totalPages} and nothing else,
	// and totalPages was declared as the count field -- so three documents at
	// ten per page reported three rather than one. That is worse than an
	// invented field: the name is real and the number is plausible, so a
	// client looping while page <= totalPages asks for two pages that do not
	// exist and reads them as empty results rather than as a mistake.
	//
	// An empty set is nought pages here. Providers differ about whether it is
	// nought or one, and nought is the reading that stops a loop rather than
	// sending it after a page with nothing in it.
	PagesField string `yaml:"pages_field"`
	// EntryField makes each entry in the collection that one field's value
	// rather than the whole record.
	//
	// Plenty of APIs answer a listing with an array of identifiers and keep
	// the object for the fetch beside it. DynamoDB's ListTables sends
	// TableNames as an array of strings; SQS's ListQueues sends QueueUrls the
	// same way. Both Recipes emitted arrays of objects under those names, so
	// a client doing TableNames.forEach(name => describe(name)) received
	// objects and called describe([object Object]).
	//
	// It belongs on the route rather than the Recipe, because the listing and
	// the fetch disagree by design: DescribeTable answers with the table and
	// ListTables answers with its name.
	EntryField string `yaml:"entry_field"`
	// LinkHeader advertises the next page in an RFC 5988 Link response
	// header rather than in the body.
	//
	// Five providers modelled here page that way and it is the mechanism
	// their own documentation leads with: GitHub, Ably, WordPress, Greenhouse
	// and Buildkite. Buildkite's says it plainly -- "the pagination
	// information can be found in the Link HTTP response header" -- and Ably
	// pages by nothing else at all.
	//
	// Without it the page size works and the next page is unreachable, which
	// is the quietest way for a listing to be wrong: one page comes back, it
	// is a correct page, and the loop that should have asked for the second
	// one has nothing to follow.
	//
	// Only next is emitted. Providers also advertise prev, first and last,
	// and last needs a total this does not have -- so a client that follows
	// next walks the whole collection here, and one that reads last finds
	// nothing. That is stated rather than guessed at.
	LinkHeader bool `yaml:"link_header"`
	// PageField and LimitField name properties echoing the page number and
	// the page size the request asked for.
	//
	// A constant cannot do this job, and putting one there is worse than
	// leaving the field out. Algolia answers every search with the page it
	// served and the page size it used, and the Recipe declared them as the
	// constants 0 and 20 -- so a client that asked for page 3 was told it was
	// looking at page 0, by a field whose entire purpose is to say where you
	// are. Paging code that trusts the response rather than its own counter
	// reads that as "still on the first page" forever.
	PageField  string `yaml:"page_field"`
	LimitField string `yaml:"limit_field"`
	// CountAsString sends the counts as strings. Docusign does, and code that
	// compares totalSetSize to a number never matches, so emitting a number
	// here would quietly fix a bug the caller has to handle.
	CountAsString bool `yaml:"count_as_string"`
	// HasMoreField names a boolean saying whether more pages remain. The
	// envelope style always sends has_more because Stripe does; other styles
	// send one only when the Recipe says so.
	HasMoreField string `yaml:"has_more_field"`
	// OmitWhenEmpty leaves the collection key out entirely when there is
	// nothing to send, rather than sending an empty array.
	//
	// SQS does this, and it is the difference between a consumer that waits on
	// an idle queue and one that throws. ReceiveMessage with nothing to give
	// answers 200 with no Messages key at all, so
	// `for (const m of response.Messages)` fails on the quietest possible
	// input. An emulator sending [] is the helpful kind of wrong: every test
	// passes and the first quiet minute in production does not.
	OmitWhenEmpty bool `yaml:"omit_when_empty"`
	// FinalField names a field sent only on the last page of a list, and left
	// out of every page before it.
	//
	// Google Calendar sends nextSyncToken this way and Microsoft Graph sends
	// @odata.deltaLink. The two tokens are not interchangeable: a page token
	// resumes the listing you are in the middle of, and a sync token starts a
	// later incremental one. Only one of them is ever present, so code that
	// reads whichever it finds on the first response and calls it "the token"
	// stores a page token, and the next sync either replays from the
	// beginning or fails with an error that names neither field.
	//
	// Sending it on every page would be the helpful kind of wrong. The
	// caller's storage logic would work locally against any list short enough
	// to fit one page, and break on the first calendar busy enough to need
	// two.
	FinalField string `yaml:"final_field"`
	// CompleteField names a boolean saying the opposite: that no pages remain.
	//
	// Salesforce sends done, and false is the interesting value. Modelling it
	// as a negated has_more would be a lie about the field's name, and
	// modelling it as has_more would invert its meaning, so a query that
	// matched more rows than it returned would claim to be finished. Code that
	// ignores done processes a prefix of its own result set and is never told.
	CompleteField string `yaml:"complete_field"`
	// Fields are constants added to a list response only. Notion stamps
	// {"object": "list"} on a collection and {"object": "page"} on a single
	// page, so the two cannot share one set of envelope constants.
	Fields map[string]any `yaml:"fields"`
	// EntryStyle wraps each item in the collection under the resource's own
	// name when set to "wrapped". Chargebee answers a subscription list with
	// {"list": [{"subscription": {...}}]}, so a client reads
	// list[0].subscription.id and anyone indexing straight into the item
	// finds nothing.
	EntryStyle string `yaml:"entry_style"`
	// CollapseSingle sends a collection of one as the object rather than as a
	// list of one.
	//
	// Tradier documents it in its own words: if you have a single order, it
	// will be returned as a JSON obj/dict whereas multiple orders will be
	// returned as an array. Every API that grew out of XML does this, because
	// a single child element and a repeated one are the same thing there and
	// are not the same thing in JSON.
	//
	// It is the dangerous half of an axis this format already had. Xero sends
	// one invoice as a list of one, which resource.array describes, and that
	// is the safe direction: a client written for the list keeps working. This
	// is the other way round, and a client written against a fixture with two
	// records in it crashes the first time production has one.
	CollapseSingle bool `yaml:"collapse_single"`
}

// Resource is an object type the provider exposes.
type Resource struct {
	// Collection is the plural name the provider wraps lists in, e.g. "orders"
	// for an order. Declared rather than derived: guessing English plurals is
	// exactly the kind of cleverness that produces a fake which is subtly
	// wrong for "person", "category" or "status".
	Collection string `yaml:"collection"`
	ID         ID     `yaml:"id"`
	// Alias names a second field a path may address this resource by.
	//
	// Jira answers /issue/10001 and /issue/PLAT-42 with the same issue. The
	// two identifiers are not interchangeable underneath: the numeric one is
	// permanent and the key changes when the issue moves project, so anything
	// that stored the readable one holds a dangling reference and gets no
	// error saying so. An emulator accepting only the identifier would reject
	// half the calls that work against the real API.
	Alias  string           `yaml:"alias"`
	Fields map[string]Field `yaml:"fields"`
	// Constants are fields the provider always sends with a fixed value, such
	// as Stripe's object discriminator and livemode flag. Unlike a default they
	// cannot be overridden by the caller, because the provider does not let you
	// override them either. Applications really do branch on these.
	Constants map[string]any `yaml:"constants"`
}

// ID describes how the provider mints identifiers. Getting this right matters
// more than it looks: applications routinely parse or prefix-match IDs.
type ID struct {
	// Style is one of: prefixed (cus_abc123), numeric (1, 2, 3), timestamp
	// (1767225600.000100, which is how Slack identifies a message), opaque
	// (a bare random string, which is what SendGrid returns as a message id)
	// uuid (Notion, and most APIs designed after about 2015), hex (Intercom,
	// and anything whose identifiers came out of MongoDB) or digits (Discord
	// snowflakes, and any provider whose ids are long numeric strings that must
	// not be parsed as numbers).
	// Empty means prefixed.
	Style  string `yaml:"style"`
	Prefix string `yaml:"prefix"`
	Length int    `yaml:"length"`
	// Field is the property the provider returns the identifier in. Empty means
	// "id". Twilio calls it "sid" everywhere, and code that reads response.id
	// against Twilio gets nothing at all. A dotted name nests, which is how
	// Contentful keeps the identifier at sys.id.
	//
	// "-" means the provider does not echo it at all. Some resources are
	// addressed by a key that appears only in the path: Marqeta's balance is
	// fetched at /v3/balances/{token} and the body that comes back carries no
	// token anywhere. Cauldron still keys the record internally, because it
	// has to be found somehow, but emitting an identifier the provider never
	// sends would put a field on the wire that real code cannot rely on.
	Field string `yaml:"field"`
	// Type is the JSON type the identifier travels as: "string" (the default)
	// or "number".
	//
	// Identifiers are minted and stored as strings, because that is the only
	// form every style shares and the only form a path parameter arrives in.
	// What they are on the wire is a separate question, and the two answers
	// disagree more often than they agree. GitHub sends an issue id as the
	// number 1. HubSpot sends a contact id as the string "1". Meilisearch
	// sends a task uid as a number and Jira sends an issue id as a string.
	//
	// It is not cosmetic. id === 1 fails against a string, typeof id ===
	// "number" fails, and a schema declaring "type": "integer" rejects the
	// response outright. An emulator answering with a string where the
	// provider answers with a number commits the exact class of bug it exists
	// to catch.
	//
	// The default stays "string" because changing it silently would rewrite
	// the wire shape of every shipped Recipe at once, and each one has to be
	// checked against its provider rather than assumed.
	Type string `yaml:"type"`
	// CarriedBy names the field that holds the identifier when the provider
	// does not send it under a name of its own.
	//
	// Dwolla is HAL: there is no id property on anything, and identity lives
	// in _links.self.href with the identifier as its last segment. So the
	// record is addressable and the id is genuinely absent, which are two
	// things that are usually not true together, and a Recipe that only said
	// field: "-" would be claiming a listing hands back records nothing can
	// identify.
	//
	// It is a declaration rather than a behaviour: nothing reads it at
	// runtime. What it does is answer the question a reader of the Recipe
	// asks first -- if there is no id, how do I address one of these -- and
	// let the validator tell a described absence from an undescribed one.
	CarriedBy string `yaml:"carried_by"`
}

// Filter is a query parameter that narrows a listing to records whose field
// matches, and usually narrows it whether or not the caller asked.
type Filter struct {
	// Param is the query parameter's name.
	Param string `yaml:"param"`
	// Field is the record field it matches against.
	Field string `yaml:"field"`
	// Default is the value applied when the parameter is absent. Empty means
	// the filter only applies when the caller supplies it, which is the less
	// interesting half: a filter nobody asked for is the one that surprises
	// people.
	Default string `yaml:"default"`
	// All is the value that turns the filter off, for providers that have one.
	// GitHub and Alpaca both spell it "all". Empty means there is no way to
	// ask for everything, which is itself worth knowing.
	All string `yaml:"all"`
	// Values expands a parameter value into the set of field values it
	// covers, for the filters whose vocabulary is not the field's.
	//
	// Alpaca's order listing takes status=open, and "open" is not a status
	// any order holds. It is a bucket: new and partially_filled are open,
	// filled and canceled are closed, and nothing in the record says which
	// bucket its status belongs to. A filter that matched the word literally
	// would hide every partially filled order, which is precisely the order
	// that most needs to be visible, since it is a real position.
	//
	// A value with no entry here matches itself, so most filters need none.
	Values map[string][]string `yaml:"values"`
}

// Field is a single attribute on a resource.
type Field struct {
	// Type is string, integer, boolean, timestamp (a Unix integer in seconds,
	// which is what Stripe and Twilio send), timestamp_ms (milliseconds, which
	// is what Clerk and most JavaScript-first APIs send) or datetime (an RFC
	// 3339 string, which is what GitHub, HubSpot and most newer APIs send).
	// The difference is not cosmetic: one parses as a number and the other
	// does not, and a factor of a thousand puts a date in 1970 or in the year
	// 55000.
	// list is a sequence of anything, and it exists because several providers
	// send one where a client expects an object. WooCommerce's meta_data is an
	// array of {id, key, value} objects, so order.meta_data.some_key finds
	// undefined and the value is sitting in the array under a key that has to
	// be searched for. Declaring the field as a list says that on purpose;
	// leaving it untyped would emit the same bytes while claiming nothing.
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
	Default  any    `yaml:"default"`
	// Stamped decides whether a timestamp or datetime field is filled in
	// automatically when the caller does not supply it. Nil means yes, which
	// is right for created_at and updated_at.
	//
	// Set it to false for a field whose absence is the meaningful state: a
	// Webflow site that has never been published has no lastPublished, and a
	// Typeform response that was abandoned has no submitted_at. Stamping
	// those makes the emulator claim an event happened that did not, which is
	// the kind of infidelity a test can never catch.
	Stamped *bool `yaml:"stamped"`
	// NullWhenUnset sends the field as null rather than leaving it out.
	//
	// Absent and null are different on the wire and providers disagree about
	// which they use, so the format has to be able to say both. Bandwidth
	// leaves errorCode out of the messages that worked, so a client testing
	// for null never sees one. Alpaca sends every timestamp on every order and
	// leaves the ones that have not happened as null, so a client testing for
	// the key's existence finds it and reads nothing out of it. Each of those
	// is a real bug in code written against the other assumption.
	//
	// It implies no stamping, because a field with a value is not unset.
	NullWhenUnset bool `yaml:"null_when_unset"`
	// In nests this field under a sub-object on the wire. HubSpot puts every
	// business attribute under "properties" and leaves only id, timestamps and
	// archived at the top level, so a client reads contact.properties.email.
	// The store stays flat; only the shape on the wire changes, and requests
	// are flattened back on the way in.
	//
	// "-" nests it nowhere: the record holds the field and the wire never
	// carries it. A route's scope needs this. A partition that lives in the
	// path has to be a field, because that is how the record is partitioned,
	// and most providers do not repeat it in the body -- Fly does not send
	// app_name on a machine, Hetzner does not send its collection on a point,
	// Tradier does not say which account an order is in.
	//
	// Before this the only way to say so was a route's returns naming every
	// other field, which was twenty-three names to hide one on Fly, and which
	// says nothing at all when the same resource is served by two routes. An
	// audit found 115 scope fields across 37 Recipes going onto the wire with
	// no case mentioning them; some of those providers really do echo the
	// partition and each one has to be read before it is changed, so this is
	// the tool for the ones that have been.
	In string `yaml:"in"`
	// As is the name this field takes on the wire, when it differs from the
	// name it is stored under.
	//
	// Nesting alone was not enough. A field's own name is the key inside the
	// sub-object, so two fields could not share a key under different parents:
	// a resource wanting both title.rendered and content.rendered had to call
	// one of them something else, and what it got called leaked onto the wire.
	// Thirty-one fields across six Recipes ended up emitting amount.amount_value
	// where Adyen sends amount.value, and title.title_rendered where WordPress
	// sends title.rendered. Every conformance case about them passed, because
	// they asserted the shape the emulator produced.
	As string `yaml:"as"`
}

// WireName is the key this field takes in a response.
func (f Field) WireName(name string) string {
	if f.As != "" {
		return f.As
	}

	return name
}

// Route binds an HTTP method and path to an operation on a resource.
type Route struct {
	Method   string `yaml:"method"`
	Path     string `yaml:"path"`
	Resource string `yaml:"resource"`
	// Operation is one of: create, get, list, update, delete.
	Operation string `yaml:"operation"`
	// Scope names the path parameters that partition this resource, e.g.
	// [owner, repo] for /repos/{owner}/{repo}/issues. Scoped requests only
	// ever see records whose matching fields agree, and creates stamp them.
	//
	// Path parameters left out of scope are ignored, which is how an API
	// version segment like /admin/api/{version}/orders stays a path parameter
	// without becoming a filter.
	Scope []string `yaml:"scope"`
	// Status is the success status this route returns. Empty means 200. GitHub
	// answers a create with 201, Stripe with 200, and a client checking for one
	// exact code is not being unreasonable.
	Status int `yaml:"status"`
	// Fields are constants this route adds to its response body, on top of
	// whatever the Recipe-wide response constants already put there.
	//
	// A one-path API needs them. ShipHero puts request_id and complexity
	// beside each connection -- data.orders.complexity, data.products.complexity
	// -- so the key depends on which query was asked, and a Recipe-wide
	// constant would stamp the orders metadata onto a products response. A
	// dotted name nests, the same way the Recipe-wide ones do.
	Fields map[string]any `yaml:"fields"`
	// Selects disambiguates several routes that share one path by what the
	// request body asks for.
	//
	// A GraphQL API is one path and one method, so the path cannot say which
	// route should answer. What can is the query itself: a request naming
	// `orders` wants the orders route and one naming `products` wants that
	// one. Selects holds the root field to look for, and the route matches
	// only when the body's query mentions it.
	//
	// This does not parse GraphQL and does not pretend to. It looks for the
	// field named as a whole word -- a substring match sent `viewer` queries
	// to a route selecting `me`, because "name" contains it -- which is
	// enough to pick a fixture and is exactly the bargain every Recipe here
	// already makes: model what comes back, not how the provider decided it.
	//
	// Seven providers were unreachable without this -- Linear, Monday, Attio,
	// New Relic, Railway, ShipHero, and half of Fly.io -- and each had been
	// recorded as its own judgement call rather than as one missing feature.
	Selects string `yaml:"selects"`
	// IDFrom says where the identifier comes from when it is not a path
	// parameter: "query:channel" or "body:channel". A body name may be
	// dotted, because a provider that puts the identifier in the body does
	// not always put it at the top of one -- DynamoDB's GetItem takes
	// {"Key": {"id": {"S": "..."}}}, which is body:Key.id.S. Slack and every other
	// RPC-shaped API put it in the query string or the body, and without this
	// the format could only describe APIs that happen to be RESTful.
	//
	// "auth" is the third case and it is not a location: it says the request
	// carries no identifier at all and the provider answers about whoever the
	// credentials belong to. GitHub's /user, Stripe's /v1/account, Slack's
	// auth.test and Backblaze's b2_authorize_account are all this shape, and
	// none of them could be described before, because both other forms name a
	// place to read from and the whole point of these routes is that there is
	// nothing to read.
	IDFrom string `yaml:"id_from"`
	// EmptyBody sends no body at all. SendGrid accepts a send with 202 and
	// nothing else, and a client that calls .json() on that response throws.
	// An emulator that helpfully returns an object hides the bug.
	EmptyBody bool `yaml:"empty_body"`
	// Headers are response headers this route sets. "{id}" is replaced with
	// the record's identifier, which is how SendGrid hands back the message id
	// a client needs to correlate a later event with the send.
	Headers    map[string]string `yaml:"headers"`
	Pagination Pagination        `yaml:"pagination"`
	// Filters are query parameters that narrow a listing, and the reason they
	// exist is the default rather than the filtering.
	//
	// GitHub's issue listing answers with open issues unless you ask for
	// otherwise. Alpaca's order listing answers with open orders. Both are
	// documented and both are forgotten, and the failure they produce is the
	// worst-shaped one there is: a client places an order, the order fills, the
	// client lists its orders and sees nothing, and concludes the order never
	// existed. Nothing errored. The list was correct.
	//
	// An emulator that returns everything is being helpful in the direction
	// that hides this. Cauldron did exactly that until this existed, and
	// GitHub's own Recipe had a closed issue in its fixture that the listing
	// returned and the real API would not.
	Filters []Filter `yaml:"filters"`
	// Beside names other resources whose records travel in the same response
	// body, each under its own collection name.
	//
	// One endpoint, several collections, and the fact that it is one endpoint
	// is the point. GoCardless Bank Account Data answers a request for
	// transactions with a booked array and a pending array in one body: the
	// same purchase appears in pending first and booked later, with a
	// different identifier, so code that merges the two arrays counts it
	// twice. Describing that as two endpoints would lose the thing worth
	// describing, and describing only one of the arrays would answer with a
	// shape no bank sends.
	//
	// The scope applies to all of them, because they are one request. Paging
	// applies only to the route's own resource, because that is the one the
	// cursor refers to.
	Beside []string `yaml:"beside"`
	// LookupBy names the field the value from IDFrom is matched against, for
	// the routes that address a record by something that is not its
	// identifier.
	//
	// SQS deletes a message by the receipt handle from a receive, and a
	// receipt handle is deliberately not a message id: it is issued per
	// receive, two consumers holding two handles for the same message is
	// normal, and a handle from an earlier receive is stale. Anything keyed
	// by a natural key -- an email address, an external reference, a slug --
	// has the same shape.
	//
	// IDFrom says where the value comes from; this says what it is compared
	// with. Without it a handle was looked up as though it were an id, found
	// nothing, and every delete failed.
	LookupBy string `yaml:"lookup_by"`
	// MatchesHeader names request headers whose values pick this route, for
	// the APIs where the path does not say which operation you meant.
	//
	// The AWS JSON protocol is the reason it exists: every operation is a
	// POST to the root and the operation is named in X-Amz-Target. Without a
	// way to route on that, the three AWS Recipes here encoded the operation
	// in the path instead -- /ListSecrets, /tables, /queues -- and served
	// URLs AWS does not have. A client can be written entirely against those
	// paths, pass every test, and be entirely wrong.
	//
	// It is the same shape as selects, which tells GraphQL routes apart by a
	// word in the query body: one path, several routes, distinguished by
	// something that is not the path. A route declaring it beats an
	// equally-scoring route that declares nothing, so a Recipe can have a
	// fallback for the operations it does not model.
	MatchesHeader map[string]string `yaml:"matches_header"`
	// List overrides the Recipe-wide list envelope for this route.
	//
	// A provider's listings do not always share a shape. Clerk's users and
	// sessions answer with bare arrays and its organisations with
	// {data, total_count}; Algolia's browse carries a cursor its search does
	// not have. A Recipe-wide envelope makes one of those wrong, and the
	// wrongness is the expensive kind: code written against the emulator
	// reads response.data.map(...) and receives an array from the provider,
	// where .data is undefined.
	//
	// Only the fields set here are overridden; the rest are inherited. A
	// string field set to "-" is cleared rather than inherited, which is how
	// a route says the provider sends nothing there.
	List *ListResponse `yaml:"list"`
	// Returns limits the response to the named fields, for the routes that
	// answer with less than the record they touched.
	//
	// Jira's create hands back an id, a key and a URL and none of them is the
	// issue, so anything reading created.fields.summary gets undefined and a
	// suite asserting on the create response is asserting on almost nothing.
	// Plenty of APIs do this and it is always the same surprise, because the
	// convention everywhere else is that a create echoes what you sent.
	//
	// Echoing the whole record would be the helpful kind of wrong: the caller
	// would read fields back that the provider never sends, locally, for as
	// long as the test suite is the only thing calling it.
	Returns []string `yaml:"returns"`
	// DeletedBody says what a delete answers with, for the providers that
	// answer with something.
	//
	// Empty means no body and a 204, which is what most providers do and what
	// this used to do for none of them. Every delete fabricated Stripe's
	// receipt — an id, an object discriminator and deleted: true — using keys
	// no Recipe declares, on 31 of 35 routes whose providers send nothing at
	// all. So await response.json() succeeded locally and threw
	// SyntaxError: Unexpected end of JSON input in production, and code
	// branching on response.deleted === true was reading undefined against the
	// real API.
	//
	// "receipt" is that Stripe shape, for the providers it is actually true
	// of. "record" answers with the deleted object, for the providers that
	// hand it back.
	DeletedBody string `yaml:"deleted_body"`
	// Error names a failure from the Recipe's own table that this route always
	// answers with, whatever the request. It is how a retired endpoint is
	// described.
	//
	// Jira's old search path answers 410 Gone to the thousands of integrations
	// still calling it, and 410 rather than 404 is the entire message: the
	// path was right, the endpoint is gone, and retrying will not help. An
	// emulator that let the path fall through to its unknown-route handler
	// would answer 404, and a client branching on the difference would take
	// the wrong branch locally and the right one in production, which is the
	// hardest kind of disagreement to notice.
	//
	// A route declaring one needs no resource and no operation, because it
	// never reaches either.
	Error string `yaml:"error"`
}

// Pagination describes how a list endpoint pages.
type Pagination struct {
	// Style is one of: cursor, offset, page.
	Style string `yaml:"style"`
	Limit int    `yaml:"limit"`
	// LimitParam names the query parameter carrying the page size, for the
	// providers that do not call it "limit". Google Calendar calls it
	// maxResults, GitHub calls it per_page, Salesforce does not accept one at
	// all.
	//
	// "-" says the provider accepts no name for it, which is different from
	// leaving it empty: empty falls back to reading "limit", and "-" reads
	// nothing and keeps the declared page size. Datadog's event listing fixes
	// the page at a thousand and takes no size parameter.
	//
	// It matters more than it looks. An emulator that only understands "limit"
	// ignores the size the caller asked for and answers with its own default,
	// and for a fixture of four records that default is the whole collection.
	// The response has no next page in it, so the paging loop the client
	// carefully wrote executes exactly once and every test of it passes
	// without ever taking the branch. The first collection large enough to
	// page is in production.
	//
	// Declaring it also makes "limit" inert, which is the faithful part:
	// Google does not accept limit, and an emulator that quietly honours both
	// spellings lets a typo work locally.
	LimitParam string `yaml:"limit_param"`
	// CursorParam names the query parameter carrying the position to resume
	// from, for the providers that call it neither cursor nor starting_after.
	// Google calls it pageToken and Shopify calls it page_info.
	CursorParam string `yaml:"cursor_param"`
	// FirstPage is the number the provider gives its first page, for the page
	// style. Empty means one.
	//
	// Providers disagree, and the disagreement is invisible: Algolia,
	// Elasticsearch and everything shaped like them count from nought, so
	// page 1 is the second page. Read as though it were the first, a client
	// asking for page 1 is handed page 0 again -- the same record twice, no
	// error, and a loop that either never terminates or quietly returns
	// duplicates. That is the off-by-one-page bug positionOf already warned
	// about, in the direction nobody had checked.
	FirstPage *int `yaml:"first_page"`
	// In is where the parameters travel: "query" (the default) or "body".
	//
	// A listing reached by POST usually carries its paging in the JSON body,
	// and reading it from the query string means reading nothing at all. What
	// that produced was worse than an error: the caller's limit was ignored,
	// so the first request answered with the entire collection and no next
	// page, and a paging loop written against it ran exactly once and looked
	// correct. Dropbox shipped that way, and its own conformance case sent
	// ?limit=1 -- a parameter Dropbox does not read -- because the case was
	// written against what came out rather than against the provider.
	//
	// A dotted name nests, because a provider that puts paging in the body
	// often puts it inside something: Plaid's count and offset live under
	// options.
	In string `yaml:"in"`
}

// Webhooks describes what the provider sends back, and how it signs it.
type Webhooks struct {
	Events  []string `yaml:"events"`
	Signing Signing  `yaml:"signing"`
	// Payload is the envelope the provider wraps the changed record in.
	//
	// Empty keeps Cauldron's default, which is Stripe's shape: an id, a type,
	// a created timestamp and the record under data.object. That default was
	// fine while Stripe was the only Recipe and became a quiet lie as the
	// collection grew, because almost nobody else uses it. Adyen wraps every
	// notification in an array of NotificationRequestItem and reports success
	// as the string "true", so a truth test on it is also true for "false" —
	// which is exactly the sort of thing this project exists to reproduce, and
	// could not be expressed at all until a Recipe could describe its own
	// envelope.
	//
	// Four placeholders are substituted anywhere inside the template:
	//
	//   {event}       the event name
	//   {id}          the delivery identifier
	//   {created}     the sandbox clock as a Unix timestamp, as a number
	//   {created_iso} the same instant as RFC 3339 text
	//
	// A string value of exactly "{object}" is replaced by the record. A map
	// *key* of "{object}" merges the record's fields into that map instead,
	// which is what Adyen needs: its notification item carries the payment's
	// fields alongside eventCode and success rather than under them.
	Payload map[string]any `yaml:"payload"`
}

// Signing describes webhook payload signing.
type Signing struct {
	// Scheme is one of: hmac-sha256, none.
	Scheme string `yaml:"scheme"`
	Header string `yaml:"header"`
	Secret string `yaml:"secret"`
}

// Error is a named failure mode that `cauldron fault` can inject.
type Error struct {
	Status int    `yaml:"status"`
	Code   string `yaml:"code"`
	// Type is the provider's error category, which is often a much smaller set
	// than the codes. Stripe has four types and dozens of codes, and client
	// libraries switch on the type. Empty falls back to the code, which is
	// wrong often enough that every Recipe should set it.
	Type    string            `yaml:"type"`
	Message string            `yaml:"message"`
	Headers map[string]string `yaml:"headers"`
	// Fields are extra body properties this failure carries, merged over the
	// Recipe-wide ones. Dropbox describes each failure with its own nested
	// union, so a single set of constants would make every error claim to be
	// the same one.
	Fields map[string]any `yaml:"fields"`
}

// Fixture is a named seed dataset: resource name to a list of records.
type Fixture map[string][]map[string]any

// indexedSegment matches a path segment naming an array position, e.g. to[0].
var indexedSegment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*(\[\d+\])+$`)

// plainSegment matches an ordinary object key.
//
// Hyphens are allowed because JSON keys have them: npm sends dist-tags and
// plenty of providers send header-shaped names. This rule refused the first
// hyphenated key it ever met, which is what a pattern written against the
// Recipes that happened to exist will do. A dot is still refused, because the
// runtime splits on it and a dot inside a segment would silently mean two
// levels rather than one.
var plainSegment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*$`)

var (
	namePattern    = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	datePattern    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	validSchemes    = []string{"bearer", "basic", "header", "query", "none"}
	validOperations = []string{"create", "get", "list", "update", "delete"}
	validPagination = []string{"", "cursor", "offset", "page"}
	// The categories a Recipe may declare. Kept short on purpose: a taxonomy
	// with sixty words in it is a taxonomy nobody can hold in their head, and
	// the point of this field is that somebody looking for a payments
	// emulator finds all of them at once.
	validCapabilities = []string{
		"payments", "banking", "accounting", "tax", "payroll",
		"email", "sms", "chat", "push", "voice",
		"auth", "identity", "crm", "support", "marketing", "brokerage",
		"commerce", "shipping", "storage", "database", "queue",
		"search", "cdn", "hosting", "observability", "analytics",
		"flags", "ci", "vcs", "issues", "docs",
		"calendar", "files", "media", "ai", "signing",
		"scheduling", "hr", "forms", "cms", "infrastructure",
	}
	validDeletedBody = []string{"", "receipt", "record"}
	validSigning     = []string{"", "none", "hmac-sha256"}
	validListStyles  = []string{"", "envelope", "bare", "wrapped", "map"}
	validErrStyles   = []string{"", "nested", "flat", "list", "string_list", "text"}
	validCodeTypes   = []string{"", "string", "number"}
	validIDTypes     = []string{"", "string", "number"}
	validIDStyles    = []string{"", "prefixed", "numeric", "timestamp", "opaque", "uuid", "hex", "digits"}
	validFieldTypes  = []string{"", "string", "integer", "number", "boolean", "timestamp", "timestamp_ms", "datetime", "list"}
	// The types Cauldron fills in from the sandbox clock, and therefore the
	// only ones a stamped declaration can affect.
	timeFieldTypes = []string{"timestamp", "timestamp_ms", "datetime"}

	knownMethods = map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true,
		"DELETE": true, "HEAD": true, "OPTIONS": true,
	}
)

// Load reads and validates a Recipe from a YAML file.
func Load(path string) (*Recipe, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(contents)
}

// Parse decodes and validates a Recipe from YAML bytes.
func Parse(contents []byte) (*Recipe, error) {
	var r Recipe

	decoder := yaml.NewDecoder(strings.NewReader(string(contents)))
	decoder.KnownFields(true)

	if err := decoder.Decode(&r); err != nil {
		return nil, fmt.Errorf("recipe is not valid YAML: %w", err)
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	return &r, nil
}

// ValidationError collects every problem with a Recipe, so an author sees all
// of them at once instead of fixing one and rerunning.
type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	if len(e.Problems) == 1 {
		return "invalid recipe: " + e.Problems[0]
	}

	return fmt.Sprintf("invalid recipe:\n  - %s", strings.Join(e.Problems, "\n  - "))
}

// Validate checks the Recipe is internally consistent.
func (r *Recipe) Validate() error {
	var problems []string

	add := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	}

	if r.Name == "" {
		add("recipe name is required")
	} else if !namePattern.MatchString(r.Name) {
		add("recipe name %q must be lowercase words separated by hyphens", r.Name)
	}

	if r.Version == "" {
		add("version is required")
	} else if !versionPattern.MatchString(r.Version) {
		add("version %q must look like 1.2.3", r.Version)
	}

	if r.Upstream.API == "" {
		add("upstream.api is required so the Recipe records which API version it targets")
	}

	if r.Capability == "" {
		add("capability is required: a hundred Recipes are only findable by what they do")
	} else if !contains(validCapabilities, r.Capability) {
		add("capability %q must be one of %s", r.Capability, strings.Join(validCapabilities, ", "))
	}

	if r.Auth.Scheme != "" && !contains(validSchemes, r.Auth.Scheme) {
		add("auth.scheme %q must be one of %s", r.Auth.Scheme, strings.Join(validSchemes, ", "))
	}

	if r.Auth.Credential != "" && r.Auth.Credential != "username" && r.Auth.Credential != "password" {
		add("auth.credential %q must be username or password", r.Auth.Credential)
	}

	if r.Auth.Credential != "" && r.Auth.Scheme != "basic" {
		add("auth.credential only applies to the basic scheme")
	}

	if r.Auth.Scheme == "query" && r.Auth.Param == "" {
		add("auth.param is required when the scheme is query")
	}

	if r.Auth.Param != "" && r.Auth.Scheme != "query" {
		add("auth.param only applies to the query scheme")
	}

	if len(r.Resources) == 0 {
		add("at least one resource is required")
	}

	for _, name := range sortedKeys(r.Resources) {
		resource := r.Resources[name]

		if !contains(validIDStyles, resource.ID.Style) {
			add("resource %q has id.style %q, which must be one of %s", name, resource.ID.Style, strings.Join(validIDStyles[1:], ", "))
		}

		if resource.ID.CarriedBy != "" {
			if _, declared := resource.Fields[resource.ID.CarriedBy]; !declared {
				add("resource %q says its identifier is carried by %q, which is not a field on it", name, resource.ID.CarriedBy)
			}

			// Carried by something means not sent under its own name. A
			// resource that sends both is not describing this shape, and the
			// declaration would read as though it were.
			if resource.ID.Field != "-" {
				add("resource %q says its identifier is carried by %q and also sends it as %q; carried_by is for the providers that send it in one place only",
					name, resource.ID.CarriedBy, identifierFieldName(resource.ID.Field))
			}
		}

		if !contains(validIDTypes, resource.ID.Type) {
			add("resource %q has id.type %q, which must be one of %s", name, resource.ID.Type, strings.Join(validIDTypes[1:], ", "))
		}

		// A number on the wire has to be a number in the store too, or the
		// response disagrees with the declaration and the disagreement shows
		// up as a string where a number was promised. The styles that mint
		// something with letters in it cannot be numbers, and saying so here
		// is cheaper than finding out from a response.
		if resource.ID.Type == "number" && !contains([]string{"numeric", "digits"}, resource.ID.Style) {
			add("resource %q says its identifier is a number and mints it with the %s style, which does not produce one",
				name, styleName(resource.ID.Style))
		}

		// digits exists for identifiers that are numeric and must not be
		// parsed as numbers, because they exceed what a JavaScript number can
		// hold. Declaring one a number is asking for the rounding bug the
		// style was added to prevent.
		if resource.ID.Type == "number" && resource.ID.Style == "digits" {
			add("resource %q mints a long numeric string and declares it a number, which is the rounding bug the digits style exists to avoid",
				name)
		}

		// A prefix is required for the prefixed style precisely because
		// applications prefix-match identifiers. Providers whose ids carry no
		// prefix say so by declaring the opaque style, rather than leaving it
		// ambiguous whether the omission was deliberate.
		if (resource.ID.Style == "" || resource.ID.Style == "prefixed") && resource.ID.Prefix == "" {
			add("resource %q must declare an id.prefix, or use id.style opaque if the provider's identifiers have none", name)
		}

		// A resource with no fields is almost always an unfinished one, and
		// the exception is narrow enough to name: a receipt. Knock's trigger
		// answers with a workflow_run_id and nothing else, and describing
		// that as a resource with one invented field would put a property on
		// the wire the provider never sends. A receipt is created and never
		// read back, so if anything reads this resource it is not one.
		if len(resource.Fields) == 0 && len(resource.Constants) == 0 && !onlyCreated(r, name) {
			add("resource %q has no fields", name)
		}

		for _, field := range sortedKeys(resource.Fields) {
			spec := resource.Fields[field]

			if spec.As != "" && spec.In == "" {
				add("resource %q field %q sets as without in, and a top-level field is already named by its key", name, field)
			}

			// A field that never reaches the wire cannot be renamed on the
			// way there, and cannot be null there either.
			if spec.In == "-" {
				if spec.As != "" {
					add("resource %q field %q is not sent and also declares as, so it names a wire field that does not exist", name, field)
				}

				if spec.NullWhenUnset {
					add("resource %q field %q is not sent and also declares null_when_unset, which is a claim about a key nobody receives", name, field)
				}
			}

			if spec.NullWhenUnset && spec.Default != nil {
				add("resource %q field %q is null when unset and also has a default, so it is never unset", name, field)
			}

			if spec.NullWhenUnset && spec.Required {
				add("resource %q field %q is null when unset and also required, so it is never unset", name, field)
			}

			// Every segment of an in: path has to be something the runtime can
			// walk. A segment it cannot parse became a literal key: a field
			// nested under "to[0]" was emitted under a property actually
			// spelled "to[0]", which is a shape no provider sends, produced in
			// silence, and invisible to a conformance suite with no case
			// naming it. That has now happened three times in different
			// disguises, so it is a rule rather than a habit.
			for _, segment := range splitNonEmpty(spec.In) {
				// "-" is the whole value rather than a path: the field is not
				// sent, so there is no segment to be well formed.
				if spec.In == "-" {
					continue
				}

				if plainSegment.MatchString(segment) || indexedSegment.MatchString(segment) {
					continue
				}

				add("resource %q field %q nests under %q, and the segment %q is not a name or a name with an array index, so it would be sent as a property literally spelled that",
					name, field, spec.In, segment)
			}

			// A field named parent_thing nested under parent emits
			// parent.parent_thing, which no provider sends. It happened
			// thirty-one times across six Recipes before anything checked,
			// and every conformance case about them passed because they
			// asserted the shape the emulator produced. A provider that
			// really does repeat the parent says so with an explicit as.
			// A field named the same as its parent emits parent.parent, which
			// is usually an accident: the author meant the field to be the
			// parent object itself and wrote in: rather than leaving it off.
			// Four Recipes were written with that shape in one sitting.
			//
			// It is not always an accident. ClickUp really does send a status
			// object with a status inside it, so an explicit as: says the
			// provider means it, exactly as the sibling rule below allows.
			if spec.In != "" && spec.As == "" && strings.EqualFold(field, spec.In) {
				add("resource %q field %q nests under %q and would be sent as %s.%s; drop the in: if the field is the parent object, or set as: %s if the provider really sends that",
					name, field, spec.In, spec.In, field, field)
			}

			if spec.In != "" && spec.As == "" && strings.HasPrefix(strings.ToLower(field), strings.ToLower(spec.In)+"_") {
				add("resource %q field %q nests under %q and would be sent as %s.%s, which repeats the parent; set as: %s, or set as explicitly if the provider really sends that",
					name, field, spec.In, spec.In, field, field[len(spec.In)+1:])
			}

			if spec.Type == "list" && spec.Default != nil {
				if _, ok := spec.Default.([]any); !ok {
					add("resource %q field %q is a list, but its default is not a sequence", name, field)
				}
			}

			if !contains(validFieldTypes, spec.Type) {
				add("resource %q field %q has type %q, which must be one of %s",
					name, field, spec.Type, strings.Join(validFieldTypes[1:], ", "))
			}

			// Only time fields are ever filled in automatically, so declaring
			// stamped on anything else does nothing and reads as though it
			// does something.
			if spec.Stamped != nil && !contains(timeFieldTypes, spec.Type) {
				add("resource %q field %q declares stamped, which only applies to %s",
					name, field, strings.Join(timeFieldTypes, ", "))
			}
		}
	}

	if len(r.Routes) == 0 {
		add("at least one route is required")
	}

	seen := map[string]bool{}

	for i, route := range r.Routes {
		where := fmt.Sprintf("route %d (%s %s)", i+1, route.Method, route.Path)

		if route.Method == "" {
			add("%s: method is required", where)
		} else if route.Method != strings.ToUpper(route.Method) {
			add("%s: method must be upper case", where)
		}

		if !strings.HasPrefix(route.Path, "/") {
			add("%s: path must start with /", where)
		}

		// Selects is part of the identity: a GraphQL Recipe is several routes
		// on one path, told apart by what the query asks for, and two of them
		// selecting the same field would still be the duplicate this rule
		// exists to catch. The headers a route matches on are part of it for
		// the same reason -- an AWS Recipe is every operation on POST /,
		// told apart by X-Amz-Target.
		key := route.Method + " " + route.Path + " " + route.Selects
		for _, name := range sortedKeys(route.MatchesHeader) {
			key += " " + name + "=" + route.MatchesHeader[name]
		}
		if seen[key] {
			add("%s: duplicate route", where)
		}
		seen[key] = true

		// A route that only ever fails reaches no operation and touches no
		// resource, so requiring either would be asking for a declaration that
		// cannot mean anything.
		if route.Error != "" {
			if _, ok := r.Errors[route.Error]; !ok {
				add("%s: unknown error %q", where, route.Error)
			}

			if route.Operation != "" || route.Resource != "" {
				add("%s: declares error %q and an operation or resource, and it can only ever do the first", where, route.Error)
			}

			continue
		}

		if route.Operation == "" {
			add("%s: operation is required", where)
		} else if !contains(validOperations, route.Operation) {
			add("%s: operation %q must be one of %s", where, route.Operation, strings.Join(validOperations, ", "))
		}

		if route.Resource == "" {
			add("%s: resource is required", where)
		} else if _, ok := r.Resources[route.Resource]; !ok {
			add("%s: unknown resource %q", where, route.Resource)
		}

		if !contains(validDeletedBody, route.DeletedBody) {
			add("%s: deleted_body %q must be one of %s", where, route.DeletedBody, strings.Join(validDeletedBody[1:], ", "))
		}

		if route.DeletedBody != "" && route.Operation != "delete" {
			add("%s: declares deleted_body on a %s, which never deletes anything", where, route.Operation)
		}

		// Offset and page numbering are only paging if the runtime reads the
		// parameter the caller sends, and there is no universal spelling: one
		// provider says per_page, another PageSize, another count. Without
		// both names the runtime reads "limit" and the style's own word, which
		// is right for some and wrong for plenty, and the wrongness is
		// invisible — the page size is ignored, one full page comes back, and
		// the caller's paging loop runs once and passes.
		//
		// 146 routes declared a style with no names, which was harmless while
		// nothing read the style and became a claim the moment something did.
		//
		// "-" is the third thing a name can be, beside a spelling and an
		// omission: the provider accepts no name for this at all. Datadog's
		// event listing fixes the page at a thousand and takes no size
		// parameter, so honouring an invented one would let code that sends
		// it work here and be ignored in production.
		if s := route.Pagination.Style; s == "offset" || s == "page" {
			if route.Pagination.LimitParam == "" || route.Pagination.CursorParam == "" {
				add("%s: declares %s paging without naming the provider's parameters, so the runtime would read `limit` and %q, which is a guess; set limit_param and cursor_param", where, s, s)
			}
		}

		if len(route.Returns) > 0 {
			if route.EmptyBody {
				add("%s: declares returns and empty_body, and there is no body for the fields to be in", where)
			}

			if resource, ok := r.Resources[route.Resource]; ok {
				for _, field := range route.Returns {
					// Constants count. They are stamped onto the record and
					// trimmed with everything else, so a route that answers
					// with one has to name it. Braze stamps "message":
					// "success" on every object, and leaving it out of a
					// trimmed response drops the only field a Braze client
					// checks.
					_, constant := resource.Constants[field]

					if _, declared := resource.Fields[field]; !declared && !constant && field != "id" {
						// The identifier is held as "id" whatever the provider
						// calls it on the wire, and returns names the record's
						// own keys rather than the rendered ones. Writing the
						// wire name here is the obvious mistake, so say so
						// instead of reporting an unknown field.
						if field == resource.ID.Field {
							add("%s: returns %q, which is what the identifier is called on the wire; name it \"id\", because the trim runs before the rename",
								where, field)

							continue
						}

						add("%s: returns %q, which is not a field on resource %q", where, field, route.Resource)
					}
				}
			}
		}

		for _, name := range route.Scope {
			if name == "id" {
				add("%s: id cannot be a scope parameter", where)
				continue
			}

			if !strings.Contains(route.Path, "{"+name+"}") {
				add("%s: scope %q does not appear in the path", where, name)
			}

			if resource, ok := r.Resources[route.Resource]; ok {
				if _, declared := resource.Fields[name]; !declared {
					add("%s: scope %q is not a field on resource %q", where, name, route.Resource)
				}
			}
		}

		if in := route.Pagination.In; in != "" && in != "query" && in != "body" {
			add("%s: pagination in %q must be query or body", where, in)
		}

		if !contains(validPagination, route.Pagination.Style) {
			add("%s: pagination.style %q is not supported", where, route.Pagination.Style)
		}

		if route.IDFrom != "" {
			source, name, ok := strings.Cut(route.IDFrom, ":")

			switch {
			case route.IDFrom == "auth":
				// Names no location, because there is nothing to read. The
				// only thing that can be wrong with it is a path that also
				// carries an identifier, which would mean the request does
				// carry one after all.
				if strings.Contains(route.Path, "{id}") {
					add("%s: id_from auth says the request carries no identifier, and the path carries one", where)
				}
			case !ok || name == "":
				add("%s: id_from %q must look like query:channel, body:channel or auth", where, route.IDFrom)
			case source != "query" && source != "body":
				add("%s: id_from source %q must be query, body or auth", where, source)
			case strings.Contains(route.Path, "{id}"):
				add("%s: id_from and an {id} path parameter cannot both apply", where)
			}
		}

		if route.Operation != "list" && route.Operation != "create" &&
			route.IDFrom == "" && !strings.Contains(route.Path, "{id}") {
			add("%s: a %s needs an {id} in the path or an id_from", where, route.Operation)
		}

		// A hidden identifier is defensible on a resource fetched one at a
		// time by a key that lives in the path. It is not defensible on a
		// collection: a page of records with nothing to tell them apart
		// cannot be what the provider sends, because nobody could address
		// the second one.
		for i, f := range route.Filters {
			fwhere := fmt.Sprintf("%s: filters[%d]", where, i)

			if route.Operation != "list" {
				add("%s: filters only narrow a listing, and this is a %s", fwhere, route.Operation)
			}

			if f.Param == "" || f.Field == "" {
				add("%s: needs both a param and a field", fwhere)

				continue
			}

			resource, ok := r.Resources[route.Resource]
			if !ok {
				continue
			}

			if _, declared := resource.Fields[f.Field]; !declared {
				add("%s: filters on %q, which is not a field on resource %q", fwhere, f.Field, route.Resource)

				continue
			}

			if f.All != "" && f.All == f.Default {
				add("%s: the value that turns the filter off is also its default, so it never applies", fwhere)
			}

			for value, options := range f.Values {
				if len(options) == 0 {
					add("%s: expands %q into nothing, so asking for it would hide every record", fwhere, value)
				}

				if f.All != "" && value == f.All {
					add("%s: expands %q, which is also the value that turns the filter off", fwhere, value)
				}
			}

			if f.Default == "" {
				continue
			}

			// A default that names a bucket has to be checked against the
			// bucket's members rather than against the word, because no
			// record's field ever holds the word.
			covered := map[string]bool{}

			if options, grouped := f.Values[f.Default]; grouped {
				for _, option := range options {
					covered[option] = true
				}
			} else {
				covered[f.Default] = true
			}

			// A default that excludes nothing is a default nobody can see. It
			// is the same rule as a scope: the behaviour is only described
			// when a fixture holds a record the filter must hide, because
			// otherwise the listing looks identical with the filter and
			// without it, and a mutation removing it passes.
			excluded := false

			for _, resources := range r.Fixtures {
				for _, record := range resources[route.Resource] {
					if value, present := record[f.Field]; present && !covered[fmt.Sprint(value)] {
						excluded = true
					}
				}
			}

			if !excluded {
				add("%s: defaults to %q and no fixture holds a %s the default would hide, so the listing looks the same with the filter and without it",
					fwhere, f.Default, route.Resource)
			}
		}

		// A receipt has no fields, so nothing filters the request on the way
		// in and the whole payload comes back out. That is the helpful kind
		// of wrong: a client reads its own request back and believes the
		// provider confirmed it, locally, for as long as the test suite is
		// the only thing calling.
		if route.Operation == "create" && len(route.Returns) == 0 &&
			len(r.Resources[route.Resource].Fields) == 0 {
			add("%s: resource %q declares no fields, so this create would echo the request body back; set returns to say what the provider actually sends",
				where, route.Resource)
		}

		for _, name := range route.Beside {
			bwhere := fmt.Sprintf("%s: beside %q", where, name)

			if route.Operation != "list" {
				add("%s: only a listing carries other collections, and this is a %s", bwhere, route.Operation)
			}

			if name == route.Resource {
				add("%s: is the route's own resource, so it would be written twice into the same body", bwhere)

				continue
			}

			beside, ok := r.Resources[name]
			if !ok {
				add("%s: is not a resource in this Recipe", bwhere)

				continue
			}

			// Without its own name it lands wherever the route's resource
			// landed, and one collection overwrites the other in silence.
			if beside.Collection == "" {
				add("%s: needs a collection name of its own, or it would overwrite the collection beside it", bwhere)

				continue
			}

			if beside.Collection == r.Resources[route.Resource].Collection {
				add("%s: names the same collection as %q, so one would overwrite the other", bwhere, route.Resource)
			}

			// The scope is applied to every collection in the body, so a
			// resource that cannot be partitioned by it would be returned
			// whole to whoever asked for one slice of it.
			for _, field := range route.Scope {
				if _, declared := beside.Fields[field]; !declared {
					add("%s: has no %q field, so the scope this route applies would not narrow it", bwhere, field)
				}
			}
		}

		// A hidden identifier on a listing is usually wrong, and the
		// exception is a positional array: a collection where position is
		// the identity because there is nothing else. Cohere's embeddings
		// are exactly that, and it is the trap rather than an oversight --
		// the only thing tying a vector to its input is the index, so a
		// client that filters or reorders inputs anywhere in the pipeline
		// pairs the wrong vector with the wrong document.
		//
		// The test is whether anything fetches one on its own. If a route
		// does, the identifier exists and hiding it from the listing is
		// inconsistent. If none does, position is all there is and saying so
		// is the honest description.
		if route.Operation == "list" && r.Resources[route.Resource].ID.Field == "-" &&
			r.Resources[route.Resource].ID.CarriedBy == "" && fetchedByID(r, route.Resource) {
			add("%s: resource %q hides its identifier and %s fetches one by id, so the listing withholds an identifier that exists; name the field that carries it with carried_by, or stop hiding it",
				where, route.Resource, fetchRoute(r, route.Resource))
		}
	}

	for i, c := range r.Conformance {
		if c.Arm != "" {
			if _, ok := r.Errors[c.Arm]; !ok {
				add("conformance[%d] %q: arms %q, which is not in the errors table", i, c.Name, c.Arm)
			}

			// Arming something and then expecting a status it does not
			// produce means the fault did nothing, and the case passes while
			// proving the opposite of what it says. That is the same class of
			// mistake as a case asserting its own request back.
			//
			// This used to compare against 400, on the reasoning that an
			// armed thing is a failure and a failure is a 4xx or a 5xx. Not
			// every armed thing is a failure: Snowflake's SQL API answers the
			// same endpoint with a 202 and a statement handle when the query
			// is slow, which is an alternate path rather than an error, and a
			// case arming it has to expect the 202 it installs. Comparing
			// against the armed entry's own status catches the real mistake
			// -- arm a 429, expect a 200 -- and permits the honest one.
			if armed, ok := r.Errors[c.Arm]; ok && c.Expect.Status > 0 && c.Expect.Status != armed.Status {
				add("conformance[%d] %q: arms %q, which answers %d, and expects %d, so what it installed changed nothing",
					i, c.Name, c.Arm, armed.Status, c.Expect.Status)
			}
		}

		if !c.Expect.NoBody {
			continue
		}

		if len(c.Expect.Body) > 0 || len(c.Expect.Matches) > 0 || len(c.Expect.Absent) > 0 {
			add("conformance[%d] %q: expect.no_body claims the response is empty, so there is nothing for body, matches or absent to look at", i, c.Name)
		}
	}

	// A field name the Recipe chooses is only a claim if a case asserts it
	// where the value exists. Asserting its absence on a last page holds
	// whatever the field happens to be called, so a Recipe could declare
	// cursor_field: next_page_token, be renamed to anything, and no case
	// would notice.
	//
	// Twenty-one of these were shipped before anything checked, across twenty
	// Recipes. Two of them turned out to have no successful list case at all:
	// every case touching the collection was checking a failure, so nothing
	// asserted what a working listing looks like.
	for _, declared := range []struct{ what, name string }{
		{"cursor_field", r.Responses.List.CursorField},
		{"count_field", r.Responses.List.CountField},
		{"has_more_field", r.Responses.List.HasMoreField},
		{"complete_field", r.Responses.List.CompleteField},
		{"final_field", r.Responses.List.FinalField},
	} {
		if declared.name == "" || assertsName(r.Conformance, declared.name) {
			continue
		}

		add("responses.list.%s is %q and no conformance case asserts that name where the value exists, so renaming it would break nothing",
			declared.what, declared.name)
	}

	if r.Responses.List.HasMoreField != "" && r.Responses.List.CompleteField != "" {
		add("responses.list declares both has_more_field and complete_field, which are the same flag with opposite senses")
	}

	// The whole point of a final field is that it is not the cursor. Naming
	// them the same thing would emit one key that changes meaning depending on
	// which page you are looking at, which is worse than either field alone.
	if f := r.Responses.List.FinalField; f != "" && f == r.Responses.List.CursorField {
		add("responses.list declares cursor_field and final_field as the same name %q, so one key would mean a page token on one page and a sync token on the next", f)
	}

	if r.Responses.Error.Key == "-" && r.Responses.Error.Style != "list" {
		add("responses.error.key is \"-\", which removes the envelope, but that only means anything for the list style")
	}

	if r.Responses.Error.Key == "-" && len(r.Responses.Error.Fields) > 0 {
		add("responses.error.key is \"-\", so there is no envelope for responses.error.fields to sit beside")
	}

	if r.Auth.Pattern != "" {
		if _, err := regexp.Compile(r.Auth.Pattern); err != nil {
			add("auth.pattern is not a valid regular expression: %v", err)
		}

		if len(r.Auth.Keys) > 0 {
			add("auth declares both keys and a pattern, and only one can decide whether a credential is accepted")
		}
	}

	if !contains(validCodeTypes, r.Responses.Error.CodeType) {
		add("responses.error.code_type %q must be one of %s", r.Responses.Error.CodeType, strings.Join(validCodeTypes[1:], ", "))
	}

	if !contains(validErrStyles, r.Responses.Error.Style) {
		add("responses.error.style %q must be one of %s", r.Responses.Error.Style, strings.Join(validErrStyles[1:], ", "))
	}

	if r.Responses.List.EntryStyle != "" && r.Responses.List.EntryStyle != "wrapped" {
		add("responses.list.entry_style %q must be wrapped", r.Responses.List.EntryStyle)
	}

	if !contains(validListStyles, r.Responses.List.Style) {
		add("responses.list.style %q must be one of %s", r.Responses.List.Style, strings.Join(validListStyles[1:], ", "))
	}

	if r.Responses.List.Style == "wrapped" && r.Responses.List.Key != "" {
		// The key is a fallback for resources that do not name their own
		// collection. When every one of them does, it is unreachable, and an
		// unreachable declaration reads as a description of the provider
		// while describing nothing. Five Recipes shipped with one before this
		// rule existed, each found by mutating it and watching nothing fail.
		covered := len(r.Resources) > 0

		for _, name := range sortedKeys(r.Resources) {
			if r.Resources[name].Collection == "" {
				covered = false
			}
		}

		if covered {
			add("responses.list.key is %q and every resource names its own collection, so nothing reads it",
				r.Responses.List.Key)
		}
	}

	if r.Responses.List.Style == "wrapped" && r.Responses.List.Key == "" {
		for _, name := range sortedKeys(r.Resources) {
			if r.Resources[name].Collection == "" {
				add("resource %q needs a collection name, or responses.list.key must be set, because the list style is wrapped", name)
			}
		}
	}

	if !contains(validSigning, r.Webhooks.Signing.Scheme) {
		add("webhooks.signing.scheme %q is not supported", r.Webhooks.Signing.Scheme)
	}

	if r.Webhooks.Signing.Scheme == "hmac-sha256" && r.Webhooks.Signing.Header == "" {
		add("webhooks.signing.header is required when signing is enabled")
	}

	for _, event := range r.Webhooks.Events {
		if strings.TrimSpace(event) == "" {
			add("webhooks.events contains an empty event name")
		}
	}

	for _, name := range sortedKeys(r.Errors) {
		e := r.Errors[name]

		if e.Status < 100 || e.Status > 599 {
			add("error %q has status %d, which is not a valid HTTP status", name, e.Status)
		}
	}

	for _, header := range sortedKeys(r.RequiredHeaders) {
		required := r.RequiredHeaders[header]

		if name := required.Error; name != "" {
			if _, ok := r.Errors[name]; !ok {
				add("required_headers[%s] names error %q, which is not declared", header, name)
			}
		}

		for _, method := range required.Methods {
			if !knownMethods[strings.ToUpper(method)] {
				add("required_headers[%s] limits to method %q, which is not an HTTP method", header, method)
			}
		}
	}

	for _, fixtureName := range sortedKeys(r.Fixtures) {
		for _, resourceName := range sortedKeys(r.Fixtures[fixtureName]) {
			resource, ok := r.Resources[resourceName]
			if !ok {
				add("fixture %q seeds unknown resource %q", fixtureName, resourceName)
				continue
			}

			// A route that reads its identifier from the credentials answers
			// with the one record there is. Two of them and the runtime picks
			// the first, which is the seeding order, which is not a decision
			// anybody made.
			if len(r.Fixtures[fixtureName][resourceName]) > 1 {
				for _, route := range r.Routes {
					if route.Resource == resourceName && route.IDFrom == "auth" {
						add("fixture %q seeds %d records of %q, and %s %s answers about whoever asked, so there is nothing to choose between them with",
							fixtureName, len(r.Fixtures[fixtureName][resourceName]), resourceName, route.Method, route.Path)

						break
					}
				}
			}

			// A list field declared as one and then seeded with a scalar would
			// serve the scalar and pass every case written about it, so the
			// declaration has to mean something.
			//
			// An explicit null is allowed, because several providers send one
			// and it is a distinct state from both an absent field and an
			// empty array. AssemblyAI sends utterances: null on a transcript
			// that succeeded without speaker labels, so a rule that refused
			// null would block a true claim.
			for i, record := range r.Fixtures[fixtureName][resourceName] {
				// The identifier, which is not among the declared fields and
				// so is checked on its own.
				idField := resource.ID.Field
				if idField == "" || idField == "-" {
					idField = "id"
				}

				// An identifier the provider never echoes cannot disagree
				// with anything a client can see, so its shape is free.
				// Cohere keys embeddings e1 and e2 to find them again and
				// emits neither.
				if seeded, present := record[idField]; present && resource.ID.Field != "-" {
					if why, ok := seededID(resource.ID, seeded); !ok {
						add("fixture %q record %d of %q seeds an id that %s, so seeded records and created ones would not have the same shape",
							fixtureName, i, resourceName, why)
					}
				}

				for _, field := range sortedKeys(record) {
					// A key the resource does not declare is dropped on the
					// way in, so a fixture can set it, a reader can believe
					// it, and the emulator never sends it. That silence is
					// the problem: it fooled a mutation into passing, which
					// means it could fool a conformance case the same way.
					//
					// id is always allowed, and so is whatever the Recipe
					// renamed it to, because the identifier is not declared
					// among the fields.
					if _, declared := resource.Fields[field]; !declared {
						if _, constant := resource.Constants[field]; !constant &&
							field != "id" && field != resource.ID.Field {
							add("fixture %q record %d of %q sets %q, which is not a field on that resource, so it would be dropped in silence",
								fixtureName, i, resourceName, field)
						}
					}

					if resource.Fields[field].Type != "list" || record[field] == nil {
						if seeded, ok := seededAs(resource.Fields[field].Type, record[field]); !ok {
							add("fixture %q record %d of %q declares %q as %s and seeds %s, and nothing coerces it, so %s is what goes on the wire",
								fixtureName, i, resourceName, field, resource.Fields[field].Type, seeded, seeded)
						}

						continue
					}

					if _, isList := record[field].([]any); !isList {
						add("fixture %q record %d of %q sets list field %q to something that is neither a sequence nor null",
							fixtureName, i, resourceName, field)
					}
				}
			}
		}
	}

	for i, c := range r.Conformance {
		where := fmt.Sprintf("conformance case %d", i+1)
		if c.Name != "" {
			where = fmt.Sprintf("conformance case %q", c.Name)
		} else {
			add("%s: name is required", where)
		}

		if c.Source == "" {
			add("%s: source is required, so a reader can check the claim against the provider", where)
		}

		if c.Verified != "" && !datePattern.MatchString(c.Verified) {
			add("%s: verified %q must be a date like 2026-08-15", where, c.Verified)
		}

		if c.Request.Method == "" {
			add("%s: request.method is required", where)
		} else if c.Request.Method != strings.ToUpper(c.Request.Method) {
			add("%s: request.method must be upper case", where)
		}

		if !strings.HasPrefix(c.Request.Path, "/") {
			add("%s: request.path must start with /", where)
		}

		if len(c.Request.Form) > 0 && len(c.Request.JSON) > 0 {
			add("%s: a request sends form or json, not both", where)
		}

		if c.Fixture != "" {
			if _, ok := r.Fixtures[c.Fixture]; !ok {
				add("%s: unknown fixture %q", where, c.Fixture)
			}
		}

		if c.Expect.Status == 0 {
			add("%s: expect.status is required", where)
		}

		for field, pattern := range c.Expect.Matches {
			if _, err := regexp.Compile(pattern); err != nil {
				add("%s: expect.matches[%s] is not a valid regular expression: %v", where, field, err)
			}
		}

		for header, pattern := range c.Expect.HeaderMatches {
			if _, err := regexp.Compile(pattern); err != nil {
				add("%s: expect.header_matches[%s] is not a valid regular expression: %v", where, header, err)
			}
		}

		if c.Expect.BodyMatches != "" {
			if _, err := regexp.Compile(c.Expect.BodyMatches); err != nil {
				add("%s: expect.body_matches is not a valid regular expression: %v", where, err)
			}
		}

		// An absence is a claim too: "another project's issues are not visible"
		// is exactly the kind of thing worth asserting, and it has no positive
		// half. Leaving it out of this check rejected real cases.
		// A case whose every body assertion repeats a value it sent proves
		// nothing about the emulator: it asserts its own request came back.
		// Four of these were written and only caught by deliberately breaking
		// the Recipe and noticing the case still passed, so the check is here
		// rather than in somebody's memory.
		//
		// One non-echoed claim is enough to clear it, because then something
		// the emulator decided is under test. matches and absent count, since
		// neither can be satisfied by echoing.
		if w := c.Expect.Webhook; w != nil {
			if w.None && (w.Event != "" || len(w.Body) > 0 || len(w.Matches) > 0 || len(w.Absent) > 0) {
				add("%s: expect.webhook.none claims nothing was emitted, so there is nothing for event, body, matches or absent to look at", where)
			}

			if !w.None && w.Event == "" && len(w.Body) == 0 && len(w.Matches) == 0 && len(w.Absent) == 0 {
				add("%s: expect.webhook asserts nothing, which is not the same as claiming none was emitted", where)
			}

			if w.Event != "" && !contains(r.Webhooks.Events, w.Event) {
				add("%s: expects webhook %q, which the Recipe does not declare", where, w.Event)
			}
		}

		if echoesOnly(c) {
			add("%s: every assertion repeats a value the request sent, so the case passes whatever the emulator does; assert something it decides, or add a matches or absent claim", where)
		}

		// A webhook claim counts. It is a claim about what the request caused,
		// which is evidence of exactly the kind this rule exists to demand.
		if c.Expect.Status < 400 && !c.Expect.NoBody && c.Expect.BodyMatches == "" && c.Expect.Webhook == nil &&
			len(c.Expect.Body) == 0 && len(c.Expect.Matches) == 0 &&
			len(c.Expect.Headers) == 0 && len(c.Expect.HeaderMatches) == 0 && len(c.Expect.Absent) == 0 &&
			len(c.Expect.AbsentHeaders) == 0 {
			add("%s: a case that asserts nothing about the response is not evidence of anything", where)
		}
	}

	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}

	return nil
}

// Verified reports how many conformance cases were observed against the real
// API, and how many rest on documentation alone. The distinction is the whole
// value of the suite, so it is reported rather than averaged away.
func (r *Recipe) Verified() (observed, documented int) {
	for _, c := range r.Conformance {
		if c.Verified != "" {
			observed++
			continue
		}

		documented++
	}

	return observed, documented
}

// Events returns the webhook event names this Recipe can emit.
func (r *Recipe) Events() []string {
	out := make([]string, len(r.Webhooks.Events))
	copy(out, r.Webhooks.Events)
	sort.Strings(out)

	return out
}

// fetchedByID reports whether any route fetches one record of a resource on
// its own, which is what makes an identifier something a client can use.
func fetchedByID(r *Recipe, resource string) bool {
	for _, route := range r.Routes {
		if route.Resource != resource {
			continue
		}

		if route.Operation == "get" || route.Operation == "update" || route.Operation == "delete" {
			return true
		}
	}

	return false
}

// fetchRoute names the first route that fetches one, for the message.
func fetchRoute(r *Recipe, resource string) string {
	for _, route := range r.Routes {
		if route.Resource != resource {
			continue
		}

		if route.Operation == "get" || route.Operation == "update" || route.Operation == "delete" {
			return route.Method + " " + route.Path
		}
	}

	return "another route"
}

// onlyCreated reports whether every route touching a resource creates it.
//
// Such a resource is a receipt: the provider hands back an identifier and
// nothing else, and there is nowhere to go and read it. Knock's trigger is
// one. A resource that is read anywhere is not, and a resource with no fields
// that is read anywhere is describing a response nobody can use.
func onlyCreated(r *Recipe, resource string) bool {
	touched := false

	for _, route := range r.Routes {
		if route.Resource != resource {
			continue
		}

		touched = true

		if route.Operation != "create" {
			return false
		}
	}

	return touched
}

// identifierFieldName renders the property an identifier goes out under, for a
// message, naming the default rather than reporting an empty string.
func identifierFieldName(field string) string {
	if field == "" {
		return "id"
	}

	return field
}

// styleName renders an id style for a message, naming the default rather than
// reporting an empty string nobody wrote.
func styleName(style string) string {
	if style == "" {
		return "prefixed"
	}

	return style
}

func contains(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if candidate == needle {
			return true
		}
	}

	return false
}

// sortedKeys keeps iteration deterministic, so validation problems always
// appear in the same order.
// seededAs reports whether a fixture value matches the type its field
// declares, and names what it is when it does not.
//
// A field's type is documentation. Nothing in the runtime reads it: the value
// a fixture seeds is the value that goes on the wire, with whatever type YAML
// gave it. So a Recipe can declare a field a string, seed it with something
// that is not one, and serve the wrong type for as long as no case looks.
//
// The reason to check it is YAML 1.1 rather than carelessness. An unquoted
// off, no, on or yes is a boolean, and those are all real values of real
// string fields: a DigitalOcean droplet's status is one of new, active, off
// and archive, and the fixture that seeded off to say a machine was powered
// down served false to everything that read it. The comment beside it
// explained that a cost report counting running machines would miss that
// droplet, and no case asserted on it, so the claim and the wire disagreed
// from the day it was written.
//
// Only the unambiguous directions are checked. A number field seeded with a
// whole number is fine, and a timestamp is a number whichever width it lands
// in.
func seededAs(declared string, value any) (string, bool) {
	if value == nil {
		return "", true
	}

	name := "something else"

	switch value.(type) {
	case bool:
		name = "a boolean"
	case string:
		name = "a string"
	case int, int64, uint64:
		name = "a whole number"
	case float64, float32:
		name = "a number"
	case []any:
		name = "a sequence"
	case map[string]any:
		name = "a mapping"
	}

	switch declared {
	case "string", "datetime":
		_, ok := value.(string)

		return name, ok
	case "boolean":
		_, ok := value.(bool)

		return name, ok
	case "integer", "number", "timestamp", "timestamp_ms":
		switch value.(type) {
		case int, int64, uint64, float64, float32:
			return name, true
		}

		return name, false
	}

	return name, true
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// ValidAuthSchemes returns the credential schemes a Recipe may declare.
//
// Exported so the runtime's test suite can assert that every scheme the
// validator accepts is one the handler actually checks. The two are separate
// pieces of code that have to agree and nothing else makes them: adding a
// scheme here without adding a case there would silently authorise every
// request against every Recipe using it.
func ValidAuthSchemes() []string {
	out := make([]string, len(validSchemes))
	copy(out, validSchemes)

	return out
}

// splitNonEmpty splits a dotted path, and returns nothing for an empty one so
// a field that does not nest is not asked to justify its parent.
func splitNonEmpty(path string) []string {
	if path == "" {
		return nil
	}

	return strings.Split(path, ".")
}

// seededID reports whether a fixture's explicit identifier has the shape the
// same resource would mint. It returns what is wrong with it, and true when
// nothing is.
//
// A fixture that seeds an id in a shape the generator would never produce puts
// two shapes in one collection: the seeded records look one way and anything
// created during a run looks another. Client code that parses or
// prefix-matches an id -- which is common enough that ID exists as a concept
// here at all -- then works against the fixtures and fails against the
// records it made itself.
//
// It was a surviving mutation that found this. Monday declares ten-digit ids
// and seeds ten-digit ids, and changing the declaration to six broke nothing:
// every case reads a seeded record, so the declared shape was never on the
// wire and could have drifted from the provider without a single case
// noticing.
func seededID(id ID, value any) (string, bool) {
	text, ok := value.(string)
	if !ok {
		// A numeric id is checked as a number elsewhere; nothing to say here.
		return "", true
	}

	if text == "" {
		return "is empty", false
	}

	switch id.Style {
	case "uuid":
		if !looksUUID(text) {
			return "is not a UUID", false
		}

	case "hex":
		body, ok := strings.CutPrefix(text, id.Prefix)
		if !ok {
			return "does not start with " + id.Prefix, false
		}

		if !all(body, isHexDigit) {
			return "is not hexadecimal", false
		}

		if id.Length > 0 && len(body) != id.Length {
			return fmt.Sprintf("is %d characters, not the declared %d", len(body), id.Length), false
		}

	case "digits":
		body, ok := strings.CutPrefix(text, id.Prefix)
		if !ok {
			return "does not start with " + id.Prefix, false
		}

		if !all(body, isDecimalDigit) {
			return "is not all digits", false
		}

		// A snowflake is a number written as text, and a number does not
		// begin with a zero.
		if body[0] == '0' {
			return "begins with a zero, which a number written as text does not", false
		}

		if id.Length > 0 && len(body) != id.Length {
			return fmt.Sprintf("is %d digits, not the declared %d", len(body), id.Length), false
		}

	case "numeric":
		if !all(text, isDecimalDigit) {
			return "is not a number", false
		}

	case "", "prefixed":
		if id.Prefix != "" && !strings.HasPrefix(text, id.Prefix) {
			return "does not start with " + id.Prefix, false
		}
	}

	return "", true
}

func looksUUID(text string) bool {
	runs := strings.Split(text, "-")
	if len(runs) != 5 {
		return false
	}

	for i, want := range []int{8, 4, 4, 4, 12} {
		if len(runs[i]) != want || !all(runs[i], isHexDigit) {
			return false
		}
	}

	return true
}

func all(text string, predicate func(byte) bool) bool {
	for i := 0; i < len(text); i++ {
		if !predicate(text[i]) {
			return false
		}
	}

	return true
}

func isDecimalDigit(b byte) bool { return b >= '0' && b <= '9' }

func isHexDigit(b byte) bool {
	return isDecimalDigit(b) || (b >= 'a' && b <= 'f')
}

// IsPath reports whether a name is a well-formed dotted path rather than a
// literal key that happens to contain a dot.
//
// The distinction is not academic. Dropbox names a field ".tag" -- the leading
// dot is part of the name, not a separator -- so treating every dotted name as
// a path turns it into an object under an empty key. A path is at least two
// segments and every one of them is a name.
func IsPath(name string) bool {
	segments := strings.Split(name, ".")
	if len(segments) < 2 {
		return false
	}

	for _, segment := range segments {
		if !plainSegment.MatchString(segment) && !indexedSegment.MatchString(segment) {
			return false
		}
	}

	return true
}

// FirstPageNumber is the number this provider gives its first page. One unless
// the Recipe says otherwise.
func (p Pagination) FirstPageNumber() int {
	if p.FirstPage == nil {
		return 1
	}

	return *p.FirstPage
}

// ListFor returns the list envelope a route answers with: the Recipe-wide one,
// with the route's own overrides applied.
//
// Empty means inherit and "-" means clear, so a route can both add a field the
// Recipe does not declare and remove one it does. A boolean can only be turned
// on, because an unset boolean and a false one are the same value in YAML and
// guessing which was meant is how a Recipe ends up asserting something nobody
// wrote.
func (r Recipe) ListFor(route Route) ListResponse {
	spec := r.Responses.List
	if route.List == nil {
		return spec
	}

	override := func(into *string, with string) {
		switch with {
		case "":
		case "-":
			*into = ""
		default:
			*into = with
		}
	}

	override(&spec.Style, route.List.Style)
	override(&spec.Key, route.List.Key)
	override(&spec.CursorField, route.List.CursorField)
	override(&spec.CountField, route.List.CountField)
	override(&spec.PageField, route.List.PageField)
	override(&spec.LimitField, route.List.LimitField)
	override(&spec.HasMoreField, route.List.HasMoreField)
	override(&spec.EntryStyle, route.List.EntryStyle)
	override(&spec.EntryField, route.List.EntryField)
	override(&spec.PagesField, route.List.PagesField)

	if route.List.LinkHeader {
		spec.LinkHeader = true
	}

	if route.List.CountAsString {
		spec.CountAsString = true
	}

	if route.List.OmitWhenEmpty {
		spec.OmitWhenEmpty = true
	}

	return spec
}

// GuessedPagination counts the routes whose paging the runtime has to guess
// at: a declared page size with neither a style nor a parameter name beside
// it.
//
// The runtime then reads "limit", which is right for some providers and wrong
// for plenty, and the wrongness is invisible -- the page size is ignored, one
// full page comes back, and the caller's paging loop runs once and passes. A
// Recipe in that state is making a claim nobody checked, and the point of
// counting them is that the number should be visible rather than buried.
//
// Naming either the style or the parameter is what marks a route as checked,
// because neither is a name anybody writes down by accident.
func (r Recipe) GuessedPagination() int {
	guessed := 0

	for _, route := range r.Routes {
		if route.Pagination.Limit == 0 {
			continue
		}

		if route.Pagination.Style == "" && route.Pagination.LimitParam == "" {
			guessed++
		}
	}

	return guessed
}
