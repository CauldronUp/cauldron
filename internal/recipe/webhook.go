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
	Payload map[string]any `yaml:"payload"`
}

// Signing describes webhook payload signing.
type Signing struct {
	// Scheme is one of: hmac-sha256, none.
	Scheme string `yaml:"scheme"`
	Header string `yaml:"header"`
	Secret string `yaml:"secret"`
}
