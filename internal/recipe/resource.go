// Resources: what a provider has, what identifies it, and what it sends.

package recipe

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
