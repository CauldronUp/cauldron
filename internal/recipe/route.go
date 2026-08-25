// Routes and paging: what a provider answers, at which path, a page at a time.

package recipe

// Route binds an HTTP method and path to an operation on a resource.
type Route struct {
	Method   string `yaml:"method"`
	Path     string `yaml:"path"`
	Resource string `yaml:"resource"`
	// NotFound names the error to raise when this route is asked for a record
	// that is not there, instead of resource_missing.
	//
	// One provider can answer a missing thing two ways, and the npm registry
	// does. A package that does not exist answers {"error":"Not found"}; a
	// version that does not exist on a package that does answers the bare JSON
	// string "version not found: 99.99.99". Same status, same registry, and
	// nothing but the route to tell the emulator which is which.
	NotFound string `yaml:"not_found"`
	// Emits names the webhook event this route fires, for the providers whose
	// event names are not resource.created.
	//
	// The runtime emits resource.action and nothing else, so a Recipe declaring
	// Freshdesk's ticket_create, Bitbucket's repo:push or Zoom's meeting
	// .started -- all of them the provider's real names -- had those events
	// declared and never fired. Creating a record produced no webhook at all,
	// silently, and the only way to see one was to ask for it by hand.
	//
	// Naming the event here is what connects a change to the notification a
	// provider would actually send.
	//
	// Most of the collection is now wired. What is left declared and unfired
	// is mostly not wirable this way: an event like Freshdesk's
	// ticket_status_change or Recurly's renewed_subscription_notification is
	// not what an update route does, it is what one particular change to one
	// field does. EmitsWhen is for those.
	Emits string `yaml:"emits"`
	// EmitsWhen names events that fire only when a particular field changes,
	// rather than on every write.
	//
	// A great many declared events are this shape and no other: Freshdesk
	// sends ticket_status_change when a ticket's status moves and stays quiet
	// when its subject is edited, and ClickUp does the same for
	// taskStatusUpdated. Hanging those off the update route unconditionally
	// would be worse than leaving them unfired, because an application would
	// see the event on every edit here and on almost none in production --
	// the emulator would teach the handler to run when it will not.
	//
	// A list, because one route can owe several: the same Freshdesk update
	// answers for status and priority separately. It composes with Emits, so
	// a route may send an unconditional ticket_updated and a conditional
	// ticket_status_change from the same write, which is what Freshdesk does.
	EmitsWhen []ChangeEmit `yaml:"emits_when"`
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
	// SelectsBody does the same job as Selects and looks anywhere in the
	// request body rather than only in a GraphQL query.
	//
	// It is a separate field on purpose. Selects is fed the `query` property
	// of a GraphQL envelope and nothing else, and seven Recipes depend on
	// that narrowness: a marker word that today matches only inside a query
	// would start matching variable names and argument values if the search
	// were widened underneath them.
	//
	// What it is for is the providers whose response shape depends on what
	// the request asked for rather than on where it was sent. Gemini answers
	// a blocked prompt with a 200 and no candidates array at all, and a
	// permitted one with candidates and no block reason -- one path, one
	// method, two shapes, and nothing outside the body to tell them apart.
	// That was written up in the backlog as unservable before this existed.
	//
	// The match is the same whole-word one Selects makes, over the raw body
	// rather than a parsed field, because no emulator here can decide which
	// answer a model would have given and the marker is how a fixture says
	// which one to serve.
	SelectsBody string `yaml:"selects_body"`
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
	// MatchesQuery names query parameters whose values pick this route, for
	// the APIs where asking for more changes the shape of what comes back.
	//
	// Clover is the reason it exists. An order carries no line items unless
	// the request asks for them with ?expand=lineItems, so the same path
	// answers two different shapes and the compact one looks like an order
	// with nothing in it. Asana does the same with opt_fields, and its Recipe
	// says in as many words that Cauldron could not express it.
	//
	// A declared value matches when it appears among the comma-separated
	// members of the parameter, because that is what both providers send:
	// ?expand=lineItems,payments asks for two things and a route selecting
	// either one should answer. Equality is the single-member case of the
	// same rule.
	//
	// It is the third spelling of one idea -- selects reads the body,
	// matches_header reads a header, this reads the query string -- and a
	// route declaring any of them beats an equally-scoring route that
	// declares nothing, so a Recipe can have a compact fallback and an
	// expanded route above it.
	MatchesQuery map[string]string `yaml:"matches_query"`
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
	// Envelope overrides how this route wraps a single object, for the
	// providers that do not wrap every resource the same way.
	//
	// responses.resource is one setting for a whole Recipe, and two providers
	// have already needed it to be two. Datadog wraps a created event under
	// "event" with a status beside it and wraps nothing around a created
	// monitor. Vercel wraps a single domain under "domain" and wraps neither a
	// project nor a deployment. Both were written down as gaps rather than
	// modelled, because saying it for one resource said it for all of them.
	//
	// Empty inherits the Recipe's, and "-" clears it, the same way a route's
	// list override works.
	Envelope *ResourceResponse `yaml:"envelope"`
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
	//
	// Three more, each written because a Recipe was left wrong without them
	// and the note saying so is still in its file:
	//
	// "flagged" is the receipt without Stripe's discriminator, so a provider
	// that calls it something else can supply the name as a constant.
	// Intercom sends {type: contact, id, deleted: true}, and only the object
	// key was ever Stripe's.
	//
	// "id" is the identifier alone. Cloudflare answers {"result": {"id": ...}}
	// once its envelope is on.
	//
	// "empty" is an object with nothing in it. Asana answers {"data": {}},
	// which is not the same as no body at all: a client calling .json()
	// succeeds against one and throws on the other, which is exactly the kind
	// of difference this format exists to record.
	DeletedBody string `yaml:"deleted_body"`
	// DeletedKey names the key the identifier arrives under, for the "id"
	// body. Datadog answers a monitor delete with deleted_monitor_id rather
	// than id, so code reading response.deleted_monitor_id finds nothing
	// unless the key can be said.
	DeletedKey string `yaml:"deleted_key"`
	// IDAs renames the identifier on this route alone, for the providers that
	// call it one thing here and another everywhere else.
	//
	// Documenso answers a document create with documentId and sends id on
	// every other route. Modelling that as an ordinary create answered with a
	// document, which taught a one-step flow that does not exist -- so the
	// route was removed and the gap written down until this existed.
	IDAs string `yaml:"id_as"`
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
	// MaxLimit is the largest page a provider will serve, for the ones that
	// cap it and answer with less rather than refusing.
	//
	// Printify is the reason. Its own description says "default: 10, maximum:
	// 10", so a client asking for a hundred orders is answered with ten and
	// not told -- and a paging loop that stops when it receives fewer records
	// than it asked for stops on the first page. A shop with four hundred
	// orders reports ten, and nothing errored.
	//
	// Without this the declared Limit is only a default and the caller always
	// wins, so a Recipe could describe the cap in a comment and not serve it.
	// Zero means the provider serves whatever is asked for, which is what
	// every Recipe written before this assumed.
	MaxLimit int `yaml:"max_limit"`
	// OverLimit names the failure a route answers with when the caller asks
	// for a bigger page than MaxLimit, for the providers that refuse instead
	// of trimming.
	//
	// Both answers are common and they are not interchangeable. Printify
	// trims: ask for a hundred, receive ten, hear nothing. Shopware's entity
	// route refuses: ask for two hundred and fifty, receive
	// FRAMEWORK__QUERY_LIMIT_EXCEEDED and a 400 naming the ceiling. A client
	// written against one is broken against the other in opposite
	// directions -- one silently under-reads a collection, the other throws
	// on a request it thought was fine.
	//
	// Shopware is also why this belongs on the route rather than the Recipe.
	// It does both, on two listings of the same resource: /store-api/product
	// refuses and /store-api/product-listing/{categoryId} trims at the same
	// hundred, because the second is the storefront's own listing and runs
	// through a processor that calls min() on its way past.
	//
	// Empty keeps the trimming behaviour, which is what every Recipe written
	// before this assumed.
	OverLimit string `yaml:"over_limit"`
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

// FirstPageNumber is the number this provider gives its first page. One unless
// the Recipe says otherwise.
func (p Pagination) FirstPageNumber() int {
	if p.FirstPage == nil {
		return 1
	}

	return *p.FirstPage
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

// UnstatedPagination counts the listings that say nothing about paging at all.
//
// The runtime pages every listing: a route with no page size declared is given
// ten and reads "limit", exactly as a route declaring a size with no name is.
// GuessedPagination cannot see these, because it starts from a declared page
// size -- so the figure it reports has always been the smaller half of its own
// justification. Sixty routes page by a parameter nobody named; another
// hundred and eight page by a parameter and a page size nobody named.
//
// Nothing is truncated by it today, because no fixture behind one of these
// holds more than ten records. That is the reason it stayed invisible and not
// a reason it is fine: the claim is about the provider, and the fixture is not
// the provider. A listing the Recipe describes as unpaged answers at most ten
// and offers a cursor, and the first collection large enough to notice is not
// going to be one of ours.
//
// Counted apart rather than folded in, because they are not the same
// omission. One Recipe looked at paging and did not finish; the other has not
// looked.
func (r Recipe) UnstatedPagination() int {
	unstated := 0

	for _, route := range r.Routes {
		if route.Operation != "list" {
			continue
		}

		p := route.Pagination
		if p.Limit == 0 && p.Style == "" && p.LimitParam == "" {
			unstated++
		}
	}

	return unstated
}

// ChangeEmit is one conditional emission: an event, and the field whose change
// triggers it.
type ChangeEmit struct {
	Event string `yaml:"event"`
	// Field is compared before and after the write. Absent from the request
	// means unchanged, and a write that sets a field to the value it already
	// held is not a change either -- providers key these events off the
	// transition, not off the request naming the field.
	Field string `yaml:"field"`
}
