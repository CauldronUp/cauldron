// Webhooks: what a provider sends when nobody asked.

package recipe

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
	// any rather than a map, because a payload is not always an object.
	// SendGrid's Event Webhook posts an array and batches several events into
	// one delivery, HubSpot does the same, and QuickBooks sends an array of
	// change notifications -- three providers whose defining behaviour could
	// not be described while this was a map, and all three would have arrived
	// as a single object under Stripe's envelope.
	Payload any `yaml:"payload"`
}

// Signing describes webhook payload signing.
type Signing struct {
	// Scheme is one of: hmac-sha256, none.
	Scheme string `yaml:"scheme"`
	Header string `yaml:"header"`
	Secret string `yaml:"secret"`
	// Format is how the digest is wrapped and what string it is taken over.
	//
	// Every Recipe declaring hmac-sha256 used to get Stripe's shape --
	// t=<unix>,v1=<hex> over "<unix>.<body>" -- which is right for Stripe and
	// wrong for the other seventy-three. The reason it matters more than an
	// envelope does is that nobody parses a signature by hand: an application
	// passes the header to the provider's own SDK, and a signature in the
	// wrong shape fails that check every time. A fake that hands an
	// application a signature its verifier rejects is worse than one that
	// sends no signature at all, which is what the code here said while doing
	// the opposite.
	//
	// One of:
	//
	//   stripe        t=<unix>,v1=<hex>   over "<unix>.<body>"   (the default)
	//   hex           <hex>               over the body
	//   prefixed-hex  sha256=<hex>        over the body
	//   base64        <base64>            over the body
	//   v0-hex        v0=<hex>            over "v0:<unix>:<body>"
	//
	// Named for the shape rather than for a provider, because the shapes are
	// shared: Zoom's signature is Slack's, deliberately, and calling that one
	// "slack" would be the same mistake in miniature as giving every Recipe
	// Stripe's. stripe keeps its name only because no second provider here
	// wants it, and would lose it if one did.
	//
	// Empty means stripe, so a Recipe that has not been looked at keeps the
	// shape it had rather than silently changing.
	Format string `yaml:"format"`
	// TimestampHeader is the header the signed timestamp travels in, for the
	// providers whose signature covers a timestamp the value does not carry.
	//
	// Slack signs "v0:<ts>:<body>" and sends the timestamp in
	// X-Slack-Request-Timestamp; Zoom does the same in
	// x-zm-request-timestamp. Without it a verifier has the signature and no
	// way to reconstruct what was signed, so the delivery carries a value
	// that cannot be checked -- which is the failure the whole signing
	// surface exists to avoid, arrived at from a different direction.
	//
	// Stripe needs none of this because its timestamp is inside the value.
	TimestampHeader string `yaml:"timestamp_header"`
}
