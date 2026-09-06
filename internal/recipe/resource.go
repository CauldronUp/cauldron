// Resources: what a provider has, what identifies it, and what it sends.

package recipe

import "strings"

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
	Alias string `yaml:"alias"`
	// VersionField names a field the provider keeps as an optimistic lock: a
	// number that moves on every write, which a caller has to quote back
	// before the provider will accept the next one.
	//
	// commercetools is the reason. Every resource it serves carries a
	// version, every update body is {version, actions} rather than a
	// document, and writing over a version that is not the current one is
	// refused with the current one in the reply so the retry can be
	// scripted. Without a way to say that, a Recipe describing such an API
	// serves an emulator that takes any write at all -- and the code written
	// against it passes every test, because a test suite is the one place
	// where nothing else is writing.
	//
	// That is the failure this is for. Ignoring the version is invisible
	// until two things touch one record at once, and then it is a silent
	// overwrite rather than an error: the later write wins and the earlier
	// one is gone, with nothing logged anywhere.
	VersionField string `yaml:"version_field"`
	// VersionConflict names the failure a stale write is refused with.
	VersionConflict string `yaml:"version_conflict"`
	// VersionMissing names the failure a write carrying no version at all is
	// refused with, for the providers that require one.
	//
	// Separate from VersionConflict because providers separate them, and
	// commercetools does: a stale version is a 409 that hands back the
	// current one, and an absent version is a 400 about a required field.
	// A client that retries on 409 and gives up on 400 needs them to be
	// different, and folding the two together here would teach it that every
	// rejected write is worth retrying.
	//
	// Empty lets a write with no version through, which is what a provider
	// that treats the field as optional does.
	VersionMissing string           `yaml:"version_missing"`
	Fields         map[string]Field `yaml:"fields"`
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
	// OtherPrefixes are prefixes a record may legitimately carry that this
	// Recipe does not mint, for providers whose identifiers do not all have
	// one shape.
	//
	// Auth0 is the reason. A user_id encodes the connection the user came
	// from: auth0|abc is a database user, google-oauth2|123 signed in with
	// Google, samlp|... came from an enterprise connection. Code that parses
	// the identifier assuming auth0| breaks on the first social login, which
	// is the first thing the Auth0 Recipe says. Its fixture holds a Google
	// user on purpose.
	//
	// Without this the fixture and the declaration have to disagree, and the
	// only ways to settle it are both lies: drop the social user, which
	// deletes the trap, or drop the prefix, which claims Auth0 mints bare
	// strings. Minting still uses Prefix -- there is one shape a new record
	// takes -- and these are the others a real account already contains.
	OtherPrefixes []string `yaml:"other_prefixes"`
	Length        int      `yaml:"length"`
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
	// Pattern is the shape an identifier has to have for the provider to look
	// it up at all, as a regular expression anchored by the Recipe.
	//
	// Declaring it says the provider checks before it searches, which is a
	// distinction the collection could not make: every absence was a 404.
	// Squarespace documents both answers on one route -- 404 "The requested
	// Order was not found" for an id that could exist and does not, and 400
	// "The id is not in the expected format" for one that could not -- and
	// Stripe, Intercom and everything else built on ObjectIds behave the same
	// way.
	//
	// It matters because the two failures are not interchangeable to the code
	// receiving them. A 404 is a fact about the account: the order was
	// deleted, or belongs to somebody else, and retrying will not help. A 400
	// is a fact about the caller: an id from the wrong provider, a truncated
	// string, an empty variable interpolated into the path. An emulator that
	// answers 404 to both teaches an application to treat its own bugs as
	// missing data -- and the test that proves the handler works asks for
	// "nonexistent", which is exactly the id that does not behave this way.
	//
	// Empty means the provider looks up whatever it is given, which is the
	// majority and stays the default.
	Pattern string `yaml:"pattern"`
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
//
// A field of type map is free-form: it accepts whatever keys the caller sends
// and answers with them. Stripe's metadata is the reason it exists -- arbitrary
// key-value pairs the provider stores and echoes without knowing what they
// mean -- and it has to be declared rather than assumed, because a create that
// echoes any field it is sent cannot tell a Recipe's model from a typo.
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

// UnassertedField counts the record field names no conformance case asserts.
//
// The format already refuses this for the two envelopes. responses.list.* and
// responses.error.* have to be asserted somewhere or the Recipe does not load,
// on the plain grounds that a name nothing checks is a name nobody would notice
// changing -- rename it and every case stays green.
//
// The records inside those envelopes had no such rule, and they are the larger
// half by an order of magnitude: the envelope has three or four names in it and
// a resource has thirty. They are also the half a client spends its time in. A
// Recipe can declare that a customer carries default_currency and delinquent and
// invoice_prefix, pass every case it has, and be describing three names it made
// up, because nothing ever looked inside the collection.
//
// Fields the wire never carries are not counted. in: "-" is how a Recipe says
// the record holds the field and the response does not send it, which is a
// statement about absence -- asking for evidence of it would be asking a case
// to assert something the Recipe has already said cannot appear.
//
// The name looked for is the one on the wire, so a field stored under one name
// and sent under another is credited by a case asserting what the response
// actually carries.
//
// There is already a rule that looks like this one and is not. A fixture may
// only set fields the resource declares, so renaming a declared field breaks
// every fixture holding it -- which guards that the field *exists*. It says
// nothing about where the field appears in the response. Nesting one under an
// object with `in:`, or renaming what it is sent as with `as:`, changes the
// path a client reads and leaves every case green. Tried on Strava: moving
// average_speed under a metrics object passed all fourteen cases, and failed
// the moment one of them asserted the path.
func (r Recipe) UnassertedField() int {
	unasserted := 0

	for what, resource := range r.Resources {
		for name, field := range resource.Fields {

			if field.In == "-" {
				continue
			}

			if r.assertsFieldOf(what, wirePath(name, field)) {
				continue
			}

			unasserted++
		}
	}

	return unasserted
}

// assertsFieldOf reports whether some successful case about this resource names
// the field.
//
// The first version of the counter asked whether any case in the file asserted
// the name anywhere, which is the right question for an envelope -- there is one
// per Recipe -- and the wrong one for a record. Nutritionix declares name on
// four resources and a case about one of them credited all four. Across the
// corpus, 476 names were credited by a case about something else entirely,
// which is the same over-generosity every other counter here has been caught in.
//
// A resource served by no route of its own is the exception. It appears only
// nested inside another's response -- Harvest's client, a time entry's task --
// so an assertion anywhere in the file is an assertion about it, and demanding
// a case on a route it does not have would be demanding the impossible.
//
// The case has to succeed, for the same reason a paging parameter's does: a
// refusal never reached the response, so it says nothing about its shape.
// wirePath is where a field actually sits in the response: the object it nests
// into, then the name it is sent under.
//
// The leaf alone is not enough once a field is renamed. Front's message carries
// author_id declared in: author and as: id, so the response has author.id and
// the message's own id is a different field at the top level. Looking for "id"
// credited author_id with every assertion of any id on that route -- which is
// precisely the collision that as: exists to make possible, so the field that
// most needs telling apart was the one this could not tell apart.
func wirePath(name string, field Field) string {
	wire := name
	if field.As != "" {
		wire = field.As
	}

	if field.In != "" && field.In != "-" {
		return field.In + "." + wire
	}

	return wire
}

func (r Recipe) assertsFieldOf(resource, wire string) bool {
	own := false

	for _, route := range r.Routes {
		if route.Resource == resource {
			own = true

			break
		}
	}

	for _, c := range r.Conformance {
		if c.Arm != "" {
			continue
		}

		if status := c.Expect.Status; status != 0 && (status < 200 || status >= 300) {
			continue
		}

		if own && !r.aboutResource(c, resource) {
			continue
		}

		if strings.Contains(wire, ".") {
			if assertsPath([]Case{c}, wire) {
				return true
			}

			continue
		}

		if assertsName([]Case{c}, wire) {
			return true
		}
	}

	return false
}

// aboutResource reports whether a case's request goes to a route serving this
// resource.
func (r Recipe) aboutResource(c Case, resource string) bool {
	for _, route := range r.Routes {
		if route.Resource != resource || c.Request.Method != route.Method {
			continue
		}

		if samePath(c.Request.Path, route.Path) {
			return true
		}
	}

	return false
}

// UnservedField counts the declared fields the emulator can never send.
//
// UnassertedField counts names with no evidence behind them. This is the harder
// version of the same question: names the fake will not produce at all. The
// Recipe says a customer carries default_currency, a fixture holds a customer,
// and the response omits the key entirely -- so code written against Cauldron
// reads undefined and finds out what the provider really sends the first time
// it talks to one. That is the failure this whole project exists to prevent,
// and it was going unmeasured while the weaker version of it was counted.
//
// The rule is the sandbox's own, read off the shaping it does. A declared field
// reaches the wire when a fixture record sets it, when it has a default, when
// it is null_when_unset, or when its type is one the sandbox stamps from the
// clock and the Recipe has not said stamped: false. Nothing else appears.
//
// Only resources some fixture actually holds are counted. A resource nothing
// seeds has no record for any of its fields to be missing from, and the listing
// counters already say so in plainer words.
//
// And a field some case asserts is not counted, whatever this rule thinks it
// knows: a green case asserting the name has seen it. Fourteen fields arrive
// that way, nearly all of them written by a request and echoed back -- Square's
// payment note, SingleStore's admin password -- so this is a strict subset of
// UnassertedField, which is what the report line says it is.
func (r Recipe) UnservedField() int {
	unserved := 0

	for what, resource := range r.Resources {
		set, held := r.fixtureFields(what)
		if !held {
			continue
		}

		for name, field := range resource.Fields {
			if field.In == "-" || set[name] {
				continue
			}

			if field.Default != nil || field.NullWhenUnset {
				continue
			}

			if stamps(field.Type) && (field.Stamped == nil || *field.Stamped) {
				continue
			}

			// A field some case asserts has been seen on the wire, whatever
			// this rule thinks. A write case is the usual way: Square's
			// payment note is sent in the request body and echoed back, so
			// the field reaches a client without any fixture setting it.
			// Counting those would be claiming the emulator cannot do
			// something a green case shows it doing.
			if r.assertsFieldOf(what, wirePath(name, field)) {
				continue
			}

			unserved++
		}
	}

	return unserved
}

// fixtureFields is every field name some fixture record of this resource sets,
// and whether any fixture holds one at all.
func (r Recipe) fixtureFields(resource string) (map[string]bool, bool) {
	set := map[string]bool{}
	held := false

	for _, fixture := range r.Fixtures {
		for _, record := range fixture[resource] {
			held = true

			for name := range record {
				set[name] = true
			}
		}
	}

	return set, held
}

// stamps reports whether the sandbox fills a field of this type from its clock
// when nothing else has.
func stamps(kind string) bool {
	switch kind {
	case "timestamp", "timestamp_ms", "datetime", "timestamp_ms_string", "msdate":
		return true
	}

	return false
}

// assertsPath reports whether a case asserts something at this nested path,
// which is the whole chain contiguously and not merely its last name.
func assertsPath(cases []Case, chain string) bool {
	want := strings.Split(chain, ".")

	for _, c := range cases {
		paths := make([]string, 0, len(c.Expect.Body)+len(c.Expect.Matches))
		for path := range c.Expect.Body {
			paths = append(paths, path)
		}

		for path := range c.Expect.Matches {
			paths = append(paths, path)
		}

		for _, path := range paths {
			segments := splitFieldPathForMatch(path)

			for start := 0; start+len(want) <= len(segments); start++ {
				same := true

				for i, name := range want {
					if segments[start+i] != name {
						same = false

						break
					}
				}

				if same {
					return true
				}
			}
		}
	}

	return false
}

// splitFieldPathForMatch splits an assertion path into names, dropping any
// [n] index so items[0].author.id reads as items, author, id.
func splitFieldPathForMatch(path string) []string {
	var out []string

	for _, segment := range strings.Split(path, ".") {
		if at := strings.IndexByte(segment, '['); at >= 0 {
			segment = segment[:at]
		}

		if segment != "" {
			out = append(out, segment)
		}
	}

	return out
}
