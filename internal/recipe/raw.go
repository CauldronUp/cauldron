// Answering with bytes rather than with a record.

package recipe

// RawBody is a response a route sends verbatim, for the providers that do not
// send JSON.
//
// This exists because refusing it was making Recipes lie about their own
// shape. arXiv answers Atom XML and nothing else -- not with a query
// parameter, not with an Accept header, not ever -- and the only mechanism
// that could write non-JSON bytes was the error table. So its Recipe declared
// thirteen routes, every one of them an error, including the eleven that
// succeed; the one resource it declared could never be routed to, because
// routing it would have rendered an XML feed as {"title": "..."}. It was the
// only Recipe of three hundred and fifteen with no routed resource, and it
// passed its whole suite while saying "error" about eleven things that work.
//
// Healthchecks arrived at the same place independently. Its ping answers
// 200 OK (not found) as text/plain, which is the normal, only behaviour of the
// endpoint, and it too had to be declared among the faults.
//
// The distinction this restores is worth stating plainly: an error table
// describes what a provider does when something is wrong, and a reader is
// entitled to read it that way. A raw body describes what a provider sends,
// full stop. Canada Post, PrestaShop and the older USPS generation are all
// waiting behind the same wall, and a client parses the response -- so for
// those providers, the parse is the integration and JSON is simply the wrong
// answer.
type RawBody struct {
	// ContentType is the media type to send.
	//
	// It has no default, and that is deliberate. Serving XML as
	// application/json is the failure mode NCBI's efetch demonstrates in
	// production -- a header that is simply false, with a client calling
	// .json() on the strength of it -- and a format that guessed here would
	// reproduce the bug rather than describe it.
	ContentType string `yaml:"content_type"`
	// Text is the body, byte for byte as it was recorded.
	Text string `yaml:"text"`
	// Empty says this route really does answer with nothing.
	//
	// Without it an empty Text is a Recipe that forgot to write one down, and
	// validation refuses it. A 204, or a ping that acknowledges and says
	// nothing, is a real answer and has to be distinguishable from an
	// omission -- the same reason routes already carry empty_body.
	Empty bool `yaml:"empty"`
}
