// Response envelopes: the shapes a provider wraps its answers in.

package recipe

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
	// CursorNull sends the cursor field as null on the last page rather than
	// leaving it out.
	//
	// Absent and null are different on the wire, and for a paging loop the
	// difference decides whether it terminates. Metronome's customer listing
	// declares next_page required and nullable and its own example shows
	// "next_page": null on the last page; Notion's next_cursor has the same
	// shape. A loop written as `while (body.next_page !== undefined)` stops
	// against a provider that omits the key and runs for ever against one
	// that nulls it, and a loop written the other way round fails in exactly
	// the opposite circumstances.
	//
	// Omitting stays the default, for the same reason CursorField is opt in
	// at all: sending a field the provider does not send is the more
	// dangerous of the two mistakes.
	CursorNull bool `yaml:"cursor_null"`
	// CursorURL says the cursor field carries an address rather than a token,
	// and which kind: "absolute" for a whole URL, "path" for the path and
	// query alone.
	//
	// The difference is not cosmetic. Salesforce sends a path because its
	// clients join it to the instance URL they authenticated against, so an
	// absolute address would be joined to that and produce nonsense -- the
	// same concatenation bug as a token, arrived at from the other side.
	//
	// Eight Recipes describe their paging pointer as a full URL and emitted an
	// opaque cursor, so the fake taught the mistake the Recipe warned about.
	// Merge's own comment put it exactly: "both full URLs rather than opaque
	// cursors ... a client concatenating a base URL to next builds a URL that
	// does not exist" -- and a client written against a token does precisely
	// that concatenation.
	//
	// The URL is this request with its position moved on, which is what the
	// Link header already renders, so a Recipe saying so gets the same value
	// in its body.
	CursorURL string `yaml:"cursor_url"`
	// CountField names a property carrying how many records matched in total,
	// which is not the same as how many are on this page. Zendesk sends one and
	// a pagination UI cannot be built without it.
	CountField string `yaml:"count_field"`
	// CountMeans says what the count field counts, for the providers where it
	// is not how many records matched.
	//
	// Empty is the whole matching set, which is what every Recipe before this
	// assumed and what nearly every provider sends. Two other quantities
	// arrive under the same name:
	//
	// "page" is the length of the page in front of you. Shopware sends this
	// by default, from a field called total, because computing a real total
	// costs a second query and it does not run one unless asked. So a shop
	// with four hundred products answers a ten-record page with total: 10 --
	// a number that is not wrong about anything except the question it looks
	// like it is answering. A client that stops when it has read total
	// records reads one page; a client that divides total by the page size
	// finds one page; and neither errors, because ten really is a number of
	// products.
	//
	// "lookahead" is a bounded count: the provider fetches a few pages past
	// this one, counts what it found, and stops. Shopware's next-pages mode
	// reads limit * 6 + 1 rows and reports how many came back, so the same
	// four hundred products report 61. That is neither the page nor the
	// total, and it is the most misleading of the three, because it is large
	// enough to look real.
	//
	// The distinction cannot be left to the fixture. A fixture small enough
	// to fit on one page makes all three modes agree, which is exactly why a
	// Recipe can describe one and serve another and no case notices.
	CountMeans string `yaml:"count_means"`
	// CountLookahead is how many pages a lookahead count reaches, including
	// the one being served. Shopware's is six.
	//
	// The count is limit * CountLookahead + 1 where the collection is larger
	// than that, and the real total where it is not. The extra row is the
	// provider's sentinel: its presence is how the shop knows there is
	// anything beyond the window, and it is why the number ends in a 1 rather
	// than landing on a page boundary.
	CountLookahead int `yaml:"count_lookahead"`
	// PageCountField names a property carrying how many records are on this
	// page, for the providers that send that beside the total rather than
	// instead of it.
	//
	// count_means exists because Shopware sends one number and it is
	// sometimes the page; this exists because commercetools sends both, and
	// the format could previously say only one of them. Its listings carry
	// count -- "actual number of results returned" -- next to total, which is
	// the whole matching set and which its own description warns is an
	// estimate rather than a strongly consistent figure.
	//
	// Naming them apart is the honest half of the same problem. A provider
	// that sends both has told the caller which is which; the trap is only
	// there when one name has to carry both meanings.
	PageCountField string `yaml:"page_count_field"`
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
	// PrevLink adds a rel="prev" beside the next link, for the providers that
	// send one.
	//
	// Not implied by LinkHeader, because providers disagree and the
	// disagreement is the whole point of asking. GitHub's last page carries a
	// Link header holding rel="prev" and no next; Basecamp's own README
	// describes rel="next" alone, so its last page carries no header at all.
	// A client that stops when the header is missing works against Basecamp
	// and never terminates against GitHub.
	//
	// Only offset and page numbering can have one. A cursor names a position
	// the caller was handed and cannot be arithmetic'd backwards.
	PrevLink bool `yaml:"prev_link"`
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

// ListFor returns the list envelope a route answers with: the Recipe-wide one,
// with the route's own overrides applied.
//
// Empty means inherit and "-" means clear, so a route can both add a field the
// Recipe does not declare and remove one it does. A boolean can only be turned
// on, because an unset boolean and a false one are the same value in YAML and
// guessing which was meant is how a Recipe ends up asserting something nobody
// wrote.
// EnvelopeFor is how this route wraps a single object.
//
// The Recipe's own setting unless the route overrides it. Empty inherits and
// "-" clears, so a Recipe that wraps everything can say that one route does
// not -- which is the shape Datadog and Vercel both have, wrapping some of
// their resources and not others.
func (r Recipe) EnvelopeFor(route Route) ResourceResponse {
	spec := r.Responses.Resource
	if route.Envelope == nil {
		return spec
	}

	switch route.Envelope.Style {
	case "":
	case "-":
		spec.Style = ""
	default:
		spec.Style = route.Envelope.Style
	}

	switch route.Envelope.Key {
	case "":
	case "-":
		spec.Key = ""
	default:
		spec.Key = route.Envelope.Key
	}

	// A boolean can only be turned on, the same as everywhere else a route
	// narrows a Recipe-wide setting: there is no false to distinguish from
	// unset, and clearing the style is how a route says it wraps nothing.
	if route.Envelope.Array {
		spec.Array = true
	}

	return spec
}

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
	override(&spec.CountMeans, route.List.CountMeans)
	override(&spec.PageCountField, route.List.PageCountField)
	override(&spec.PageField, route.List.PageField)
	override(&spec.LimitField, route.List.LimitField)
	override(&spec.HasMoreField, route.List.HasMoreField)
	override(&spec.EntryStyle, route.List.EntryStyle)
	override(&spec.EntryField, route.List.EntryField)
	override(&spec.PagesField, route.List.PagesField)

	if route.List.LinkHeader {
		spec.LinkHeader = true
	}

	if route.List.PrevLink {
		spec.PrevLink = true
	}

	if route.List.CountLookahead > 0 {
		spec.CountLookahead = route.List.CountLookahead
	}

	if route.List.CountAsString {
		spec.CountAsString = true
	}

	if route.List.OmitWhenEmpty {
		spec.OmitWhenEmpty = true
	}

	if route.List.CursorNull {
		spec.CursorNull = true
	}

	// The five below were declared and dropped. A route-level list override
	// merged eighteen of the envelope's fields and silently ignored the rest,
	// so a Recipe saying collapse_single on one route had it read, validated
	// and thrown away -- and no conformance case could tell the difference,
	// because the emulator simply behaved as though the line were not there.
	//
	// This is the fifth declared-and-ignored field closed in this collection,
	// after a list key the collection name already supplied, auth.header on
	// the bearer scheme, a declared Content-Type overwritten by a default, and
	// an empty object expectation that asserted nothing. The shape is always
	// the same: the file says something, the runtime does something else, and
	// nothing in between notices.
	//
	// Open Library found it. /api/books answers one book as the object under
	// the caller's own query string and no books as {}, which is
	// collapse_single and omit_when_empty on one route -- and only the second
	// of those was reaching the runtime.
	if route.List.CollapseSingle {
		spec.CollapseSingle = true
	}

	override(&spec.CursorURL, route.List.CursorURL)
	override(&spec.FinalField, route.List.FinalField)
	override(&spec.CompleteField, route.List.CompleteField)

	for name, value := range route.List.Fields {
		if spec.Fields == nil {
			spec.Fields = map[string]any{}
		}

		spec.Fields[name] = value
	}

	return spec
}
