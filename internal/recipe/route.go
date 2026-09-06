// Routes and paging: what a provider answers, at which path, a page at a time.

package recipe

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Route binds an HTTP method and path to an operation on a resource.
type Route struct {
	// Derived marks a route that came from the provider's own description
	// rather than from a recorded response, which Augment sets and a Recipe
	// file cannot. It is the difference between an observation and an
	// un-contradicted guess, and a user must never have to work out which
	// kind of answer they are looking at.
	Derived  bool   `yaml:"-"`
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
	// Raw is a body this route sends verbatim, for a provider that does not
	// send JSON. A route with one answers no resource and runs no
	// operation: it is the recorded bytes and nothing else.
	Raw *RawBody `yaml:"raw"`
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
	// Public marks a route the provider answers without a credential, on a
	// Recipe whose other routes need one.
	//
	// Auth is one setting for a whole Recipe, and nineteen Recipes here have
	// had to write down that their provider disagrees. Cognito serves a
	// genuinely public key set beside SigV4-gated operations. FireHydrant,
	// which hit this first, has one public route it could not model.
	// Checkly's runtimes list answers 200 with nothing presented. Ashby's job
	// board is public and the rest of Ashby is not. Chroma's whole retired
	// surface is public. Coinbase's public market data sits on the same hosts
	// as its private endpoints. Appsmith has exactly one.
	//
	// Every one of those had the same two choices: declare the route and
	// misrepresent it as needing a credential, or leave it out and describe a
	// provider as smaller than it is. Both are wrong in the direction that
	// matters, because a caller writing against the emulator would learn to
	// send a credential where none is wanted, or would not find the route at
	// all.
	//
	// This exempts the credential and nothing else. A required header is a
	// separate contract -- a provider can want its version header on a public
	// route, and Cognito does -- so those are still checked. Routing,
	// method resolution and the error table are unaffected.
	//
	// It reads two ways, because providers mean two different things by open.
	//
	//   public: true          the credential is not examined at all
	//   public: when-absent   exempt only a request that presented nothing
	//
	// The first is FireHydrant's, and it was verified rather than assumed: its
	// public route answers byte-identically with a junk bearer token attached
	// and with no header at all, so it genuinely ignores the credential.
	//
	// The second is football-data's, and the difference is not academic. Its
	// competitions list answers with nothing presented and refuses a wrong
	// token with a 400. Under the first reading a fake would wave that wrong
	// token through, which teaches a caller their bad credential is fine on
	// exactly the route they are most likely to try first. Under the second
	// the absent case is served, a valid key still works, and anything else is
	// judged as it would be anywhere else on the host.
	Public PublicMode `yaml:"public"`
	// Auth overrides the Recipe-wide credential for this route.
	//
	// Auth is one setting for a whole Recipe, and that is wrong for any
	// provider serving two surfaces from one Recipe. Twenty-four Recipes have
	// written the same paragraph: Coinbase's two hosts want different
	// credentials and no field records which route came from which; Mezmo's
	// two hosts resolve routing and the credential in opposite orders;
	// Healthchecks wanted to say "this route needs a key and this one, on a
	// different host, never does" and could only say it in prose.
	//
	// Public already covers the last of those, and only that one. This covers
	// the rest: a route may name its own scheme, header, keys and ordering,
	// and anything it does not name it inherits.
	//
	// Inherits field by field, which is worth stating because it did not at
	// first: this replaced the Recipe's credential wholesale while the comment
	// beside it said otherwise. The two readings agree for every Recipe that
	// scopes a credential to a route today, because all of them write the
	// whole thing out, and they disagree in the worst available direction the
	// moment one does not. An Auth with no scheme accepts every request, so a
	// route naming only its own refusal sentence would have switched its own
	// credential off while reading as though it had tightened it.
	//
	// The prefix is the one field where empty does not mean inherit, and the
	// validator refuses a route that leaves it unsaid while moving to another
	// scheme. A prefix belongs to a carrier: "Bearer " means nothing to an
	// X-Api-Key header, and inheriting it there refuses credentials the
	// provider accepts. Write prefix: "-" for a carrier taking the bare
	// secret, which is what clear already means for a list key and an
	// identifier field.
	//
	// Nil means inherit, which is every route in every Recipe shipping today,
	// so no existing behaviour and no existing fingerprint moves. That is the
	// same bargain List and Envelope already make one level down.
	//
	// One half of those twenty-four this does not fix, and the limit is worth
	// knowing before reaching for it. A route's AfterRouting governs requests
	// that reach that route, and ordering is mostly a claim about what happens
	// when routing FAILS -- an unrouted path, a wrong method. Neither of those
	// matches a route, so neither has a route to take an ordering from, and
	// the Recipe's own decides. Mezmo's two hosts disagreeing about that stays
	// in prose.
	//
	// Fixing that would mean deciding which route a request that matched none
	// of them should have been judged by, which is a guess. Recording it is
	// better than guessing.
	Auth *Auth `yaml:"auth"`
}

// PublicMode is how far a route's credential exemption reaches.
type PublicMode struct {
	// Always exempts every request from the credential check.
	Always bool
	// WhenAbsent exempts only a request that presented no credential.
	WhenAbsent bool
}

// Exempts reports whether this mode excuses a request that reached the given
// verdict.
func (p PublicMode) Exempts(v Verdict) bool {
	switch {
	case p.Always:
		return true
	case p.WhenAbsent:
		return v == Absent
	}

	return false
}

// Declared reports whether the route claims any exemption at all.
func (p PublicMode) Declared() bool {
	return p.Always || p.WhenAbsent
}

// UnmarshalYAML accepts the bare boolean and the named mode.
//
// The boolean came first and four Recipes ship it, so it keeps meaning what it
// meant. Anything else has to be a mode this understands: a typo silently
// reading as "not public" would turn a declared exemption into a gated route,
// which fails as a 401 on a path the provider serves to anyone.
func (p *PublicMode) UnmarshalYAML(value *yaml.Node) error {
	var flag bool
	if err := value.Decode(&flag); err == nil {
		p.Always = flag

		return nil
	}

	var named string
	if err := value.Decode(&named); err != nil {
		return err
	}

	switch named {
	case "always":
		p.Always = true
	case "when-absent":
		p.WhenAbsent = true
	default:
		return fmt.Errorf("route public %q must be true, always, or when-absent", named)
	}

	return nil
}

// Pagination describes how a list endpoint pages.
type Pagination struct {
	// Style is one of: cursor, offset, page, none.
	//
	// none says the provider serves the whole collection: no page size, no
	// position, no pointer to a next page. It exists because silence could not
	// say it. A listing declaring nothing is paged anyway -- ten records,
	// reading "limit" -- which is a claim about the provider that 344 routes
	// across 222 Recipes never made, and could not be told apart from nobody
	// having looked yet.
	//
	// Several providers' own descriptions settle it outright, declaring no
	// query parameters at all on a listing: OpenAI's /v1/models, xAI's,
	// Perplexity's, Supabase's projects, Upstash's Redis databases, Turso's
	// organizations, Redis Cloud's subscriptions. Reading that and writing
	// nothing down keeps the guess.
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
	// MayOvershoot says a page can carry more records than the size asked
	// for, for the providers whose page size is advice rather than a
	// contract.
	//
	// Missive documents it in one sentence -- "A page may return more
	// [items] than limit" -- on both of its listings. MaxLimit and OverLimit
	// describe the two ways a provider can serve *less* than was asked for,
	// and there was no way at all to say it serves more.
	//
	// It matters because of the loop everybody writes. `while len(page) ==
	// limit: fetch(next)` terminates on a page of 26 when the caller asked
	// for 25, which is the first page, so the walk stops after one page and
	// reports whatever it found as the whole collection. Nothing errors and
	// nothing is short -- the client received *more* data than it asked for
	// and concluded there was no more.
	//
	// Declared, this emulator serves one extra record on every page that is
	// not the last, so that loop breaks here rather than in production.
	//
	// This and MayUndershoot are not exclusive, though they were refused
	// together for the first few hours they existed. Modern Treasury settles
	// it in one sentence about one endpoint: "the actual number of records
	// returned may be less than, equal to, or more than the requested
	// amount." When both are set they take turns by position -- the first
	// page overshoots, every page after it undershoots -- because doing both
	// at once would cancel out and demonstrate neither.
	MayOvershoot bool `yaml:"may_overshoot"`
	// MayUndershoot says a page can carry fewer records than the size asked
	// for without being the last page, for the providers that fill a page
	// from a source that may not have enough ready.
	//
	// Onfleet is the reason. Its task listing "will return up to 64 tasks but
	// may return fewer", and it takes no size parameter at all, so a caller
	// cannot even ask for a number it could compare against.
	//
	// Two more turned up the same afternoon, written by different companies
	// for different purposes, which is the argument that this is a shape
	// rather than a quirk. Sumo Logic, on a log search: "the number of
	// messages returned may be less than the limit." Brave Search, on web
	// results: "the actual number of results returned may be less than
	// count."
	//
	// This is the same bug as MayOvershoot from the other side, and it is the
	// commoner spelling: `while len(page) == limit` and `if len(page) <
	// limit: break` both stop early, mid-collection, on a page that was
	// merely thin. What ends a walk at these providers is the cursor going
	// absent, and nothing else.
	//
	// Declared, this emulator serves one fewer record on every page that is
	// not the last, and never fewer than one -- a page of zero would end the
	// walk for a different and untrue reason.
	//
	// See MayOvershoot for what happens when both are declared.
	MayUndershoot bool `yaml:"may_undershoot"`
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
// Naming the style is not naming the parameters, and this counted it as though
// it were for months. Front declares cursor paging on three listings and names neither
// parameter, and Front's own description says it takes "limit" and
// "page_token" -- so the runtime read "cursor", Front ignored it, and every
// request came back on page one. Ninety-four routes across forty Recipes were
// in that state and neither this counter nor UnstatedPagination could see any
// of them: one required an empty style and the other required an empty
// everything.
func (r Recipe) GuessedPagination() int {
	guessed := 0

	for _, route := range r.Routes {
		p := route.Pagination

		// A listing that says it does not page has nothing to guess at, and
		// one that says nothing at all is counted by UnstatedPagination.
		if p.Style == "none" || (p.Limit == 0 && p.Style == "") {
			continue
		}

		// Either name left blank is a name the runtime supplies: "limit" for
		// the page size, and the style's own word for the position. "-" is not
		// blank -- it is a Recipe saying the provider accepts no name, which is
		// a decision rather than a gap.
		if p.LimitParam == "" || p.CursorParam == "" {
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
